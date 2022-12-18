package main

import (
	"github.com/gin-gonic/gin"
	"xyz.xeonds/xdu-planet/controller"
)

func initRouter(r *gin.Engine) {
	// public directory is used to serve static resources
	r.Static("/", "./")

	//路由分组
	apiRouter := r.Group("/api/v1")
	apiRouter.GET("/feed/", controller.FetchFeed)
}
