package router

import (
	"time"
	"wiki_bit/app/api/configure"
	"wiki_bit/app/api/geo"
	"wiki_bit/app/api/oauth"
	"wiki_bit/app/api/openapi"
	"wiki_bit/app/api/participant"
	"wiki_bit/app/api/system"
	"wiki_bit/middleware"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 路由配置
func InitRouter(zapLogger *zap.Logger) *gin.Engine {

	engine := gin.New()
	//gin中间件加载
	engine.Use(middleware.Cors())
	engine.Use(middleware.Secure())
	engine.Use(middleware.Language())
	engine.Use(middleware.Ginzap(zapLogger, time.RFC3339, true))
	engine.Use(middleware.RecoveryWithZap(zapLogger, true))

	apiGroup := engine.Group("api/v1")
	{
		// 验证以太坊地址
		apiGroup.GET("/participate/verification", participant.Verification)
		// 领取糖果
		//apiGroup.POST("/participate/create", participant.Create)
		// 完善手机/邮箱
		apiGroup.POST("/participate/updateCommunication", participant.UpdateCommunication)
		// 获取认购详细信息
		apiGroup.GET("/participate/info", participant.GetInfo)
		// 获取首页数据
		apiGroup.GET("/participate/homeInfo", participant.GetHomeInfo)
		// 登录
		apiGroup.POST("/participate/login", participant.Login)
		// h5-领取类型
		apiGroup.POST("/participate/receiveType", participant.ReceiveType)
		// h5-充币
		//apiGroup.POST("/participate/recharge", participant.Recharge)
		// geo
		apiGroup.GET("/geo/getWithIpToLocation", geo.GetWithIpToLocation)

		// oauth
		//apiGroup.GET("/oauth/smsSend", oauth.SmsSend)
		//apiGroup.POST("/oauth/validateCode", oauth.ValidateCode)
		//apiGroup.POST("/oauth/sendEmailCode", oauth.SendEmailCode)
		//apiGroup.POST("/oauth/validateEmailCode", oauth.ValidateEmailCode)
		//apiGroup.GET("/oauth/receiveSmsSend", oauth.ReceiveSmsSend)
		//apiGroup.POST("/oauth/receiveValidateCode", oauth.ReceiveValidateCode)
		//apiGroup.POST("/oauth/receiveSendEmailCode", oauth.ReceiveSendEmailCode)
		//apiGroup.POST("/oauth/receiveValidateEmailCode", oauth.ReceiveValidateEmailCode)
		apiGroup.POST("/oauth/userInfo", oauth.UserInfo)

		// system-验证码
		apiGroup.GET("/getCaptcha", system.GenerateCaptchaHandler)
		apiGroup.POST("/checkCaptcha", system.CheckCode)
		// system-上传图片
		apiGroup.POST("/uploadByBase64", system.UploadByBase64)
		apiGroup.GET("/download", system.Download)

		// configure
		apiGroup.GET("/getConfigure", configure.GetConfigure)

		// query 给市场部查询数据统计使用
		apiGroup.GET("/query", participant.Query)

		// 给wikiglobal调用谷歌地图经纬度中转用
		apiGroup.POST("/googleMap/findCoordinateByAddress", openapi.FindCoordinateByAddress)

		// 给bit调用清理数据
		apiGroup.GET("/clear", participant.Clear)
	}

	return engine
}
