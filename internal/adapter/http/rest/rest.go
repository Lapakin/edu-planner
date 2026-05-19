package rest

import (
	"github.com/gin-gonic/gin"
)

func RespondJSON(c *gin.Context, statusCode int, data any) {
	c.JSON(statusCode, data)
}

func RespondError(c *gin.Context, statusCode int, message error) {
	c.JSON(statusCode, gin.H{"error": message.Error()})
}
