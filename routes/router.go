package routes

import (
	"github.com/gin-gonic/gin"

	implementations "leanmeal/api/Implementations"
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

	mfaController := &MfaController{
		TotpService: &implementations.TimeBasedService{},
		AccountsRepository: repositories.Accounts{
			Storage: r.Storage,
		},
		MfaRepository: repositories.MfaRepository{
			Storage: r.Storage,
		},
		EmailService: &implementations.MailService{
			From:     r.Configuration.GetKey("from").(string),
			Password: r.Configuration.GetKey("smtp-password").(string),
			SkipSSl:  r.Configuration.GetKey("ssl").(bool),
			Smtp:     r.Configuration.GetKey("smtp").(string),
			Port:     r.Configuration.GetKey("port").(int),
		},
	}

	passwordsController := &PasswordsController{
		PasswordService: r.PasswordService,
		PasswordRepository: repositories.PasswordRepository{
			Storage: r.Storage,
		},
		EventRepository: repositories.EventRepository{
			Storage: r.Storage,
		},
	}
	secretsController := &SecretsController{
		SecretReoistory: repositories.SecretsRepository{
			Storage: r.Storage,
		},
		EventRepository: repositories.EventRepository{
			Storage: r.Storage,
		},
	}
	identitiesController := &IdentitiesController{
		IdentityRepository: repositories.IdentityRepository{
			Storage: r.Storage,
		},
		EventsRepository: repositories.EventRepository{
			Storage: r.Storage,
		},
	}
	eventsController := &EventsController{
		EventsRepository: repositories.EventRepository{
			Storage: r.Storage,
		},
	}
	totpController := &TotpController{
		EventsRepository: repositories.EventRepository{
			Storage: r.Storage,
		},
		TotpRepository: repositories.TotpRepository{
			Storage: r.Storage,
		},
		TotpService: &implementations.TimeBasedService{},
	}

	authController.Init(r.V1)
	passwordsController.Init(r.V1, r.AuthMiddlewhere)
	secretsController.Init(r.V1, r.AuthMiddlewhere)
	identitiesController.Init(r.V1, r.AuthMiddlewhere)
	eventsController.Init(r.V1, r.AuthMiddlewhere)
	totpController.Init(r.V1, r.AuthMiddlewhere)
	mfaController.Init(r.V1, r.AuthMiddlewhere)
}
