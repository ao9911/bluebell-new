package router

import (
	"strconv"

	"github.com/ao9911/go-matrix/response"
	"github.com/gin-gonic/gin"

	v1 "github.com/ao9911/bluebell-new/api/post/v1"
	"github.com/ao9911/bluebell-new/pkg/ecode"
)

func CreatePost(c *gin.Context) {
	var param v1.CreatePostRequest
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
	data, err := srv.CreatePost(c, uid, &param)
	if err != nil {
		HandleError(c, err)
		return
	}
	// 返回结果
	response.JSONSuccess(c, data)
}

func GetPostDetail(c *gin.Context) {
	// 获取参数并转换
	postID, err := strconv.ParseInt(c.Param("post_id"), 10, 64)
	if err != nil || postID <= 0 {
		response.JSONFail(c, ecode.RequestErr, nil)
		return
	}
	// 业务处理
	data, err := srv.GetPostDetail(c, postID)
	if err != nil {
		HandleError(c, err)
		return
	}
	// 返回结果
	response.JSONSuccess(c, data)
}

func GetPostList(c *gin.Context) {
	// 获取分页参数
	page, size := getPageInfo(c)
	// 业务处理
	data, err := srv.GetPostList(c, page, size)
	if err != nil {
		HandleError(c, err)
		return
	}
	// 返回结果
	response.JSONSuccess(c, data)
}
