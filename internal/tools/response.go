package tools

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

//定义反馈至用户的错误码及信息,统一编写供handler调用 避免重复编码
const (
	CodeSuccess      = 0
	CodeFail         = 1
	CodeUnknownError = -1
	CodeSessionError = 40000
)

var MsgCodeMap = map[int]string{
	CodeSuccess:      "success",
	CodeFail:         "fail",
	CodeUnknownError: "unknow error",
	CodeSessionError: "session error",
}

func ResponseWithCode(c *gin.Context, msgCode int, msg interface{}, data interface{}) {
	//如果没有指定消息,则从map中查询对应的错误码
	if msg == nil {
		if val, ok := MsgCodeMap[msgCode]; ok {
			msg = val
		} else {
			msg = MsgCodeMap[-1]
		}
	}

	//终止后续的处理器执行并返回,对于所有的请求,都回复200状态码,但是信息主体中传入自定义的错误信息
	c.AbortWithStatusJSON(http.StatusOK, gin.H{
		"code":    msgCode,
		"message": msg,
		"data":    data,
	})

}

func FailWithMsg(c *gin.Context, msg interface{}) {
	ResponseWithCode(c, CodeFail, msg, nil)
}
func SuccessWithMsg(c *gin.Context, msg interface{}, data interface{}) {
	ResponseWithCode(c, CodeSuccess, msg, data)
}
