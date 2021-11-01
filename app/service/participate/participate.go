package participate

import (
	"errors"
	"github.com/k0kubun/pp"
	"go.uber.org/zap"
	"time"
	"wiki_bit/boot/log"
	"wiki_bit/library/convert/xstring"
	aes "wiki_bit/library/encryption"
	"wiki_bit/library/fmt"
	"wiki_bit/library/geo"
	"wiki_bit/library/i18nresponse"
	"wiki_bit/services/wikibit"

	log4 "github.com/jeanphorn/log4go"

	"github.com/gin-gonic/gin"

	"wiki_bit/app/model/participant"
	"wiki_bit/app/model/record"
	"wiki_bit/config"
	"wiki_bit/library/constant"
	"wiki_bit/library/ids"
)

//type CreateBitConsumer struct{}

// Create Service
func Create(c *gin.Context, participant participant.Participant) (result participant.Participant, err error) {

	// 校验ip是否已经超过领取最大值
	count, err := participant.FindCountByIp()
	if err != nil {
		return result, errors.New("500")
	}
	if count >= 50 {
		return result, errors.New("1010017")
	}

	participant.Language = i18nresponse.GetLang(c)
	participant.Ip = c.ClientIP()
	participant.WaitIntegral = 100
	participant.CountryCode = geo.GetCountryCode(c.ClientIP())

	// 查询此地址是否认购过
	if result, err = GetInfo(participant.EthAddress); err != nil {
		return result, err
	}

	if result.Id != 0 {
		return result, errors.New("1010009")
	}

	// 查询当前剩余币数量是否满足认购
	recordRequest, err := record.Record{}.Statistics()
	if err != nil {
		return result, err
	}

	// 验证邀请链接的合法性
	if participant.Event != "" {
		p2, err := participant.FindByShareUrl(config.Conf().Url.Website + participant.Event)
		if err != nil {
			return result, err
		}
		if p2.Id == 0 {
			return result, errors.New("1010010")
		}
	}

	// 验证手机号或者邮箱是否已经参与过活动
	p3, err := participant.FindByCommunication()
	if err != nil {
		return result, err
	}
	if p3.Id != 0 {
		if participant.Type == 1 {
			return result, errors.New("1010011")
		} else if participant.Type == 2 {
			return result, errors.New("1010012")
		}
	}

	// 如果已经发放完毕则直接返回错误
	if recordRequest.Received >= constant.TotalAmountOfCurrencyIssued || (constant.TotalAmountOfCurrencyIssued-recordRequest.Received < 100) {
		return result, errors.New("1010014")
	}

	if (participant.Event == "" && recordRequest.SurplusNum <= 100 && recordRequest.IncTotal < constant.TotalAmountOfCurrencyIssued) || (participant.Event != "" && recordRequest.SurplusNum <= 150 && recordRequest.IncTotal < constant.TotalAmountOfCurrencyIssued) {
		// 增发货币
		_, err = record.Record{}.AddCurrency()
		if err != nil {
			return result, err
		}
	}

	// 生成邀请链接
	idWorker, err := ids.NewIdWorker(10)
	if err != nil {
		return result, errors.New("500")
	}
	code, err := idWorker.ShortId()
	if err != nil {
		return result, errors.New("500")
	}
	participant.ShareUrl = config.Conf().Url.Website + code

	if result, err = participant.Create(recordRequest); err != nil {
		return result, err
	}

	return result, nil
}

// UpdateCommunication Service
func UpdateCommunication(participant participant.Participant) (result participant.Participant, err error) {

	// 生成十位唯一码
	idWorker, err := ids.NewIdWorker(10)
	if err != nil {
		return result, err
	}
	code, err := idWorker.ShortId()
	if err != nil {
		return result, err
	}
	participant.ShareUrl = config.Conf().Url.Website + code

	if result, err = participant.UpdateCommunication(); err != nil {
		return result, err
	}

	return result, nil
}

// GetInfo 认购信息service
func GetInfo(ethAddress string) (p participant.Participant, err error) {

	p, err = participant.Participant{EthAddress: ethAddress}.FindOne()
	if err != nil {
		return participant.Participant{}, err
	}

	return p, nil
}

// Login 登录service
func Login(l participant.LoginRequest) (resp participant.Participant, err error) {

	resp, err = participant.Participant{Type: l.Type, Communication: l.Communication, AreaCode: l.AreaCode, EthAddress: l.EthAddress}.FindByCommunicationOrEthAddress()
	if err != nil {
		return participant.Participant{}, err
	}

	return resp, nil
}

// GetHomeInfo 获取首页数据 service
func GetHomeInfo() (recordRequest record.RecordRequest, err error) {

	recordRequest, err = record.Record{}.Statistics()
	if err != nil {
		return recordRequest, err
	}

	return recordRequest, nil
}

// Query
func Query(type1, startTime, endTime, channelId int64, areaCode string) (queryRequest []participant.QueryResponse, err error) {

	queryRequest, err = participant.QueryRequest{}.Query(type1, startTime, endTime, channelId, areaCode)
	if err != nil {
		return queryRequest, err
	}

	return queryRequest, nil
}

