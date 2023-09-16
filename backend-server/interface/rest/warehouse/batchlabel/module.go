package batchlabel

import (
	"github.com/gin-gonic/gin"
	"github.com/vamika-digital/wms-api-server/core/business/warehouse/service"
)

type BatchLabelRestModule struct {
	Handler *BatchLabelRestHandler
}

func NewBatchLabelRestModule(batchlabelService service.BatchLabelService) *BatchLabelRestModule {
	batchlabelHandler := NewBatchLabelHandler(batchlabelService)
	return &BatchLabelRestModule{Handler: batchlabelHandler}
}

func (m *BatchLabelRestModule) RegisterRoutes(r *gin.RouterGroup) {
	batchlabelGroup := r.Group("/batchlabels")
	{
		batchlabelGroup.POST("", m.Handler.CreateBatchLabel)
		batchlabelGroup.GET("", m.Handler.GetAllBatchLabels)
		batchlabelGroup.POST("/bulkdelete", m.Handler.DeleteBatchLabelByIDs)
		batchlabelGroup.GET("/:id", m.Handler.GetBatchLabelByID)
		batchlabelGroup.PUT("/:id", m.Handler.UpdateBatchLabel)
		batchlabelGroup.DELETE("/:id", m.Handler.DeleteBatchLabel)
		batchlabelGroup.POST("/:id/generate-stickers", m.Handler.GenerateBatchLabelStickers)

		// GenerateBatchLabelStickers(batchlabelID int64, form *batchlabel.BatchLabelStickersCreateDto)
	}
}
