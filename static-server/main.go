package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()

	r.Static("/", "./")

	r.Run(":8080")
}
