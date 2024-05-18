package routes

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	implementations "leanmeal/api/Implementations"
	"leanmeal/api/dtos"
	"leanmeal/api/interfaces"
	"leanmeal/api/middlewhere"
	"leanmeal/api/repositories"
)

type IdentitiesController struct {
	IdentityRepository repositories.IdentityRepository
	EventsRepository   repositories.EventRepository
	KeyManager         interfaces.KeyManager
}

func (ic *IdentitiesController) all(ctx *gin.Context) {
	ic.IdentityRepository.Storage.Open()
	identities, err := ic.IdentityRepository.All()
	ic.IdentityRepository.Storage.Close()

	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"Message": "Bad Request"})
		return
	}

	ctx.JSON(http.StatusOK, identities)
}

func (ic *IdentitiesController) types(ctx *gin.Context) {
	ic.IdentityRepository.Storage.Open()
	keyTypes, err := ic.IdentityRepository.GetKeyTypes()
	ic.IdentityRepository.Storage.Close()

	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"Message": "Bad Request"})
		return
	}

	ctx.JSON(http.StatusOK, keyTypes)

}

func (ic *IdentitiesController) add(ctx *gin.Context) {
	account := ctx.MustGet("ID")

	request := &dtos.CreateIdentity{}
	if err := ctx.BindJSON(request); err != nil {
		fmt.Println("could not bind request body")
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
		return
	}

	if request.KeyType > 2 {
		ctx.JSON(http.StatusBadRequest, gin.H{"Message": "Bad request parameter."})
	}

	switch request.KeyType {
	case 1:
		ic.KeyManager = &implementations.ED25519Key{}
	case 2:
		ic.KeyManager = &implementations.RSAKey{}
	}

	public, private, err := ic.KeyManager.Generate(request.KeySize)

	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"Message": "Bad Request"})
		return
	}

	ic.IdentityRepository.Storage.Open()
	created, err := ic.IdentityRepository.Add(*request, &public, &private, account.(uuid.UUID))
	if err != nil {
		ic.IdentityRepository.Storage.Close()
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
		return
	}

	event, err := ic.EventsRepository.Add(dtos.CreateEvent{
		TypeId:      1,
		Description: "New key pair generated, with name: " + request.Name,
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
	var public []byte
	var private []byte

	oldIdentity, err := ic.IdentityRepository.GetInternal(request.Id)

	if err != nil {
		fmt.Println("Record doesn't exist")
		ic.EventsRepository.Storage.Close()
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
		return
	}

	public, err = base64.StdEncoding.DecodeString(oldIdentity.PublicKey)
	if err != nil {
		ic.EventsRepository.Storage.Close()
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
		return
	}
	private, err = base64.StdEncoding.DecodeString(oldIdentity.PrivateKey)
	if err != nil {
		ic.EventsRepository.Storage.Close()
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
		return
	}

	switch request.KeyType {
	case 1:
		ic.KeyManager = &implementations.ED25519Key{}
	case 2:
		ic.KeyManager = &implementations.RSAKey{}
	}

	if request.RegenerateKey {
		public, private, err = ic.KeyManager.Generate(request.KeySize)

		if err != nil {
			fmt.Println(err)
			ic.EventsRepository.Storage.Close()
			ctx.JSON(http.StatusBadRequest, gin.H{"Message": "Bad Request"})
			return
		}
	}

	updated := ic.IdentityRepository.Update(request, &public, &private)

	if !updated {
		fmt.Println("Failed to update entity, aborting")
		ic.EventsRepository.Storage.Close()
		ctx.JSON(http.StatusBadRequest, gin.H{"Message": "Bad Request"})
		return
	}

	event, err := ic.EventsRepository.Add(dtos.CreateEvent{
		TypeId:      2,
		Description: "Key updated, name set to: " + oldIdentity.Name,
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
	controller := r.Group("identities")
	controller.Use(a.Authorize())

	controller.GET("all", ic.all)
	controller.GET("types", ic.types)
	controller.POST("add", ic.add)
	controller.POST("update", ic.update)

}
