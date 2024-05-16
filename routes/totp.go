package routes

import (
	"github.com/gin-gonic/gin"

	"leanmeal/api/middlewhere"
)

type TotpController struct {
}

func (tc *TotpController) Init(r *gin.RouterGroup, authMiddlewhere *middlewhere.AuthenticationMiddlewhere) {
	controller := r.Group("passwords")
	controller.Use(authMiddlewhere.Authorize())

}
