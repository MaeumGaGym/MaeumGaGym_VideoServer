package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"pokabook/go-file-server/route"
	"pokabook/go-file-server/utils"
)

func main() {
	r := gin.Default()

	r.Use(cors.Default())
	route.SetupRoutes(r)

	utils.InitRabbitMQ()
	utils.InitMinio()

	r.Run(":9999")
}
