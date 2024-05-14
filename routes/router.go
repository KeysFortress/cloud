package routes

import (
	"github.com/gin-gonic/gin"

	"leanmeal/api/interfaces"
	"leanmeal/api/middlewhere"
	"leanmeal/api/repositories"
)

type ApplicationRouter struct {
	Configuration   interfaces.Configuration
	Storage         interfaces.Storage
	PasswordService interfaces.PasswordService
	AuthMiddlewhere *middlewhere.AuthenticationMiddlewhere
	Jwt             interfaces.JwtService
	V1              *gin.RouterGroup
}

func (r *ApplicationRouter) Init() {

	authController := &AuthenticationController{
		JwtService:    r.Jwt,
		Configuration: r.Configuration,
		AccountRepository: repositories.Accounts{
			Storage: r.Storage,
		},
		AccessKeysRepository: repositories.AccessKeysRepository{
			Storage: r.Storage,
		},
	}
	passwordsController := &PasswordsController{
		PasswordService: r.PasswordService,
		PasswordRepository: repositories.PasswordRepository{
			Storage: r.Storage,
		},
	}
	secretsController := &SecretsController{
		SecretReoistory: repositories.SecretsRepository{
			Storage: r.Storage,
		},
	}
	identitiesController := &IdentitiesController{
		IdentityRepository: repositories.IdentityRepository{
			Storage: r.Storage,
		},
	}
	eventsController := &EventsController{
		EventsRepository: repositories.EventRepository{
			Storage: r.Storage,
		},
	}

	authController.Init(r.V1)
	passwordsController.Init(r.V1, r.AuthMiddlewhere)
	secretsController.Init(r.V1, r.AuthMiddlewhere)
	identitiesController.Init(r.V1, r.AuthMiddlewhere)
	eventsController.Init(r.V1, r.AuthMiddlewhere)
}
