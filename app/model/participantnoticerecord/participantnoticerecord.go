package participantnoticerecord

import (
	"go.uber.org/zap"
	orm "wiki_bit/boot/db/mysql"
	"wiki_bit/boot/log"
)

type ParticipantNoticeRecord struct {
	Id            int64  `json:"id" xorm:"autoincr"`
	EthAddress    string `json:"eth_address" xorm:"eth_address"`
	Type          int    `json:"type" xorm:"type"`
	AreaCode      string `json:"area_code" xorm:"area_code"`
	Communication string `json:"communication" xorm:"communication"`
	CountryCode   string `json:"country_code" xorm:"country_code"`
	Result        int64  `json:"result" xorm:"result"`
	RetCode       int64  `json:"ret_code" xorm:"ret_code"`
	Message       string `json:"message" xorm:"message"`
	SessionNo     string `json:"session_no" xorm:"session_no"`
	CreateAt      int64  `json:"create_at" xorm:"created"`
}

// Create 创建
func (p *ParticipantNoticeRecord) Create() (participantNoticeRecord *ParticipantNoticeRecord, err error) {

	if _, err = orm.Engine.Cols("id", "eth_address", "type", "area_code", "communication", "country_code", "result", "ret_code", "message", "session_no", "create_at").InsertOne(p); err != nil {
		log.Logger().Error("model participant Create InsertOne err：", zap.Error(err))
		return nil, err
	}
	return p, nil
}
