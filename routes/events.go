package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"leanmeal/api/dtos"
	"leanmeal/api/middlewhere"
)

type EventsController struct {
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

func (ec *EventsController) Init(r *gin.RouterGroup, authMiddlewhere *middlewhere.AuthenticationMiddlewhere) {
	controller := r.Group("passwords")
	controller.Use(authMiddlewhere.Authorize())

	controller.GET("all", ec.all)
	controller.GET("event-types", ec.types)
	controller.POST("add", ec.add)
}
