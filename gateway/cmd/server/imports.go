package main

import (
	"github.com/sdq-codes/claimclam/internal/controllers"
	"github.com/sdq-codes/claimclam/internal/events"
	"github.com/sdq-codes/claimclam/internal/repositories"
	"github.com/sdq-codes/claimclam/internal/services"
	httptransport "github.com/sdq-codes/claimclam/internal/transport/http"
	"gorm.io/gorm"
)

func AllImports(db *gorm.DB) *httptransport.Server {
	userRepositories := repositories.NewUserRepository(db)
	e := events.New()
	userService := services.NewUserService(userRepositories, e)
	userController := controllers.NewUserController(userService, e)
	podcastService := services.NewPodcastService(e)
	podcastController := controllers.NewPodcastController(podcastService)
	httpServer := httptransport.New(userController, podcastController, e)
	return httpServer
}
