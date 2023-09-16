package labelsticker

import (
	"github.com/gin-gonic/gin"
	"github.com/vamika-digital/wms-api-server/core/business/warehouse/service"
)

type LabelStickerRestModule struct {
	Handler *LabelStickerRestHandler
}

func NewLabelStickerRestModule(labelstickerService service.LabelStickerService) *LabelStickerRestModule {
	labelstickerHandler := NewLabelStickerHandler(labelstickerService)
	return &LabelStickerRestModule{Handler: labelstickerHandler}
}

func (m *LabelStickerRestModule) RegisterRoutes(r *gin.RouterGroup) {
	labelstickerGroup := r.Group("/labelstickers")
	{
		labelstickerGroup.GET("", m.Handler.GetAllLabelStickers)
		labelstickerGroup.GET("/:id", m.Handler.GetLabelStickerByID)
	}
}
