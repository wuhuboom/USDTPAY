package common

var (
	VerifyErrCode             = 1002  //校验错误Code
	ErrReturnCode             = -101  //全局返回错误的信息
	SuccessReturnCode         = 2000  //全局返回正常的数据
	IpLimitWaring             = -102  //全局限制警告  ip限制
	IllegalityCode            = -103  //全局非法请求警告,一般是没有token
	NeedGoogleBind            = -104  //管理员登录需要绑定谷歌验证码
	TokenExpire               = -105  //管理员或者用户的token过期
	NoHavePermission          = -106  //没有权限
	MysqlErr                  = -107  //数据的一些报错
	TaskClearing              = 300   //任务结算
	ReturnOldOrderCode        = 20001 //返回已经获取的账单
	NoBank                    = 400
	SystemMinWithdrawal       = 401
	NoEnoughMoney             = -888
	StatusNotFound            = 200404
	StatusInternalServerError = 200500
)
