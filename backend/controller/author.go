package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthorPage(c *gin.Context) {
	c.HTML(http.StatusOK, "author.html", gin.H{
		"author": feed.Author,
	})
}
