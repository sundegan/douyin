package controller

import (
	"douyin-server/dao"
	"douyin-server/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

// 这块仍是官方demo内容，现在已经可以查数据库
// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
// test data: username=zhanglei, password=douyin
var usersLoginInfo = map[string]dao.User{
	"zhangleidouyin": {
		Id:            1,
		Name:          "zhanglei",
		FollowCount:   10,
		FollowerCount: 5,
		IsFollow:      true,
	},
}

type UserLoginResponse struct {
	Response
	Token int64 `json:"token"`
}

type UserResponse struct {
	Response
	User dao.User `json:"user"`
}

// LoginLimit 中间件，限制注册登录操作过于频繁。
func LoginLimit(c *gin.Context) {
	ipAddress := c.ClientIP()
	ok := service.LoginLimit(ipAddress)
	if !ok {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "操作过于频繁，请稍后再试"},
		})
		c.Abort()
	}
}

// VerifyToken 验证token中间件，成功会将用户id写入gin上下文中，否则会直接拦截请求
func VerifyToken(c *gin.Context) {
	token, ok := c.GetQuery("token")
	// 误用中间件，直接跳过
	if !ok {
		return
	}

	ipAddress := c.ClientIP()
	// 错误可忽略
	times, _ := dao.RDB.Get(dao.Ctx, ipAddress).Int64()
	if times > 10 {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "操作过于频繁，请稍后再试"})
		c.Abort()
		return
	}

	_id, err := dao.RDB.Get(dao.Ctx, token).Result()
	if err == nil {
		// token续期
		dao.RDB.Expire(dao.Ctx, token, 12*time.Hour)
	} else {
		// token不存在，记录该ip此次访问
		dao.RDB.Set(dao.Ctx, ipAddress, times+1, time.Minute)
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "登陆已过期，请重新登陆"})
		c.Abort()
		return
	}

	id, err := strconv.ParseInt(_id, 10, 64)
	if err != nil {
		log.Println("出现无法解析成64位整数的token")
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "登陆已过期，请重新登陆"})
		c.Abort()
		return
	}

	//将id写入gin上下文中
	c.Set("id", id)
}

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	id, err := service.Register(username, password)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
	} else {
		token := service.CreateToken(id)
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			Token:    token,
		})
	}
}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	id, err := service.Login(username, password)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
	} else {
		token := service.CreateToken(id)
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			Token:    token,
		})
	}
}

func UserInfo(c *gin.Context) {
	_id, ok := c.Get("id")
	if !ok {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "获取用户id失败，请重试"})
		return
	}

	id, ok := _id.(int64)
	if !ok {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "获取用户id失败，请重试"})
		return
	}

	user, err := service.UserInfo(id)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
	} else {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 0},
			User:     user,
		})
	}
}
