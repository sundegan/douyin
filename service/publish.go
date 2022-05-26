package service

import (
	"douyin-server/dao"
)

const staticRouter = "http://172.25.151.124:3000/static/"

func Publish(filename, title string, user_id int64) error {

	v := dao.Video{
		AuthorId: user_id,
		PlayUrl:  staticRouter + "videos/" + filename,
		CoverUrl: staticRouter + "covers/" + filename,
	}
	if err := dao.DB.Create(&v).Error; err != nil {
		return err
	}

	return nil
}
