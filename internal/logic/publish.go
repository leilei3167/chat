package logic

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/leilei3167/chat/config"
	"github.com/leilei3167/chat/internal/db"
	"github.com/leilei3167/chat/internal/tools"
	"github.com/pkg/errors"
	"github.com/rcrowley/go-metrics"
	"github.com/rpcxio/rpcx-etcd/serverplugin"
	"github.com/sirupsen/logrus"
	"github.com/smallnest/rpcx/server"
	"strings"
	"time"
)

var RedisClient *redis.Client
var RedisSessClient *redis.Client

func (logic *Logic) InitPublishRedisClient() error {
	//根据配置文件构建redis的配置选项
	redisOpt := db.RedisOption{
		Address:  config.Conf.Common.CommonRedis.RedisAddress,
		Password: config.Conf.Common.CommonRedis.RedisPassword,
		Db:       config.Conf.Common.CommonRedis.Db,
	}
	RedisClient = db.GetRedisInstance(redisOpt)
	//测试连接
	if _, err := RedisClient.Ping(context.Background()).Result(); err != nil {
		return errors.Wrap(err, "redis无法ping通")
	}
	RedisSessClient = RedisClient //共用
	return nil
}

func (logic *Logic) InitRpcServer() (err error) {
	var network, addr string
	//单进程多端口
	rpcAddrList := strings.Split(config.Conf.Logic.LogicBase.RpcAddress, ",")
	for _, bind := range rpcAddrList {
		//先验证地址,tcp@1270.0.0.1:8080
		if network, addr, err = tools.ParseNetwork(bind); err != nil {
			logrus.Panicf("InitLogicRpc ParseNetwork error : %s", err.Error())
		}
		logrus.Printf("logic start run at-->%s:%s", network, addr)
		//每个地址开启服务,此处忽略了开启Serve的错误处理,应该可以使用errorgroup来处理多端口监听的情况?
		go logic.createRpcServer(network, addr)
	}
	return
}

func (logic *Logic) createRpcServer(network string, addr string) {
	//创建rpcx服务端,etcd插件
	s := server.NewServer()
	addRegistryPlugin(s, network, addr)

	//注册服务,元数据使用本机的uuid
	err := s.RegisterName(config.Conf.Common.CommonEtcd.ServerPathLogic, new(RpcLogic),
		fmt.Sprintf("%s", logic.ServerID))
	if err != nil {
		logrus.Fatalf("注册rpc服务错误:%s", err.Error())
	}
	//TODO:是如何下线的
	s.RegisterOnShutdown(func(s *server.Server) {
		s.UnregisterAll() //在优雅退出时,将注册的服务下线
	})
	s.Serve(network, addr)
}

func addRegistryPlugin(s *server.Server, network, addr string) {
	r := &serverplugin.EtcdV3RegisterPlugin{
		ServiceAddress: network + "@" + addr,
		EtcdServers:    []string{config.Conf.Common.CommonEtcd.Host},
		BasePath:       config.Conf.Common.CommonEtcd.BasePath,
		Metrics:        metrics.NewRegistry(),
		UpdateInterval: time.Second * 30,
	}
	err := r.Start() //测试连接
	if err != nil {
		logrus.Fatal("添加etcd插件错误:", err)
	}
	s.Plugins.Add(r)
}
