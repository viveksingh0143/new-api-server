package auth

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/vamika-digital/wms-api-server/config/authconfig"
	"github.com/vamika-digital/wms-api-server/core/base/dto"
	"github.com/vamika-digital/wms-api-server/core/base/validators"
	"github.com/vamika-digital/wms-api-server/core/business/auth/dto/auth"
	"github.com/vamika-digital/wms-api-server/core/business/auth/service"
)

type AuthRestHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthRestHandler {
	return &AuthRestHandler{authService: authService}
}

func (handler *AuthRestHandler) LoginHandler(c *gin.Context) {
	var credentials = &auth.LoginUserDto{}
	if err := c.ShouldBindJSON(&credentials); err != nil {
		log.Printf("%+v\n", err)
		c.JSON(http.StatusBadRequest, dto.GetErrorRestResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	if err := validators.Validate.Struct(credentials); err != nil {
		log.Printf("%+v\n", err)
		errors := validators.GetAllErrors(err, credentials)
		c.JSON(http.StatusBadRequest, dto.GetErrorRestResponse(http.StatusBadRequest, "Please fill the form correctly", errors))
		return
	}

	validUser, err := handler.authService.ValidateCredentials(credentials.Username, credentials.Password, credentials.LoginVia)
	if err != nil {
		log.Printf("%+v\n", err)
		c.JSON(http.StatusUnauthorized, dto.GetErrorRestResponse(http.StatusUnauthorized, "Credentials not matched", nil))
		return
	}

	if !validUser.Status.IsEnable() {
		c.JSON(http.StatusUnauthorized, dto.GetErrorRestResponse(http.StatusUnauthorized, "User account has been deactivated", nil))
		return
	}

	accessToken, err := handler.authService.GenerateAccessToken(validUser)
	if err != nil {
		log.Printf("%+v\n", err)
		c.JSON(http.StatusInternalServerError, dto.GetErrorRestResponse(http.StatusInternalServerError, "Issue at generating access token", nil))
		return
	}

	refreshToken, err := handler.authService.GenerateRefreshToken(validUser, credentials.RememberMe)
	if err != nil {
		log.Printf("%+v\n", err)
		c.JSON(http.StatusInternalServerError, dto.GetErrorRestResponse(http.StatusInternalServerError, "Issue at generating refresh token", nil))
		return
	}

	allPermissions, err := handler.authService.GetAllPermissions(validUser)
	if err != nil {
		log.Printf("%+v\n", err)
		c.JSON(http.StatusInternalServerError, dto.GetErrorRestResponse(http.StatusInternalServerError, "Issue at getting all permissions", nil))
		return
	}

	c.JSON(http.StatusOK, dto.RestResponse{
		Status: http.StatusOK,
		Data: auth.LoginTokenDto{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			Name:         validUser.Name,
			StaffID:      validUser.StaffID.String,
			Permissions:  allPermissions,
		},
	})
}

func (handler *AuthRestHandler) RefreshTokenHandler(c *gin.Context) {
	var refreshTokenDto auth.RefreshTokenDto
	if err := c.ShouldBindQuery(&refreshTokenDto); err != nil {
		log.Printf("%+v\n", err)
		c.JSON(http.StatusBadRequest, dto.GetErrorRestResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	claims := &jwt.StandardClaims{}
	tkn, err := jwt.ParseWithClaims(refreshTokenDto.RefreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return authconfig.Cfg.SecretKey, nil
	})
	if err != nil || !tkn.Valid {
		log.Printf("%+v\n", err)
		c.JSON(http.StatusUnauthorized, dto.GetErrorRestResponse(http.StatusUnauthorized, "Invalid refresh token", nil))
		return
	}

	validUser, err := handler.authService.GetUserById(claims.Subject)
	if err != nil {
		log.Printf("%+v\n", err)
		c.JSON(http.StatusNotFound, dto.GetErrorRestResponse(http.StatusUnauthorized, "Account not exists", nil))
		return
	}

	if !validUser.Status.IsEnable() {
		c.JSON(http.StatusUnauthorized, dto.GetErrorRestResponse(http.StatusUnauthorized, "User account has been deactivated", nil))
		return
	}

	accessToken, err := handler.authService.GenerateAccessToken(validUser)
	if err != nil {
		log.Printf("%+v\n", err)
		c.JSON(http.StatusInternalServerError, dto.GetErrorRestResponse(http.StatusUnauthorized, "Issue at generating access token", nil))
		return
	}

	refreshToken, err := handler.authService.GenerateRefreshToken(validUser, true)
	if err != nil {
		log.Printf("%+v\n", err)
		c.JSON(http.StatusInternalServerError, dto.GetErrorRestResponse(http.StatusUnauthorized, "Issue at generating refresh token", nil))
		return
	}

	allPermissions, err := handler.authService.GetAllPermissions(validUser)
	if err != nil {
		log.Printf("%+v\n", err)
		c.JSON(http.StatusInternalServerError, dto.GetErrorRestResponse(http.StatusInternalServerError, "Issue at getting all permissions", nil))
		return
	}

	c.JSON(http.StatusOK, dto.RestResponse{
		Status: http.StatusOK,
		Data: auth.LoginTokenDto{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			Name:         validUser.Name,
			StaffID:      validUser.StaffID.String,
			Permissions:  allPermissions,
		},
	})
}
