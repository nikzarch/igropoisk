package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"io"
	"log"
	"os"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalln("Error loading .env file")
	}
	r := gin.New()
	f, _ := os.Create("log.txt")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	r.Use(gin.LoggerWithWriter(gin.DefaultWriter))
	r.Use(gin.Recovery())
	r.Use(cors.Default()) //temp
}
