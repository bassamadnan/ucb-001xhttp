package main

import (
	"github.com/gin-gonic/gin"
)

func GetDummyEndpoint(c *gin.Context) {
	resp := map[string]string{"hello": "world"}
	c.JSON(200, resp)
}

func main() {
	api := gin.Default()
	api.GET("/dummy", GetDummyEndpoint)
	api.Run(":5000")
}
