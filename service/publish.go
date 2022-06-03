package service

import (
	"douyin-server/dao"
)

const staticRouter = "http://10.0.2.2:8080/static/"

func Publish(filename, title string, user_id int64) error {

	v := dao.Video{
		AuthorId: user_id,
		Title:    title,
		PlayUrl:  staticRouter + "videos/" + filename,
		CoverUrl: staticRouter + "covers/" + filename,
	}
	if err := dao.DB.Create(&v).Error; err != nil {
		return err
	}

	return nil
}
