package dao

import (
	"context"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB

	RDB        *redis.Client
	RDB_FOLLOW *redis.Client // db1,存放关注列表
	RDB_FANS   *redis.Client // db2,存放粉丝列表

	Ctx context.Context
)

// InitDB 数据库连接配置
func InitDB() {
	var err error

	// MySQL连接
	dsn := "sundegan:SunDeGan1998@tcp(47.107.99.130:3306)/" +
		"douyin?charset=utf8&interpolateParams=true&parseTime=True&loc=Local"

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		PrepareStmt: true,
	})
	if err != nil {
		panic("数据库连接失败,error:" + err.Error())
	}

	_ = DB.AutoMigrate(&User{}, &Video{})

	// Redis连接
	RDB = redis.NewClient(&redis.Options{
		Addr:     "47.107.99.130:6379",
		Password: "",
		DB:       0,
	})
	// 使用DB1作为关注列表数据库
	RDB_FOLLOW = redis.NewClient(&redis.Options{
		Addr:     "47.107.99.130:6379",
		Password: "",
		DB:       1,
	})
	// 使用DB2作为粉丝列表数据库
	RDB_FANS = redis.NewClient(&redis.Options{
		Addr:     "47.107.99.130:6379",
		Password: "",
		DB:       2,
	})

	Ctx = context.Background()
}
