package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"leanmeal/api/dtos"
	"leanmeal/api/middlewhere"
	"leanmeal/api/repositories"
)

type SecretsController struct {
	SecretReoistory repositories.SecretsRepository
}

func (sc *SecretsController) all(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "OK")
}

func (sc *SecretsController) add(ctx *gin.Context) {
	userId := ctx.MustGet("UUID")
	_ = userId

	request := &dtos.IncomingSecretsRequest{}
	if err := ctx.BindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
		return
	}

	ctx.JSON(http.StatusOK, "OK")
}

func (sc *SecretsController) update(ctx *gin.Context) {
	userId := ctx.MustGet("UUID")
	_ = userId

	request := &dtos.IncomingSecretsUpdateRequest{}

	if err := ctx.BindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
	}

	ctx.JSON(http.StatusOK, "OK")
}

func (sc *SecretsController) remove(ctx *gin.Context) {
	userId := ctx.MustGet("UUID")
	_ = userId

	request := &uuid.UUID{}

	if err := ctx.BindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
		return
	}

	ctx.JSON(http.StatusOK, "OK")
}

func (sc *SecretsController) Init(r *gin.RouterGroup, am *middlewhere.AuthenticationMiddlewhere) {
	controller := r.Group("Secrets")
	controller.Use(am.Authorize())

	controller.GET("all", sc.all)
	controller.POST("add", sc.add)
	controller.POST("update", sc.update)
	controller.DELETE("remove", sc.remove)
}
