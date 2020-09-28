package main

import (
	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func main() {
	router = gin.Default()

	MongoInit()

	router.LoadHTMLGlob("assets/html/*")
	router.Static("/assets", "D:/Golang Tasks/Code/Server/assets/")

	InitialiseRouters()

	err := router.Run(":8181")
	if err != nil {
		println(err.Error())
		return
	}

}
