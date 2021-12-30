package main

import (
	"github.com/gin-gonic/gin"
	"netthrough.server/controllers"
)

func main() {
	r := gin.Default()
	r.POST("/register", controllers.Register)
	r.POST("/statuscheck", controllers.StatusCheck)
	r.POST("/readdata", controllers.ReadData)
	r.POST("/writedata", controllers.WriteData)
	r.POST("/unregister", controllers.UnRegister)
	panic(r.Run(":5001"))
}
