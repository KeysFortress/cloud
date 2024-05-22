package routes

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/lib/pq"

	implementations "leanmeal/api/Implementations"
	"leanmeal/api/dtos"
	"leanmeal/api/interfaces"
	"leanmeal/api/repositories"
	"leanmeal/api/utils"
)

type AuthenticationController struct {
	AuthenticationService interfaces.AuthenticationService
	AccountRepository     repositories.Accounts
	AccessKeysRepository  repositories.AccessKeysRepository
	JwtService            interfaces.JwtService
	Configuration         interfaces.Configuration
}

func (ac *AuthenticationController) beginRequest(ctx *gin.Context) {

	email := ctx.Param("email")

	if email == "" {
		ctx.JSON(http.StatusBadRequest, "Bad request")
		return
	}

	ac.AccessKeysRepository.Storage.Open()
	userExists, err := ac.AccountRepository.UserExists(email)
	ac.AccessKeysRepository.Storage.Close()

	if err != nil {
		ctx.JSON(http.StatusBadRequest, "Bad request")
		return
	}

	fmt.Println(email)
	data := ac.AuthenticationService.GetMessage(&userExists.Email, &userExists.Id)
	ctx.JSON(http.StatusOK, data)
}

func (ac *AuthenticationController) finishRequest(ctx *gin.Context) {
	request := &dtos.FinishAuthResponse{}
	if err := ctx.BindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
		return
	}

	accountId := ac.AuthenticationService.GetRequestById(request.Uuid)
	if accountId == uuid.Nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Request expired or canceled!"})
		return
	}

	ac.AccessKeysRepository.Storage.Open()
	keys := ac.AccessKeysRepository.GetAccountKeys(accountId)
	ac.AccessKeysRepository.Storage.Close()
	isVerified, err := ac.AuthenticationService.VerifySignature(*request, &keys)

	if err != nil || isVerified == uuid.Nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad Signature!"})
		return
	}

	approvedKey, _ := ac.AuthenticationService.ExchangeCodeForPublicKey(request.Uuid)

	token := ac.JwtService.IssueMfaToken(isVerified.String(), approvedKey)
	ctx.JSON(http.StatusOK, token)
}

func (ac *AuthenticationController) waitLogin(ctx *gin.Context) {
	code := ctx.Param("code")

	if code == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"Message": "Bad request"})
		return
	}

	id, err := utils.ParseUUID(code)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Message": "Invalid request, bad id"})
		return
	}

	var getToken uuid.UUID
	var valid bool
	valid = true

	for valid {
		getToken, valid = ac.AuthenticationService.ExchangeCodeForToken(id)

		if !valid {
			ctx.JSON(http.StatusBadRequest, gin.H{"Message": "Invalid request, token has expired"})
			return
		}

		if getToken != uuid.Nil {
			approvedKey, _ := ac.AuthenticationService.ExchangeCodeForPublicKey(id)

			jwtToken := ac.JwtService.IssueMfaToken(getToken.String(), approvedKey)
			ctx.JSON(http.StatusOK, jwtToken)
			return
		}
	}
}

func (ac *AuthenticationController) pickMethod(ctx *gin.Context) {

}

func (ac *AuthenticationController) Init(r *gin.RouterGroup) {
	domain := ac.Configuration.GetKey("domain")

	if domain == nil {
		panic("Domain is not set in the configuration file")
	}

	ac.AuthenticationService = &implementations.AuthenticationService{
		Domain: domain.(string),
	}

	go ac.AuthenticationService.Start()

	r.GET("/begin-request/:email", ac.beginRequest)
	r.GET("/login/:code", ac.waitLogin)
	r.GET("/pick-method/:id", ac.pickMethod)
	r.POST("/finish-request", ac.finishRequest)
	r.POST("/perform-method/:id", ac.pickMethod)
}
