package record

import (
	"errors"
	"time"
	orm "wiki_bit/boot/db/mysql"
	"wiki_bit/boot/log"
	"wiki_bit/library/constant"
	"wiki_bit/library/convert/xint64"

	"go.uber.org/zap"
)

type Record struct {
	ID               int64 `json:"id" xorm:"id"`
	Inc              int64 `json:"inc" xorm:"inc"`
	IncNum           int64 `json:"inc_num" xorm:"inc_num"`
	IncTotal         int64 `json:"inc_total" xorm:"inc_total"`
	BuyBackStartTime int64 `json:"buy_back_start_time" xorm:"buy_back_start_time"`
	BuyBackEndTime   int64 `json:"buy_back_end_time" xorm:"buy_back_end_time"`
	CreateAt         int64 `json:"create_at" xorm:"created"`
}

// 统计返回
type RecordRequest struct {
	ID               int64 `json:"id"`
	Inc              int64 `json:"inc"`
	IncNum           int64 `json:"inc_num"`
	IncTotal         int64 `json:"inc_total"`
	BuyBackStartTime int64 `json:"buy_back_start_time"`
	BuyBackEndTime   int64 `json:"buy_back_end_time"`
	CurrentTime      int64 `json:"current_time"`
	Received         int64 `json:"received"`    // 已领取
	PartakeNum       int64 `json:"partake_num"` // 参与人数
	SurplusNum       int64 `json:"surplus_num"` // 剩余币数
	CreateAt         int64 `json:"create_at" xorm:"created"`
}

func (r Record) Statistics() (recordRequest RecordRequest, err error) {

	if recordRequest, err = r.FindOne(); err != nil {
		return RecordRequest{}, err
	}

	// 统计已领取和参与人数
	sql := "SELECT SUM(wait_integral + already_integral) received,count(1) partake_num FROM participant"
	result, err := orm.Engine.QueryString(sql)
	if err != nil {
		log.Logger().Error("record Statistics QueryString err：", zap.Error(err))
		return RecordRequest{}, errors.New("500")
	}

	for _, res := range result {
		for k, v := range res {
			if k == "received" {
				recordRequest.Received = xint64.StrToInt64(v)
			} else if k == "partake_num" {
				recordRequest.PartakeNum = xint64.StrToInt64(v)
			}
		}
	}
	recordRequest.SurplusNum = recordRequest.Inc*10000 - recordRequest.Received
	recordRequest.CurrentTime = time.Now().Unix()

	return recordRequest, nil
}

// 增发货币
func (r Record) AddCurrency() (val bool, err error) {

	var recordRequest RecordRequest
	if recordRequest, err = r.FindOne(); err != nil {
		return false, err
	}

	if recordRequest.IncTotal >= constant.TotalAmountOfCurrencyIssued {
		return false, errors.New("1010014")
	}

	recordRequest.Inc = recordRequest.Inc + 1
	recordRequest.IncTotal = recordRequest.IncTotal + recordRequest.IncNum

	if _, err = orm.Engine.Table("record").Cols("id", "inc", "inc_num", "inc_total", "buy_back_start_time", "buy_back_end_time", "create_at").InsertOne(&recordRequest); err != nil {
		log.Logger().Error("record AddCurrency InsertOne err：", zap.Error(err))
		return false, errors.New("500")
	}

	return true, nil
}

// 根据条件查单条
func (r Record) FindOne() (recordRequest RecordRequest, err error) {

	if _, err := orm.Engine.Table("record").Cols("id", "inc", "inc_num", "inc_total", "buy_back_start_time", "buy_back_end_time", "create_at").OrderBy("id desc").Get(&recordRequest); err != nil {
		log.Logger().Error("record FindOne Get err：", zap.Error(err))
		return RecordRequest{}, errors.New("500")
	}

	return recordRequest, nil
}
