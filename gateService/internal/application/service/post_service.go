package service

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"gateService/internal/domain/entity"
	"gateService/internal/domain/repository"
	"gateService/internal/infrastructure/config"
	"gateService/internal/interfaces/dto"
	"gateService/pkg/mq/nsqpool"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// PostServiceImpl 实现帖子相关的业务逻辑
type PostServiceImpl struct {
	storageConfig             *config.StorageConfig
	postRepository            repository.PostRepository            // 帖子仓储接口
	postTagRepository         repository.PostTagRepository         // 帖子标签仓储接口
	postTagRelationRepository repository.PostTagRelationRepository // 帖子标签关系仓储接口
	postCommentRepository     repository.PostCommentRepository     // 帖子评论仓储接口
	userRepository            repository.UserRepository            // 用户仓储接口
	producerPool              *nsqpool.ProducerPool
}

// NewPostServiceImpl 创建PostServiceImpl的新实例
// 参数:
// - postRepository: 帖子仓储接口,用于帖子的增删改查
// - postTagRepository: 帖子标签仓储接口,用于标签的管理
// - postTagRelationRepository: 帖子标签关系仓储接口,用于维护帖子和标签的关联
// - postCommentRepository: 帖子评论仓储接口,用于评论的管理
// - userRepository: 用户仓储接口,用于获取用户信息
// 返回:
// - *PostServiceImpl: 初始化后的PostServiceImpl实例
func NewPostServiceImpl(storageConfig *config.StorageConfig, postRepository repository.PostRepository, postTagRepository repository.PostTagRepository, postTagRelationRepository repository.PostTagRelationRepository, postCommentRepository repository.PostCommentRepository, userRepository repository.UserRepository, producerPool *nsqpool.ProducerPool) *PostServiceImpl {
	return &PostServiceImpl{
		storageConfig:             storageConfig,
		postRepository:            postRepository,
		postTagRepository:         postTagRepository,
		postTagRelationRepository: postTagRelationRepository,
		postCommentRepository:     postCommentRepository,
		userRepository:            userRepository,
		producerPool:              producerPool,
	}
}

// CreatePost 创建新帖子,包括帖子内容和标签
// 参数:
// - ctx: 上下文,用于传递请求上下文
// - post: 创建帖子请求DTO,包含帖子标题、内容、标签等信息
// 返回:
// - *dto.CreatePostResponse: 创建帖子响应DTO,包含帖子ID、创建时间等信息
// - error: 创建过程中的错误信息,如果创建成功则返回nil
func (s *PostServiceImpl) CreatePost(ctx context.Context, post *dto.CreatePostRequest) (*dto.CreatePostResponse, error) {
	// 构建帖子实体
	postEntity := s.buildPostEntity(post)

	tx, err := s.postRepository.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("开启事务失败: %v", err)
	}

	committed := false
	defer func() {
		if !committed {
			tx.Rollback()
		}
	}()

	// 创建帖子
	postID, err := s.createPostEntity(ctx, tx, postEntity)
	if err != nil {
		return nil, fmt.Errorf("创建帖子失败: %v", err)
	}

	// 创建标签
	tagIDs, err := s.createPostTags(ctx, tx, post.Tags)
	if err != nil {
		return nil, fmt.Errorf("创建标签失败: %v", err)
	}

	// 创建关联关系
	err = s.createPostTagRelations(ctx, tx, postID, tagIDs)
	if err != nil {
		return nil, fmt.Errorf("创建关联关系失败: %v", err)
	}

	// 更新分类帖子数量
	err = s.postRepository.UpdateCategoryCountTx(ctx, tx, post.CategoryID, 1)
	if err != nil {
		return nil, fmt.Errorf("更新分类帖子数量失败: %v", err)
	}

	// 处理图片
	if len(post.Images) > 0 {
		postImages := make([]*entity.PostImage, 0, len(post.Images))
		for _, base64Image := range post.Images {
			// 将Base64图片数据保存到文件系统或对象存储
			imageURL, err := s.savePostImage(postID, base64Image)
			if err != nil {
				return nil, fmt.Errorf("保存图片失败: %v", err)
			}
			postImages = append(postImages, &entity.PostImage{
				PostID:   postID,
				ImageURL: imageURL,
			})
		}
		// 批量插入图片记录
		err = s.postRepository.CreatePostImagesTx(ctx, tx, postImages)
		if err != nil {
			return nil, fmt.Errorf("处理图片失败: %v", err)
		}
	}

	// 提交事务
	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("提交事务失败: %v", err)
	}
	committed = true

	return &dto.CreatePostResponse{
		Code:    200,
		Message: "帖子创建成功",
	}, nil
}

