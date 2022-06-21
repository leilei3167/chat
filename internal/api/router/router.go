package router

import (
	"github.com/gin-gonic/gin"
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
	userGroup.Use(CheckSessionID()) //都必须检查session
	userGroup.POST("/login")        //登录逻辑
	userGroup.POST("/register")     //注册逻辑
	{
		userGroup.POST("/checkAuth")
		userGroup.POST("/logout")
	}

}
func initPushRouter(r *gin.Engine) {

}

type FormCheckSessionId struct { //每一个请求 都必须附带上Token
	AuthToken string `form:"authToken" json:"authToken" binding:"required"`
}

func CheckSessionID() gin.HandlerFunc {
	return func(c *gin.Context) {
		/*	var formCheckSessionId FormCheckSessionId
			if err := c.ShouldBindBodyWith(&formCheckSessionId, binding.JSON); err != nil {
				tools.ResponseWithCode(c, tools.CodeSessionError, nil, nil)
				return
			}*/
		//获取到session后,和数据库进行比对,通过rpc调用logic的方法来验证
		//TODO:验证session
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
