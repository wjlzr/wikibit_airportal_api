package responseerror

import (
	"fmt"
	orm "wiki_bit/boot/db/mysql"
)

type Error struct {
	Code int64  `json:"code"`
	ZhCN int64  `json:"zh-CN"`
	ZhHK string `json:"zh-HK"`
	ZhTW string `json:"zh-TW"`
	En   string `json:"en"`
	Vi   string `json:"vi"`
	Th   string `json:"th"`
}

func GetError(code, lang string) string {

	sql := fmt.Sprintf("SELECT * FROM `error` WHERE code = '%s'", code)
	result, _ := orm.Engine.QueryString(sql)
	for _, val := range result {
		for k, v := range val {
			if k == lang {
				return v
			}
		}
	}
	return ""
}
