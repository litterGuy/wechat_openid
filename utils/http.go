package utils

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func Get(url string, params map[string]interface{}, headers ...map[string]string) ([]byte, error) {
	var rawQuery string
	for key, value := range params {
		rawQuery += fmt.Sprintf("&%s=%v", key, value)
	}
	if rawQuery != "" {
		if strings.Contains(url, "?") {
			url += rawQuery
		} else {
			url += strings.Replace(rawQuery, "&", "?", 1)
		}
	}
	return callHttp(http.MethodGet, url, params, headers...)
}

func PostJson(url string, params map[string]interface{}, headers ...map[string]string) ([]byte, error) {
	header := make(map[string]string)
	if len(headers) > 0 {
		header = headers[0]
	} else {
		headers = make([]map[string]string, 1)
	}
	if _, ok := header["Content-Type"]; !ok {
		header["Content-Type"] = "application/json"
		headers[0] = header
	}
	return callHttp(http.MethodPost, url, params, headers...)
}

func PostForm(url string, params map[string]interface{}, headers ...map[string]string) ([]byte, error) {
	header := make(map[string]string)
	if len(headers) > 0 {
		header = headers[0]
	} else {
		headers = make([]map[string]string, 1)
	}
	if _, ok := header["Content-Type"]; !ok {
		header["Content-Type"] = "application/x-www-form-urlencoded"
		headers[0] = header
	}
	return callHttp(http.MethodPost, url, params, headers...)
}

//Http客户端
func callHttp(method string, url string, params map[string]interface{}, headers ...map[string]string) ([]byte, error) {
	var bodyreader io.Reader
	if len(headers) > 0 {
		header := headers[0]
		if contentType, ok := header["Content-Type"]; ok {
			if contentType == "application/x-www-form-urlencoded" {
				var r http.Request
				r.ParseForm()
				for key, val := range params {
					r.Form.Add(key, fmt.Sprintf("%v", val))
				}
				bodyreader = strings.NewReader(r.Form.Encode())
			} else {
				body, _ := json.Marshal(params)
				bodyreader = bytes.NewReader(body)
			}
		}
	} else {
		body, _ := json.Marshal(params)
		bodyreader = bytes.NewReader(body)
	}
	//组装请求信息

	//创建HttpRequest
	req, err := http.NewRequest(strings.ToUpper(method), url, bodyreader)
	if err != nil {
		return nil, err
	}
	//req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	//req.Header.Add("Content-Type", "application/json")
	if len(headers) > 0 {
		for k, v := range headers[0] {
			req.Header.Add(k, v)
		}
	}

	//创建HttpClient并发起请求
	client := http.Client{
		Transport: &http.Transport{
			DisableKeepAlives:   true, //true:不同HTTP请求之间TCP连接的重用将被阻止（http1.1默认为长连接，此处改为短连接）
			MaxIdleConnsPerHost: 512,  //控制每个主机下的最大闲置连接数目
		},
		Timeout: time.Second * 60, //Client请求的时间限制,该超时限制包括连接时间、重定向和读取response body时间;Timeout为零值表示不设置超时
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	//解析响应信息
	reply, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

//Https客户端
func callHttps(crtPath string, method string, url string, params map[string]interface{}, headers ...map[string]string) ([]byte, error) {
	//组装请求信息
	body, _ := json.Marshal(params)

	//创建HttpRequest
	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	//req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	//req.Header.Add("Content-Type", "application/json")
	if len(headers) > 0 {
		for k, v := range headers[0] {
			req.Header.Add(k, v)
		}
	}

	//获取客户端证书
	crts := x509.NewCertPool()
	if crtPath == "" {
		crtPath = "ca.crt"
	}
	crt, err := ioutil.ReadFile(crtPath)
	if err != nil {
		return nil, fmt.Errorf("获取客户端https证书出错，err: %v", err)
	}
	crts.AppendCertsFromPEM(crt)

	//创建HttpClient并发起请求
	client := http.Client{
		Transport: &http.Transport{
			DisableKeepAlives:   true,                       //true:不同HTTP请求之间TCP连接的重用将被阻止（http1.1默认为长连接，此处改为短连接）
			MaxIdleConnsPerHost: 512,                        //控制每个主机下的最大闲置连接数目
			TLSClientConfig:     &tls.Config{RootCAs: crts}, //添加证书
		},
		Timeout: time.Second * 60, //Client请求的时间限制,该超时限制包括连接时间、重定向和读取response body时间;Timeout为零值表示不设置超时
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	//解析响应信息
	reply, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return reply, nil
}
