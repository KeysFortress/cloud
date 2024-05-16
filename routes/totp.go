package routes

import (
	"github.com/gin-gonic/gin"

	"leanmeal/api/middlewhere"
	"leanmeal/api/repositories"
)

type TotpController struct {
	TotpRepository   repositories.TotpRepository
	EventsRepository repositories.EventRepository
}

func (tc *TotpController) Init(r *gin.RouterGroup, authMiddlewhere *middlewhere.AuthenticationMiddlewhere) {
	controller := r.Group("passwords")
	controller.Use(authMiddlewhere.Authorize())

}