// buildPostEntity 根据请求DTO构建帖子实体
// 参数:
// - post: 创建帖子请求DTO,包含帖子的基本信息
// 返回:
// - *entity.Post: 构建的帖子实体对象
func (s *PostServiceImpl) buildPostEntity(post *dto.CreatePostRequest) *entity.Post {
	return &entity.Post{
		UserID:     post.UserID,
		CategoryID: post.CategoryID,
		Title:      post.Title,
		Content:    post.Content,
	}
}

// SavePostImage 保存帖子图片
func (s *PostServiceImpl) savePostImage(postID int64, base64Image string) (string, error) {
	// 处理可能的 Data URI 格式
	base64Data := base64Image
	if len(base64Image) > 5 && base64Image[:5] == "data:" {
		// 分离 MIME 类型和 base64 数据
		parts := strings.Split(base64Image, ",")
		if len(parts) != 2 {
			return "", fmt.Errorf("无效的 Data URI 格式")
		}
		base64Data = parts[1]
	}

	// 解码base64图片数据
	imageData, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return "", fmt.Errorf("解码base64图片数据失败: %w", err)
	}

	// 检查图片大小
	if len(imageData) > s.storageConfig.Post.MaxSize {
		return "", fmt.Errorf("图片大小超过限制: %d > %d", len(imageData), s.storageConfig.Post.MaxSize)
	}

	// 检查图片类型
	contentType := http.DetectContentType(imageData)
	validType := false
	for _, allowedType := range s.storageConfig.Post.AllowedTypes {
		if contentType == allowedType {
			validType = true
			break
		}
	}
	if !validType {
		return "", fmt.Errorf("不支持的图片类型: %s", contentType)
	}

	// 生成文件名
	ext := ""
	switch contentType {
	case "image/jpeg":
		ext = ".jpg"
	case "image/png":
		ext = ".png"
	case "image/gif":
		ext = ".gif"
	}
	fileName := fmt.Sprintf("%d_%d%s", postID, time.Now().UnixNano(), ext)

	// 确保存储目录存在
	if err := os.MkdirAll(s.storageConfig.Post.Path, 0755); err != nil {
		return "", fmt.Errorf("创建存储目录失败: %w", err)
	}

	// 保存图片文件
	filePath := filepath.Join(s.storageConfig.Post.Path, fileName)
	if err := os.WriteFile(filePath, imageData, 0644); err != nil {
		return "", fmt.Errorf("保存图片文件失败: %w", err)
	}

	// 返回图片URL (使用 filepath.ToSlash 统一转换为前向斜杠)
	return filepath.ToSlash(filepath.Join(s.storageConfig.Post.URL, fileName)), nil
}

// createPostEntity 将帖子实体保存到数据库
// 参数:
// - ctx: 上下文,用于传递请求上下文
// - post: 要创建的帖子实体
// 返回:
// - int64: 创建成功后的帖子ID
// - error: 创建过程中的错误信息
func (s *PostServiceImpl) createPostEntity(ctx context.Context, tx *sql.Tx, post *entity.Post) (int64, error) {
	return s.postRepository.CreatePostTx(ctx, tx, post)
}

// createPostTags 创建帖子的标签
// 参数:
// - ctx: 上下文,用于传递请求上下文
// - tagNames: 标签名称列表,需要创建的标签名称
// 返回:
// - []int: 创建成功的标签ID列表
// - error: 创建过程中的错误信息
func (s *PostServiceImpl) createPostTags(ctx context.Context, tx *sql.Tx, tagNames []string) ([]int, error) {
	postTags := make([]*entity.PostTag, 0, len(tagNames))
	for _, tag := range tagNames {
		postTags = append(postTags, &entity.PostTag{Name: tag})
	}
	return s.postTagRepository.CreatePostTagTx(ctx, tx, postTags)
}

