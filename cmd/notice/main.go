package main

import (
	fmt2 "fmt"
	"github.com/goinggo/mapstructure"
	"github.com/k0kubun/pp"
	"os"
	"sync"
	"time"
	"wiki_bit/app/model/invite"
	"wiki_bit/app/model/oauth"
	"wiki_bit/app/model/participant"
	"wiki_bit/app/model/participantnoticerecord"
	oauth2 "wiki_bit/app/service/oauth"
	orm "wiki_bit/boot/db/mysql"
	"wiki_bit/boot/log"
	"wiki_bit/config"
	"wiki_bit/library/convert/xint"
	"wiki_bit/library/convert/xstring"
	"wiki_bit/library/fmt"
	"wiki_bit/library/geo"
	"wiki_bit/services/eth"
	"wiki_bit/services/ucloud"
)

var (
	//single     = 100000
	//totalCount = 187079800

	single     = 1000
	totalCount = 100000
)

func init() {
	//初始化配置文件
	config.LoadConfig()
	//初始化log
	log.Init("logs")
	//mysql初始化
	_ = orm.Init(
		config.Conf().MySQL.DriverName,
		config.Conf().MySQL.Dsn,
		config.Conf().MySQL.MaxOpenConns,
		config.Conf().MySQL.MaxIdleConns,
	)
}

func main() {
	//wg := sync.WaitGroup{}
	//count := math.Ceil(float64(totalCount) / float64(single))
	//for i := 0; i < int(count); i++ {
	//	wg.Add(1)
	//	go query3(int(single), int(i*single))
	//	time.Sleep(2 * time.Second)
	//}
	//wg.Wait()
	query3(10, 0)
	_, _ = pp.Println("###########数据筛查结束#############")
	_, _ = pp.Println("#######结束时间：" + time.Now().Format("2006-01-02 15:04:05") + "#########")
}

func sms() {

	// 查询需要通知的账号
	sql := "select * from participant where area_code in(select a.area_code from (select area_code,count(1) as \"sum\" from participant where is_cheat = 2 and channel = 2 and area_code != \"0086\" and type = 1 and area_code != \"\" and convert_type = 2 GROUP BY area_code HAVING sum >= 100) as a) and is_cheat = 2 and channel = 2 and area_code = \"0033\" and type = 1 and convert_type = 2 ORDER BY id DESC"
	result, err := orm.Engine.QueryString(sql)

	if err != nil {
		fmt2.Printf("participant Find Err:%s", err.Error())
		os.Exit(-1)
	}
	pp.Println(result)
	for _, v := range result {
		fmt2.Printf("ID：%s\n", v["id"])
		result, err := ucloud.SmsSend("("+v["area_code"]+")", v["communication"])
		if err != nil {
			if err.Error() == "100000" {
				// 记录发送失败的数据
				var p participant.Participant
				if err = mapstructure.Decode(v, &p); err != nil {
					fmt2.Printf("map to struct err：%s\n", err.Error())
				}
				record(p, result.Message, result.SessionNo, result.RetCode, 2)
			} else {
				_, _ = pp.Println("未知错误意外终止")
				_, _ = pp.Println(err.Error())
				//os.Exit(-1)
			}
		}
		pp.Println(result)
	}
}

