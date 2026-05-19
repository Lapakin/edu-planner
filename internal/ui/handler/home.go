package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func SetPassword(c *gin.Context) {
	token := c.Query("token")
	c.HTML(http.StatusOK, "set-password.html", gin.H{"token": token})
}