// createPostTagRelations 创建帖子和标签的关联关系
// 参数:
// - ctx: 上下文,用于传递请求上下文
// - postID: 帖子ID,要关联的帖子
// - tagIDs: 标签ID列表,要关联的标签
// 返回:
// - error: 创建关联关系过程中的错误信息
func (s *PostServiceImpl) createPostTagRelations(ctx context.Context, tx *sql.Tx, postID int64, tagIDs []int) error {
	postTagRelations := make([]*entity.PostTagRelation, 0, len(tagIDs))
	for _, tagID := range tagIDs {
		postTagRelations = append(postTagRelations, &entity.PostTagRelation{
			PostID: postID,
			TagID:  tagID,
		})
	}
	return s.postTagRelationRepository.CreatePostTagRelationTx(ctx, tx, postTagRelations)
}

// UpdatePost 更新帖子信息
// 参数:
// - ctx: 上下文,用于传递请求上下文
// - post: 要更新的帖子实体,包含更新后的信息
// 返回:
// - error: 更新过程中的错误信息
func (s *PostServiceImpl) UpdatePost(ctx context.Context, post *entity.Post) error {
	return s.postRepository.UpdatePost(ctx, post)
}

// DeletePost 删除指定帖子
// 参数:
// - ctx: 上下文,用于传递请求上下文
// - postID: 要删除的帖子ID
// 返回:
// - error: 删除过程中的错误信息
func (s *PostServiceImpl) DeletePost(ctx context.Context, postID int64) error {
	return s.postRepository.DeletePost(ctx, postID)
}

// GetPostByID 根据ID获取帖子详细信息
// 参数:
// - ctx: 上下文,用于传递请求上下文
// - postID: 要查询的帖子ID
// 返回:
// - *entity.Post: 查询到的帖子实体信息
// - error: 查询过程中的错误信息
func (s *PostServiceImpl) GetPostByID(ctx context.Context, postID int64) (*entity.Post, error) {
	return s.postRepository.GetPostByID(ctx, postID)
}

// GetPostCommentsByPostID 获取帖子的一级评论列表
// 参数:
// - ctx: 上下文,用于传递请求上下文
// - request: 获取评论请求,包含帖子ID、分页参数和当前用户ID等信息
// 返回:
// - *dto.GetCommentsResponse: 评论列表响应,包含评论列表、总页数等信息
// - error: 获取评论过程中的错误信息
func (s *PostServiceImpl) GetPostCommentsByPostID(ctx context.Context, request *dto.GetPostCommentsRequest) (*dto.GetPostCommentsResponse, error) {
	// 1. 获取评论基本数据
	// 根据帖子ID和分页参数获取评论列表
	comments, err := s.postCommentRepository.GetPostCommentsByPostID(ctx, request.ID, request.Page, request.PageSize)
	if err != nil {
		return nil, fmt.Errorf("获取帖子评论失败: %v", err)
	}

	// 2. 获取总页数
	// 用于前端分页展示
	totalPage, err := s.postCommentRepository.GetCommentTotalPage(ctx, request.ID, request.PageSize)
	if err != nil {
		return nil, fmt.Errorf("获取评论总页数失败: %v", err)
	}

	// 3. 获取评论作者信息
	// 3.1 提取所有评论的作者ID
	userIDs := make([]int, 0, len(comments))
	for _, comment := range comments {
		userIDs = append(userIDs, comment.UserID)
	}

	// 3.2 批量获取作者信息
	users, err := s.userRepository.GetUsersByIDs(ctx, &userIDs)
	if err != nil {
		return nil, fmt.Errorf("获取用户信息失败: %v", err)
	}

	// 4. 获取当前用户的点赞状态
	// 查询当前用户对这些评论的点赞记录
	likeMap, err := s.postCommentRepository.GetPostCommentLikesByUserID(ctx, request.UserID)
	if err != nil {
		return nil, fmt.Errorf("获取用户点赞状态失败: %v", err)
	}

	// 5. 组装评论数据
	commentItems := make([]*dto.CommentItem, 0, len(comments))
	for i, comment := range comments {
		// 构建每条评论的完整信息
		commentItems = append(commentItems, &dto.CommentItem{
			ID:      comment.CommentID, // 评论ID
			Content: comment.Content,   // 评论内容
			// 评论作者信息
			Author: dto.UserInfo{
				ID:        (*users)[i].UserID,
				Username:  (*users)[i].Username,
				AvatarURL: (*users)[i].AvatarURL,
			},
			CreatedAt: comment.CreatedAt,                  // 评论创建时间
			LikeCount: comment.LikeCount,                  // 点赞数
			IsLiked:   (*likeMap)[comment.CommentID] == 1, // 当前用户是否点赞
			Replies:   make([]dto.ReplyItem, 0),           // 初始化子评论列表
			ReplyNum:  comment.ReplyCount,                 // 子评论数量
		})
	}

	// 6. 返回响应数据
	return &dto.GetPostCommentsResponse{
		Code:      200,          // 成功状态码
		TotalPage: totalPage,    // 总页数
		Comments:  commentItems, // 评论列表
	}, nil
}

