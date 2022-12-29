package main

import (
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"xyz.xeonds/xdu-planet/controller"
)

func main() {
	r := gin.Default()
	initRouter(r)
	crontab := cron.New(cron.WithSeconds())
	crontab.AddFunc("* * */4 * * *", controller.FetchFeed)
	crontab.Start()

	panic(r.Run(":8192"))
}
