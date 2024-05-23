package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	implementations "leanmeal/api/Implementations"
	"leanmeal/api/interfaces"
	"leanmeal/api/middlewhere"
	"leanmeal/api/routes"
)

func main() {
	// DI Service registration
	config := implementations.Configuration{}
	config.Load()
	connectionString := config.GetKey("ConnectionString").(string)

	jwt := implementations.JwtService{}
	jwt.Secret = config.GetKey("jwt-key").(string)
	jwt.Issuer = config.GetKey("jwt-issuer").(string)

	storage := implementations.Storage{
		ConnectionString: connectionString,
	}
	passwordService := implementations.PasswordService{}

	initializationService := implementations.Initialization{
		Storage: storage,
	}
	// Middlewhere setups
	authMiddlewhere := middlewhere.AuthenticationMiddlewhere{
		JwtService: &jwt,
	}
	cors := middlewhere.Cors()

	if !initializationService.Initialized() {
		initializationService.Database()
		initializationService.Seed()
	}

	//Init the server
	startServer(&config, &storage, &passwordService, &jwt, authMiddlewhere, &cors)
}

func startServer(configuration interfaces.Configuration, storage interfaces.Storage, passwordService interfaces.PasswordService,
	jwt interfaces.JwtService, authMiddlewhere middlewhere.AuthenticationMiddlewhere, cors *gin.HandlerFunc) {
	port := configuration.GetKey("Port").(string)

	router := gin.New()
	router.Use(*cors)
	v1 := router.Group("/v1")
	appRouter := routes.ApplicationRouter{
		Configuration:   configuration,
		Storage:         storage,
		PasswordService: passwordService,
		AuthMiddlewhere: &authMiddlewhere,
		Jwt:             jwt,
		V1:              v1,
	}

	appRouter.Init()

	srv := &http.Server{
		Addr:    port,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	got := <-quit
	fmt.Println(got)
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	connectionDone := <-ctx.Done()
	fmt.Println(connectionDone)
	log.Println("Server exiting")
}
