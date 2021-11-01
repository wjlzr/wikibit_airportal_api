package configure

import (
	"errors"
	"go.uber.org/zap"
	orm "wiki_bit/boot/db/mysql"
	"wiki_bit/boot/log"
)

type Configure struct {
	Id               int64 `json:"id" xorm:"id"`
	EnableMailbox    int64 `json:"enable_mailbox" xorm:"enable_mailbox"`
	BuyBackStartTime int64 `json:"buy_back_start_time" xorm:"buy_back_start_time"`
	BuyBackEndTime   int64 `json:"buy_back_end_time" xorm:"buy_back_end_time"`
	CreateAt         int64 `json:"create_at" xorm:"created"`
}

func (r *Configure) FindOne() (configure Configure, err error) {
	if _, err = orm.Engine.Where("id = ?", 1).Get(&configure); err != nil {
		log.Logger().Error("model configure FindOne Get err：", zap.Error(err))
		return Configure{}, errors.New("500")
	}

	return configure, nil
}
