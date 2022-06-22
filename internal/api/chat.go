package api

import (
	"context"
	"fmt"
	"github.com/leilei3167/chat/config"
	"github.com/leilei3167/chat/internal/api/router"
	"github.com/leilei3167/chat/internal/api/rpc"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Chat struct {
}

func New() *Chat { //统一各层服务开启方式,创建一个空结构体,初始化以及运行都在Run内部
	return &Chat{}
}

func (c *Chat) Run() {
	//初始化本层所需服务
	rpc.InitLogicRpcClient()
	//初始化rpc客户端,用于api层调用logic进行上线注册;注册路由开启服务

	r := router.Register()
	apiConfig := config.Conf.Api
	port := apiConfig.ApiBase.ListenPort
	logrus.Println("开始监听端口:", port)
	srv := &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: r} //可以有更多配置选项,如tls等

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Errorf("start listen : %s\n", err)
		}
	}()

	//优雅退出
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit
	logrus.Println("Shutdown Server ...")

	ctx, cancle := context.WithTimeout(context.Background(), time.Second*10) //10s的超时关闭时间
	defer cancle()
	if err := srv.Shutdown(ctx); err != nil {
		logrus.Errorf("Server Shutdown:%v", err)
	}
	logrus.Println("服务已关闭")
	os.Exit(0)
}
