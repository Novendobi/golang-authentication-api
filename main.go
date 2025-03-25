package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)


func main(){
	err := godotenv.Load()
	if err != nil {
		log.Fatal("ERROR: Could not load enviroment, create a '.env' file", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router := gin.Default()
	router.GET("home/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Hello... Welcome to our API",
		})
	})
	router.Run()
}