package router

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Lapakin/edu-planner/internal/middleware"
	"github.com/Lapakin/edu-planner/internal/user-management/handler"
	"github.com/Lapakin/edu-planner/internal/user-management/service"
)

func NewRouter(svc *service.Services) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	apiV1 := r.Group("/api/v1")
	apiV1.Handle(http.MethodPost, "/login", handler.NewLoginHandler(svc.AuthSvc))
	apiV1.Handle(http.MethodPost, "/set-password", handler.NewSetPasswordHandler(svc.AuthSvc))
	apiV1.Handle(http.MethodPost, "/reset-password", handler.NewResetPasswordHandler(svc.AuthSvc))

	protected := apiV1.Group("/")
	protected.Use(middleware.JWTMiddleware())

	protected.Handle(http.MethodPost, "/auth/invite", handler.NewCreateInviteHandler(svc.AuthSvc))

	protected.Handle(http.MethodGet, "/users/:userId", handler.NewGetUserByIDHandler(svc.UserSvc))
	protected.Handle(http.MethodPost, "/users", handler.NewCreateUsersHandler(svc.UserSvc))
	protected.Handle(http.MethodGet, "/users", handler.NewFetchUsersHandler(svc.UserSvc))
	protected.Handle(http.MethodPut, "/users", handler.NewUpdateUsersHandler(svc.UserSvc))
	protected.Handle(http.MethodPost, "/users/attach", handler.NewAttachUsersHandler(svc.UserSvc))
	protected.Handle(http.MethodPost, "/users/delete", handler.NewDeleteUsersHandler(svc.UserSvc))
	protected.Handle(http.MethodPost, "/users/:userId/activate", handler.NewActivateUserHandler(svc.UserSvc))
	protected.Handle(http.MethodPost, "/users/:userId/deactivate", handler.NewDeactivateUserHandler(svc.UserSvc))
	protected.Handle(http.MethodPost, "/users/:userId/reset-invite", handler.NewResetInviteHandler(svc.AuthSvc))
	protected.Handle(http.MethodPost, "/users/:userId/generate-reset-password-link", handler.NewGenerateResetPasswordLinkHandler(svc.AuthSvc))

	r.HEAD("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	return r
}
