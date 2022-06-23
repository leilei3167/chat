package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/leilei3167/chat/internal/api/rpc"
	"github.com/leilei3167/chat/internal/proto"
	"github.com/leilei3167/chat/internal/tools"
)

// FormLogin 定义登录所需的结构体
type FormLogin struct {
	UserName string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var formLogin FormLogin
	//接收参数
	if err := c.ShouldBindBodyWith(&formLogin, binding.JSON); err != nil {
		//出错
		tools.FailWithMsg(c, err.Error())
		return
	}

	req := &proto.LoginRequest{
		Name:     formLogin.UserName,
		Password: formLogin.Password,
	}
	//使用全局的RpcLogicObj调用登录方法(将Login作为函数而不是方法,应该也是一样的)
	code, authToken, msg := rpc.RpcLogicObj.Login(req)
	if code == tools.CodeFail || authToken == "" {
		tools.FailWithMsg(c, msg)
		return
	}
	tools.SuccessWithMsg(c, "登陆成功", authToken) //将Token反馈给前端
}

// FormRegister 提供注册
type FormRegister struct {
	UserName string `form:"userName" json:"userName" binding:"required"`
	Password string `form:"passWord" json:"passWord" binding:"required"`
}

// Register 提供注册的处理器
func Register(c *gin.Context) {
	var formRegister FormRegister
	if err := c.ShouldBindBodyWith(&formRegister, binding.JSON); err != nil {
		tools.FailWithMsg(c, err.Error())
		return
	}

	req := &proto.RegisterRequest{
		Name:     formRegister.UserName,
		Password: formRegister.Password,
	}
	code, authToken, msg := rpc.RpcLogicObj.Register(req)
	if code == tools.CodeFail || authToken == "" {
		tools.FailWithMsg(c, msg)
		return
	}
	tools.SuccessWithMsg(c, "register success", authToken)
}

// FormCheckAuth 用于获取token用于验证用户状态
type FormCheckAuth struct {
	AuthToken string `form:"authToken" json:"authToken" bingding:"required"`
}

func CheckAuth(c *gin.Context) {
	var formAuth FormCheckAuth
	if err := c.ShouldBindBodyWith(&formAuth, binding.JSON); err != nil {
		tools.FailWithMsg(c, err.Error())
		return
	}
	//构建rpc请求
	token := formAuth.AuthToken
	req := &proto.CheckAuthRequest{AuthToken: token}
	code, id, name := rpc.RpcLogicObj.CheckAuth(req)
	if code == tools.CodeFail {
		tools.FailWithMsg(c, "auth fail")
		return
	}
	var result = map[string]interface{}{
		"userId":   id,
		"userName": name,
	}
	tools.SuccessWithMsg(c, "auth success", result)
}

// FormLogout 退出登录,将该token删除即可
type FormLogout struct {
	AuthToken string `form:"authToken" json:"authToken" binding:"required"`
}

func Logout(c *gin.Context) {
	var formLogout FormLogout
	if err := c.ShouldBindBodyWith(&formLogout, binding.JSON); err != nil {
		tools.FailWithMsg(c, err.Error())
		return
	}
	authToken := formLogout.AuthToken
	logoutReq := &proto.LogoutRequest{
		AuthToken: authToken,
	}
	code := rpc.RpcLogicObj.Logout(logoutReq)
	if code == tools.CodeFail {
		tools.FailWithMsg(c, "logout fail!")
		return
	}
	tools.SuccessWithMsg(c, "logout ok!", nil)
}
