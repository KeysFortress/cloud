package middlewhere

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"leanmeal/api/interfaces"
	"leanmeal/api/utils"
)

type UnsignedResponse struct {
	Message interface{} `json:"message"`
}

type AuthenticationMiddlewhere struct {
	JwtService interfaces.JwtService
}

func (aw *AuthenticationMiddlewhere) Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		jwtToken, err := extractBearerToken(c.GetHeader("Authorization"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, UnsignedResponse{
				Message: err.Error(),
			})
			return
		}

		result := aw.JwtService.ValidateToken(jwtToken)
		if !result {
			c.AbortWithStatusJSON(http.StatusBadRequest, UnsignedResponse{
				Message: "Invalid access token",
			})
			return
		}

		id := aw.JwtService.ExtractValue(jwtToken, "id")

		userId, err := utils.ParseUUID(id.(string))
		if err != nil {
			fmt.Print("User id is not a valid UUID")
			c.AbortWithStatusJSON(http.StatusBadRequest, UnsignedResponse{
				Message: "Invalid access token",
			})
			return
		}
		c.Set("ID", userId)
		deviceKey := aw.JwtService.ExtractValue(jwtToken, "deviceKey")
		c.Set("DeviceKey", deviceKey)

		c.Copy().Next()
	}
}

func extractBearerToken(header string) (string, error) {
	if header == "" {
		return "", errors.New("bad header value given")
	}

	jwtToken := strings.Split(header, " ")
	if len(jwtToken) != 2 {
		return "", errors.New("incorrectly formatted authorization header")
	}

	return jwtToken[1], nil
}
