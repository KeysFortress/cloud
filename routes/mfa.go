package routes

import (
	"github.com/gin-gonic/gin"

	"leanmeal/api/middlewhere"
)

type MfaController struct {
}

func (m *MfaController) pickMethod(ctx *gin.Context) {
	return
}

func (m *MfaController) performMethod(ctx *gin.Context) {
	return
}

func (m *MfaController) all(ctx *gin.Context) {
	return
}

func (m *MfaController) add(ctx *gin.Context) {
	return
}

func (m *MfaController) delete(ctx *gin.Context) {
	return
}

func (m *MfaController) Init(r *gin.RouterGroup, a *middlewhere.AuthenticationMiddlewhere) {
	controller := r.Group("mfa")
	controller.Use(a.AuthorizeMFA())

	controller.GET("pick-method", m.pickMethod)
	controller.POST("perform-method", m.performMethod)

	authorizedController := r.Group("user-mfa")
	authorizedController.Use(a.Authorize())

	authorizedController.GET("all", m.all)
	authorizedController.POST("add", m.add)
	authorizedController.DELETE("delete", m.delete)
}
