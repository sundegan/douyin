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
	RdbFollow  *redis.Client // 存放关注列表
	RdbFans    *redis.Client // 存放粉丝列表
)

// Redis数据库编号
const (
	numTokenDB = iota
	numLoginCacheDB
	numFollowListDB
	numFollowerListDB
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

	Ctx = context.Background()

	err = DB.AutoMigrate(&User{}, &Video{}, &Favorite{}, &Comment{})
	log.Println(err)

	RDB = redis.NewClient(&redis.Options{
		Addr:     "192.168.200.128:7000",
		Password: "zxc05020519",
		DB:       numTokenDB,
	})

	LoginCache = cache.New(&cache.Options{
		Redis: redis.NewClient(&redis.Options{
			Addr:     "192.168.200.128:7000",
			Password: "zxc05020519",
			DB:       numLoginCacheDB,
		}),
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})

	// 关注列表数据库
	RdbFollow = redis.NewClient(&redis.Options{
		Addr:     "192.168.200.128:7000",
		Password: "zxc05020519",
		DB:       numFollowListDB,
	})
	// 粉丝列表数据库
	RdbFans = redis.NewClient(&redis.Options{
		Addr:     "192.168.200.128:7000",
		Password: "zxc05020519",
		DB:       numFollowerListDB,
	})
}
