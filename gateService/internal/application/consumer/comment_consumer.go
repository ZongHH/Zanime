package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"gateService/internal/domain/entity"
	"gateService/internal/domain/repository"
	"gateService/internal/infrastructure/middleware/websocket"
	"gateService/pkg/mq/nsqpool"
	"log"
	"strconv"

	"github.com/redis/go-redis/v9"
)

// CommentConsumer 负责处理评论消息的消费者
// 包含评论存储仓库和NSQ消费者池
type CommentConsumer struct {
	postRepository        repository.PostRepository        // 帖子持久化存储接口
	postCommentRepository repository.PostCommentRepository // 评论持久化存储接口
	userRepository        repository.UserRepository
	commentConsumerPool   *nsqpool.ConsumerPool // NSQ消费者池实例
	likeConsumerPool      *nsqpool.ConsumerPool // NSQ消费者池实例
	websocketManager      *websocket.Manager
}

// NewCommentConsumer 构造函数，创建新的评论消费者实例
// 参数: postCommentRepository - 评论存储仓库实现
func NewCommentConsumer(postRepository repository.PostRepository, postCommentRepository repository.PostCommentRepository, userRepository repository.UserRepository, websocketManager *websocket.Manager) *CommentConsumer {
	return &CommentConsumer{
		postRepository:        postRepository,
		postCommentRepository: postCommentRepository,
		userRepository:        userRepository,
		websocketManager:      websocketManager,
	}
}

// persistComment 持久化评论消息到数据库
// 参数:
//   - ctx: 上下文对象，用于传递截止时间和取消信号
//   - msg: 从消息队列接收的原始消息字节
//
// 返回值:
//   - error: 处理过程中遇到的错误，包含详细的错误信息
func (c *CommentConsumer) persistComment(ctx context.Context, msg []byte) error {
	// 反序列化消息到领域实体
	postComment := &entity.PostComment{}
	err := json.Unmarshal(msg, postComment)
	if err != nil {
		return fmt.Errorf("反序列化评论失败: %v", err)
	}

	if postComment.Level == 2 {
		parentID, err := c.postCommentRepository.GetCommentVirtualID(ctx, postComment.UserID, *postComment.ParentID)
		if err != nil && err != redis.Nil {
			return fmt.Errorf("获取虚拟父ID失败: %v", err)
		}
		if err != redis.Nil {
			postComment.ParentID = &parentID
		}
		rootID, err := c.postCommentRepository.GetCommentVirtualID(ctx, postComment.UserID, *postComment.RootID)
		if err != nil && err != redis.Nil {
			return fmt.Errorf("获取虚拟根ID失败: %v", err)
		}
		if err != redis.Nil {
			postComment.RootID = &rootID
		}
	}

	virtualID := postComment.CommentID

	// 开启事务
	tx, err := c.postCommentRepository.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("开启事务失败: %v", err)
	}

	committed := false
	defer func() {
		if !committed {
			tx.Rollback()
		}
	}()

	// 在事务中持久化评论
	err = c.postCommentRepository.CreatePostCommentTx(ctx, tx, postComment)
	if err != nil {
		return fmt.Errorf("持久化评论失败: %v", err)
	}

	// 在事务中更新帖子评论量
	err = c.postRepository.UpdateCommentCountTx(ctx, tx, postComment.PostID, 1)
	if err != nil {
		return fmt.Errorf("更新帖子评论量失败: %v", err)
	}

	if postComment.Level == 2 && postComment.ToUserID != nil {
		user, err := c.userRepository.GetUserByID(ctx, postComment.UserID)
		if err != nil {
			return fmt.Errorf("获取用户信息失败: %v", err)
		}

		// 在事务中更新回复数量
		err = c.postCommentRepository.UpdateReplyCountTx(ctx, tx, *postComment.RootID, 1)
		if err != nil {
			return fmt.Errorf("更新回复数量失败: %v", err)
		}

		// 发送websocket消息
		replyComment := &CommentMessage{
			MsgType:      "WS_MESSAGE",
			SendUserName: user.Username,
			Content:      postComment.Content,
		}
		replyCommentJson, err := json.Marshal(replyComment)
		if err != nil {
			return fmt.Errorf("序列化回复评论失败: %v", err)
		}
		toUserID := strconv.Itoa(*postComment.ToUserID)
		_ = c.websocketManager.SendMessage(toUserID, replyCommentJson)
	}

	err = c.postCommentRepository.SetCommentVirtualID(ctx, postComment.UserID, virtualID, postComment.CommentID)
	if err != nil {
		return fmt.Errorf("缓存虚拟ID失败: %v", err)
	}

	// 提交事务
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("提交事务失败: %v", err)
	}
	committed = true

	return nil
}

