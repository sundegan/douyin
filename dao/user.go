package dao

import (
	"context"
	"fmt"
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

// BeforeSave 进行延迟双删第一删
func (u *User) BeforeSave(tx *gorm.DB) (err error) {
	u.deleteFromCache()
	return nil
}

// AfterSave 进行延迟双删第二删
func (u *User) AfterSave(tx *gorm.DB) (err error) {
	go func() {
		time.Sleep(200 * time.Millisecond)
		u.deleteFromCache()
	}()
	return nil
}

// BeforeUpdate 进行延迟双删第一删
func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	u.deleteFromCache()
	return nil
}

// AfterUpdate 进行延迟双删第二删
func (u *User) AfterUpdate(tx *gorm.DB) (err error) {
	go func() {
		time.Sleep(200 * time.Millisecond)
		u.deleteFromCache()
	}()
	return nil
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
			TTL:   30 * time.Second,
		})
		if err != nil {
			log.Println("用户登录缓存失败:", err)
		}
	}

std:
	u.saveIntoCache()

	// 无论是否成功写缓存，继续完成事务
	return nil
}

// AfterCreate 创建完后进行写缓存
func (u *User) AfterCreate(tx *gorm.DB) (err error) {
	u.saveIntoCache()

	// 无论是否成功写缓存，继续完成事务
	return nil
}

// 将除密码和随机盐外的字段放入缓存中
func (u *User) saveIntoCache() {
	if u.Name != "" {
		err := UserCache.Set(&cache.Item{
			Key:   fmt.Sprintf("%d:Name", u.Id),
			Value: u.Name,
			TTL:   time.Minute,
		})
		if err != nil {
			log.Println("用户信息缓存失败:", err)
		}
	}

	if u.TotalFavorited != 0 {
		err := UserCache.Set(&cache.Item{
			Key:   fmt.Sprintf("%d:TotalFavorited", u.Id),
			Value: strconv.FormatInt(u.TotalFavorited, 10),
			TTL:   time.Minute,
		})
		if err != nil {
			log.Println("用户信息缓存失败:", err)
		}
	}

	if u.FavoriteCount != 0 {
		err := UserCache.Set(&cache.Item{
			Key:   fmt.Sprintf("%d:FavoriteCount", u.Id),
			Value: strconv.FormatInt(u.FavoriteCount, 10),
			TTL:   time.Minute,
		})
		if err != nil {
			log.Println("用户信息缓存失败:", err)
		}
	}
}

func (u *User) deleteFromCache() {
	_ = UserCache.Delete(context.Background(), fmt.Sprintf("%d:Name", u.Id))
	_ = UserCache.Delete(context.Background(), fmt.Sprintf("%d:TotalFavorited", u.Id))
	_ = UserCache.Delete(context.Background(), fmt.Sprintf("%d:FavoriteCount", u.Id))
}
