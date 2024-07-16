package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/kawijayaa/pintas/api"
	"github.com/kawijayaa/pintas/frontend"
)

func main() {
	godotenv.Load()
	server := gin.Default()

	frontend.Routes(server)
	api.Routes(server)

	server.Run(":8080")
}
