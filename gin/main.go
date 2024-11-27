package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()
	r.GET("ping", func(c *gin.Context) {
		c.JSON(200, "Hello World")
	})
	r.Run(":8080")
}