package router

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Lapakin/edu-planner/internal/ui/handler"
	"github.com/Lapakin/edu-planner/internal/ui/service"
)

func NewRouter(svc *service.Services) *gin.Engine {
	r := gin.New()

	r.LoadHTMLGlob("internal/ui/template/*.html")
	r.Static("/static", "internal/ui/static")

	r.GET("/", handler.Index)
	r.GET("/set-password", handler.SetPassword)

	api := r.Group("/api")
	api.Any("/auth/*path", handler.AuthProxy(svc))
	api.Any("/syllabus/*path", handler.SyllabusProxy(svc))

	r.HEAD("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	return r
}
