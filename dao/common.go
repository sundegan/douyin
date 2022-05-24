package dao

import "time"

type Video struct {
	Id            int64     `json:"id,omitempty"`
	CreatedAt     time.Time `json:"-" gorm:"index"`
	AuthorId      int64     `json:"-"`
	Author        User      `json:"author"`
	PlayUrl       string    `json:"play_url,omitempty" gorm:"type:varchar(100)"`
	CoverUrl      string    `json:"cover_url,omitempty" gorm:"type:varchar(100)"`
	FavoriteCount int64     `json:"favorite_count,omitempty"`
	CommentCount  int64     `json:"comment_count,omitempty"`
	IsFavorite    bool      `json:"is_favorite,omitempty"`
}

type Comment struct {
	Id         int64  `json:"id,omitempty"`
	User       User   `json:"user"`
	Content    string `json:"content,omitempty"`
	CreateDate string `json:"create_date,omitempty"`
}

type User struct {
	Id            int64  `json:"id,omitempty" gorm:"primaryKey"`
	Name          string `json:"name,omitempty" gorm:"type:varchar(32); index"`
	Pwd           string `json:"-" gorm:"type:char(60)"`
	Salt          string `json:"-" gorm:"type:char(4)"`
	FollowCount   int64  `json:"follow_count,omitempty"`
	FollowerCount int64  `json:"follower_count,omitempty"`
	IsFollow      bool   `json:"is_follow,omitempty"`
}
