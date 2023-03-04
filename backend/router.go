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

	// APIs
	apiRouter := r.Group("/api/v1")
	apiRouter.GET("/feed", controller.FetchRawFeed)             // Get feed url list
	apiRouter.GET("/feed/:feed", controller.FetchRawFeed)       // Get feed info
	apiRouter.PUT("/comment", controller.FetchRawFeed)          // Send comment by article ID
	apiRouter.GET("/comment/:comment", controller.FetchRawFeed) // Get comment by article ID

	// Pages
	r.GET("/", controller.HomePage)         // Home page
	r.GET("/author", controller.AuthorPage) // Author list
	r.GET("/about", controller.AboutPage)   // About site
}
