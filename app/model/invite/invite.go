package invite

import (
	"errors"
	"wiki_bit/boot/log"

	"go.uber.org/zap"
	"xorm.io/xorm"
)

type Invite struct {
	Id                int64  `json:"id" xorm:"id"`
	InviteEthAddress  string `json:"invite_eth_address" xorm:"invite_eth_address"`
	InvitedEthAddress string `json:"invited_eth_address" xorm:"invited_eth_address"`
	CreateAt          int64  `json:"create_at" xorm:"created"`
}

type Test struct {
	Aa int64 `json:"id" xorm:"id"`
}

// Create 创建邀请人映射信息
func (i Invite) Create(session *xorm.Session) (Invite, error) {

	if _, err := session.InsertOne(i); err != nil {
		log.Logger().Error("model invite Create InsertOne err：", zap.Error(err))
		return Invite{}, errors.New("500")
	}

	return i, nil
}
