package oauth

import (
	"fmt"
	"time"
	"wiki_bit/app/model/oauth"
	"wiki_bit/boot/db/redis"
	"wiki_bit/boot/log"
	"wiki_bit/config"
	"wiki_bit/library/validatecode"
	"wiki_bit/services/email"

	"go.uber.org/zap"
)

// SendEmailCode 发送邮箱验证码service
func SendEmailCode(req oauth.SendEmailCodeRequest) (val bool, err error) {

	validateCodeStr := validatecode.GenValidateCode(4)

	emailMessage := email.NewEmailMessage("WikiBit", config.Conf().Email.ContentType, "Verification code："+validateCodeStr, "", []string{req.Email}, []string{})

	if val, err = email.NewEmailClient(emailMessage).SendMessage(); !val || err != nil {
		fmt.Printf("\n邮箱验证码发送错误-邮箱：%s \n", req.Email)
		log.Logger().Error("service oauth SendEmailCode SendMessage Err：", zap.Error(err))
		return false, err
	}

	if err = redis.Client().Set("email_validate_code_"+req.Email, validateCodeStr, 300*time.Second).Err(); err != nil {
		log.Logger().Error("service oauth SendEmailCode Redis Set Err：", zap.Error(err))
		return false, err
	}

	return true, nil
}

// ValidateEmailCode 验证邮箱验证码service
func ValidateEmailCode(email, code string) (val bool, err error) {

	rcode, err := redis.Client().Get("email_validate_code_" + email).Result()
	if err != nil || rcode == "" || rcode != code {
		log.Logger().Error("service oauth ValidateEmailCode Redis Get Err：", zap.Error(err))
		return false, err
	}

	// 验证完删除验证码
	redis.Client().Del("email_validate_code_" + email)

	return true, nil
}

// SendEmailNotice 发送邮箱通知消息service
func SendEmailNotice(req oauth.SendEmailCodeRequest) (val bool, err error) {

	emailMessage := email.NewEmailMessage("WikiBit", config.Conf().Email.ContentType, "亲爱的用户，WikiBit已全部发放于区块天眼APP中，请尽快领取，领取步骤：https://www.wikibit.cc/ 如有疑问，请添加官方号wikibit002咨询。", "", []string{req.Email}, []string{})

	if val, err = email.NewEmailClient(emailMessage).SendMessage(); !val || err != nil {
		fmt.Printf("\n邮箱验证码发送错误-邮箱：%s \n", req.Email)
		log.Logger().Error("service oauth SendEmailCode SendMessage Err：", zap.Error(err))
		return false, err
	}

	return true, nil
}
