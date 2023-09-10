package auth

import (
	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
)

type LoginUserDto struct {
	Username   string                    `json:"username" validate:"required"`
	Password   string                    `json:"password" validate:"required,min=3"`
	RememberMe bool                      `json:"remember_me"`
	LoginVia   *customtypes.LoginViaEnum `json:"login_via" validate:"omitempty,LoginViaEnum=STAFF_ID EMAIL USERNAME"`
}
