package main

import (
	"douyin-server/dao"
	"douyin-server/service"
	"github.com/gin-gonic/gin"
	"log"
	"math/rand"
	"time"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	rand.Seed(time.Now().UnixNano())

	dao.InitDB()
	service.InitUser()

	r := gin.Default()

	initRouter(r)

	err := r.Run(":80") // http默认端口
	if err != nil {
		panic(err)
	}

	err = dao.RdbToken.Close()
	log.Fatal(err.Error())
}