// GetPostCommentsByRootID 获取指定根评论的子评论列表
// 参数:
// - ctx: 上下文,用于传递请求上下文
// - request: 获取评论请求,包含根评论ID、分页参数和当前用户ID等信息
// 返回:
// - []*dto.ReplyItem: 子评论列表数据
// - error: 获取子评论过程中的错误信息
func (s *PostServiceImpl) GetPostCommentsByRootID(ctx context.Context, request *dto.GetPostCommentsRequest) ([]*dto.ReplyItem, error) {
	// 1. 从数据库获取指定根评论下的子评论列表
	comments, err := s.postCommentRepository.GetPostCommentsByRootID(ctx, request.RootID, request.Page, request.PageSize)
	if err != nil {
		return nil, fmt.Errorf("获取二级评论失败: %v", err)
	}

	// 2. 收集评论相关的用户ID
	// userIDs: 评论作者的用户ID列表
	// toUserIDs: 被回复用户的ID列表
	userIDs := make([]int, 0)
	toUserIDs := make([]int, 0)
	for _, comment := range comments {
		userIDs = append(userIDs, comment.UserID)
		toUserIDs = append(toUserIDs, *comment.ToUserID)
	}

	// 3. 批量获取用户信息
	// 3.1 获取评论作者信息
	users, err := s.userRepository.GetUsersByIDs(ctx, &userIDs)
	if err != nil {
		return nil, fmt.Errorf("获取用户信息失败: %v", err)
	}
	// 3.2 获取被回复用户信息
	toUsers, err := s.userRepository.GetUsersByIDs(ctx, &toUserIDs)
	if err != nil {
		return nil, fmt.Errorf("获取被回复用户信息失败: %v", err)
	}

	// 4. 获取当前用户对这些评论的点赞状态
	likeMap, err := s.postCommentRepository.GetPostCommentLikesByUserID(ctx, request.UserID)
	if err != nil {
		return nil, fmt.Errorf("获取用户点赞状态失败: %v", err)
	}

	// 5. 组装回复数据
	replys := make([]*dto.ReplyItem, len(comments))
	for i, comment := range comments {
		// 构建每条回复的完整信息
		replys[i] = &dto.ReplyItem{
			ID:      comment.CommentID, // 回复ID
			Content: comment.Content,   // 回复内容
			// 回复作者信息
			Author: dto.UserInfo{
				ID:        (*users)[i].UserID,
				Username:  (*users)[i].Username,
				AvatarURL: (*users)[i].AvatarURL,
			},
			CreatedAt: comment.CreatedAt,                  // 回复创建时间
			LikeCount: comment.LikeCount,                  // 点赞数
			IsLiked:   (*likeMap)[comment.CommentID] == 1, // 当前用户是否点赞
			// 被回复用户信息
			ReplyTo: dto.UserInfo{
				Username: (*toUsers)[i].Username,
			},
		}
	}

	return replys, nil
}

// SubmitComment 提交评论或回复
// 参数:
// - ctx: 上下文,用于传递请求上下文
// - request: 提交评论的请求参数,包含评论内容、帖子ID、评论ID等信息
// 返回:
// - *dto.SubmitPostCommentResponse: 提交评论的响应数据,包含状态码等信息
// - error: 提交评论过程中的错误信息
func (s *PostServiceImpl) SubmitComment(ctx context.Context, request *dto.SubmitPostCommentRequest) (*dto.SubmitPostCommentResponse, error) {
	comment := &entity.PostComment{
		PostID:    request.PostID,
		CommentID: request.CommentID,
		UserID:    request.UserID,
		Content:   request.Content,
		CreatedAt: request.CreatedAt,
		Level:     1,
		Status:    1,
	}

	requestJson, err := json.Marshal(comment)
	if err != nil {
		return nil, fmt.Errorf("序列化评论参数失败: %v", err)
	}
	err = s.producerPool.Publish(ctx, "comment_channel", requestJson)
	if err != nil {
		return nil, fmt.Errorf("发布评论到消息队列失败: %v", err)
	}
	return &dto.SubmitPostCommentResponse{
		Code: 200, // 200表示成功
	}, nil
}