func (c *CommentConsumer) updateCommentLike(ctx context.Context, msg []byte) error {
	// 解析消息内容
	var commentLike entity.PostCommentLike
	err := json.Unmarshal(msg, &commentLike)
	if err != nil {
		return fmt.Errorf("解析点赞消息失败: %v", err)
	}

	commentID, err := c.postCommentRepository.GetCommentVirtualID(ctx, commentLike.UserID, commentLike.CommentID)
	if err != nil && err != redis.Nil {
		return fmt.Errorf("获取虚拟ID失败: %v", err)
	}
	if err != redis.Nil {
		commentLike.CommentID = commentID
	}

	tx, err := c.postCommentRepository.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("开启事务失败: %v", err)
	}

	committed := false
	defer func() {
		if !committed {
			tx.Rollback()
		}
	}()

	// 更新评论点赞状态
	err = c.postCommentRepository.InsertCommentLikeTx(ctx, tx, &commentLike)
	if err != nil {
		return fmt.Errorf("更新点赞状态失败: %v", err)
	}

	// 更新评论点赞数
	var count int
	if commentLike.Status == 1 {
		count = 1
	} else {
		count = -1
	}
	err = c.postCommentRepository.UpdateLikeCountTx(ctx, tx, commentLike.CommentID, count)
	if err != nil {
		return fmt.Errorf("更新点赞数量失败: %v", err)
	}

	// 提交事务
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("提交事务失败: %v", err)
	}
	committed = true

	return nil
}

// Start 启动评论消费者
// 初始化NSQ消费者池，注册回调函数，并启动消费
func (c *CommentConsumer) Start() {
	c.startCommentConsumer()
	c.startLikeConsumer()
}

func (c *CommentConsumer) startCommentConsumer() {
	consumerPool, err := nsqpool.NewConsumerPool(&nsqpool.ConsumerOptions{
		Topic:    "comment_channel", // 订阅的NSQ主题
		Channel:  "comment_channel", // 消费者通道名称
		PoolSize: 5,                 // 消费者池大小（并发处理数）
	})
	if err != nil {
		log.Fatalf("创建评论消费者池失败: %v\n", err)
	}
	c.commentConsumerPool = consumerPool

	consumerPool.RegisterCallback(c.persistComment)

	err = consumerPool.Start()
	if err != nil {
		log.Fatalf("启动评论消费者池失败: %v\n", err)
	}
}

func (c *CommentConsumer) startLikeConsumer() {
	consumerPool, err := nsqpool.NewConsumerPool(&nsqpool.ConsumerOptions{
		Topic:    "comment_like_channel", // 订阅的NSQ主题
		Channel:  "comment_like_channel", // 消费者通道名称
		PoolSize: 5,                      // 消费者池大小（并发处理数）
	})
	if err != nil {
		log.Fatalf("创建点赞消费者池失败: %v\n", err)
	}
	c.likeConsumerPool = consumerPool

	consumerPool.RegisterCallback(c.updateCommentLike)

	err = consumerPool.Start()
	if err != nil {
		log.Fatalf("启动点赞消费者池失败: %v\n", err)
	}
}

// Stop 停止评论消费者
// 优雅关闭NSQ消费者池，释放资源
func (c *CommentConsumer) Stop() {
	c.commentConsumerPool.Stop()
	c.likeConsumerPool.Stop()
}
