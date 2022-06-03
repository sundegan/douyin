package dao

import (
	jsoniter "github.com/json-iterator/go"
	"time"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

type Video struct {
	IsFavorite    bool      `json:"is_favorite,omitempty" gorm:"-"`
	Id            int64     `json:"id,omitempty" gorm:"primaryKey"`
	AuthorId      int64     `json:"-"`
	FavoriteCount int64     `json:"favorite_count,omitempty"`
	CommentCount  int64     `json:"comment_count,omitempty"`
	Title         string    `json:"title,omitempty" gorm:"type:varchar(100)"`
	PlayUrl       string    `json:"play_url,omitempty" gorm:"type:varchar(100)"`
	CoverUrl      string    `json:"cover_url,omitempty" gorm:"type:varchar(100)"`
	CreatedAt     time.Time `json:"-" gorm:"index"`

	Author User `json:"author"`
}

type Comment struct {
	Id         int64  `json:"id,omitempty"`
	User       User   `json:"user"`
	Content    string `json:"content,omitempty"`
	CreateDate string `json:"create_date,omitempty"`
}

type User struct {
	// id、密码、随机盐字段在返回给用户时应屏蔽
	IsFollow       bool   `json:"is_follow,omitempty" gorm:"-"`
	Id             int64  `json:"id,omitempty" gorm:"primaryKey"`
	FollowCount    int64  `json:"follow_count,omitempty"`
	FollowerCount  int64  `json:"follower_count,omitempty"`
	TotalFavorited int64  `json:"total_favorited,omitempty"`
	FavoriteCount  int64  `json:"favorite_count,omitempty"`
	Salt           string `json:"salt,omitempty" gorm:"type:char(4)"`
	Name           string `json:"name,omitempty" gorm:"type:varchar(32); index"`
	Pwd            string `json:"pwd,omitempty" gorm:"type:char(60)"`
}
