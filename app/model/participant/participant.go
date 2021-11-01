package participant

import (
	"errors"
	"fmt"
	"github.com/k0kubun/pp"
	"wiki_bit/app/model/invite"
	"wiki_bit/app/model/record"
	orm "wiki_bit/boot/db/mysql"
	"wiki_bit/boot/log"
	"wiki_bit/config"
	"wiki_bit/library/convert/xtime"

	mapstructure1 "github.com/mitchellh/mapstructure"
	"go.uber.org/zap"
	"xorm.io/xorm"
)

type Participant struct {
	Id                int64  `json:"id" xorm:"autoincr"`
	EthAddress        string `json:"eth_address" xorm:"eth_address"`
	Ip                string `json:"ip" xorm:"ip"`
	Type              int    `json:"type" xorm:"type"`
	AreaCode          string `json:"area_code" xorm:"area_code"`
	Communication     string `json:"communication" xorm:"communication"`
	WaitIntegral      int64  `json:"wait_integral" xorm:"wait_integral"`
	AlreadyIntegral   int64  `json:"already_integral" xorm:"already_integral"`
	ShareUrl          string `json:"share_url" xorm:"share_url"`
	RewardIntegral    int64  `json:"reward_integral" xorm:"reward_integral"`
	InviteNum         int64  `json:"invite_num" xorm:"invite_num"`
	Over20Reward      int64  `json:"over_20_reward" xorm:"over_20_reward"`
	ChristmasReward   int64  `json:"christmas_reward" xorm:"christmas_reward"`
	Language          string `json:"language" xorm:"language"`
	Event             string `json:"event" xorm:"-"`
	Channel           int    `json:"channel" xorm:"channel"`
	ChannelId         int64  `json:"channel_id" xorm:"channel_id"`
	CountryCode       string `json:"country_code" xorm:"country_code"`
	ConvertType       int64  `json:"convert_type" xorm:"convert_type"`
	IsCheat           int64  `json:"is_cheat" xorm:"is_cheat"`
	CreateAt          int64  `json:"create_at" xorm:"created"`
	BindCommunication string `json:"bind_communication" xorm:"bind_communication"`
	BindAreaCode      string `json:"bind_area_code" xorm:"bind_area_code"`
	UserId            string `json:"user_id" xorm:"user_id"`
	ScreenWithEth     int64  `json:"screen_with_eth" xorm:"screen_with_eth"`
}

type LoginRequest struct {
	Type          int    `json:"type"`
	EthAddress    string `json:"eth_address"`
	AreaCode      string `json:"area_code"`
	Communication string `json:"communication"`
}

type QueryRequest struct {
	Type      int    `json:"type"`
	AreaCode  string `json:"area_code"`
	ChannelId int64  `json:"channel_id"`
	StartTime int64  `json:"start_time"`
	EndTime   int64  `json:"end_time"`
}

type QueryResponse struct {
	AreaCode  string `json:"areacode"`
	UserCount string `json:"usercount"`
	CoinCount string `json:"coincount"`
}

type ReceiveType struct {
	Type       int64  `json:"type"` // 1 不合法用户 2 已参与并已领取完-不可以继续领取 3 已参与还可以领取 4 未参与
	Url        string `json:"url"`
	Received   int64  `json:"received"`
	EthAddress string `json:"eth_address"`
	Quota      int64  `json:"quota"`
}

type RechargeRequest struct {
	UserId     string `json:"user_id" binding:"required"`
	EthAddress string `json:"eth_address" binding:"required"`
	//Type          int64  `json:"type" binding:"required"` // 1 匹配上 2 未匹配上
	AreaCode      string `json:"area_code"`
	Communication string `json:"communication" binding:"required"`
}

