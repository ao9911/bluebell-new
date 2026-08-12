package router

import (
	"github.com/ao9911/go-matrix/response"
	"github.com/gin-gonic/gin"

	v1 "github.com/ao9911/bluebell-new/api/post/v1"
	"github.com/ao9911/bluebell-new/pkg/ecode"
)

func PostVote(c *gin.Context) {
	var param v1.VotePostRequest
	// 获取参数&参数校验
	if err := c.ShouldBindJSON(&param); err != nil {
		response.JSONFail(c, ecode.RequestErr, nil)
		return
	}
	uid, ok := getCurrentUserID(c)
	if !ok {
		response.JSONFail(c, ecode.NeedLogin, nil)
		return
	}
	// 业务处理
	if err := srv.PostVote(c, uid, &param); err != nil {
		HandleError(c, err)
		return
	}
	// 返回结果
	response.JSONSuccess(c, nil)
}
