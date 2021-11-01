package main

import (
	fmt2 "fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/k0kubun/pp"
	"github.com/xluohome/phonedata"
	"os"
	"wiki_bit/app/model/participant"
	orm "wiki_bit/boot/db/mysql"
	"wiki_bit/boot/log"
	"wiki_bit/config"
	"wiki_bit/library/convert/xstring"
)

func init() {
	// 初始化配置文件
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
	//geo1("13501848687")
	var p []participant.Participant
	err := orm.Engine.Where("area_code = ? and type = ? and channel = ? and is_cheat = ? ", "0086", 1, 2, 2).Desc("id").Find(&p)
	if err != nil {
		fmt2.Printf("participant Find Err:%s", err.Error())
		os.Exit(-1)
	}
	f := excelize.NewFile()
	_ = f.SetCellValue("Sheet1", "A1", "区号")
	_ = f.SetCellValue("Sheet1", "B1", "手机号")
	_ = f.SetCellValue("Sheet1", "C1", "IP")
	i := 1
	for _, v := range p {
		fmt2.Printf("ID: %d\n", v.Id)
		//code := geo.GetCityCode(v.Ip)
		//fmt2.Printf("CODE: %s\n", code)
		//if code == "Shanghai" {
		info, err := phonedata.Find(v.Communication)
		if err != nil {
			continue
		}
		if info.City == "上海" {
			i = i + 1
			num := xstring.Int64ToString(int64(i))
			_ = f.SetCellValue("Sheet1", "A"+num, v.AreaCode)
			_ = f.SetCellValue("Sheet1", "B"+num, v.Communication)
			_ = f.SetCellValue("Sheet1", "C"+num, v.Ip)
		}
		//}
	}
	// 根据指定路径保存文件
	if err := f.SaveAs("shanghai.xlsx"); err != nil {
		_, _ = pp.Println("生成excel失败")
		_, _ = pp.Println(err.Error())
	}
	_, _ = pp.Println("筛查已结束")
	fmt2.Printf("总筛查行数: %d\n", len(p))
}

func geo1(phone string) {
	info, err := phonedata.Find(phone)
	if err != nil {
		pp.Println(err)
	}
	pp.Println(info)
}
