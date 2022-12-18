package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()
	initRouter(r)
	panic(r.Run(":8192"))
}
