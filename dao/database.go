package dao

import (
	"context"
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"time"
)

var (
	DB *gorm.DB

	RDB        *redis.Client
	Ctx        context.Context
	LoginCache *cache.Cache
)

// Redis数据库编号
const (
	numTokenDB = iota
	numLoginCacheDB
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

	err = DB.AutoMigrate(&User{}, &Video{})
	log.Println(err)

	RDB = redis.NewClient(&redis.Options{
		Addr:     "192.168.200.128:7000",
		Password: "zxc05020519",
		DB:       numTokenDB,
	})
	Ctx = context.Background()

	LoginCache = cache.New(&cache.Options{
		Redis: redis.NewClient(&redis.Options{
			Addr:     "192.168.200.128:7000",
			Password: "zxc05020519",
			DB:       numLoginCacheDB,
		}),
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})
}