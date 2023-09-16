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
		containerGroup.GET("find-one", m.Handler.GetOneContainerByCodeAndType)
		containerGroup.POST("mark-container-full", m.Handler.MarkContainerFullByCode)
	}
}
