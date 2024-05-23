package routes

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pquerna/otp"

	"leanmeal/api/dtos"
	"leanmeal/api/interfaces"
	"leanmeal/api/middlewhere"
	"leanmeal/api/repositories"
	"leanmeal/api/utils"
)

type MfaController struct {
	TotpService        interfaces.TimeBasedService
	AccountsRepository repositories.Accounts
	MfaRepository      repositories.MfaRepository
	EmailService       interfaces.MailService
	JwtService         interfaces.JwtService
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

	_, err = m.MfaRepository.Add(secret, 2, id.(uuid.UUID), sql.NullString{
		String: account.Email,
	})

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Message": "Failed to save new secret, aborted operation!"})
		return
	}

	ctx.JSON(http.StatusOK, secret)
}

func (m *MfaController) pickMethod(ctx *gin.Context) {
	id := ctx.MustGet("ID")

	if id == "" {
		fmt.Println("Failed to extract user id aborting")
		ctx.JSON(http.StatusBadRequest, gin.H{"Message": "Bad Request"})
		return
	}
	typeIdData := ctx.Param("type")

	typeId, err := utils.ParseInt(typeIdData)

	if err != nil {
		fmt.Println("Failed to parse type id aborting")
		ctx.JSON(http.StatusBadRequest, gin.H{"Message": "Bad Request"})
		return
	}

	if typeId == 2 {
		ctx.JSON(http.StatusOK, gin.H{"Message": "Please check your authenticator application"})
		return
	}

	m.MfaRepository.Storage.Open()
	emails, err := m.MfaRepository.GetForUserByType(id.(uuid.UUID), typeId)
	m.MfaRepository.Storage.Close()

	if err != nil {
		result := "Failed to retrive emails methods for user " + id.(string)
		ctx.JSON(http.StatusBadRequest, gin.H{"Message": result})
		return

	}

	for _, maiData := range emails {
		code, err := m.TotpService.GenerateTOTPCode(maiData.Value, 30, otp.AlgorithmMD5)
		if err != nil {
			result := "Bad Request failed to generate Code for email" + maiData.Address.String
			ctx.JSON(http.StatusBadRequest, gin.H{"Message": result})
			return
		}

		emailSent, err := m.EmailService.SendMessage(maiData.Address.String, "MFA Verification", "Your mfa verification code is "+code)
		if err != nil || !emailSent {
			result := "Bad Request failed to send for email" + maiData.Address.String
			ctx.JSON(http.StatusBadRequest, gin.H{"Message": result})
			return
		}

	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Please check your email"})
}

func (m *MfaController) performMethod(ctx *gin.Context) {
	id := ctx.MustGet("ID")
	deviceKey := ctx.MustGet("DeviceKey")
	request := &dtos.MfaCodeRequest{}
	if err := ctx.BindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
		return
	}

	m.MfaRepository.Storage.Open()
	methods, err := m.MfaRepository.GetForUser(id.(uuid.UUID))
	m.MfaRepository.Storage.Close()

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Message": "Bad Request"})
	}

	var valid bool
	for _, method := range methods {
		valid, err = m.TotpService.VerifyTOTP(request.Code, method.Value, 30, otp.AlgorithmMD5)

		if err != nil {
			fmt.Println(err)
		}

		if valid {
			break
		}
	}

	if !valid {
		ctx.JSON(http.StatusBadRequest, gin.H{"Message": "Failed to verify method, code doesn't match"})
		return
	}

	token := m.JwtService.IssueToken("user", id.(string), deviceKey.(string))
	ctx.JSON(http.StatusOK, token)
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
	controller.GET("pick-method/:type", m.pickMethod)
	controller.POST("perform-method", m.performMethod)

	authorizedController := r.Group("user-mfa")
	authorizedController.Use(a.Authorize())

	authorizedController.GET("all", m.all)
	authorizedController.POST("add", m.add)
	authorizedController.DELETE("delete", m.delete)
}
