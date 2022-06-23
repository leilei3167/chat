package connect

import (
	"github.com/leilei3167/chat/config"
	"github.com/sirupsen/logrus"
	"runtime"
)

type Connect struct {
	ServerId string
}

func New() *Connect {
	return new(Connect)
}

func (c *Connect) Run() {
	ConnConf := config.Conf.Connect
	runtime.GOMAXPROCS(ConnConf.ConnectBucket.CpuNum) //桶的数量和cpu一致

	//connect层 既要作为rpc客户端,同时也要作为服务端
	//1.初始化rpc客户端(调用logic层)
	if err := c.InitLogicRpcClient(); err != nil {
		logrus.Panicf("InitLogicRpcClient err:%s", err.Error())
	}
	//2.初始化rpc服务端(供task层调用)

}