func query(limit, skip int, group *sync.WaitGroup) {
	defer group.Done()
	var p []participant.Participant
	err := orm.Engine.Where("channel = ? and is_cheat = ? and type = ?", 2, 2, 1).Desc("id").Limit(limit, skip).Find(&p)
	if err != nil {
		fmt2.Printf("participant Find Err:%s", err.Error())
		os.Exit(-1)
	}

	for _, v := range p {
		pp.Println("=======" + xstring.Int64ToString(v.Id) + "======")
		code := geo.GetCountryCode(v.Ip)
		pp.Println("CODE：", code)
		if code == "CN" {
			if v.AreaCode != "0086" {
				fmt2.Printf("符合要求：%d\n", v.Id)
				fmt2.Printf("IP：%s\n", v.Ip)
				fmt2.Printf("phone：%s\n", v.AreaCode+v.Communication)
				_, err = orm.Engine.Table("participant").Where("id = ?", v.Id).Update(map[string]interface{}{"is_cheat": 1})
			}
		} else {
			if v.AreaCode == "0086" {
				fmt2.Printf("符合要求：%d\n", v.Id)
				fmt2.Printf("IP：%s\n", v.Ip)
				fmt2.Printf("phone：%s\n", v.AreaCode+v.Communication)
				_, err = orm.Engine.Table("participant").Where("id = ?", v.Id).Update(map[string]interface{}{"is_cheat": 1})
			}
		}
	}
}

func query2(limit, skip int) {
	//defer group.Done()
	var p []participant.Participant
	err := orm.Engine.Where("channel = ? and is_cheat =?", 2, 2).Desc("id").Limit(limit, skip).Find(&p)
	if err != nil {
		fmt2.Printf("participant Find Err:%s", err.Error())
		os.Exit(-1)
	}

	for _, v := range p {
		fmt2.Printf("审查：%s\n", v.EthAddress)
		var ids []int64
		ids = append(ids, v.Id)
		var i []invite.Invite
		if err = orm.Engine.Where("invite_eth_address = ?", v.EthAddress).Find(&i); err != nil {
			fmt2.Printf("invite Find Err:%s", err.Error())
		}
		num := 0
		// 记录此邀请者以及他邀请的所有人的ID 方便后续操作
		for _, v1 := range i {
			var p2 participant.Participant
			if _, err = orm.Engine.Where("eth_address = ?", v1.InvitedEthAddress).Get(&p2); err != nil {
				fmt2.Printf("participant Get Err:%s", err.Error())
			}
			if p2.ScreenWithEth == 3 {
				num = num + 1
			}
			ids = append(ids, p2.Id)
		}
		// 如果邀请的虚假用户大于等于10个 那么邀请者及所有被邀请者都判定为虚假用户
		if num >= 10 {
			_, _ = pp.Println("符合要求的所有用户ID")
			_, _ = pp.Println(ids)
			//_, err = orm.Engine.Table("participant").In("id", ids).Update(map[string]interface{}{"is_cheat": 1})
		}
	}
}

func query3(limit, skip int) {
	//defer group.Done()
	var p []participant.Participant
	err := orm.Engine.Where("channel = ? and is_cheat = ? and invite_num >= ?", 2, 2, 10).Desc("id").Limit(limit, skip).Find(&p)
	if err != nil {
		fmt2.Printf("participant Find Err:%s", err.Error())
		os.Exit(-1)
	}

	for _, v := range p {
		fmt2.Printf("审查：%s\n", v.EthAddress)
		var ids []int64
		ids = append(ids, v.Id)
		var i []invite.Invite
		if err = orm.Engine.Where("invite_eth_address = ?", v.EthAddress).Find(&i); err != nil {
			fmt2.Printf("invite Find Err:%s", err.Error())
		}
		num := 0
		// 记录此邀请者以及他邀请的所有人的ID 方便后续操作
		for _, v1 := range i {
			var p2 participant.Participant
			if _, err = orm.Engine.Where("eth_address = ? and area_code = ? and (communication like ? or communication like ? or communication like ? or communication like ? or communication like ? or communication like ? or communication like ? or communication like ? or communication like ?)", v1.InvitedEthAddress, "0086", "141%", "144%", "145%", "146%", "162%", "165%", "167%", "1740%", "184%").Get(&p2); err != nil {
				fmt2.Printf("participant Get Err:%s", err.Error())
			}
			if p2.Id != 0 {
				num = num + 1
			}
			ids = append(ids, p2.Id)
		}
		// 如果邀请的虚假用户大于等于10个 那么邀请者及所有被邀请者都判定为虚假用户
		if num >= 10 {
			_, _ = pp.Println("符合要求的所有用户ID")
			_, _ = pp.Println(ids)
			_, err = orm.Engine.Table("participant").In("id", ids).Update(map[string]interface{}{"is_cheat": 1})
		}
	}
}

