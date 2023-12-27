package main

import (
	"github.com/gin-gonic/gin"
	"pokabook/go-file-server/route"
	"pokabook/go-file-server/utils"
)

func main() {
	r := gin.Default()

	route.SetupRoutes(r)

	utils.InitRabbitMQ()
	utils.InitMinio()

	r.Run(":9999")
}
