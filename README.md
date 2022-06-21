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
- 由router包负责路由,中间件处理器的注册;Cors跨域管理的中间件全局使用,检查Session的中间件在与User相关的路由组中必须添加
- 统一制定错误码,以及反馈至前端的消息结构;反馈至前端的错误码统一为200,只是返回的数据中还会包含自定义的错误码
- 优雅退出逻辑,在协程中开启ListenAndServe,之后设置监听signal的缓存为1的chan,当收到退出信号时,调用http.ShutDown来实现优雅退出

## Logic层