// Create 创建
func (p Participant) Create(request record.RecordRequest) (participant Participant, err error) {

	// 以下所有操作受事务控制 保证原子性 唯一性
	session := orm.Engine.NewSession()
	defer session.Close()
	if err := session.Begin(); err != nil {
		log.Logger().Error("participant Create 开启事务 err：", zap.Error(err))
		return Participant{}, errors.New("500")
	}

	// 判断此人是否是被邀请
	if p.Event != "" && request.SurplusNum > 150 {
		invitees, val, err := p.updateByInviter(session)
		if err != nil || !val {
			_ = session.Rollback()
			return Participant{}, err
		}
		// 每个被邀请人做渠道号邀请记录
		if invitees.Channel == 1 {
			p.ChannelId = invitees.Id
		} else {
			p.ChannelId = invitees.ChannelId
		}
	}

	// 插入
	if _, err = session.Cols("id", "eth_address", "ip", "type", "area_code", "communication", "wait_integral", "already_integral", "share_url", "reward_integral", "invite_num", "over_20_reward", "language", "channel_id", "country_code", "create_at").InsertOne(&p); err != nil {
		log.Logger().Error("model participant Create InsertOne err：", zap.Error(err))
		_ = session.Rollback()
		return Participant{}, errors.New("500")
	}

	_ = session.Commit()

	return p, nil
}

// Update 更新邀请人的相关信息
func (p Participant) updateByInviter(session *xorm.Session) (Participant, bool, error) {

	var participant Participant

	// 先查询邀请人信息
	if _, err := session.Cols("id", "eth_address", "ip", "area_code", "communication", "wait_integral", "already_integral", "share_url", "reward_integral", "invite_num", "over_20_reward", "language", "channel", "channel_id", "create_at").Where("BINARY share_url = ?", config.Conf().Url.Website+p.Event).Get(&participant); err != nil {
		log.Logger().Error("model participant updateByInviter Get err：", zap.Error(err))
		return Participant{}, false, err
	}

	// 根据不同情况处理不同业务
	if participant.Channel == 1 {
		participant.InviteNum = participant.InviteNum + 1
	} else {
		if participant.InviteNum == 19 && participant.Over20Reward == 0 {
			participant.InviteNum = participant.InviteNum + 1
			participant.Over20Reward = 500
			participant.RewardIntegral = participant.RewardIntegral + 50
			// 活动当天
			if xtime.Second() >= config.Conf().Activity.StartDate && xtime.Second() <= config.Conf().Activity.EndDate {
				participant.WaitIntegral = participant.WaitIntegral + 1050
				participant.ChristmasReward = 500
			} else {
				participant.WaitIntegral = participant.WaitIntegral + 550
			}
		} else if participant.InviteNum >= 20 && participant.Over20Reward > 0 {
			// 不做处理直接返回
			return Participant{}, true, nil
		} else if participant.InviteNum < 19 && participant.Over20Reward == 0 {
			participant.InviteNum = participant.InviteNum + 1
			participant.WaitIntegral = participant.WaitIntegral + 50
			participant.RewardIntegral = participant.RewardIntegral + 50
		}
	}

	if _, err := session.Where("share_url = ?", config.Conf().Url.Website+p.Event).Cols("invite_num", "wait_integral", "reward_integral", "over_20_reward", "channel_id", "christmas_reward").Update(&participant); err != nil {
		log.Logger().Error("model participant update Update err：", zap.Error(err))
		return Participant{}, false, errors.New("500")
	}

	// 记录邀请人
	_, err := invite.Invite{InviteEthAddress: participant.EthAddress, InvitedEthAddress: p.EthAddress}.Create(session)
	if err != nil {
		return Participant{}, false, err
	}

	return participant, true, nil
}

// FindOne 查询单条
func (p Participant) FindOne() (participant Participant, err error) {

	if _, err = orm.Engine.Cols("id", "eth_address", "ip", "area_code", "communication", "type", "wait_integral", "already_integral", "share_url", "reward_integral", "invite_num", "over_20_reward", "christmas_reward", "language", "user_id", "is_cheat", "create_at").Where("eth_address = ?", p.EthAddress).Get(&participant); err != nil {
		log.Logger().Error("model participant FindOne Get err：", zap.Error(err))
		return Participant{}, errors.New("500")
	}

	return participant, nil
}

