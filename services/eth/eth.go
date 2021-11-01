package eth

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/k0kubun/pp"
	"go.uber.org/zap"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
	"wiki_bit/boot/log"
	"wiki_bit/library/convert/xstring"
	"wiki_bit/services"
)

const gateway = "http://www.tokenview.com:8088/eth/address"

type ethResponse struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
	Data data   `json:"data"`
}

type data struct {
	Type          string `json:"type"`
	Network       string `json:"network"`
	Hash          string `json:"hash"`
	Spend         string `json:"spend"`
	Receive       string `json:"receive"`
	Balance       string `json:"balance"`
	NormalTxCount int64  `json:"normalTxCount"`
	AddressType   string `json:"addressType"`
	Nonce         int64  `json:"nonce"`
	Txs           []txs  `json:"txs"`
}

type txs struct {
	Type          string `json:"type"`
	Network       string `json:"network"`
	BlockNo       int64  `json:"block_no"`
	Height        int64  `json:"height"`
	Index         int64  `json:"index"`
	Time          int64  `json:"time"`
	Txid          string `json:"txid"`
	Fee           string `json:"fee"`
	Confirmations int64  `json:"confirmations"`
	From          string `json:"from"`
	To            string `json:"to"`
	Nonce         int64  `json:"nonce"`
	GasPrice      int64  `json:"gasPrice"`
	GasLimit      int64  `json:"gasLimit"`
	Value         string `json:"value"`
	GasUsed       int64  `json:"gasUsed"`
}

type ip struct {
	Ip   string `json:"ip"`
	Port int64  `json:"port"`
}

func Screen(ethAddress string) (ethResponse, error) {

	num := services.GenerateRangeNum(0, 399)

	ip := read()[num]
	uri, err := url.Parse(ip)
	fmt.Printf("IP：%s\n", ip)
	// 检测ip是否有效
	strArr := strings.Split(ip, ":")
	ipstr := strings.TrimLeft(strArr[1], "/")
	val := tcpGather(ipstr, strArr[2])
	if val == false {
		return ethResponse{}, errors.New("5000")
	}

	client := http.Client{
		Transport: &http.Transport{
			// 设置代理
			Proxy: http.ProxyURL(uri),
		},
	}
	request, err := client.Get(fmt.Sprintf(gateway+"/%s", ethAddress))

	//request, err := services.Request(http.MethodGet, fmt.Sprintf(gateway+"/%s", ethAddress), nil)
	if err != nil {
		log.Logger().Error("eth Screen 请求 err：", zap.Error(err))
		return ethResponse{}, err
	}
	defer request.Body.Close()
	data, _ := ioutil.ReadAll(request.Body)
	var v ethResponse
	_ = json.Unmarshal(data, &v)
	if v.Code == 404 {
		log.Logger().Info("eth Screen 查询response：", zap.Any("response", v))
		return v, errors.New("100000")
	} else if v.Code != 404 && v.Code != 1 {
		return v, errors.New("200000")
	}
	return v, nil

	//content, _ := services.ResponseHandle(request)
	//var v ethResponse
	//_ = json.Unmarshal(content, &v)
	//if v.Code == 404 {
	//	log.Logger().Info("eth Screen 查询response：", zap.Any("response", v))
	//	return &v, errors.New("100000")
	//} else if v.Code != 404 && v.Code != 1 {
	//	return &v, errors.New("200000")
	//}
	//return &v, nil
}

// telnet 检测iP加端口
func tcpGather(ip, port string) bool {
	address := net.JoinHostPort(ip, port)
	fmt.Printf("开始检测IP：%s\n", address)
	// 3 秒超时
	conn, err := net.DialTimeout("tcp", address, 1*time.Second)
	if err != nil {
		fmt.Printf("IP检测ERR：%s\n", err.Error())
		return false
	} else {
		if conn != nil {
			_ = conn.Close()
			return true
		} else {
			return false
		}
	}
}

func read() [400]string {
	data, err := ioutil.ReadFile("D:/project/wikibit_airportal_api/test/json.json")
	if err != nil {
		_, _ = pp.Println("打开文件错误：", err.Error())
		os.Exit(-1)
	}

	var ips []ip
	var ippool [400]string
	_ = json.Unmarshal([]byte(data), &ips)
	for k, v := range ips {
		ippool[k] = "http://" + v.Ip + ":" + xstring.Int64ToString(v.Port)
	}
	return ippool
}
