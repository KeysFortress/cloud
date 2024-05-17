package routes

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"leanmeal/api/dtos"
	"leanmeal/api/interfaces"
	"leanmeal/api/middlewhere"
	"leanmeal/api/repositories"
	"leanmeal/api/utils"
)

type TotpController struct {
	TotpRepository   repositories.TotpRepository
	EventsRepository repositories.EventRepository
	TotpService      interfaces.TimeBasedService
}

func (tc *TotpController) all(ctx *gin.Context) {
	tc.TotpRepository.Storage.Open()
	all, err := tc.TotpRepository.All()
	tc.TotpRepository.Storage.Close()

	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"Message": "Bad Request"})
		return
	}

	ctx.JSON(http.StatusOK, all)
}

func (tc *TotpController) content(ctx *gin.Context) {
	contentId := ctx.Param("id")

	id, err := utils.ParseUUID(contentId)

	if err != nil {
		fmt.Println("Wrong parameter format")
		ctx.JSON(http.StatusBadRequest, gin.H{"Message": "Bad Request"})
		return
	}

	tc.TotpRepository.Storage.Open()
	secret, err := tc.TotpRepository.Content(id)
	tc.TotpRepository.Storage.Close()

	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"Message": "Bad Request"})
		return
	}

	ctx.JSON(http.StatusOK, secret)
}

func (tc *TotpController) types(ctx *gin.Context) {
	tc.TotpRepository.Storage.Open()
	types, err := tc.TotpRepository.GetCodeTypes()
	tc.TotpRepository.Storage.Close()

	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"Message": "Bad Request"})
		return
	}

	ctx.JSON(http.StatusOK, types)
}

func (tc *TotpController) algorithms(ctx *gin.Context) {

	tc.TotpRepository.Storage.Open()
	algorithms, err := tc.TotpRepository.GetAlgorithms()
	tc.TotpRepository.Storage.Close()

	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"Message": "Bad Request"})
		return
	}

	ctx.JSON(http.StatusOK, algorithms)
}

func (tc TotpController) code(ctx *gin.Context) {
	id := ctx.Param("id")
	uuid, err := utils.ParseUUID(id)

	if err != nil {
		fmt.Println("Bad parameter")
		ctx.JSON(http.StatusBadRequest, gin.H{"Message": "Bad Request"})
		return
	}

	result := "--- ---"
	tc.TotpRepository.Storage.Open()
	secret, err := tc.TotpRepository.GetInternal(uuid)
	tc.TotpRepository.Storage.Close()

	if err != nil {
		fmt.Println("Record doesn't exist")
		tc.TotpRepository.Storage.Close()
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
		return
	}

	switch secret.Type {
	case 1:
		code, err := tc.TotpService.GenerateHOTP(secret.Secret)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"Message": "Bad Request"})
			return
		}
		var codes string
		for _, c := range code {
			codes += c
		}
		result = codes
	default:
		currentAlgorithm := utils.ParseAlgorithm(secret.Type)
		code, err := tc.TotpService.GenerateTOTP(secret.Secret, secret.Validity, currentAlgorithm)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"Message": "Bad Request"})
			return
		}

		result = code
	}

	ctx.JSON(http.StatusOK, result)
}

func (tc *TotpController) add(ctx *gin.Context) {
	request := &dtos.CreateTimeBasedCode{}
	if err := ctx.BindJSON(request); err != nil {
		fmt.Println("could not bind request body")
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
		return
	}

	tc.TotpRepository.Storage.Open()
	created, err := tc.TotpRepository.Add(*request)
	if err != nil {
		tc.TotpRepository.Storage.Close()
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
		return
	}

	event, err := tc.EventsRepository.Add(dtos.CreateEvent{
		TypeId:      1,
		Description: "New time based password created, for account: " + request.Email + " with address: " + request.Website,
		CreatedAt:   time.Now().UTC(),
	})
	if err != nil {
		fmt.Println("Failed to create event abording")
		tc.TotpRepository.Storage.Close()
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
		return
	}

	tc.TotpRepository.Storage.Close()

	fmt.Println(event)
	ctx.JSON(http.StatusOK, created)
}

func (tc *TotpController) update(ctx *gin.Context) {
	id := ctx.MustGet("ID")
	_ = id

	request := &dtos.UpdateTimeBasedCode{}
	if err := ctx.BindJSON(request); err != nil {
		fmt.Println("could not bind request body")
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
		return
	}

	tc.TotpRepository.Storage.Open()
	oldPassword, err := tc.TotpRepository.Get(request.Id)
	if err != nil {
		fmt.Println("Record doesn't exist")
		tc.TotpRepository.Storage.Close()
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
		return
	}

	updated := tc.TotpRepository.Update(request)

	if !updated {
		fmt.Println("Failed to update entity, aborting")
		tc.TotpRepository.Storage.Close()
		ctx.JSON(http.StatusBadRequest, gin.H{"Message": "Bad Request"})
		return
	}

	event, err := tc.EventsRepository.Add(dtos.CreateEvent{
		TypeId:      2,
		Description: "TOTP updated, from account: " + oldPassword.Email + " with address: " + oldPassword.Website + "to account: " + request.Email + " with address: " + request.Website,
		CreatedAt:   time.Now().UTC(),
	})

	if err != nil {

		fmt.Println("Entity was updated, but failed to create an event.")
		tc.TotpRepository.Storage.Close()
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})

		return
	}

	fmt.Println(event)
	tc.TotpRepository.Storage.Close()
	ctx.JSON(http.StatusOK, updated)
}

func (tc *TotpController) Init(r *gin.RouterGroup, authMiddlewhere *middlewhere.AuthenticationMiddlewhere) {
	controller := r.Group("totp")
	controller.Use(authMiddlewhere.Authorize())

	controller.GET("all", tc.all)
	controller.GET("content/:id", tc.content)
	controller.GET("types", tc.types)
	controller.GET("algorithms", tc.algorithms)
	controller.GET("code/:id", tc.code)
	controller.POST("add", tc.add)
	controller.POST("update", tc.update)
}
