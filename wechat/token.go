package wechat

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"wechat_openid/utils"
)

func GetAccessToken(appid, appsecret string) (string, error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s", appid, appsecret)
	temp, err := utils.Get(url, nil)
	if err != nil {
		log.Println(err)
		return "", err
	}
	rt := make(map[string]interface{})
	err = json.Unmarshal(temp, &rt)
	if err != nil {
		log.Println(err)
		return "", err
	}
	if token, ok := rt["access_token"]; ok {
		return token.(string), nil
	} else {
		return "", errors.New("token获取错误")
	}
}
