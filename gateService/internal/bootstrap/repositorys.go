package bootstrap

import (
	"gateService/internal/domain/repository"
	"gateService/internal/infrastructure/database"
)

// repositories 结构体包含所有数据访问层的仓储实例
type repositories struct {
	// UserRepo 用户仓储,处理用户数据的持久化
	UserRepo repository.UserRepository
	// PostRepo 帖子仓储,处理帖子内容的存储和查询
	PostRepo repository.PostRepository
	// PostTagRepo 帖子标签仓储,管理帖子标签数据
	PostTagRepo repository.PostTagRepository
	// PostTagRelationRepo 帖子-标签关系仓储,维护帖子和标签的多对多关系
	PostTagRelationRepo repository.PostTagRelationRepository
	// PostCommentRepo 帖子评论仓储,处理帖子评论数据
	PostCommentRepo repository.PostCommentRepository
	// ProgressRepo 进度仓储,记录和更新视频观看进度
	ProgressRepo repository.ProgressRepository
	// CommentRepo 评论仓储,处理评论数据,同时使用MySQL和Redis
	CommentRepo repository.CommentRepository
	// VideoRepo 视频仓储,管理视频元数据,使用MySQL存储基本信息,Redis缓存热点数据
	VideoRepo repository.VideoRepository
	// ProductRepo 商品仓储,处理商品信息的CRUD操作
	ProductRepo repository.ProductRepository
	// OrderRepo 订单仓储,管理订单相关数据
	OrderRepo repository.OrderRepository
}

// initRepositories 初始化所有仓储实例
// 参数:
//   - bases: 包含数据库连接等基础设施组件
//
// 返回:
//   - 初始化完成的repositories实例
func initRepositories(bases *bases) *repositories {
	return &repositories{
		// 初始化用户仓储,仅使用MySQL
		UserRepo: database.NewUserRepositoryImpl(bases.DB.GetDB()),
		// 初始化帖子仓储,仅使用MySQL
		PostRepo: database.NewPostRepositoryImpl(bases.DB.GetDB()),
		// 初始化帖子标签仓储,仅使用MySQL
		PostTagRepo: database.NewPostTagRepositoryImpl(bases.DB.GetDB()),
		// 初始化帖子-标签关系仓储,仅使用MySQL
		PostTagRelationRepo: database.NewPostTagRelationRepositoryImpl(bases.DB.GetDB()),
		// 初始化帖子评论仓储,仅使用MySQL
		PostCommentRepo: database.NewPostCommentRepositoryImpl(bases.DB.GetDB(), bases.RDB.GetRDB()),
		// 初始化进度仓储,仅使用MySQL
		ProgressRepo: database.NewProgressRepositoryImpl(bases.DB.GetDB()),
		// 初始化评论仓储,同时使用MySQL和Redis
		CommentRepo: database.NewCommentRepositoryImpl(bases.DB.GetDB(), bases.RDB.GetRDB()),
		// 初始化视频仓储,同时使用MySQL和Redis
		VideoRepo: database.NewVideoRepositoryImpl(bases.DB.GetDB(), bases.RDB.GetRDB()),
		// 初始化商品仓储,仅使用MySQL
		ProductRepo: database.NewProductRepositoryImpl(bases.DB.GetDB()),
		// 初始化订单仓储,仅使用MySQL
		OrderRepo: database.NewOrderRepositoryImpl(bases.DB.GetDB()),
	}
}
