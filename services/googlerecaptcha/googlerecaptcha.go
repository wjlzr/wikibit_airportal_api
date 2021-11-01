package googlerecaptcha

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const siteVerifyURL = "https://www.google.com/recaptcha/api/siteverify"
const secret = "6Lc0ef8ZAAAAAAWAHRf39RWyoYh-v6KIQNm6jfFn"

type SiteVerifyResponse struct {
	Success     bool      `json:"success"`
	Score       float64   `json:"score"`
	Action      string    `json:"action"`
	ChallengeTS time.Time `json:"challenge_ts"`
	Hostname    string    `json:"hostname"`
	ErrorCodes  []string  `json:"error-codes"`
}

// 校验
func CheckRecaptcha(response, communication string) (val bool, err error) {
	req, err := http.NewRequest(http.MethodPost, siteVerifyURL, nil)
	if err != nil {
		return false, err
	}

	// Add necessary request parameters.
	q := req.URL.Query()
	q.Add("secret", secret)
	q.Add("response", response)
	req.URL.RawQuery = q.Encode()

	// Make request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	// Decode response.
	var body SiteVerifyResponse
	if err = json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return false, err
	}

	// Check recaptcha verification success.
	if !body.Success {
		fmt.Printf("Google人机验证失败 通信方式：%s 错误：%+v\n", communication, body.ErrorCodes)
		return false, err
	}

	// Check response score.
	if body.Score < 0 {
		fmt.Printf("Google人机验证得分过低 通信方式：%s 分数：%+v\n", communication, body.Score)
		return false, nil
	}

	return true, nil
}
