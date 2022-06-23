package router

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/leilei3167/chat/internal/api/handler"
	"github.com/leilei3167/chat/internal/api/rpc"
	"github.com/leilei3167/chat/internal/proto"
	"github.com/leilei3167/chat/internal/tools"
)

// Register 初始化路由,注册处理器,中间件等
func Register() *gin.Engine {
	r := gin.Default()
	r.Use(CorsMid())

	initUserRouter(r)

	return r
}

//初始化用户交互相关的处理器,如登录等
func initUserRouter(r *gin.Engine) {
	userGroup := r.Group("/user")
	//userGroup.Use(CheckSessionID())               //都必须检查session
	userGroup.POST("/login", handler.Login)       //登录逻辑
	userGroup.POST("/register", handler.Register) //注册逻辑
	{
		userGroup.POST("/checkAuth", handler.CheckAuth) //和中间件调用的rpc服务一致
		userGroup.POST("/logout", handler.Logout)
	}

}

//push主要处理消息的推送和接收
func initPushRouter(r *gin.Engine) {
	//TODO
}

type FormCheckSessionId struct { //每一个请求 都必须附带上Token这个form表单数据
	AuthToken string `form:"authToken" json:"authToken" binding:"required"`
}

func CheckSessionID() gin.HandlerFunc {
	return func(c *gin.Context) {
		var formCheckSessionId FormCheckSessionId
		if err := c.ShouldBindBodyWith(&formCheckSessionId, binding.JSON); err != nil {
			tools.ResponseWithCode(c, tools.CodeSessionError, nil, nil)
			return
		}
		//获取到session后,和数据库进行比对,通过rpc调用logic的方法来验证
		//对于每一个经过的请求,都要检查其session!(rpc调用logic)
		authToken := formCheckSessionId.AuthToken
		req := &proto.CheckAuthRequest{
			AuthToken: authToken,
		}
		code, userId, userName := rpc.RpcLogicObj.CheckAuth(req)
		if code == tools.CodeFail || userId <= 0 || userName == "" {
			c.Abort()
			tools.ResponseWithCode(c, tools.CodeSessionError, nil, nil)
			return
		}
		c.Next()
		return
	}
}

func CorsMid() gin.HandlerFunc {
	return func(c *gin.Context) {
		//将请求添加以下的Header,便于支持跨域
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
		c.Header("Access-Control-Allow-Methods", "GET, OPTIONS, POST, PUT, DELETE")
		c.Set("content-type", "application/json")
		c.Next()
	}
}
