package main

import (
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"netthrough.server/controllers"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard
	r := gin.Default()
	r.POST("/register", controllers.Register)
	r.POST("/statuscheck", controllers.StatusCheck)
	r.POST("/readdata", controllers.ReadData)
	r.POST("/writedata", controllers.WriteData)
	r.POST("/unregister", controllers.UnRegister)
	panic(r.Run(":5001"))
}
