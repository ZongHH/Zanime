package router

func (c *Controller) setupAuthRoutes() {
	// 认证API
	c.engine.POST("/api/signInfo", c.userHandler.Register)
	c.engine.POST("/api/sendVerificationCode", nil)
	c.engine.POST("/api/loginInfo", c.userHandler.Login)
}
