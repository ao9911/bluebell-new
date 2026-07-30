package router

import (
	"github.com/ao9911/go-matrix/response"
	"github.com/gin-gonic/gin"

	v1 "github.com/ao9911/bluebell-new/api/user/v1"
	"github.com/ao9911/bluebell-new/pkg/ecode"
)

func SignUp(c *gin.Context) {
	var param v1.SignupRequest
	// 获取参数&参数校验
	if err := c.ShouldBindJSON(&param); err != nil {
		response.JSONFail(c, ecode.RequestErr, nil)
		return
	}
	// 业务处理
	data, err := srv.SignUp(c, &param)
	if err != nil {
		HandleError(c, err)
		return
	}
	// 返回结果
	response.JSONSuccess(c, data)
}

func Login(c *gin.Context) {
	var param v1.LoginRequest
	// 获取参数&参数校验
	if err := c.ShouldBindJSON(&param); err != nil {
		response.JSONFail(c, ecode.RequestErr, nil)
		return
	}
	// 业务处理
	data, err := srv.Login(c, &param)
	if err != nil {
		HandleError(c, err)
		return
	}
	// 返回结果
	response.JSONSuccess(c, data)
}

func RefreshToken(c *gin.Context) {
	var param v1.RefreshTokenRequest
	// 获取参数&参数校验
	if err := c.ShouldBindJSON(&param); err != nil {
		response.JSONFail(c, ecode.RequestErr, nil)
		return
	}
	// 业务处理
	data, err := srv.RefreshToken(c, &param)
	if err != nil {
		HandleError(c, err)
		return
	}
	// 返回结果
	response.JSONSuccess(c, data)
}
