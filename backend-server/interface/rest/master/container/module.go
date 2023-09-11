package container

import (
	"github.com/gin-gonic/gin"
	"github.com/vamika-digital/wms-api-server/core/business/master/service"
)

type ContainerRestModule struct {
	Handler *ContainerRestHandler
}

func NewContainerRestModule(containerService service.ContainerService) *ContainerRestModule {
	containerHandler := NewContainerHandler(containerService)
	return &ContainerRestModule{Handler: containerHandler}
}

func (m *ContainerRestModule) RegisterRoutes(r *gin.RouterGroup) {
	containerGroup := r.Group("/containers")
	{
		containerGroup.POST("", m.Handler.CreateContainer)
		containerGroup.GET("", m.Handler.GetAllContainers)
		containerGroup.GET("/container-code-info", m.Handler.GetContainerCodeInfoDto)
		containerGroup.POST("/bulkdelete", m.Handler.DeleteContainerByIDs)
		containerGroup.GET("/:id", m.Handler.GetContainerByID)
		containerGroup.PUT("/:id", m.Handler.UpdateContainer)
		containerGroup.DELETE("/:id", m.Handler.DeleteContainer)
	}
}
