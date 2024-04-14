package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"leanmeal/api/interfaces"
	"leanmeal/api/models"
	"leanmeal/api/repositories"
)

type IngridientsControlller struct {
	Storage              interfaces.Storage
	ingridientRepository repositories.IngridientRepository
}

func (ic *IngridientsControlller) add(ctx *gin.Context) {
	request := &models.Ingridient{}
	if err := ctx.BindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
		return
	}
	ic.ingridientRepository.OpenConnection(&ic.Storage)
	ingridient, err := ic.ingridientRepository.Add(request)
	ic.ingridientRepository.Close()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
	}

	ctx.JSON(http.StatusOK, ingridient)
}

func (ic *IngridientsControlller) Init(r *gin.RouterGroup) {
	r.POST("/ingridients/add", ic.add)
}
