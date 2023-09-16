package store

import (
	"github.com/gin-gonic/gin"
	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
	"github.com/vamika-digital/wms-api-server/core/business/admin/dto/user"
)

type StoreFilterDto struct {
	Query     string                 `form:"query"`
	ID        int64                  `form:"id"`
	Code      string                 `form:"code"`
	Name      string                 `form:"name"`
	StoreType string                 `form:"store_type"`
	Status    customtypes.StatusEnum `form:"status"`
	Owner     *user.UserMinimalDto   `form:"owner"`
}

func (s *StoreFilterDto) Bind(ctx *gin.Context) error {
	value := ctx.DefaultQuery("status", "")
	return s.Status.FromString(value)
}
