package main

import (
	"douyin-server/dao"
	"github.com/gin-gonic/gin"
)

func main() {
	dao.InitDB()

	r := gin.Default()

	initRouter(r)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
