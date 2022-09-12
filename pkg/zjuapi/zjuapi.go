package zjuapi

import (
	"net/http"
	"net/http/cookiejar"
	"time"
)

type ZJUAPIClient struct {
	HttpClient *http.Client
}

func NewClient() *ZJUAPIClient {
	jar, _ := cookiejar.New(nil)
	return &ZJUAPIClient{
		// 这里是为了考试页面设置了差不多合适的 timeout
		// 校内网络环境下，大概 8~9s 左右加载完成
		HttpClient: &http.Client{Jar: jar, Timeout: time.Second * 12},
	}
}
