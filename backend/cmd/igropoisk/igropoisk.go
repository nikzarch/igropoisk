package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"igropoisk_backend/internal/db"
	"igropoisk_backend/internal/game"
	"igropoisk_backend/internal/middleware"
	"igropoisk_backend/internal/review"
	"igropoisk_backend/internal/user"
	"io"
	"log"
	"os"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalln("Error loading .env file")
	}
	db := db.GetConnection()
	defer db.Close()

	userRepo := user.NewPostgresRepository(db)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)

	gameRepo := game.NewPostgresRepository(db)
	gameService := game.NewService(gameRepo)
	gameHandler := game.NewHandler(gameService)

	reviewRepo := review.NewPostgresRepository(db)
	reviewService := review.NewService(reviewRepo, gameService)
	reviewHandler := review.NewHandler(reviewService)
	r := gin.New()
	f, _ := os.Create("log.txt")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	r.Use(gin.LoggerWithWriter(gin.DefaultWriter))
	r.Use(gin.Recovery())
	r.Use(cors.Default()) //temp

	api := r.Group("api")
	{
		api.POST("register", userHandler.HandleRegistration)
		api.POST("login", userHandler.HandleLogin)
		api.GET("games/:id/reviews", reviewHandler.GetReviewsByGameId)
	}
	authorizedApi := r.Group("api", middleware.AuthMiddleware())
	{
		authorizedApi.GET("games/:id", gameHandler.GetGameByID)
		authorizedApi.GET("games", gameHandler.GetAllGames)
		authorizedApi.POST("games", gameHandler.AddGame)
		authorizedApi.DELETE("games/:id", gameHandler.DeleteGameByID)

		authorizedApi.POST("games/:id/reviews", reviewHandler.AddReview)
	}

	r.Run("localhost:" + os.Getenv("PORT"))
}