// SubmitPostReply 提交帖子回复
// 参数:
// - ctx: 上下文,用于传递请求上下文
// - request: 提交回复的请求参数,包含回复内容、帖子ID、评论ID等信息
// 返回:
// - *dto.SubmitPostReplyResponse: 提交回复的响应数据,包含状态码等信息
// - error: 提交回复过程中的错误信息
func (s *PostServiceImpl) SubmitPostReply(ctx context.Context, request *dto.SubmitPostReplyRequest) (*dto.SubmitPostReplyResponse, error) {
	reply := &entity.PostComment{
		PostID:    request.PostID,
		CommentID: request.CommentID,
		UserID:    request.UserID,
		Content:   request.Content,
		CreatedAt: request.CreatedAt,
		ParentID:  request.ParentID,
		RootID:    request.RootID,
		ToUserID:  request.ToUserID,
		Level:     2,
		Status:    1,
	}

	requestJson, err := json.Marshal(reply)
	if err != nil {
		return nil, fmt.Errorf("序列化回复参数失败: %v", err)
	}
	err = s.producerPool.Publish(ctx, "comment_channel", requestJson)
	if err != nil {
		return nil, fmt.Errorf("发布回复到消息队列失败: %v", err)
	}
	return &dto.SubmitPostReplyResponse{
		Code: 200, // 200表示成功
	}, nil
}

// CommentLike 处理评论点赞/取消点赞的请求
func (s *PostServiceImpl) CommentLike(ctx context.Context, request *dto.CommentLikeRequest) (*dto.CommentLikeResponse, error) {
	// 构建点赞实体
	like := &entity.PostCommentLike{
		UserID:    request.UserID,
		CommentID: request.CommentID,
		Status:    1,
	}
	if !request.Status {
		like.Status = 0 // 取消点赞
	}

	// 序列化点赞数据
	requestJson, err := json.Marshal(like)
	if err != nil {
		return nil, fmt.Errorf("序列化点赞参数失败: %v", err)
	}

	// 发布点赞消息到消息队列
	err = s.producerPool.Publish(ctx, "comment_like_channel", requestJson)
	if err != nil {
		return nil, fmt.Errorf("发布点赞到消息队列失败: %v", err)
	}

	return &dto.CommentLikeResponse{
		Code: 200, // 200表示成功
	}, nil
}

// GetPostCategoryList 获取帖子分类列表
func (s *PostServiceImpl) GetPostCategoryList(ctx context.Context, request *dto.GetCategoryListRequest) (*dto.GetCategoryListResponse, error) {
	// 从数据库获取分类列表
	categories, err := s.postRepository.GetPostCategoryList(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取分类列表失败: %v", err)
	}

	// 构建响应数据
	categoryItems := make([]*dto.PostCategory, 0, len(categories))
	for _, category := range categories {
		categoryItems = append(categoryItems, &dto.PostCategory{
			ID:        category.CategoryID,
			Name:      category.Name,
			Icon:      category.Icon,
			PostCount: category.PostCount,
		})
	}

	return &dto.GetCategoryListResponse{
		Code:       200,           // 200表示成功
		Categories: categoryItems, // 分类列表
	}, nil
}

