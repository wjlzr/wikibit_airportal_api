package validate

import "regexp"

// 校验以太坊地址准确性
func ValidateAddress(address string) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	if re.MatchString(address) {
		return true
	}
	return false
}

// 校验邮箱
func VerifyEmailFormat(email string) bool {

	pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}
