package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
	"wechat_openid/conf"
	"wechat_openid/utils"
	"wechat_openid/wechat"
)

func main() {
	//读取配置文件
	if err := conf.Init(); err != nil {
		log.Println("conf.Init() error(%v)", err)
		panic(err)
	}
	defer func() {
		if err := recover(); err != nil {
			log.Println("exception is %v", err)
		}
	}()

	//添加主要业务
	err := writeOpenId()
	if err != nil {
		log.Println(err)
	} else {
		log.Println("写入完成")
		fmt.Println("写入完成")
	}
	time.Sleep(5 * time.Second)
}

func init() {
	file := filepath.Join("", "wechat_openid.log")
	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		panic(err)
	}
	log.SetOutput(logFile) // 将文件设置为log输出的文件
	log.SetPrefix("[wechat_openid]")
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.LUTC)
	return
}

func writeOpenId() error {
	token, err := wechat.GetAccessToken(conf.Conf.AppId, conf.Conf.AppSecret)
	if err != nil {
		return err
	}
	flag := true
	nextopenid := ""
	result := ""
	for flag {
		userlist, err := wechat.GetUsersOpenIds(token, nextopenid)
		if err != nil {
			return err
		}
		if len(userlist.Data.Openid) == 0 {
			break
		}
		if len(userlist.NextOpenId) == 0 {
			flag = false
		} else {
			nextopenid = userlist.NextOpenId
		}
		result += strings.Join(userlist.Data.Openid, "\n")
	}
	err = utils.WriteFile("openid.txt", result)
	return err
}
