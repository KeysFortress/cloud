package routes

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"leanmeal/api/dtos"
	"leanmeal/api/middlewhere"
	"leanmeal/api/repositories"
	"leanmeal/api/utils"
)

type SecretsController struct {
	SecretReoistory repositories.SecretsRepository
	EventRepository repositories.EventRepository
}

func (sc *SecretsController) all(ctx *gin.Context) {

	sc.SecretReoistory.Storage.Open()
	secrets, err := sc.SecretReoistory.All()
	sc.SecretReoistory.Storage.Close()

	if err != nil {
		fmt.Println("Failed to fetch secrets")
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"Message": "Bad Request"})
		return
	}

	ctx.JSON(http.StatusOK, secrets)
}

func (sc *SecretsController) add(ctx *gin.Context) {
	id := ctx.MustGet("ID")

	request := &dtos.IncomingSecretsRequest{}
	if err := ctx.BindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
		return
	}

	sc.SecretReoistory.Storage.Open()
	created, err := sc.SecretReoistory.Add(*request, id.(uuid.UUID))

	if err != nil {
		fmt.Println("Failed to create a secret")
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"Message": "Bad Request"})
		return
	}

	event, err := sc.EventRepository.Add(dtos.CreateEvent{
		TypeId:      4,
		Description: "New secret created, for account: " + request.Email + " with address: " + request.Website,
		CreatedAt:   time.Now().UTC(),
	})
	if err != nil {
		fmt.Println("Failed to create event abording")
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
		return
	}

	sc.SecretReoistory.Storage.Close()
	fmt.Println(event)

	ctx.JSON(http.StatusOK, created)
}

func (sc *SecretsController) content(ctx *gin.Context) {
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

	sc.SecretReoistory.Storage.Open()
	password, err := sc.SecretReoistory.Content(passwordId)
	sc.SecretReoistory.Storage.Close()

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Message": "Bad Request"})
		return
	}

	ctx.JSON(http.StatusOK, password)
}

func (sc *SecretsController) update(ctx *gin.Context) {

	request := &dtos.IncomingSecretsUpdateRequest{}

	if err := ctx.BindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
	}

	sc.SecretReoistory.Storage.Open()
	oldSecret, err := sc.SecretReoistory.Get(request.Id)

	if err != nil {
		fmt.Println("Record doesn't exist")
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
		return
	}

	event, err := sc.EventRepository.Add(dtos.CreateEvent{
		TypeId:      5,
		Description: "Password updated, from account: " + oldSecret.Email + " with address: " + oldSecret.Website + "to account: " + request.Email + " with address: " + request.Website,
		CreatedAt:   time.Now().UTC(),
	})

	if err != nil {
		fmt.Println("Failed to create an event aborting")
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
		return
	}

	fmt.Println(event)

	updated := sc.SecretReoistory.Update(request)
	sc.SecretReoistory.Storage.Close()

	ctx.JSON(http.StatusOK, updated)
}

func (sc *SecretsController) remove(ctx *gin.Context) {
	id := ctx.Request.FormValue("id")
	uuid, err := utils.ParseUUID(id)

	if err != nil {
		fmt.Println("id is not in a valid format")
		ctx.JSON(http.StatusBadRequest, gin.H{"Message": "Bad Request"})
		return
	}

	sc.SecretReoistory.Storage.Open()
	event, err := sc.EventRepository.Add(dtos.CreateEvent{
		TypeId:      6,
		Description: "Secret with id: " + uuid.String() + "has been removed",
		CreatedAt:   time.Now().UTC(),
	})

	if err != nil {
		fmt.Println("Failed to create an event aborting")
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
		return
	}
	fmt.Println(event)

	sc.SecretReoistory.Delete(uuid)
	sc.SecretReoistory.Storage.Close()

	ctx.JSON(http.StatusOK, "OK")
}

func (sc *SecretsController) Init(r *gin.RouterGroup, am *middlewhere.AuthenticationMiddlewhere) {
	controller := r.Group("secrets")
	controller.Use(am.Authorize())

	controller.GET("all", sc.all)
	controller.GET("content/:id", sc.content)
	controller.POST("add", sc.add)
	controller.POST("update", sc.update)
	controller.DELETE("remove", sc.remove)
}
