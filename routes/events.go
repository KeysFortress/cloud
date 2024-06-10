package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"leanmeal/api/dtos"
	"leanmeal/api/middlewhere"
	"leanmeal/api/repositories"
)

type EventsController struct {
	EventsRepository repositories.EventRepository
}

func (ec *EventsController) all(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "Do something"})
}

func (ec *EventsController) types(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "Do something"})
}

func (ec *EventsController) add(ctx *gin.Context) {
	deviceId := ctx.MustGet("DeviceKey") //public key of the device
	_ = deviceId

	request := &dtos.CreateEventRequest{}
	if err := ctx.BindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Do something"})
}

func (ec *EventsController) eventsByType(ctx *gin.Context) {
	var request int
	if err := ctx.Bind(request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Message": "Event Type doesn't exist"})
		return
	}

	ec.EventsRepository.Storage.Open()
	defer ec.EventsRepository.Storage.Close()

	events, err := ec.EventsRepository.GetByType(request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Message": "Internal Error"})
		return
	}

	ctx.JSON(http.StatusOK, events)
}

func (ec *EventsController) Init(r *gin.RouterGroup, authMiddlewhere *middlewhere.AuthenticationMiddlewhere) {
	controller := r.Group("events")
	controller.Use(authMiddlewhere.Authorize())

	controller.GET("all", ec.all)
	controller.GET("events-by-type/:type", ec.eventsByType)
	controller.GET("event-types", ec.types)
	controller.POST("add", ec.add)
}