// FindOneByUserId 查询单条
func (p Participant) FindOneByUserId() (participant Participant, err error) {

	if _, err = orm.Engine.Cols("id", "user_id", "is_cheat").Where("user_id = ?", p.UserId).Get(&participant); err != nil {
		log.Logger().Error("model participant FindOneByUserId Get err：", zap.Error(err))
		return Participant{}, errors.New("500")
	}

	return participant, nil
}

// UpdateCommunication 更新通信
func (p Participant) UpdateCommunication() (participant Participant, err error) {

	// 先验证此号码不能重复参与
	if p.Type == 1 {
		_, err = orm.Engine.Cols("id").Where("communication = ? and area_code = ?", p.Communication, p.AreaCode).Get(&participant)
	} else {
		_, err = orm.Engine.Cols("id").Where("communication = ?", p.Communication).Get(&participant)
	}
	if err != nil {
		log.Logger().Error("model participant Create Get err：", zap.Error(err))
		return Participant{}, err
	}

	if participant.Id != 0 {
		return Participant{}, errors.New("此号码已参与过活动，不能重复参与")
	}

	if _, err = orm.Engine.Where("eth_address = ?", p.EthAddress).Cols("area_code", "communication", "type", "share_url").Update(&p); err != nil {
		log.Logger().Error("model participant update Update err：", zap.Error(err))
		return Participant{}, err
	}

	return p, nil
}

// 校验邀请链接的合法性
func (p Participant) FindByShareUrl(shareUrl string) (participant Participant, err error) {

	if _, err = orm.Engine.Cols("id", "eth_address", "ip", "area_code", "communication", "type", "wait_integral", "already_integral", "share_url", "reward_integral", "invite_num", "over_20_reward", "christmas_reward", "language", "create_at").Where("BINARY share_url = ?", shareUrl).Get(&participant); err != nil {
		log.Logger().Error("model participant FindByShareUrl Get err：", zap.Error(err))
		return Participant{}, errors.New("500")
	}

	return participant, nil
}

// 校验手机号或者邮箱的唯一性
func (p Participant) FindByCommunication() (participant Participant, err error) {

	if p.Type == 1 {
		_, err = orm.Engine.Where("communication = ? and area_code = ?", p.Communication, p.AreaCode).Get(&participant)
	} else {
		_, err = orm.Engine.Where("communication = ?", p.Communication).Get(&participant)
	}
	if err != nil {
		log.Logger().Error("model participant FindByCommunication Get err：", zap.Error(err))
		return Participant{}, errors.New("500")
	}

	return participant, nil
}

// 领取活动
func (p Participant) Receive() (convertAmount int64, ethAddress string, err error) {

	var participant Participant

	if p.Type == 1 {
		_, err = orm.Engine.Where("communication = ? and area_code = ? and eth_address = ?", p.Communication, p.AreaCode, p.EthAddress).Get(&participant)
	} else {
		_, err = orm.Engine.Where("communication = ? and eth_address = ?", p.Communication, p.EthAddress).Get(&participant)
	}
	if err != nil {
		log.Logger().Error("model participant FindByCommunicationAndEthAddress Get err：", zap.Error(err))
		return 0, "", errors.New("500")
	}

	if participant.Id == 0 {
		return 0, "", errors.New("1010025")
	}

	return participant.WaitIntegral, participant.EthAddress, nil
}

// 根据通信或者以太坊地址查询
func (p Participant) FindByCommunicationOrEthAddress() (participant Participant, err error) {

	if p.EthAddress != "" {
		_, err = orm.Engine.Where("eth_address = ?", p.EthAddress).Get(&participant)
	} else {
		_, err = orm.Engine.Where("communication = ? and area_code = ?", p.Communication, p.AreaCode).Get(&participant)
	}

	if err != nil {
		log.Logger().Error("model participant FindByCommunicationOrEthAddress Get err：", zap.Error(err))
		return Participant{}, errors.New("500")
	}

	return participant, nil
}

