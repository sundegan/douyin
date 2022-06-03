package dao

import (
	"github.com/go-redis/cache/v8"
	"gorm.io/gorm"
	"log"
	"time"
)

// AfterFind 查询完后进行写缓存
func (u *User) AfterFind(tx *gorm.DB) (err error) {
	jsonUser, err := json.Marshal(*u)
	if err != nil {
		log.Println("json编码错误：", err)
	}
	err = LoginCache.Set(&cache.Item{
		Key:   u.Name,
		Value: jsonUser,
		TTL:   time.Minute,
	})
	if err != nil {
		log.Println("用户登录缓存失败:", err)
	}
	return nil
}
