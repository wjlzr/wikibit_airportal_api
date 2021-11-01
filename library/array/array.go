package array

func StrInArray(array []string, str string) bool {
	for _, v := range array {
		if v == str {
			return true
		}
	}
	return false
}
