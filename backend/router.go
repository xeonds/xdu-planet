package main

import (
	"embed"
	"html/template"

	"github.com/gin-gonic/gin"
	"xyz.xeonds/xdu-planet/controller"
)

//go:embed template
var f embed.FS

func initRouter(r *gin.Engine) {
	r.SetHTMLTemplate(template.Must(template.New("").ParseFS(f, "template/*.html")))

	//路由分组
	apiRouter := r.Group("/api/v1")
	apiRouter.GET("/feed/", controller.FetchRawXml)

	// Home page
	r.GET("/", controller.GenPage)
}
