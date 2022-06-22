# just for learning

 [原项目地址](https://github.com/LockGit/gochat)

## 核心知识点

- web应用的构建方式,包括gin框架,session结合redis进行管理,通过中间件进行验证
- rpc应用的构建模式,此项目采用rpcx框架提供rpc服务,使用etcd做服务发现,可以了解到微服务架构的构建模式
- TCP,websocket协议的使用,自定义TCP通讯消息结构
- 消息队列在项目中的应用(此项目为简便使用的是redis,可更换kafka RMQ等)
- 各组件的分层构建方式
- 后续可拓展学习docker打包部署项目

## Common
- [x] 使用Viper管理全局配置文件 
- [x] main.go做所有组件的唯一入口,统一各层启动的代码样式
### 关键点:
- 使用sync.Once确保全局配置文件结构体只被初始化一次(单例模式)
- 全局Conf在包级别声明后在init函数内需要使用new来分配内存,否则使用viper.Unmarshal时会出现空指针异常
- 使用runtime.Caller得到config.go文件的绝对路径,配合使用filepath.Spilt得到config.go所在目录的绝对路径,以便于获取到配置文件目录
- 调用Viper的ReadInConfig将配置文件读取到内存中,后续新增加的配置模块可以使用MergeInConfig来进行追加,为了避免每一个文件都使用Viper.Get来调用配置项  
选择将viper读取到的配置选项写入到Conf结构体中(此处注意结构体必须指定mapstructure的标签)
- 使用Toml作为配置文件,单个文件->模块->字段

### Tips
- 可使用cobra进一步精简main的入口,将不同层的运行使用rootCmd的子命令来替代


## Api层
  
### 关键点

#### User相关部分
- 由router包负责路由,中间件处理器的注册;Cors跨域管理的中间件全局使用,检查Session的中间件在与User相关的路由组中必须添加
- 统一制定错误码,以及反馈至前端的消息结构;反馈至前端的错误码统一为200,只是返回的数据中还会包含自定义的错误码
- 优雅退出逻辑,在协程中开启ListenAndServe,之后设置监听signal的缓存为1的chan,当收到退出信号时,调用http.ShutDown来实现优雅退出
- 使用Gin的ShouldBindBodyWith函数来使得body可复用(但性能不高),AbortWithStatusJSON来中断handler链执行
- 在api层构建rpc包,使用单例模式构建全局的XClient,并且包内创建全局的操作实例,所有调用方法都围绕实例构建,避免任何地方直接操作XClient
- 对于api层的处理器,核心逻辑就是接收并验证前端用户传入的数据,将数据构建为rpc请求,并通过rpc包方法调用logic层进行处理,处理结果反馈至前端,  
一切的复杂逻辑处理都在logic层
- 登录或注册成功后都会得到一个token(同时logic也会将其存入redis),用于鉴别登录状态,所有请求都必须带上此Token,由checkauth中间件在每次请求  
时调用logic进行验证,验证失败(token过期)的返回登录界面
#### Push相关部分

## Logic层
1. logic作为rpc的服务端,为api和connect层提供rpc服务,logic层会频繁涉及到与数据库的交互,主要是redis的session存储,和User信息mysql的交互  
会在执行的开始将2种数据库初始化完毕
2. 与各层的rpc交互所需的model都会在proto种定义好,以及调用的状态码
3. 与数据库的交互,在db层进行初始化,在dao层通过db层进行获取
4. 关于session:将存在两种session,一种是记住登录状态的(由sess_前缀和randToken拼接而成),在注册完成登录完成后创建,另一种是上线状态的session  
   (由sess_map和userId拼接而成)
5. 上线session的组成是sess_map和userId为key,登录的随机token为value
#### 注册逻辑的实现
1. 调用dao层方法,根据userName查询是否已存在user
2. 若不存在,将user信息存入mysql,返回一个userID(主键)
3. 生成32位随机token(用于反馈至用户),组合前缀生成sessionID(用作redis的key)
4. 以session为key,用户的id和name作为value存入redis中,并设置session过期时间(需使用事务)
#### 登录逻辑的实现
1. 调用dao层方法,根据userName查询是否已存在user,以及调用bcrypt进行密码验证,不通过则返回 用户名或密码错误
2. 根据userId生成在线session,使用此session查询redis中是否已是登录状态,若是,将原session对应的条目删除
3. 生成randToken和sessionId,存入redis并设置过期时间;注意在线session的value是生成的randtoken
4. 返回randtoken
#### checkAuth方法
1. api层将获取用户请求附带的token,由此token得到sessionId,在redis中查询,如果没有相关数据,则说明token无效或过期
2. 若有效,返回sessionId的value(userId,Name)
3. 将token在postform中,还是服务端返回时直接写入cookie中?
### 关键点
- 为了使分布式系统中的所有元素都有唯一的辨识标志,每个logic层的服务都包含一个ID字段,用uuid生成唯一的识别码,避免与其他节点冲突
- 在注册rpc服务时RegisterName使用logic结构体的serverID作为metadata,用于表示每个请求发起的主机,并且使用RegisterOnShutdown注册了在shutdown时要执行的操作(如何实现优雅退出)
- 需要注册的方法都在rpc.go文件中 全部绑定在RpcLogic结构体上,所需参数,返回参数全部统一放入到proto目录中,使得调用方和被调方同时使用
- InitRpcServer时,开启协程监听多端口时,为什么不会直接退出
- 由于RPC服务调用只会返回error一个值,因此更多的响应信息应该放入到reply中,并应该定义统一的错误码常量(在config中,只有失败和成功)
- 对于数据库的操作,在db包中维护dbMap并在init中初始化数据库连接,对外提供获取数据库连接的api,需要使用到数据库的dao层依靠此api进行获取连接(通过统一方式获取db)  
dao层才是进行crud操作的地方;对于gorm,可以在初始化时设置LogMode来详细打印日志便于调试
- 数据库执行操作后返回其主键是基本意识
- GetRandomToken生成随机的字符串,使用到rand.Reader和io.ReadFull,base64.URLEncoding.EncodeToString(URL方式会避免生成/和+)
- 存入redis时,因为包含存入以及设置过期时间2步,应该使用事务来存储(事务是使用TxPipeline来实现)