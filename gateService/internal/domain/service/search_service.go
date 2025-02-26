// package service 提供了与搜索相关的业务逻辑服务
package service

import (
	"context"
	"gateService/internal/interfaces/dto"
)

// SearchService 定义了搜索服务的接口
// 提供视频搜索的基本功能和详细信息查询功能
type SearchService interface {
	// SearchVideos 搜索视频基本信息
	// 参数:
	// - ctx: 上下文信息
	// - request: 搜索请求参数,包含关键词、分页等信息
	// 返回:
	// - *dto.SearchResponse: 搜索结果响应,包含视频基本信息列表
	// - error: 搜索过程中的错误信息
	SearchVideos(ctx context.Context, request *dto.SearchRequest) (*dto.SearchResponse, error)

	// SearchVideosDetail 搜索视频详细信息
	// 参数:
	// - ctx: 上下文信息
	// - request: 详细搜索请求参数,包含视频ID等详细查询条件
	// 返回:
	// - *dto.SearchDetailResponse: 详细搜索结果响应,包含视频的完整信息
	// - error: 搜索过程中的错误信息
	SearchVideosDetail(ctx context.Context, request *dto.SearchDetailRequest) (*dto.SearchDetailResponse, error)
}
