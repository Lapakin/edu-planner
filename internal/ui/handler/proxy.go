package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/Lapakin/edu-planner/internal/ui/service"
)

func AuthProxy(svc *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.URL.Path = c.Param("path")
		svc.ProxyAuth(c.Writer, c.Request)
	}
}

func SyllabusProxy(svc *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.URL.Path = c.Param("path")
		svc.ProxySyllabus(c.Writer, c.Request)
	}
}
