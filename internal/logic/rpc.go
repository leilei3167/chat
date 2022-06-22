package logic

import (
	"context"
	"errors"
	"github.com/leilei3167/chat/config"
	"github.com/leilei3167/chat/internal/logic/dao"
	"github.com/leilei3167/chat/internal/proto"
	"github.com/leilei3167/chat/internal/tools"
	"github.com/sirupsen/logrus"
	"strconv"
	"time"
)

type RpcLogic struct {
}

//logic层所有的rpc服务的方法由RpcLogic这个结构体来管理,注册为rpc服务的函数都有要满足统一的函数签名

// Register 提供api层的注册服务逻辑
//首先检查用户名,是否已存在(mysql中查询,如有返回主键),已存在则返回错误;未存在则将信息存入数据库,返回令牌
func (r *RpcLogic) Register(ctx context.Context, args *proto.RegisterRequest, reply *proto.RegisterReply) (err error) {
	reply.Code = config.FailReplyCode
	u := new(dao.User)
	uData := u.CheckHaveUserName(args.Name) //先验证username是否已存在
	if uData.ID > 0 {
		return errors.New("this username already exists,please login")
	}
	//新用户,存入db
	u.UserName = args.Name
	hashedPassword, err := tools.HashPassword(args.Password)
	if err != nil {
		logrus.Warnf("hashing password err:%s", err.Error())
		return err
	}
	u.Password = hashedPassword
	userID, err := u.Add()
	if err != nil {
		logrus.Warnf("register err:%s", err.Error())
		return err
	}
	if userID == 0 {
		return errors.New("register userId empty!")
	}

	//存入db成功后说明注册成功,为其生成token返回
	randToken := tools.GetRandomToken(32)         //生成32位随机的token
	sessionID := tools.CreateSessionID(randToken) //组合成sesionID
	//以sessionID为key,用户的userID和userName为value,存入redis哈希表
	userDataInRedis := make(map[string]interface{})
	userDataInRedis["userId"] = userID
	userDataInRedis["userName"] = args.Name
	//-----------使用redis事务存入-----------
	tx := RedisSessClient.TxPipeline()
	tx.HMSet(context.Background(), sessionID, userDataInRedis)
	tx.Expire(context.Background(), sessionID, config.RedisBaseValidTime*time.Second)
	_, err = tx.Exec(context.Background())
	if err != nil {
		logrus.Warnf("register set redis token fail:%v", err.Error())
		return err
	}
	//--------------------------------------
	reply.Code = config.SuccessReplyCode
	reply.AuthToken = randToken
	return
}

func (r *RpcLogic) Login(ctx context.Context, args *proto.LoginRequest, reply *proto.LoginResponse) (err error) {
	reply.Code = config.FailReplyCode
	u := new(dao.User)
	userName := args.Name
	passWord := args.Password
	//检查用户名和密码
	userData := u.CheckHaveUserName(userName)
	if userData.ID == 0 {
		return errors.New("用户名或密码错误")
	}
	err = tools.CheckPassword(passWord, userData.Password)
	if err != nil {
		return errors.New("用户名或密码错误")
	}

	//验证成功,尝试创建在线session和更新登录session
	loginSession := tools.GetSessionIdByUserId(int(userData.ID)) //sess_map_%d
	randToken := tools.GetRandomToken(32)
	sessionId := tools.CreateSessionID(randToken)
	userDataInRedis := make(map[string]interface{})
	userDataInRedis["userId"] = userData.ID
	userDataInRedis["userName"] = userData.UserName
	//先检查当前用户上线状态
	oldToken, _ := RedisSessClient.Get(context.Background(), loginSession).Result()
	if oldToken != "" {
		//说明在线,将已登录的session删除
		oldSession := tools.CreateSessionID(oldToken)
		err := RedisSessClient.Del(context.Background(), oldSession).Err() //删除旧的session(下线)
		if err != nil {
			return errors.New("logout user fail!token is:" + oldToken)
		}
	}

	//----------redis事务------------
	tx := RedisSessClient.TxPipeline()
	tx.HMSet(context.Background(), sessionId, userDataInRedis)
	tx.Expire(context.Background(), sessionId, config.RedisBaseValidTime*time.Second)
	tx.Set(context.Background(), loginSession, randToken, config.RedisBaseValidTime*time.Second) //上线session会在logout时或者过期时被销毁
	_, err = tx.Exec(context.Background())
	if err != nil {
		logrus.Warnf("register set redis token fail:%v", err.Error())
		return err
	}
	//------------redis事务结束------------
	reply.Code = config.SuccessReplyCode
	reply.AuthToken = randToken
	return
}
func (r *RpcLogic) CheckAuth(ctx context.Context, args *proto.CheckAuthRequest, reply *proto.CheckAuthResponse) (err error) {
	//api中间件,每一个请求都会获取前端携带的Token调用此方法进行验证
	reply.Code = config.FailReplyCode
	authToken := args.AuthToken
	sessionId := tools.CreateSessionID(authToken) //拼接成redis的key
	var userDataMap = map[string]string{}
	userDataMap, err = RedisSessClient.HGetAll(context.Background(), sessionId).Result()
	if err != nil {
		logrus.Warnf("查询sessionID:'%s'出错:%v", sessionId, err)
		return err
	}
	if len(userDataMap) == 0 {
		logrus.Warnf("seesionID:'%s'不存在用户信息:%v", sessionId, err)
		return
	}
	intUserId, _ := strconv.Atoi(userDataMap["userId"])
	reply.UserId = intUserId
	userName, _ := userDataMap["userName"]
	reply.Code = config.SuccessReplyCode
	reply.UserName = userName
	return
}
