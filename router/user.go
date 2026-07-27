package router

import (
	"errors"

	"github.com/ao9911/go-matrix/log"
	"github.com/ao9911/go-matrix/response"
	"github.com/gin-gonic/gin"

	"github.com/ao9911/bluebell-new/api"
	"github.com/ao9911/bluebell-new/pkg/ecode"
)

func SignUp(c *gin.Context) {
	var param api.SignupRequest
	// 获取参数&参数校验
	if err := c.ShouldBindJSON(&param); err != nil {
		log.Error("SignUp invalid param failed: %v", err)
		response.JSONFail(c, ecode.RequestErr, nil)
		return
	}
	// 业务处理
	data, err := srv.SignUp(c, &param)
	if err != nil {
		if errors.Is(err, ecode.UserExist) {
			response.JSONFail(c, ecode.UserExist, nil)
			return
		}
		log.Error("SignUp failed: %v", err)
		response.JSONFail(c, ecode.ServerErr, nil)
		return
	}
	// 返回结果
	response.JSONSuccess(c, data)
}

func Login(c *gin.Context) {
	var param api.LoginRequest
	// 获取参数&参数校验
	if err := c.ShouldBindJSON(&param); err != nil {
		log.Error("Login invalid param failed: %v", err)
		response.JSONFail(c, ecode.RequestErr, nil)
		return
	}
	// 业务处理
	data, err := srv.Login(c, &param)
	if err != nil {
		if errors.Is(err, ecode.UserNotFound) {
			response.JSONFail(c, ecode.UserNotFound, nil)
			return
		}
		if errors.Is(err, ecode.InvalidPassword) {
			response.JSONFail(c, ecode.InvalidPassword, nil)
			return
		}
		log.Error("Login failed: %v", err)
		response.JSONFail(c, ecode.ServerErr, nil)
		return
	}
	// 返回结果
	response.JSONSuccess(c, data)
}
