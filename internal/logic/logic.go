// Package logic 提供核心的处理逻辑
package logic

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/leilei3167/chat/config"
	"github.com/sirupsen/logrus"
	"runtime"
)

type Logic struct {
	ServerID string
}

func New() *Logic {
	return new(Logic)
}

func (logic *Logic) Run() {
	logicConfig := config.Conf.Logic

	runtime.GOMAXPROCS(logicConfig.LogicBase.CpuNum) //TODO:如果直接设置runtime.NumCPU?(不同cpu的机器部署logic是否会出错?)
	logic.ServerID = fmt.Sprintf("logic-%s", uuid.New().String())

	//logic层需要频繁与数据库交互
	err := logic.InitPublishRedisClient()
	if err != nil {
		logrus.Panicf("logic init redis fail:%v", err)
	}
	//初始化rpc服务(单进程多端口)
	if err := logic.InitRpcServer(); err != nil {
		logrus.Panicf("logic init rpc server fail:%v", err)
	}

}
