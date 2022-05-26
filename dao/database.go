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
	dsn := "debian-sys-maint:2o2U5fxdnadNXZhq@tcp(localhost:3306)/" +
		"douyin?charset=utf8&interpolateParams=true&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		PrepareStmt: true,
	})
	if err != nil {
		panic(err)
	}

	DB.AutoMigrate(&User{}, &Video{})

	RDB = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	Ctx = context.Background()
}
