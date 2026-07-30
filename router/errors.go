package router

import (
	"errors"

	xecode "github.com/ao9911/go-matrix/ecode"
	"github.com/ao9911/go-matrix/response"
	"github.com/gin-gonic/gin"

	"github.com/ao9911/bluebell-new/pkg/ecode"
)

func HandleError(c *gin.Context, err error) {
	var ec xecode.Codes
	if errors.As(err, &ec) {
		response.JSONFail(c, ec, nil)
		return
	}
	response.JSONFail(c, ecode.ServerErr, nil)
}
