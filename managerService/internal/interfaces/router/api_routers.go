package controller

import "github.com/gin-gonic/gin"

func (c *Controller) get(path string, handlers ...gin.HandlerFunc) {
	c.engine.GET(path, handlers...)
}

func (c *Controller) setAPIRouters() {
	c.get("/api/statistics")
}