func screen(limit, skip int, group *sync.WaitGroup) {
	defer group.Done()
	var p []participant.Participant
	//count, _ := orm.Engine.Where("is_cheat = ? and channel = ? and id = ?", 2, 2, 323365).Desc("id").Limit(1000, 0).FindAndCount(&p)
	err := orm.Engine.Where("is_cheat = ? and channel = ? and screen_with_eth = ?", 2, 2, 1).Desc("id").Limit(limit, skip).Find(&p)
	if err != nil {
		pp.Println("Find Err")
		fmt2.Printf("Find Err:", err.Error())
		os.Exit(-1)
	}
	//fmt2.Printf("总数：%d\n", count)
	//numArr := services.GenerateRandomNumber(0, 99, 100)
	for _, v := range p {
		fmt2.Printf("ID：%d\n", v.Id)
		fmt2.Printf("ETH：%s\n", v.EthAddress)
		result, err := eth.Screen(v.EthAddress)
		pp.Println(result.Code)
		pp.Println(result.Msg)
		if err != nil {
			if err.Error() == "5000" {
				pp.Println("IP 无用")
				continue
			} else if err.Error() == "100000" {
				// 未查询到eth交易记录
				_, err = orm.Engine.Table("participant").Where("id = ?", v.Id).Update(map[string]interface{}{"is_cheat": 1, "screen_with_eth": 3})
			} else if err.Error() == "200000" {
				pp.Println("未知错误")
			}
		} else {
			// 查到eth交易记录
			_, err = orm.Engine.Table("participant").Where("id = ?", v.Id).Update(map[string]interface{}{"screen_with_eth": 2})
		}
		//count = count - 1
		//fmt2.Printf("========剩余未筛查数量=========：%d\n", count)
		//time.Sleep(60 * time.Second)
	}
	//pp.Println("================结束=================")
}

func query1() {
	sql := "SELECT * FROM participant WHERE ip IN ( SELECT a.ip  FROM ( SELECT ip, count( 1 ) AS \"sum\" FROM participant GROUP BY ip HAVING sum > 5 ) AS a )  and invite_num > 2"

	result, err := orm.Engine.QueryString(sql)
	if err != nil {
		pp.Println(err)
	}
	userNum := 0
	bitNum := 0
	for k, res := range result {
		fmt2.Printf("===数量====：%d\n", k)
		fmt2.Printf("ETH地址：%s\n", res["eth_address"])
		sql1 := fmt2.Sprintf("select count(p.id) as %s,sum(p.wait_integral+p.already_integral) as %s from invite i left join participant p on i.invited_eth_address = p.eth_address where i.invite_eth_address = %s and p.area_code = %s and p.type = 1 and p.channel = 2 and (p.communication LIKE '%s%s' OR p.communication LIKE '%s%s' OR p.communication LIKE '%s%s' OR p.communication LIKE '%s%s' OR p.communication LIKE '%s%s' OR p.communication LIKE '%s%s' OR p.communication LIKE '%s%s' OR p.communication LIKE '%s%s' OR p.communication LIKE '%s%s' OR p.communication LIKE '%s%s') HAVING num > 0",
			"num", "bit_num", res["eth_address"], "0086", "184", "%", "162", "%", "165", "%", "167", "%", "171", "%", "170", "%", "141", "%", "144", "%", "146", "%", "1740", "%")

		result1, err := orm.Engine.QueryString(sql1)
		//fmt2.Printf("SQL:%s\n", sql1)
		if err != nil {
			pp.Println(err)
		}

		for _, res1 := range result1 {
			aa := xint.StrToInt(res1["num"])
			//bb := xint.StrToInt(res1["bit_num"])
			fmt2.Printf("符合条件数量：%d\n", aa)
			if aa > 0 {
				userNum = userNum + 1
				bb := xint.StrToInt(res["wait_integral"]) + xint.StrToInt(res["already_integral"])
				bitNum = bitNum + bb
			}
		}
	}
	pp.Println("用户数")
	pp.Println(userNum)
	pp.Println("bi数")
	pp.Println(bitNum)
}

