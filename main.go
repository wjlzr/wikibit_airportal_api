package main

import (
	"fmt"
	"wiki_bit/boot/db/mysql"
	"wiki_bit/boot/log"
	"wiki_bit/config"
	"wiki_bit/router"

	"go.uber.org/zap"
)

func main() {
	// 初始化配置文件
	config.LoadConfig()
	//初始化log
	log.Init("logs")

	//mysql初始化
	err := mysql.Init(
		config.Conf().MySQL.DriverName,
		config.Conf().MySQL.Dsn,
		config.Conf().MySQL.MaxOpenConns,
		config.Conf().MySQL.MaxIdleConns,
	)
	if err != nil {
		log.Logger().Error(" mysql connect error", zap.Error(err))
		return
	}

	// redis
	//redis.Connect()

	// nsq consumer
	//nsq.NewNsq().InitConsumer("create_bit_topic", "create_bit_channel", &participate.CreateBitConsumer{})

	// log4
	//log4.LoadConfiguration("./example.json")

	//gin路由引擎配置
	engine := router.InitRouter(log.Logger())
	//启动服务
	if err = engine.Run(fmt.Sprintf("%s:%d", config.Conf().Application.Host, config.Conf().Application.Port)); err != nil {
		log.Logger().Error("start service error", zap.Error(err))
	}
}
