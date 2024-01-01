package main

//import (
//	httptransport "github.com/sdq-codes/ccasses/gateway/internal/transport/http"
//	"gorm.io/gorm"
//)
//
//func AllImports(db *gorm.DB) *httptransport.Server {
//	userRepositories := repositories.NewUserRepository(db)
//	e := events.New()
//	userService := services.NewUserService(userRepositories, e)
//	userController := controllers.NewUserController(userService, e)
//	podcastService := services.NewPodcastService(e)
//	podcastController := controllers.NewPodcastController(podcastService)
//	httpServer := httptransport.New(userController, podcastController, e)
//	return httpServer
//}
