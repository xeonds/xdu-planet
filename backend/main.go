package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"xyz.xeonds/xdu-planet/controller"
)

func main() {
	r := gin.Default()
	r.Use(cors.Default())
	initRouter(r)
	crontab := cron.New(cron.WithSeconds())
	crontab.AddFunc("0 15 * * * *", controller.FetchFeed)
	crontab.Start()

	panic(r.Run(":8192"))
}
