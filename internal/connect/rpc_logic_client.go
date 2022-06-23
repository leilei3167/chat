package connect

import (
	"errors"
	"github.com/leilei3167/chat/config"
	"github.com/rpcxio/libkv/store"
	etcdV3 "github.com/rpcxio/rpcx-etcd/client"
	"github.com/sirupsen/logrus"
	"github.com/smallnest/rpcx/client"
	"sync"
	"time"
)

var logicRpcClient client.XClient
var once sync.Once

type RpcConnect struct {
}

// InitLogicRpcClient 初始化客户端
func (c *Connect) InitLogicRpcClient() (err error) {
	etcdConfigOption := &store.Config{
		ClientTLS:         nil,
		TLS:               nil,
		ConnectionTimeout: time.Duration(config.Conf.Common.CommonEtcd.ConnectionTimeout) * time.Second,
		Bucket:            "",
		PersistConnection: true,
		Username:          config.Conf.Common.CommonEtcd.UserName,
		Password:          config.Conf.Common.CommonEtcd.Password,
	}
	once.Do(func() { //单例模式,确保只被创建一次
		d, e := etcdV3.NewEtcdV3Discovery(
			config.Conf.Common.CommonEtcd.BasePath,
			config.Conf.Common.CommonEtcd.ServerPathLogic,
			[]string{config.Conf.Common.CommonEtcd.Host},
			true,
			etcdConfigOption,
		)
		if e != nil {
			logrus.Fatalf("init connect rpc etcd discovery client fail:%s", e.Error())
		}
		logicRpcClient = client.NewXClient(config.Conf.Common.CommonEtcd.ServerPathLogic, client.Failtry, client.RandomSelect, d, client.DefaultOption)
	})
	if logicRpcClient == nil {
		return errors.New("get rpc client nil")
	}
	return

}
