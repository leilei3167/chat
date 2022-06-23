package proto

// LoginRequest 用于向logic发起登录调用
type LoginRequest struct {
	Name     string
	Password string
}

// LoginResponse logic对登录结果的响应
type LoginResponse struct {
	Code      int
	AuthToken string
}

// RegisterRequest 注册和登录所需的字段是一样的,为了区分
type RegisterRequest struct {
	Name     string
	Password string
}

type RegisterReply struct {
	Code      int
	AuthToken string //注册成功后会的到一个随机的令牌,同时redis中也会储存,在执行操作时进行对比
}

// CheckAuthRequest 用于验证中间件,验证正确将包含UserID和UserName作为结果
type CheckAuthRequest struct {
	AuthToken string
}

type CheckAuthResponse struct {
	Code     int
	UserId   int
	UserName string
}

type LogoutRequest struct {
	AuthToken string
}

type LogoutResponse struct { //只需返回退出成功与否
	Code int
}
