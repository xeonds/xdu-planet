package main

import (
	"flag"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"xyz.xeonds/xdu-planet/controller"
)

func main() {
	var updateDB bool
	flag.BoolVar(&updateDB, "update-db", false, "Fetch and update feed database")

	// 解析命令行参数
	flag.Parse()

	// 如果updateDB为true，则调用FetchFeed函数
	switch {
	case updateDB:
		controller.FetchFeed()
	default:
		r := gin.Default()
		r.Use(cors.Default())
		initRouter(r)
		crontab := cron.New(cron.WithSeconds())
		crontab.AddFunc("0 15 * * * *", controller.FetchFeed)
		crontab.Start()

		panic(r.Run(":8192"))
	}

}
