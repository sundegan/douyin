package dao

import (
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"time"
)

var (
	DB *gorm.DB

	LoginCache *cache.Cache
	UserCache  *cache.Cache
	RdbToken   *redis.Client
	RdbFollow  *redis.Client // 存放关注列表
	RdbFans    *redis.Client // 存放粉丝列表
)

// Redis数据库编号
const (
	numTokenDB = iota
	numLoginCacheDB
	numUserCacheDB
	numFollowListDB
	numFollowerListDB
)

func InitDB() {
	var err error
	dsn := "douyin_server:@tcp(localhost:3306)/" +
		"douyin?charset=utf8mb4&interpolateParams=true&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		PrepareStmt: true,
	})
	if err != nil {
		panic(err)
	}

	err = DB.AutoMigrate(&User{}, &Video{}, &Favorite{}, &Comment{})
	log.Println(err)

	RdbToken = redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    "mymaster",
		SentinelAddrs: []string{":17000", ":17001", ":17002"},
		DB:            numTokenDB,
	})

	LoginCache = cache.New(&cache.Options{
		Redis: redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:    "mymaster",
			SentinelAddrs: []string{":17000", ":17001", ":17002"},
			DB:            numLoginCacheDB,
		}),
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})

	UserCache = cache.New(&cache.Options{
		Redis: redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:    "mymaster",
			SentinelAddrs: []string{":17000", ":17001", ":17002"},
			DB:            numUserCacheDB,
		}),
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})

	// 关注列表数据库
	RdbFollow = redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    "mymaster",
		SentinelAddrs: []string{":17000", ":17001", ":17002"},
		DB:            numFollowListDB,
	})
	// 粉丝列表数据库
	RdbFans = redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    "mymaster",
		SentinelAddrs: []string{":17000", ":17001", ":17002"},
		DB:            numFollowerListDB,
	})
}