// GetPostsByCategoryID 获取指定分类的帖子列表
func (s *PostServiceImpl) GetPostsByCategoryID(ctx context.Context, request *dto.GetPostListRequest) (*dto.GetPostListResponse, error) {
	// 从数据库获取帖子列表
	posts, err := s.postRepository.GetPostsByCategoryID(ctx, request.CategoryID, request.Page, request.PageSize)
	if err != nil {
		return nil, fmt.Errorf("获取分类帖子列表失败: %v", err)
	}

	// 提取所有帖子的作者ID
	userIDs := make([]int, 0, len(posts))
	for _, post := range posts {
		userIDs = append(userIDs, post.UserID)
	}

	// 批量获取作者信息
	users, err := s.userRepository.GetUsersByIDs(ctx, &userIDs)
	if err != nil {
		return nil, fmt.Errorf("获取用户信息失败: %v", err)
	}

	// 获取当前用户点赞状态
	likeMap, err := s.postRepository.GetPostLikesByUserID(ctx, request.UserID)
	if err != nil {
		return nil, fmt.Errorf("获取用户点赞状态失败: %v", err)
	}

	// 获取当前用户收藏状态
	favoriteMap, err := s.postRepository.GetPostFavoritesByUserID(ctx, request.UserID)
	if err != nil {
		return nil, fmt.Errorf("获取用户收藏状态失败: %v", err)
	}

	// 构建响应数据
	postItems := make([]dto.PostResponse, 0, len(posts))
	pinnedPosts := make([]dto.PostResponse, 0)

	for i, post := range posts {
		// 获取作者信息
		author := (*users)[i]

		// 获取帖子标签
		tags, err := s.postTagRepository.GetPostTagByPostID(ctx, post.PostID)
		if err != nil {
			return nil, fmt.Errorf("获取帖子标签失败: %v", err)
		}

		// 构建帖子标签列表
		tagItems := make([]dto.PostTag, 0, len(tags))
		for _, tag := range tags {
			tagItems = append(tagItems, dto.PostTag{
				ID:   tag.TagID,
				Name: tag.Name,
			})
		}

		imageItems := make([]dto.PostImage, 0, len(post.Images))
		for _, image := range post.Images {
			imageItems = append(imageItems, dto.PostImage{
				ID:  image.ImageID,
				URL: image.ImageURL,
			})
		}

		// 构建帖子响应
		postResponse := dto.PostResponse{
			ID:      post.PostID,
			Title:   post.Title,
			Content: post.Content,
			Author: dto.PostAuthor{
				Username:  author.Username,
				AvatarURL: author.AvatarURL,
			},
			CreatedAt:     post.CreatedAt,
			Tags:          tagItems,
			IsPinned:      post.IsPinned,
			IsFeatured:    post.IsFeatured,
			ViewCount:     post.ViewCount,
			CommentCount:  post.CommentCount,
			LikeCount:     post.LikeCount,
			FavoriteCount: post.FavoriteCount,
			IsLiked:       (*likeMap)[post.PostID] == 1,
			IsFavorited:   (*favoriteMap)[post.PostID] == 1,
			Images:        imageItems,
		}

		// 区分置顶和普通帖子
		if post.IsPinned {
			pinnedPosts = append(pinnedPosts, postResponse)
		} else {
			postItems = append(postItems, postResponse)
		}
	}

	return &dto.GetPostListResponse{
		Code:        200,
		Posts:       postItems,
		PinnedPosts: pinnedPosts,
	}, nil
}

func (s *PostServiceImpl) GetPostByPostID(ctx context.Context, request *dto.GetPostByPostIDRequest) (*dto.GetPostByPostIDResponse, error) {
	// 获取帖子信息
	post, err := s.postRepository.GetPostByPostID(ctx, request.PostID)
	if err != nil {
		return nil, fmt.Errorf("获取帖子失败: %v", err)
	}

	// 获取作者信息
	author, err := s.userRepository.GetUserByID(ctx, post.UserID)
	if err != nil {
		return nil, fmt.Errorf("获取作者信息失败: %v", err)
	}

	// 获取帖子标签
	tags, err := s.postTagRepository.GetPostTagByPostID(ctx, post.PostID)
	if err != nil {
		return nil, fmt.Errorf("获取帖子标签失败: %v", err)
	}

	// 构建帖子标签列表
	tagItems := make([]dto.PostTag, 0, len(tags))
	for _, tag := range tags {
		tagItems = append(tagItems, dto.PostTag{
			ID:   tag.TagID,
			Name: tag.Name,
		})
	}

	// 获取当前用户点赞状态
	likeMap, err := s.postRepository.GetPostLikesByUserID(ctx, request.UserID)
	if err != nil {
		return nil, fmt.Errorf("获取用户点赞状态失败: %v", err)
	}

	// 获取当前用户收藏状态
	favoriteMap, err := s.postRepository.GetPostFavoritesByUserID(ctx, request.UserID)
	if err != nil {
		return nil, fmt.Errorf("获取用户收藏状态失败: %v", err)
	}

	imageItems := make([]dto.PostImage, 0, len(post.Images))
	for _, image := range post.Images {
		imageItems = append(imageItems, dto.PostImage{
			ID:  image.ImageID,
			URL: image.ImageURL,
		})
	}

	// 构建帖子响应
	postResponse := dto.PostResponse{
		ID:      post.PostID,
		Title:   post.Title,
		Content: post.Content,
		Author: dto.PostAuthor{
			ID:        author.UserID,
			Username:  author.Username,
			AvatarURL: author.AvatarURL,
		},
		CreatedAt:     post.CreatedAt,
		Tags:          tagItems,
		IsPinned:      post.IsPinned,
		IsFeatured:    post.IsFeatured,
		ViewCount:     post.ViewCount,
		CommentCount:  post.CommentCount,
		LikeCount:     post.LikeCount,
		FavoriteCount: post.FavoriteCount,
		IsLiked:       (*likeMap)[post.PostID] == 1,
		IsFavorited:   (*favoriteMap)[post.PostID] == 1,
		Images:        imageItems,
	}

	// 更新帖子浏览量
	err = s.postRepository.UpdateViewCount(ctx, request.PostID, 1)
	if err != nil {
		return nil, fmt.Errorf("更新帖子浏览量失败: %v", err)
	}

	return &dto.GetPostByPostIDResponse{
		Code: 200,
		Post: postResponse,
	}, nil
}

