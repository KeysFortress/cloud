package routes

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"leanmeal/api/dtos"
	"leanmeal/api/middlewhere"
	"leanmeal/api/repositories"
)

type IdentitiesController struct {
	IdentityRepository repositories.IdentityRepository
	EventsRepository   repositories.EventRepository
}

func (ic *IdentitiesController) all(ctx *gin.Context) {
	return
}

func (ic *IdentitiesController) types(ctx *gin.Context) {
	return
}

func (ic *IdentitiesController) add(ctx *gin.Context) {
	request := &dtos.CreateIdentity{}
	if err := ctx.BindJSON(request); err != nil {
		fmt.Println("could not bind request body")
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
		return
	}

	ic.IdentityRepository.Storage.Open()
	created, err := ic.IdentityRepository.Add(*request)
	if err != nil {
		ic.IdentityRepository.Storage.Close()
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
		return
	}

	event, err := ic.EventsRepository.Add(dtos.CreateEvent{
		TypeId:      1,
		Description: "New time based password created, for account: " + request.Email + " with address: " + request.Website,
		CreatedAt:   time.Now().UTC(),
	})
	if err != nil {
		fmt.Println("Failed to create event abording")
		ic.EventsRepository.Storage.Close()
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
		return
	}

	ic.EventsRepository.Storage.Close()

	fmt.Println(event)
	ctx.JSON(http.StatusOK, created)
}

func (ic *IdentitiesController) update(ctx *gin.Context) {
	id := ctx.MustGet("ID")
	_ = id

	request := &dtos.UpdateIdentity{}
	if err := ctx.BindJSON(request); err != nil {
		fmt.Println("could not bind request body")
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
		return
	}

	ic.EventsRepository.Storage.Open()
	oldPassword, err := ic.IdentityRepository.Get(request.Id)
	if err != nil {
		fmt.Println("Record doesn't exist")
		ic.EventsRepository.Storage.Close()
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
		return
	}

	updated := ic.IdentityRepository.Update(request)

	if !updated {
		fmt.Println("Failed to update entity, aborting")
		ic.EventsRepository.Storage.Close()
		ctx.JSON(http.StatusBadRequest, gin.H{"Message": "Bad Request"})
		return
	}

	event, err := ic.EventsRepository.Add(dtos.CreateEvent{
		TypeId:      2,
		Description: "TOTP updated, from account: " + oldPassword.Email + " with address: " + oldPassword.Website + "to account: " + request.Email + " with address: " + request.Website,
		CreatedAt:   time.Now().UTC(),
	})

	if err != nil {

		fmt.Println("Entity was updated, but failed to create an event.")
		ic.IdentityRepository.Storage.Close()
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})

		return
	}

	fmt.Println(event)
	ic.IdentityRepository.Storage.Close()
	ctx.JSON(http.StatusOK, updated)
}

func (ic *IdentitiesController) Init(r *gin.RouterGroup, a *middlewhere.AuthenticationMiddlewhere) {
	controller := r.Group("totp")
	controller.Use(a.Authorize())

	controller.GET("all", ic.all)
	controller.GET("types", ic.types)
	controller.POST("add", ic.add)
	controller.POST("update", ic.update)

}
