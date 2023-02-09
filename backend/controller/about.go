package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AboutPage(c *gin.Context) {
	c.HTML(http.StatusOK, "about.html", gin.H{})
}