// PostLike 处理帖子点赞/取消点赞的请求
func (s *PostServiceImpl) PostLike(ctx context.Context, request *dto.PostLikeRequest) (*dto.PostLikeResponse, error) {
	// 开启事务
	tx, err := s.postRepository.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("开启事务失败: %v", err)
	}

	committed := false
	defer func() {
		if !committed {
			tx.Rollback()
		}
	}()

	// 创建点赞记录
	postLike := &entity.PostLike{
		UserID: request.UserID,
		PostID: request.PostID,
		Status: request.Status,
	}

	// 插入点赞记录
	err = s.postRepository.InsertPostLikeTx(ctx, tx, postLike)
	if err != nil {
		return nil, fmt.Errorf("插入点赞记录失败: %v", err)
	}

	status := 1
	if request.Status != 1 {
		status = -1
	}

	// 更新帖子点赞数
	err = s.postRepository.UpdateLikeCountTx(ctx, tx, request.PostID, status)
	if err != nil {
		return nil, fmt.Errorf("更新帖子点赞数失败: %v", err)
	}

	// 提交事务
	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("提交事务失败: %v", err)
	}
	committed = true

	return &dto.PostLikeResponse{
		Code: 200,
	}, nil
}

func (s *PostServiceImpl) PostFavorite(ctx context.Context, request *dto.PostFavoriteRequest) (*dto.PostFavoriteResponse, error) {
	// 开启事务
	tx, err := s.postRepository.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("开启事务失败: %v", err)
	}

	committed := false
	defer func() {
		if !committed {
			tx.Rollback()
		}
	}()

	// 创建收藏记录
	postFavorite := &entity.PostFavorite{
		UserID: request.UserID,
		PostID: request.PostID,
		Status: request.Status,
	}

	// 插入收藏记录
	err = s.postRepository.InsertPostFavoriteTx(ctx, tx, postFavorite)
	if err != nil {
		return nil, fmt.Errorf("插入收藏记录失败: %v", err)
	}

	status := 1
	if request.Status != 1 {
		status = -1
	}

	// 更新帖子收藏数
	err = s.postRepository.UpdateFavoriteCountTx(ctx, tx, request.PostID, status)
	if err != nil {
		return nil, fmt.Errorf("更新帖子收藏数失败: %v", err)
	}

	// 提交事务
	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("提交事务失败: %v", err)
	}
	committed = true

	return &dto.PostFavoriteResponse{
		Code: 200,
	}, nil
}

func (s *PostServiceImpl) GetRecentPosts(ctx context.Context, request *dto.GetRecentPostsRequest) (*dto.GetRecentPostsResponse, error) {
	// 获取最近帖子列表
	posts, err := s.postRepository.GetPostsByUserID(ctx, request.UserID, request.Page, request.PageSize)
	if err != nil {
		return nil, fmt.Errorf("获取最近帖子列表失败: %v", err)
	}

	// 构建响应数据
	postItems := make([]dto.PostBriefInfo, 0, len(posts))
	for _, post := range posts {
		postItems = append(postItems, dto.PostBriefInfo{
			ID:           post.PostID,
			Title:        post.Title,
			ViewCount:    post.ViewCount,
			CommentCount: post.CommentCount,
			LikeCount:    post.LikeCount,
			CreatedAt:    post.CreatedAt,
		})
	}

	return &dto.GetRecentPostsResponse{
		Code:  200,
		Posts: postItems,
	}, nil
}
