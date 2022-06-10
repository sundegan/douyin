package dao

import (
	"github.com/go-redis/cache/v8"
	"gorm.io/gorm"
	"log"
	"strconv"
	"time"
)

type User struct {
	// id、密码、随机盐字段在返回给用户时应屏蔽
	IsFollow       bool   `json:"is_follow,omitempty" gorm:"-"`
	Id             int64  `json:"id,omitempty" gorm:"primaryKey"`
	FollowCount    int64  `json:"follow_count,omitempty" gorm:"-"`
	FollowerCount  int64  `json:"follower_count,omitempty" gorm:"-"`
	TotalFavorited int64  `json:"total_favorited,omitempty"`
	FavoriteCount  int64  `json:"favorite_count,omitempty"`
	Salt           string `json:"salt,omitempty" gorm:"type:char(4)"`
	Name           string `json:"name,omitempty" gorm:"type:varchar(32); index"`
	Pwd            string `json:"pwd,omitempty" gorm:"type:char(60)"`
}

// AfterFind 查询完后进行写缓存
func (u *User) AfterFind(tx *gorm.DB) (err error) {
	// 登录情况下，携带密码和随机盐写入登陆缓存
	_, isLogin := tx.Get("login")
	if isLogin {
		user := *u
		// 除去与密码验证无关的字段
		user.TotalFavorited = 0
		user.FavoriteCount = 0
		user.Name = ""
		jsonUser, err := json.Marshal(user)
		if err != nil {
			log.Println("json编码错误：", err)
			// 继续后续缓存
			goto std
		}

		err = LoginCache.Set(&cache.Item{
			Key:   u.Name,
			Value: jsonUser,
			TTL:   10 * time.Second,
		})
		if err != nil {
			log.Println("用户登录缓存失败:", err)
		}
	}

std:
	// 敏感数据在保存到缓存前删除
	user := *u
	user.EraseSensitiveFiled()

	jsonUser, err := json.Marshal(user)
	if err != nil {
		log.Println("json编码错误：", err)
		return nil
	}

	err = UserCache.Set(&cache.Item{
		Key:   strconv.FormatInt(u.Id, 10),
		Value: jsonUser,
		TTL:   10 * time.Second,
	})
	if err != nil {
		log.Println("用户信息缓存失败:", err)
	}
	// 无论是否成功写缓存，继续完成事务
	return nil
}

// AfterCreate 创建完后进行写缓存
func (u *User) AfterCreate(tx *gorm.DB) (err error) {
	// 敏感数据在保存到缓存前删除
	user := *u
	user.EraseSensitiveFiled()

	jsonUser, err := json.Marshal(user)
	if err != nil {
		log.Println("json编码错误：", err)
		return nil
	}

	err = UserCache.Set(&cache.Item{
		Key:   strconv.FormatInt(u.Id, 10),
		Value: jsonUser,
		TTL:   10 * time.Second,
	})
	if err != nil {
		log.Println("用户信息缓存失败:", err)
	}
	// 无论是否成功写缓存，继续完成事务
	return nil
}

func (u *User) EraseSensitiveFiled() {
	u.Pwd = ""
	u.Salt = ""
}
