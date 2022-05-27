package dao

import (
	"context"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB

	RDB *redis.Client
	Ctx context.Context
)

func InitDB() {
	var err error
	dsn := "root:zxc05020519@tcp(192.168.200.128:23306)/" +
		"douyin?charset=utf8&interpolateParams=true&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		PrepareStmt: true,
	})
	if err != nil {
		panic(err)
	}

	DB.AutoMigrate(&User{}, &Video{})

	RDB = redis.NewClient(&redis.Options{
		Addr:     "192.168.200.128:7000",
		Password: "zxc05020519",
		DB:       0,
	})
	Ctx = context.Background()
}
