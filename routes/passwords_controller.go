package routes

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"leanmeal/api/dtos"
	"leanmeal/api/interfaces"
	"leanmeal/api/middlewhere"
	"leanmeal/api/repositories"
)

type PasswordsController struct {
	PasswordRepository repositories.PasswordRepository
	PasswordService    interfaces.PasswordService
}

func (pc *PasswordsController) all(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func (pc *PasswordsController) category(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "ok"})
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
	id := ctx.MustGet("id")

	request := &dtos.IncomingPasswordRequest{}
	if err := ctx.BindJSON(request); err != nil {
		fmt.Println("could not bind request body")
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
		return
	}

	pc.PasswordRepository.Storage.Open()
	created, err := pc.PasswordRepository.Add(*request, id.(uuid.UUID))
	pc.PasswordRepository.Storage.Close()

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
	}

	ctx.JSON(http.StatusOK, created)
}

func (pc *PasswordsController) update(ctx *gin.Context) {
	id := ctx.MustGet("id")
	_ = id

	request := &dtos.IncomingPasswordUpdateRequest{}
	if err := ctx.BindJSON(request); err != nil {
		fmt.Println("could not bind request body")
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
		return
	}

	pc.PasswordRepository.Storage.Open()
	updated, err := pc.PasswordRepository.Update(request)
	pc.PasswordRepository.Storage.Close()

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
	}

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
	controller.GET("category/:id", pc.category)
	controller.POST("generate", pc.generate)
	controller.POST("add", pc.add)
	controller.POST("update", pc.update)
	controller.DELETE("remove", pc.remove)
}
