package service

import (
	"douyin-server/dao"
	"errors"
	"gorm.io/gorm"
)

func Register(username, password string) (int64, error) {
	if len(username) > 32 {
		return 0, errors.New("用户名过长，不可超过32位")
	}
	if len(password) > 32 {
		return 0, errors.New("密码过长，不可超过32位")
	}

	user := dao.User{}
	dao.DB.Where("name = ?", username).Find(&user)
	if user.Id != 0 {
		return 0, errors.New("用户已存在")
	}

	user.Name, user.Pwd = username, password
	dao.DB.Create(&user)
	return user.Id, nil
}

func Login(username, password string) (int64, error) {
	user := dao.User{}
	err := dao.DB.Where("name = ?", username).Find(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, errors.New("用户不存在")
	} else if password != user.Pwd {
		return 0, errors.New("密码错误")
	} else {
		return user.Id, nil
	}
}

func UserInfo(token string) (dao.User, error) {
	user := dao.User{}
	err := dao.DB.Where("name = ?", token).Find(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return user, errors.New("用户不存在")
	}
	return user, nil
}
