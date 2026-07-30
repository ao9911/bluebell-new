package ecode

import xecode "github.com/ao9911/go-matrix/ecode"

// Error codes.
var (
	OK         = xecode.OK
	RequestErr = xecode.RequestErr
	ServerErr  = xecode.ServerErr

	UserExist           = xecode.New(10001) // 用户已存在
	UserNotFound        = xecode.New(10002) // 用户不存在
	InvalidPassword     = xecode.New(10003) // 密码错误
	InvalidToken        = xecode.New(10004) // token无效
	NeedLogin           = xecode.New(10005) // 需要登录
	InvalidRefreshToken = xecode.New(10006) // 刷新token无效
)

// Corresponding messages.
func init() {
	xecode.Register(map[int32]string{
		0:     "成功",
		-400:  "请求错误",
		-500:  "服务器错误",
		10001: "用户名已存在",
		10002: "用户不存在",
		10003: "密码错误",
		10004: "token无效",
		10005: "需要登录",
		10006: "刷新token无效",
	})
}
