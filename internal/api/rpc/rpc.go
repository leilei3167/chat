// Package rpc 此包定义api层与其他层交互的客户端及方法
package rpc

import (
	"context"
	"github.com/leilei3167/chat/config"
	"github.com/leilei3167/chat/internal/proto"
	"github.com/rpcxio/libkv/store"
	etcdv3 "github.com/rpcxio/rpcx-etcd/client"
	"github.com/sirupsen/logrus"
	"github.com/smallnest/rpcx/client"
	"sync"
	"time"
)

//包内私有的XClient,仅用于包内方法的调用
var logicRpcClient client.XClient
var once sync.Once

type RpcLogic struct { //围绕一个实例构建调用函数
}

var RpcLogicObj *RpcLogic //外界通过rpc.RpcLogicObj来调用本包方法,在本包提前初始化好,避免调用处再创建(类似于标准库的DefaultXXX写法)

// InitLogicRpcClient 使用单例模式 创建全局的XClient
func InitLogicRpcClient() {
	once.Do(func() {
		etcdConfigOption := &store.Config{
			ClientTLS:         nil,
			TLS:               nil,
			ConnectionTimeout: time.Duration(config.Conf.Common.CommonEtcd.ConnectionTimeout) * time.Second,
			Bucket:            "",
			PersistConnection: true,
			Username:          config.Conf.Common.CommonEtcd.UserName,
			Password:          config.Conf.Common.CommonEtcd.Password,
		}
		d, err := etcdv3.NewEtcdV3Discovery(
			config.Conf.Common.CommonEtcd.BasePath,
			config.Conf.Common.CommonEtcd.ServerPathLogic,
			[]string{config.Conf.Common.CommonEtcd.Host},
			true,
			etcdConfigOption,
		)
		if err != nil {
			logrus.Fatalf("init connect rpc etcd discovery client fail:%s", err.Error())
		}
		//初始化全局rpc客户端(针对logic层的rpc服务)
		logicRpcClient = client.NewXClient(config.Conf.Common.CommonEtcd.ServerPathLogic,
			client.Failtry, client.RoundRobin, d, client.DefaultOption)
		RpcLogicObj = new(RpcLogic)
	})
	if logicRpcClient == nil {
		logrus.Fatalf("get logic rpc client nil")
	}
}

//以下就是api层和logic层的所有rpc交互

func (rpc *RpcLogic) Login(req *proto.LoginRequest) (code int, authToken string, msg string) {
	reply := new(proto.LoginResponse)
	//使用全局的XClient进行调用
	err := logicRpcClient.Call(context.Background(), "Login", req, reply)
	if err != nil { //将错误信息作为msg传回
		msg = err.Error()
	}
	code = reply.Code
	authToken = reply.AuthToken
	return
}

func (rpc *RpcLogic) Register(req *proto.RegisterRequest) (code int, authToken string, msg string) {
	reply := &proto.RegisterReply{}
	err := logicRpcClient.Call(context.Background(), "Register", req, reply)
	if err != nil {
		msg = err.Error()
	}
	code = reply.Code
	authToken = reply.AuthToken
	return
}

// CheckAuth 用于根据前端携带的authToken来验证是否有效(未过期,与数据库一致)
func (rpc *RpcLogic) CheckAuth(req *proto.CheckAuthRequest) (code int, userId int, userName string) {
	reply := &proto.CheckAuthResponse{}
	logicRpcClient.Call(context.Background(), "CheckAuth", req, reply)
	code = reply.Code
	userId = reply.UserId
	userName = reply.UserName
	return
}

func (rpc *RpcLogic) Logout(req *proto.LogoutRequest) (code int) {
	reply := &proto.LogoutResponse{}
	logicRpcClient.Call(context.Background(), "Logout", req, reply)
	code = reply.Code
	return
}