func sms1() {
	// 查询需要通知的账号
	var p []participant.Participant
	//count, err := orm.Engine.Where("communication not like ? and communication not like ? and communication not like ? and communication not like ? and communication not like ? and communication not like ? and type = ? and area_code = ? and id < ?", "184%", "162%", "165%", "167%", "171%", "170%", 1, "0086", 59748).Desc("id").FindAndCount(&p)
	err := orm.Engine.Where("communication not like ? and communication not like ? and communication not like ? and communication not like ? and communication not like ? and communication not like ? and type = ? and area_code = ? and id < ?", "184%", "162%", "165%", "167%", "171%", "170%", 1, "0086", 59748).Desc("id").Find(&p)
	//err := orm.Engine.Where("communication = ? or communication = ? or communication = ?", "13773652841", "17807566070", "17605083193").Asc("id").Limit(10, 0).Find(&p)
	if err != nil {
		fmt.Color.Println("FindAndCount Err")
		fmt.Color.Println(err.Error())
	}

	for _, v := range p {
		fmt2.Printf("ID：%d\n", v.Id)
		result, err := ucloud.SmsSend(v.AreaCode, v.Communication)
		if err != nil {
			if err.Error() == "100000" {
				// 记录发送失败的数据
				record(v, result.Message, result.SessionNo, result.RetCode, 2)
			} else {
				_, _ = pp.Println("未知错误意外终止")
				_, _ = pp.Println(err.Error())
				os.Exit(-1)
			}
		}
		pp.Println(result)
	}
}

func email() {
	// 查询需要通知的账号
	var participants []participant.Participant
	err := orm.Engine.Table("participant").Where("communication like ? or communication like ? or communication like ? and type = ? and id < ?", "%yahoo.com", "%outlook.com", "%gmail.com", 2, 324248).Desc("id").Limit(3000, 10).Find(&participants)
	//err := orm.Engine.Table("participant").Where("communication like ? or communication like ? and type = ?", "%13773652841@163.com", "%wei@wikiglobal.com", 2).Find(&participants)
	if err != nil {
		fmt.Color.Println("Email Get Err")
		fmt.Color.Println(err.Error())
		os.Exit(-1)
	}
	//_, _ = pp.Println(participants)
	//arr := []string{"13773652841@163.com", "wei@wikiglobal.com"}
	for _, v := range participants {
		fmt2.Printf("ID：%d\n", v.Id)
		var sendEmailCodeRequest oauth.SendEmailCodeRequest
		sendEmailCodeRequest.Email = v.Communication
		val, err := oauth2.SendEmailNotice(sendEmailCodeRequest)
		if val == false {
			record(v, err.Error(), "", -1, 2)
		}
	}
}

// record(p, result.Message, result.SessionNo, result.RetCode, 2)
func record(v participant.Participant, message, sessionNo string, retCode, result int64) {

	var participantNoticeRecord participantnoticerecord.ParticipantNoticeRecord

	participantNoticeRecord.Result = result
	participantNoticeRecord.EthAddress = v.EthAddress
	participantNoticeRecord.Communication = v.Communication
	participantNoticeRecord.Type = v.Type
	participantNoticeRecord.AreaCode = v.AreaCode
	participantNoticeRecord.CountryCode = v.CountryCode
	participantNoticeRecord.Message = message
	participantNoticeRecord.RetCode = retCode
	participantNoticeRecord.SessionNo = sessionNo

	if _, err := participantNoticeRecord.Create(); err != nil {
		_, _ = pp.Println("新增通知记录失败")
		_, _ = pp.Println(err.Error())
		os.Exit(1)
	}
}
