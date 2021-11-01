package services

import (
	"github.com/k0kubun/pp"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
	"wiki_bit/boot/log"
)

//统一请求分发
func Request(method, url string, body io.Reader) (request *http.Request, err error) {

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		log.Logger().Error("services http request err：", zap.Error(err))
		return request, err
	}
	if method == http.MethodPost {
		req.Header.Set("Content-Type", "application/json")
	}
	// 设置ip代理
	return req, nil
}

//返回参数统一处理
func ResponseHandle(request *http.Request) (content []byte, err error) {
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		pp.Println(err.Error())
		return content, err

	}
	defer resp.Body.Close()
	content, _ = ioutil.ReadAll(resp.Body)
	//_, _ = pp.Println("调用OpenApi接口")
	//_, _ = pp.Println(resp.Request)
	//fmt.Printf("调用OpenApi接口：%+v \n", resp.Request)
	//_, _ = pp.Println("OpenApi返回值")
	//fmt.Printf("OpenApi返回值：%s \n", string(content))
	return content, nil
}

func GenerateRangeNum(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	randNum := rand.Intn(max-min) + min
	return randNum
}

// 生成随机数
func GenerateRandomNumber(start int, end int, count int) []int {
	//范围检查
	if end < start {
		return nil
	}

	//存放结果的slice
	nums := make([]int, 0)
	//随机数生成器，加入时间戳保证每次生成的随机数不一样
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(nums) < count {
		//生成随机数
		num := r.Intn((end - start)) + start

		//查重
		exist := false
		for _, v := range nums {
			if v == num {
				exist = true
				break
			}
		}

		if !exist {
			nums = append(nums, num)
		}
	}

	return nums
}
