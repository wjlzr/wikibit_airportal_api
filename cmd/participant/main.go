package main

import (
	"github.com/k0kubun/pp"
	participant2 "wiki_bit/app/model/participant"
	orm "wiki_bit/boot/db/mysql"
	"wiki_bit/config"
	"wiki_bit/library/convert/xstring"
)

func init() {
	//初始化配置文件
	config.LoadConfig()

	//mysql初始化
	_ = orm.Init(
		config.Conf().MySQL.DriverName,
		config.Conf().MySQL.Dsn,
		config.Conf().MySQL.MaxOpenConns,
		config.Conf().MySQL.MaxIdleConns,
	)
}

func main() {

	//  更新原有已记录的币到未领取中
	//var participants []participant2.Participant
	//err := orm.Engine.Table("participant").Cols("id", "wait_integral", "already_integral").Find(&participants)
	//if err != nil {
	//	_, _ = pp.Println("全量Err：", err)
	//}
	//_, _ = pp.Println("总数量=============", len(participants))
	//for k, v := range participants {
	//	//time.Sleep(100 * 1000 * time.Microsecond)
	//	_, _ = pp.Println("key=========>", k)
	//	_, _ = pp.Println("id=========>", v.Id)
	//	notReceived := v.WaitIntegral + v.AlreadyIntegral
	//	_, err = orm.Engine.Table("participant").Where("id = ?", v.Id).Update(map[string]interface{}{"not_received": notReceived})
	//	if err != nil {
	//		_, _ = pp.Println("update -- id -> "+xstring.Int64ToString(v.Id)+" err：", err.Error())
	//	}
	//}
	//
	//_, _ = pp.Println("更新已结束")

	// 手动更新已有的真实的用户
	var participants []participant2.Participant
	err := orm.Engine.Table("participant").Cols("id", "wait_integral", "already_integral").Where("already_integral > ?", 0).Find(&participants)
	if err != nil {
		_, _ = pp.Println("全量Err：", err)
	}
	_, _ = pp.Println("=====总数量====", len(participants))
	for k, v := range participants {
		//time.Sleep(100 * 1000 * time.Microsecond)
		_, _ = pp.Println("key=========>", k)
		_, _ = pp.Println("id=========>", v.Id)
		waitIntegral := v.WaitIntegral + v.AlreadyIntegral
		_, err = orm.Engine.Table("participant").Where("id = ?", v.Id).Update(map[string]interface{}{"is_cheat": 2, "wait_integral": waitIntegral, "already_integral": 0})
		if err != nil {
			_, _ = pp.Println("update -- id -> "+xstring.Int64ToString(v.Id)+" err：", err.Error())
		}
	}

	_, _ = pp.Println("更新已结束")
}
