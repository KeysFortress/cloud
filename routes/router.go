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
	domain := r.Configuration.GetKey("domain")
	if domain == nil {
		panic("Domain is not set in the configuration file")
	}

	authService := &implementations.AuthenticationService{
		Domain: domain.(string),
	}
	go authService.Start()

	authController := &AuthenticationController{
		AuthenticationService: authService,
		JwtService:            r.Jwt,
		Configuration:         r.Configuration,
		AccountRepository: repositories.Accounts{
			Storage: r.Storage,
		},
		AccessKeysRepository: repositories.AccessKeysRepository{
			Storage: r.Storage,
		},
	}
	mfaController := &MfaController{
		TotpService: &implementations.TimeBasedService{
			Issuer: r.Configuration.GetKey("jwt-issuer").(string),
		},
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
			Port:     int(r.Configuration.GetKey("smtp-port").(float64)),
		},
		JwtService: r.Jwt,
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
	setupController := &SetupController{
		accountRepository: repositories.Accounts{
			Storage: r.Storage,
		},
		accessKeysRepository: repositories.AccessKeysRepository{
			Storage: r.Storage,
		},
		setupPath:             "v1/setup/finish",
		domain:                r.Configuration.GetKey("domain").(string),
		authenticationService: authService,
	}

	authController.Init(r.V1)
	passwordsController.Init(r.V1, r.AuthMiddlewhere)
	secretsController.Init(r.V1, r.AuthMiddlewhere)
	identitiesController.Init(r.V1, r.AuthMiddlewhere)
	eventsController.Init(r.V1, r.AuthMiddlewhere)
	totpController.Init(r.V1, r.AuthMiddlewhere)
	mfaController.Init(r.V1, r.AuthMiddlewhere)
	setupController.Init(r.V1)
}