// 根据ip查询领取次数
func (p Participant) FindCountByIp() (count int64, err error) {
	var participant Participant
	if count, err = orm.Engine.Cols(`id`).Where("ip = ?", p.Ip).Count(&participant); err != nil {
		log.Logger().Error("model participant FindCountByIp Count err：", zap.Error(err))
		return 0, errors.New("500")
	}

	return count, nil
}

// 根据ip查询领取次数
func (q QueryRequest) Query(type1, startTime, endTime, channelId int64, areaCode string) (queryResponse []QueryResponse, err error) {

	sql := "select area_code areacode,count(1) as usercount, sum(wait_integral+already_integral) as coincount from participant where create_at BETWEEN  %d and %d and channel = 2 and type = %d"
	if type1 == 2 {
		sql = "select count(1) as usercount, sum(wait_integral+already_integral) as coincount from participant where create_at BETWEEN  %d and %d and channel = 2 and type = %d"
	}

	sql = fmt.Sprintf(sql, startTime, endTime, type1)
	if channelId != 0 {
		sql = sql + " AND channel_id = %d"
		sql = fmt.Sprintf(sql, channelId)
	}
	if areaCode != "" {
		sql = sql + " AND area_code = %s"
		sql = fmt.Sprintf(sql, areaCode)
	}

	if type1 == 1 {
		sql = sql + " GROUP BY area_code"
	}

	result, err := orm.Engine.QueryString(sql)
	if err != nil {
		log.Logger().Error("participant Query QueryString err：", zap.Error(err))
		return nil, errors.New("500")
	}
	for _, res := range result {
		var r QueryResponse
		if err = mapstructure1.Decode(res, &r); err != nil {
			pp.Println(err.Error())
		}
		queryResponse = append(queryResponse, r)
	}
	return queryResponse, nil
}

func (p Participant) UpdateConvert() error {

	if _, err := orm.Engine.Where("id = ?", p.Id).Cols("convert_type", "convert_amount").Update(&p); err != nil {
		log.Logger().Error("model participant UpdateConvert Update err：", zap.Error(err))
		return err
	}

	return nil
}

// 更新
func (p Participant) UpdateByID(cols []string) error {

	if _, err := orm.Engine.Where("id = ?", p.Id).Cols(cols...).Update(&p); err != nil {
		log.Logger().Error("model participant UpdateByID Update err：", zap.Error(err))
		return errors.New("500")
	}

	return nil
}

// 校验手机号或者邮箱的唯一性
func (p Participant) FindByUserId() (participant Participant, err error) {

	_, err = orm.Engine.Where("user_id = ?", p.UserId).Get(&participant)
	if err != nil {
		log.Logger().Error("model participant FindByUserId Get err：", zap.Error(err))
		return Participant{}, errors.New("500")
	}

	return participant, nil
}

// FindOne 查询单条
func (p Participant) FindByEthAddressAndCommunication() (participant Participant, err error) {

	table := orm.Engine.Table("participant").Cols("id", "eth_address", "ip", "area_code", "communication", "type", "wait_integral", "already_integral", "share_url", "reward_integral", "invite_num", "over_20_reward", "christmas_reward", "language", "create_at", "user_id")
	table = table.Where("eth_address = ?", p.EthAddress)
	if p.Type == 1 {
		table = table.Where("area_code = ? and communication = ?", p.AreaCode, p.Communication)
	} else {
		table = table.Where("communication = ?", p.Communication)
	}
	if _, err = table.Get(&participant); err != nil {
		log.Logger().Error("model participant FindByEthAddressAndCommunication Get err：", zap.Error(err))
		return Participant{}, errors.New("500")
	}

	return participant, nil
}
