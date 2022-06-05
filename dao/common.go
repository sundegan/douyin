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
	Id         int64  `json:"id,omitempty" gorm:"primaryKey"`
	UserId     int64  `json:"-"`
	Content    string `json:"content,omitempty" gorm:"type:text"`
	CreateDate string `json:"create_date,omitempty" gorm:"type:varchar(20)"`
	VideoId    int64  `json:"video_id"`

	User User `json:"user"`
}

type Favorite struct {
	UserId  int64 `json:"-" gorm:"primaryKey;autoIncrement:false"`
	VideoId int64 `json:"-" gorm:"primaryKey;autoIncrement:false"`
}
