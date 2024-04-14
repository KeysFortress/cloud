package routes

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	implementations "leanmeal/api/Implementations"
	"leanmeal/api/dtos"
	"leanmeal/api/interfaces"
	"leanmeal/api/models"
	"leanmeal/api/repositories"
)

type AuthenticationController struct {
	AuthenticationService interfaces.AuthenticationService
	Storage               interfaces.Storage
	accountRepistory      repositories.Accounts
	accessKeysRepository  repositories.AccessKeysRepository
	JwtService            interfaces.JwtService
}

func (ac *AuthenticationController) beginRequest(ctx *gin.Context) {

	email := ctx.Param("email")

	if email == "" {
		ctx.JSON(500, "Bad request")
		return
	}

	ac.accountRepistory.OpenConnection(&ac.Storage)
	userExists := ac.accountRepistory.UserExists(email)

	if userExists == (models.Account{}) {
		ctx.JSON(500, "Bad request")
		return
	}

	ac.accountRepistory.Close()

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
	ac.accessKeysRepository.OpenConnection(&ac.Storage)
	keys := ac.accessKeysRepository.GetAccountKeys(accountId)
	ac.accessKeysRepository.Close()
	isVerified, err := ac.AuthenticationService.VerifySignature(*request, &keys)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad Signature!"})
		return
	}

	token := ac.JwtService.IssueToken("user", isVerified.String())
	ctx.JSON(http.StatusOK, gin.H{"access_token": token})
}

func (ac *AuthenticationController) createAccount(ctx *gin.Context) {
	request := &dtos.CreateAccountRequest{}

	if err := ctx.BindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
		return
	}
	ac.accessKeysRepository.OpenConnection(&ac.Storage)
	ac.accountRepistory.Storage = ac.accessKeysRepository.Storage
	created := ac.accountRepistory.CreateAccount(request)
	ac.accessKeysRepository.Add(&created, &request.PublicKey)
	ac.accessKeysRepository.Close()

	parse := created.String()
	token := ac.JwtService.IssueToken("user", parse)

	ctx.JSON(http.StatusOK, gin.H{"access_token": token})

}

func (ac *AuthenticationController) Init(r *gin.RouterGroup) {
	ac.AuthenticationService = &implementations.AuthenticationService{}
	println("initializing Authentication Controller")
	go ac.AuthenticationService.Start()

	r.GET("/begin-request/:email", ac.beginRequest)
	r.POST("/finish-request", ac.finishRequest)
	r.POST("/create-account", ac.createAccount)

}
