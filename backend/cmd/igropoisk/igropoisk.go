package main

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"igropoisk_backend/internal/auth"
	"igropoisk_backend/internal/db/elastic"
	"igropoisk_backend/internal/db/postgres"
	"igropoisk_backend/internal/game"
	"igropoisk_backend/internal/game/genre"
	"igropoisk_backend/internal/logger"
	"igropoisk_backend/internal/middleware"
	"igropoisk_backend/internal/review"
	"igropoisk_backend/internal/user"
	"log"
	"os"
)

func main() {
	auth.Init()
	postgresPool := postgres.GetPool()
	elasticClient := elastic.NewClient()
	defer postgresPool.Close()

	userRepo := user.NewPostgresRepository(postgresPool)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)

	gameRepo := game.NewPostgresRepository(postgresPool)
	genreRepo := genre.NewPostgresRepository(postgresPool)
	searchRepo := game.NewElasticRepository(elasticClient)

	gameService := game.NewService(gameRepo, genreRepo, searchRepo)
	gameHandler := game.NewHandler(gameService)

	reviewRepo := review.NewPostgresRepository(postgresPool)
	reviewService := review.NewService(reviewRepo, gameService)
	reviewHandler := review.NewHandler(reviewService)
	r := gin.New()
	err := logger.InitLogger()
	defer logger.CloseFile()
	if err != nil {
		log.Printf("failed to init logger : %s\n", err.Error())
	}
	err = searchRepo.SyncWith(context.Background(), gameRepo)
	if err != nil {
		logger.Logger.Error("Unable to sync elastic with TS",
			"error", err)
	}
	r.Use(logger.SlogMiddleware())
	r.Use(gin.Recovery())
	r.Use(cors.Default()) //temp

	api := r.Group("api")
	{
		api.POST("register", userHandler.HandleRegistration)
		api.POST("login", userHandler.HandleLogin)
		api.GET("games/:id/reviews", reviewHandler.GetReviewsByGameID)
	}
	authorizedApi := r.Group("api", middleware.AuthMiddleware())
	{
		authorizedApi.GET("games/search", gameHandler.SearchGame)
		authorizedApi.GET("games/:id", gameHandler.GetGameByID)
		authorizedApi.GET("games", gameHandler.GetAllGames)
		authorizedApi.POST("games", gameHandler.AddGame)
		authorizedApi.DELETE("games/:id", gameHandler.DeleteGameByID)

		authorizedApi.POST("games/:id/reviews", reviewHandler.AddReview)
	}

	r.Run(":" + os.Getenv("PORT"))
}