func ReceiveType(p participant.Participant) (r *participant.ReceiveType, err error) {

	// 根据phone查询详情
	result, err := p.FindByUserId()
	if err != nil {
		return nil, errors.New("500")
	}
	r = new(participant.ReceiveType)
	if result.Id == 0 {
		p.Type = 1
		result, err = p.FindByCommunication()
		if err != nil {
			return nil, errors.New("500")
		}
		if result.UserId != "" && result.UserId != p.UserId {
			// 匹配到了 但是已被其他bit账号绑定
			r.Type = 4
			return r, nil
		}
	}

	if result.Id != 0 && result.IsCheat == 1 {
		r.Type = 1
	} else if result.Id != 0 && result.ConvertType == 1 && result.WaitIntegral == 0 {
		r.Type = 2 // 已参与-已兑换-已领取完-不可以继续领取
	} else if result.Id != 0 && result.ConvertType == 1 && result.WaitIntegral > 0 {
		r.Type = 3 // 已参与-已兑换-还可以领取
	} else if result.Id != 0 && result.ConvertType == 2 && result.WaitIntegral == 0 {
		r.Type = 5 // 已参与-未兑换-已领取完-不可以继续领取
	} else if result.Id != 0 && result.ConvertType == 2 && result.WaitIntegral > 0 {
		r.Type = 6 // 已参与-未兑换-未领取完-可以继续领取
	} else if result.Id == 0 {
		r.Type = 4 // 未参与
	}

	if result.ChristmasReward == 500 {
		r.Quota = 2100
	} else {
		r.Quota = 1600
	}
	r.EthAddress = result.EthAddress
	r.Received = result.AlreadyIntegral
	return r, nil
}

func Receive(p participant.Participant) (convertAmount int64, ethAddress string, err error) {

	// 做唯一绑定验证
	result, err := p.FindByEthAddressAndCommunication()
	_, _ = pp.Println("唯一绑定验证")
	_, _ = pp.Println(result)
	fmt.Color.Println("唯一绑定验证")
	fmt.Color.Printf(result)
	if err != nil {
		return 0, "", errors.New("500")
	}
	if result.UserId != "" && result.UserId != p.UserId {
		return 0, "", errors.New("1010026")
	}

	if convertAmount, ethAddress, err = p.Receive(); err != nil {
		return 0, "", err
	}

	return convertAmount, ethAddress, nil
}

// 兑换
func Recharge(p participant.RechargeRequest) (result int64, err error) {
	fmt.Color.Println("领取请求参数")
	fmt.Color.Printf(p)
	resp, err := participant.Participant{EthAddress: p.EthAddress}.FindOne()
	if resp.Id == 0 || err != nil {
		return 0, errors.New("500")
	}

	if resp.IsCheat == 1 {
		return 0, errors.New("1010023")
	}

	if resp.WaitIntegral == 0 {
		return 0, errors.New("1010027")
	}

	if resp.UserId == "" {
		// 未绑定的情况下判断当前用户是否领过
		res, err := participant.Participant{UserId: p.UserId}.FindByUserId()
		if err != nil {
			return 0, errors.New("500")
		}
		if res.Id != 0 {
			return 0, errors.New("1010024")
		}
		// 新增绑定关系
		err = participant.Participant{Id: resp.Id, BindAreaCode: p.AreaCode, BindCommunication: p.Communication, UserId: p.UserId}.UpdateByID([]string{"bind_area_code", "bind_communication", "user_id"})
		if err != nil {
			return 0, errors.New("500")
		}
	} else {
		// 已绑定过 那么校验
		if p.UserId != resp.UserId {
			return 0, errors.New("1010026")
		}
	}

	var m wikibit.RechargeRequest
	m.UserId = p.UserId
	m.Secret = config.Conf().Wikibit.Secret
	m.AppId = config.Conf().Wikibit.AppId
	m.Money = xstring.Int64ToString(resp.WaitIntegral)
	r, err := wikibit.Recharge(m, p.Communication)
	if err != nil {
		return 0, errors.New("500")
	}

	if r.Succeed == false || r.Result == false {
		fmt.Color.Println("调用BIT返回的数据")
		fmt.Color.Printf(r)
		log.Logger().Info("service participate Recharge Err：", zap.Error(err))
		return 0, errors.New("500")
	}

	// 兑换完后修改数据
	received := resp.AlreadyIntegral + resp.WaitIntegral
	notReceived := int64(0)
	_ = participant.Participant{Id: resp.Id, AlreadyIntegral: received, WaitIntegral: notReceived, ConvertType: 1}.UpdateByID([]string{"wait_integral", "convert_type", "already_integral"})

	return resp.WaitIntegral, nil
}

// 记录日志
func Logging(p participant.Participant) {
	log4.LOGGER("wikibit").Info("id：%d", p.Id)
	log4.LOGGER("wikibit").Info("eth_address：%s", p.EthAddress)
	log4.LOGGER("wikibit").Info("type：%d", p.Type)
	log4.LOGGER("wikibit").Info("communication：%s", aes.AesEncrypt(p.Communication, config.Conf().Encryption.AesSecretKey))
	log4.LOGGER("wikibit").Info("share_url：%s", p.ShareUrl)
	log4.LOGGER("wikibit").Info("reward_integral：%d", p.RewardIntegral)
	log4.LOGGER("wikibit").Info("over_20_reward：%d", p.Over20Reward)
	log4.LOGGER("wikibit").Info("language：%s", p.Language)
	log4.LOGGER("wikibit").Info("partake_reward：%d", 100)
	log4.LOGGER("wikibit").Info("total_reward：%d", p.WaitIntegral+p.AlreadyIntegral)
	log4.LOGGER("wikibit").Info("channel_id：%d", p.ChannelId)
	log4.LOGGER("wikibit").Info("area_code：%s", p.AreaCode)
	log4.LOGGER("wikibit").Info("invite_num：%d", p.InviteNum)
	timeNow := time.Unix(p.CreateAt, 0)
	log4.LOGGER("wikibit").Info("logDate：%s", timeNow.Format("2006-01-02 15:04:05"))
	log4.Info("\n")
}

// 消费
//func (*CreateBitConsumer) HandleMessage(msg *nsq.Message) (err error) {
//	var p participant.Participant
//	_ = json.Unmarshal(msg.Body, &p)
//	if _, err = Create(p); err != nil {
//		return err
//	}
//	return nil
//}
