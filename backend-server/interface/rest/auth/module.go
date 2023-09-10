package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/vamika-digital/wms-api-server/core/business/auth/service"
)

type AuthRestModule struct {
	Handler *AuthRestHandler
}

func NewAuthRestModule(authService service.AuthService) *AuthRestModule {
	authHandler := NewAuthHandler(authService)
	return &AuthRestModule{Handler: authHandler}
}

func (m *AuthRestModule) RegisterRoutes(r *gin.RouterGroup) {
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/login", m.Handler.LoginHandler)
	}
}
