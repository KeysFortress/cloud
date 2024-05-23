package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"leanmeal/api/interfaces"
	"leanmeal/api/middlewhere"
	"leanmeal/api/repositories"
)

type MfaController struct {
	TotpService        interfaces.TimeBasedService
	AccountsRepository repositories.Accounts
	MfaRepository      repositories.MfaRepository
}

func (m *MfaController) setup(ctx *gin.Context) {

	id := ctx.MustGet("ID")

	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"Message": "Bad Request"})
		return
	}

	m.AccountsRepository.Storage.Open()
	account, err := m.AccountsRepository.GetById(id.(uuid.UUID))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Message": "Bad Request"})
		m.AccountsRepository.Storage.Close()
		return
	}

	configured, err := m.MfaRepository.IsConfigured(id.(uuid.UUID))
	m.MfaRepository.Storage.Close()

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Message": "Denied, failed to fetch existing methods, contact an administrator"})
		return
	}

	if configured {
		ctx.JSON(http.StatusBadRequest, gin.H{"Message": "Denied, mfa already configured"})
		return
	}

	secret, err := m.TotpService.GenerateTOTP(account.Email)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Message": "Bad Request"})
		return
	}

	_, err = m.MfaRepository.Add(secret, 2, id.(uuid.UUID))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Message": "Failed to save new secret, aborted operation!"})
		return
	}

	ctx.JSON(http.StatusOK, secret)
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

	controller.GET("setup", m.setup)
	controller.GET("pick-method", m.pickMethod)
	controller.POST("perform-method", m.performMethod)

	authorizedController := r.Group("user-mfa")
	authorizedController.Use(a.Authorize())

	authorizedController.GET("all", m.all)
	authorizedController.POST("add", m.add)
	authorizedController.DELETE("delete", m.delete)
}
