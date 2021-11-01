package openapi

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

//签名
//将字符串进行md5加密，并转换成32位base64值，并转换成大写
func Sign(query string) string {
	return strings.ToUpper(xmd5(query))
}

//加密签名
func CryptoSign(param, encryption string) string {
	switch encryption {
	case "HMAC-SHA256":
		return xhmac(param)
	case "MD5":
		return xmd5(param)
	}
	return ""
}

//md5加密
func xmd5(param string) string {
	h := md5.New()
	h.Write([]byte(param))
	b := h.Sum(nil)

	return hex.EncodeToString(b)
}

//hmac加密
func xhmac(param string) string {
	hash := hmac.New(sha256.New, []byte(param))
	hash.Write([]byte(param))

	return hex.EncodeToString(hash.Sum(nil))
}

//json转map
func Json2UrlValues(data interface{}) url.Values {
	t := reflect.TypeOf(data)
	v := reflect.ValueOf(data)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}

	values := url.Values{}
	for i := 0; i < v.NumField(); i++ {
		switch v.Field(i).Kind() {
		case reflect.String:
			if vv := v.Field(i).String(); vv != "" {
				if t.Field(i).Tag.Get("json") != "sign" {
					values.Set(t.Field(i).Tag.Get("json"), vv)
				}
			}
		}
	}
	return values
}

//map to urlvalues
func Map2UrlValues(data map[string]string) url.Values {
	values := url.Values{}
	for k, v := range data {
		values.Set(k, v)
	}
	return values
}

//json to map
func Json2Map(data interface{}) map[string]string {
	t := reflect.TypeOf(data)
	v := reflect.ValueOf(data)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}

	m := make(map[string]string, 0)
	for i := 0; i < v.NumField(); i++ {
		switch v.Field(i).Kind() {
		case reflect.String:
			if vv := v.Field(i).String(); vv != "" {
				if t.Field(i).Tag.Get("json") != "sign" {
					m[t.Field(i).Tag.Get("json")] = vv
				}
			}
		case reflect.Int:
			if vv := v.Field(i).Int(); vv >= 0 {
				m[t.Field(i).Tag.Get("json")] = strconv.Itoa(int(vv))
			}
		case reflect.Float64:
			if vv := v.Field(i).Float(); vv >= 0 {
				m[t.Field(i).Tag.Get("json")] = strconv.FormatFloat(vv, 'f', -1, 64)
			}
		}
	}
	return m
}
