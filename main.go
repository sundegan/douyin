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
	rand.Seed(time.Now().UnixNano())

	dao.InitDB()
	service.InitUser()

	r := gin.Default()

	initRouter(r)

	r.Run(":8080") // http默认端口

	err := dao.RDB.Close()
	log.Fatal(err.Error())
}
