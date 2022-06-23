// Package config 提供全局的配置文件管理,各组件配置文件统一集中到结构体中
package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

var once sync.Once //单例模式,确保配置文件只初始化一次
var Conf *Config

const ( //定义常量配置,rpc状态码等
	SuccessReplyCode   = 0
	FailReplyCode      = 1
	SuccessReplyMsg    = "success"
	RedisBaseValidTime = 86400
	RedisPrefix        = "gochat_"
)

// Config 汇总所有层的配置,每一个字段代表一个配置文件
type Config struct {
	Common  Common //代表共用的配置
	Connect ConnectConfig
	Logic   LogicConfig
	//Task    TaskConfig
	Api ApiConfig //api层的配置
	//Site    SiteConfig
}

func init() { //此包被导入时初始化
	Init()

}
func Init() {
	once.Do(func() { //单例模式 将全局的Conf初始化
		//设定配置文件所在位置
		env := GetEnv()
		realPath := getCurrentDir()
		configFilePath := realPath + "/" + env + "/" //得到配置文件的绝对路径

		//读取配置文件
		viper.AddConfigPath(configFilePath)

		viper.SetConfigType("toml")
		viper.SetConfigName("api") //配置文件的名字,不要加拓展名
		err := viper.ReadInConfig()
		if err != nil {
			panic(err)
		}
		viper.SetConfigName("common")
		err = viper.MergeInConfig() //合并新的配置文件
		if err != nil {
			panic(err)
		}
		viper.SetConfigName("logic")
		err = viper.MergeInConfig()
		if err != nil {
			panic(err)
		}

		Conf = new(Config) //初始化 分配内存
		if err := viper.Unmarshal(&Conf.Api); err != nil {
			panic(err)
		}
		if err := viper.Unmarshal(&Conf.Common); err != nil {
			panic(err)
		}
		if err := viper.Unmarshal(&Conf.Logic); err != nil {
			panic(err)
		}
		logrus.Debugf("%#v", Conf)
	})
}
func GetEnv() string {
	env := os.Getenv("RUN_MODE")
	if env == "" {
		env = "dev"
	}
	return env
}
func getCurrentDir() string {
	_, fileName, _, _ := runtime.Caller(1)
	path, _ := filepath.Split(fileName)
	return path[:len(path)-1] //舍弃末尾的/
}

type ApiBase struct {
	ListenPort int `mapstructure:"listenPort"` //监听地址
}

type ApiConfig struct {
	ApiBase ApiBase `mapstructure:"api-base"`
}
type Common struct {
	CommonEtcd  CommonEtcd  `mapstructure:"common-etcd"` //代表着配置文件内的模块
	CommonRedis CommonRedis `mapstructure:"common-redis"`
}

type CommonRedis struct {
	RedisAddress  string `mapstructure:"redisAddress"`
	RedisPassword string `mapstructure:"redisPassword"`
	Db            int    `mapstructure:"db"`
}

type CommonEtcd struct { //配置文件的模块字段
	Host              string `mapstructure:"host"`
	BasePath          string `mapstructure:"basePath"`
	ServerPathLogic   string `mapstructure:"serverPathLogic"` //logic层的rpc服务的名称
	ServerPathConnect string `mapstructure:"serverPathConnect"`
	UserName          string `mapstructure:"userName"`
	Password          string `mapstructure:"password"`
	ConnectionTimeout int    `mapstructure:"connectionTimeout"`
}

type LogicBase struct {
	ServerId   string `mapstructure:"serverId"`
	CpuNum     int    `mapstructure:"cpuNum"`
	RpcAddress string `mapstructure:"rpcAddress"`
	CertPath   string `mapstructure:"certPath"`
	KeyPath    string `mapstructure:"keyPath"`
}

type LogicConfig struct {
	LogicBase LogicBase `mapstructure:"logic-base"`
}

type ConnectConfig struct {
	ConnectBase                ConnectBase                `mapstructure:"connect-base"`
	ConnectRpcAddressWebSockts ConnectRpcAddressWebsockts `mapstructure:"connect-rpcAddress-websockts"`
	ConnectRpcAddressTcp       ConnectRpcAddressTcp       `mapstructure:"connect-rpcAddress-tcp"`
	ConnectBucket              ConnectBucket              `mapstructure:"connect-bucket"`
	ConnectWebsocket           ConnectWebsocket           `mapstructure:"connect-websocket"`
	ConnectTcp                 ConnectTcp                 `mapstructure:"connect-tcp"`
}

type ConnectBase struct {
	CertPath string `mapstructure:"certPath"`
	KeyPath  string `mapstructure:"keyPath"`
}

type ConnectRpcAddressWebsockts struct {
	Address string `mapstructure:"address"`
}

type ConnectRpcAddressTcp struct {
	Address string `mapstructure:"address"`
}

type ConnectBucket struct {
	CpuNum        int    `mapstructure:"cpuNum"`
	Channel       int    `mapstructure:"channel"`
	Room          int    `mapstructure:"room"`
	SrvProto      int    `mapstructure:"svrProto"`
	RoutineAmount uint64 `mapstructure:"routineAmount"`
	RoutineSize   int    `mapstructure:"routineSize"`
}

type ConnectWebsocket struct {
	ServerId string `mapstructure:"serverId"`
	Bind     string `mapstructure:"bind"`
}

type ConnectTcp struct {
	ServerId      string `mapstructure:"serverId"`
	Bind          string `mapstructure:"bind"`
	SendBuf       int    `mapstructure:"sendbuf"`
	ReceiveBuf    int    `mapstructure:"receivebuf"`
	KeepAlive     bool   `mapstructure:"keepalive"`
	Reader        int    `mapstructure:"reader"`
	ReadBuf       int    `mapstructure:"readBuf"`
	ReadBufSize   int    `mapstructure:"readBufSize"`
	Writer        int    `mapstructure:"writer"`
	WriterBuf     int    `mapstructure:"writerBuf"`
	WriterBufSize int    `mapstructure:"writeBufSize"`
}
