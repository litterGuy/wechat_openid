package wechat

import (
	"encoding/json"
	"fmt"
	"log"
	"wechat_openid/utils"
)

func GetUsersOpenIds(token string, nextopenid string) (*UserList, error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/user/get?access_token=%s", token)
	if len(nextopenid) > 0 {
		url = url + "&next_openid=" + nextopenid
	}
	temp, err := utils.Get(url, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	rt := new(UserList)
	err = json.Unmarshal(temp, rt)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return rt, nil
}

type UserList struct {
	Total      int      `json:"total"`
	Count      int      `json:"count"`
	NextOpenId string   `json:"next_openid"`
	Data       UserData `json:"data"`
}

type UserData struct {
	Openid []string `json:"openid"`
}
