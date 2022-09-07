package zjuapi

import (
	"net/http"
	"net/http/cookiejar"
)

type ZJUAPIClient struct {
	HttpClient *http.Client
}

func NewClient() *ZJUAPIClient {
	jar, _ := cookiejar.New(nil)
	return &ZJUAPIClient{
		HttpClient: &http.Client{Jar: jar},
	}
}
