package router

import (
	"strconv"

	"github.com/ao9911/bluebell-new/pkg/ecode"
	"github.com/ao9911/go-matrix/response"
	"github.com/gin-gonic/gin"
)

func GetCommunityList(c *gin.Context) {
	// 业务处理
	data, err := srv.GetCommunityList(c)
	if err != nil {
		HandleError(c, err)
		return
	}
	// 返回结果
	response.JSONSuccess(c, data)
}

func GetCommunityDetail(c *gin.Context) {
	// 获取参数并转换
	communityID, err := strconv.ParseInt(c.Param("community_id"), 10, 64)
	if err != nil || communityID <= 0 {
		response.JSONFail(c, ecode.RequestErr, nil)
		return
	}
	// 业务处理
	data, err := srv.GetCommunityDetail(c, communityID)
	if err != nil {
		HandleError(c, err)
		return
	}

	// 返回结果
	response.JSONSuccess(c, data)

}
