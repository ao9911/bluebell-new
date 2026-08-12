package router

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func getCurrentUserID(c *gin.Context) (int64, bool) {
	uid, ok := c.Get(CtxSubjectKey)
	if !ok {
		return 0, false
	}
	userID, ok := uid.(int64)
	if !ok {
		return 0, false
	}
	return userID, true
}

const (
	defaultPage = int64(1)
	defaultSize = int64(10)
	maxSize     = int64(100)
)

func getPageInfo(c *gin.Context) (int64, int64) {
	page, err := strconv.ParseInt(c.Query("page"), 10, 64)
	if err != nil || page < 1 {
		page = defaultPage
	}

	size, err := strconv.ParseInt(c.Query("size"), 10, 64)
	if err != nil || size < 1 {
		size = defaultSize
	}
	if size > maxSize {
		size = maxSize
	}

	return page, size
}
