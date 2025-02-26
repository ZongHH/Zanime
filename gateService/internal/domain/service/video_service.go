package service

import (
	"context"
	"gateService/internal/interfaces/dto"
)

// VideoService 定义了视频服务的接口
type VideoService interface {
	// GetVideoInfo 获取视频详细信息
	// ctx: 上下文信息
	// request: 包含视频ID的请求参数
	// 返回: 视频详细信息响应和可能的错误
	GetVideoInfo(ctx context.Context, request *dto.GetVideoInfoRequest) (*dto.GetVideoInfoResponse, error)

	// GetVideoLibrary 获取视频库列表
	// ctx: 上下文信息
	// request: 包含地区、年份、类型、字母、分页等筛选条件的请求参数
	// 返回: 视频库列表响应和可能的错误
	GetVideoLibrary(ctx context.Context, request *dto.GetVideoLibraryRequest) (*dto.GetVideoLibraryRespnse, error)

	// GetVideoFilters 获取视频筛选条件
	// ctx: 上下文信息
	// request: 筛选条件请求参数
	// 返回: 视频筛选条件响应和可能的错误
	GetVideoFilters(ctx context.Context, request *dto.GetVideoFiltersRequest) (*dto.GetVideoFiltersResponse, error)

	// GetVideoURL 获取视频播放地址
	// ctx: 上下文信息
	// request: 包含用户ID、视频ID和集数的请求参数
	// 返回: 视频URL响应和可能的错误
	GetVideoURL(ctx context.Context, request *dto.GetVideoURLRequest) (*dto.GetVideoURLResponse, error)

	// GetHomeAnimes 获取首页动漫列表
	// ctx: 上下文信息
	// request: 包含用户ID的请求参数
	// 返回: 首页动漫列表响应和可能的错误
	GetHomeAnimes(ctx context.Context, request *dto.GetHomeAnimesRequest) (*dto.GetHomeAnimesResponse, error)

	// Response 生成视频URL响应
	// 参数1: 状态码
	// 参数2: 视频播放地址
	// 返回: 视频URL响应对象
	Response(int, string) *dto.GetVideoURLResponse

	// AddAnimeCollection 添加动漫收藏
	// ctx: 上下文信息
	// request: 包含用户ID和视频ID的请求参数
	// 返回: 更新动漫收藏响应和可能的错误
	UpdateAnimeCollection(ctx context.Context, request *dto.UpdateAnimeCollectionRequest) (*dto.UpdateAnimeCollectionResponse, error)

	// GetAnimeCollection 获取动漫收藏列表
	// ctx: 上下文信息
	// request: 包含用户ID、页码和每页数量的请求参数
	// 返回: 动漫收藏列表响应和可能的错误
	GetAnimeCollection(ctx context.Context, request *dto.GetAnimeCollectionRequest) (*dto.GetAnimeCollectionResponse, error)

	// GetRecommend 获取相关动漫推荐
	// ctx: 上下文信息
	// request: 包含用户ID和视频ID的请求参数
	// 返回: 相关动漫推荐响应和可能的错误
	GetRecommend(ctx context.Context, request *dto.GetRecommendRequest) (*dto.GetRecommendResponse, error)
}
