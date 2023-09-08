package user

import (
	"github.com/gin-gonic/gin"
	"github.com/vamika-digital/wms-api-server/core/business/admin/service"
)

type UserRestModule struct {
	Handler *UserRestHandler
}

func NewUserRestModule(userService service.UserService) *UserRestModule {
	userHandler := NewUserHandler(userService)
	return &UserRestModule{Handler: userHandler}
}

func (m *UserRestModule) RegisterRoutes(r *gin.RouterGroup) {
	userGroup := r.Group("/users")
	{
		userGroup.POST("", m.Handler.CreateUser)
		userGroup.GET("", m.Handler.GetAllUsers)
		userGroup.GET("/:id", m.Handler.GetUserByID)
		userGroup.PUT("/:id", m.Handler.UpdateUser)
		userGroup.DELETE("/:id", m.Handler.DeleteUser)
	}
}
