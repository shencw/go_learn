package core

import (
	"github.com/gin-gonic/gin"
	"go_learn/pkg/errors"
	"go_learn/pkg/log"
	"net/http"
)

// ErrResponse 定义了发生错误时的返回消息。
// 如果引用不存在，将被省略。
type ErrResponse struct {
	// Code 业务Code码
	Code int `json:"code"`

	// Message 详细返回信息
	Message string `json:"message"`

	// Reference 返回可能对解决此错误有用的参考文档
	Reference string `json:"reference,omitempty"`
}

func WriteResponse(c *gin.Context, err error, data interface{}) {
	if err != nil {
		log.Errorf("%#+v", err)
		coder := errors.ParseCoder(err)
		c.JSON(coder.HTTPStatus(), ErrResponse{
			Code:      coder.Code(),
			Message:   coder.String(),
			Reference: coder.Reference(),
		})

		return
	}

	c.JSON(http.StatusOK, data)
}
