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
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

var RedisClient *redis.Client
var RedisSessClient *redis.Client
var once sync.Once

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
	var wg sync.WaitGroup
	exit := make(chan bool, 1)

	//单进程多端口
	rpcAddrList := strings.Split(config.Conf.Logic.LogicBase.RpcAddress, ",")
	for _, bind := range rpcAddrList {
		//先验证地址,tcp@1270.0.0.1:8080
		if network, addr, err = tools.ParseNetwork(bind); err != nil {
			logrus.Panicf("InitLogicRpc ParseNetwork error : %s", err.Error())
		}
		logrus.Printf("logic start run at-->%s:%s", network, addr)
		//每个地址开启服务,此处忽略了开启Serve的错误处理,应该可以使用errorgroup来处理多端口监听的情况?
		wg.Add(1)
		go logic.createRpcServer(network, addr, exit, &wg)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit
	logrus.Println("Shutdown Server ...")
	close(exit) //关闭后监听协程会解除阻塞读取到零值
	wg.Wait()   //等待监听全部退出
	os.Exit(0)

	return
}

func (logic *Logic) createRpcServer(network string, addr string, ch chan bool, wg *sync.WaitGroup) {
	//创建rpcx服务端,etcd插件
	defer wg.Done()
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
		err := s.UnregisterAll()
		if err != nil {
			logrus.Warnf("取消注册服务错误:%v", err)
		} //在优雅退出时,将注册的服务下线
	})
	go func() {
		err := s.Serve(network, addr)
		if err != nil {
			return
		}
	}()

	<-ch
	ctx, cancle := context.WithTimeout(context.Background(), time.Second*10)
	defer cancle()
	err = s.Shutdown(ctx)
	if err != nil {
		logrus.Warn("logic rpc 关机错误:", err)
	}
	logrus.Printf("%s 已结束监听", network+addr)
	os.Exit(0)
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

func getUserKey(userID string) string {
	return config.RedisPrefix + userID
}
