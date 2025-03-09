package bootstrap

// setAPIRouters 设置API路由
// 该方法负责配置所有API端点，将HTTP请求映射到相应的处理函数
// 参数:
//   - handlers: 包含所有处理器的结构体，用于处理不同类型的请求
func (c *Container) setAPIRouters(handlers *handlers) {
	// 获取统计数据的API端点
	// 返回系统中的用户数量、动画数量、今日播放次数和活跃用户数等统计信息
	c.Controller.GET("/api/statistics", handlers.StatisticsHandler.GetStatisticsData)

	// 获取用户行为日志的API端点
	// 返回系统中记录的用户活动日志，包括用户名、操作类型、时间等信息
	c.Controller.GET("/api/userActionLogs", handlers.UserActionLogsHandler.GetUserActionLogs)

	// 获取最新上线动漫的API端点
	// 返回系统中最新添加的动漫作品列表，包括标题、封面图片和更新时间等信息
	c.Controller.GET("/api/newAnime", handlers.StatisticsHandler.GetNewAnime)
}
