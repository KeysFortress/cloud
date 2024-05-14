package routes

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"leanmeal/api/dtos"
	"leanmeal/api/interfaces"
	"leanmeal/api/middlewhere"
	"leanmeal/api/repositories"
	"leanmeal/api/utils"
)

type PasswordsController struct {
	PasswordRepository repositories.PasswordRepository
	PasswordService    interfaces.PasswordService
	EventRepository    repositories.EventRepository
}

func (pc *PasswordsController) all(ctx *gin.Context) {
	pc.PasswordRepository.Storage.Open()
	passwords, err := pc.PasswordRepository.All()
	pc.PasswordRepository.Storage.Close()

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Messsage": "Bad Request"})
		return
	}

	ctx.JSON(http.StatusOK, passwords)
}

func (pc *PasswordsController) content(ctx *gin.Context) {
	id := ctx.Param("id")

	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"Message": "Bad Request, missing id"})
		return
	}

	passwordId, err := utils.ParseUUID(id)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Message": "Bad Request, id is not UUID"})
		return
	}

	pc.PasswordRepository.Storage.Open()
	password, err := pc.PasswordRepository.Content(passwordId)
	pc.PasswordRepository.Storage.Close()

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Message": "Bad Request"})
		return
	}

	ctx.JSON(http.StatusOK, password)
}

func (pc *PasswordsController) generate(ctx *gin.Context) {

	request := &dtos.RequestPassword{}
	if err := ctx.BindJSON(request); err != nil {
		fmt.Println("could not bind request body")
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
		return
	}

	generated, err := pc.PasswordService.GeneratePassword(request.Lenght, request.LowerCase, request.UpperCase,
		request.Unique, request.SpecialCharacters)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
		return
	}

	ctx.JSON(http.StatusOK, generated)
}

func (pc *PasswordsController) add(ctx *gin.Context) {
	id := ctx.MustGet("ID")

	request := &dtos.IncomingPasswordRequest{}
	if err := ctx.BindJSON(request); err != nil {
		fmt.Println("could not bind request body")
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
		return
	}

	pc.PasswordRepository.Storage.Open()
	created, err := pc.PasswordRepository.Add(*request, id.(uuid.UUID))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
	}

	event, err := pc.EventRepository.Add(dtos.CreateEvent{
		TypeId:      1,
		Description: "New password created, for account: " + request.Email + " with address: " + request.Website,
		CreatedAt:   time.Now().UTC(),
	})
	if err != nil {
		fmt.Println("Failed to create event abording")
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
		return
	}

	pc.PasswordRepository.Storage.Close()

	fmt.Println(event)
	ctx.JSON(http.StatusOK, created)
}

func (pc *PasswordsController) update(ctx *gin.Context) {
	id := ctx.MustGet("ID")
	_ = id

	request := &dtos.IncomingPasswordUpdateRequest{}
	if err := ctx.BindJSON(request); err != nil {
		fmt.Println("could not bind request body")
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
		return
	}

	pc.PasswordRepository.Storage.Open()
	oldPassword, err := pc.PasswordRepository.Get(request.Id)
	if err != nil {
		fmt.Println("Record doesn't exist")
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
		return
	}

	event, err := pc.EventRepository.Add(dtos.CreateEvent{
		TypeId:      2,
		Description: "Password updated, from account: " + oldPassword.Email + " with address: " + oldPassword.Website + "to account: " + request.Email + " with address: " + request.Website,
		CreatedAt:   time.Now().UTC(),
	})
	if err != nil {
		fmt.Println("Failed to create an event aborting")
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
		return
	}

	fmt.Println(event)
	updated := pc.PasswordRepository.Update(request)
	pc.PasswordRepository.Storage.Close()

	ctx.JSON(http.StatusOK, updated)
}

func (pc *PasswordsController) remove(ctx *gin.Context) {
	id := ctx.MustGet("id")
	_ = id

	request := &uuid.UUID{}

	if err := ctx.BindJSON(request); err != nil {
		fmt.Println("Could not bind id, wrong type")
		ctx.JSON(http.StatusOK, gin.H{"message": "Bad Request"})
		return
	}

	ctx.JSON(http.StatusOK, "ok")
}

func (pc *PasswordsController) Init(r *gin.RouterGroup, authMiddlewhere *middlewhere.AuthenticationMiddlewhere) {
	controller := r.Group("passwords")
	controller.Use(authMiddlewhere.Authorize())

	controller.GET("all", pc.all)
	controller.GET("content/:id", pc.content)
	controller.POST("generate", pc.generate)
	controller.POST("add", pc.add)
	controller.POST("update", pc.update)
	controller.DELETE("remove", pc.remove)
}
