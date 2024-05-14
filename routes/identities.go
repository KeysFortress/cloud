package routes

import (
	"github.com/gin-gonic/gin"

	"leanmeal/api/middlewhere"
	"leanmeal/api/repositories"
)

type IdentitiesController struct {
	IdentityRepository repositories.IdentityRepository
}

func (ic *IdentitiesController) Init(r *gin.RouterGroup, a *middlewhere.AuthenticationMiddlewhere) {

}
