package v1

import (
	"gitee.com/phper95/pkg/httpclient"
	"gitee.com/phper95/pkg/sign"
	"mail-search/config"
	"net/http"
	"net/url"
	"testing"
	"time"
)

const MailSearchHost = "http://127.0.0.1:9099"
const MailSearchUri = "/api/v1/mail-search"

var (
	ak  = "AK100523687952"
	sk  = "W1WTYvJpfeH1YpUjTpeFbEx^DnpQ&35L"
	ttl = time.Minute * 3
)

func TestProductSearch(t *testing.T) {
	params := url.Values{}
	params.Add("userid", "1")
	params.Add("email", "phper95@163.com")
	params.Add("keyword", "测试")
	params.Add("page_num", "1")
	params.Add("page_size", "10")
	authorization, date, err := sign.New(ak, sk, ttl).Generate(MailSearchUri, http.MethodGet, params)
	if err != nil {
		t.Error(err)
		return
	}
	headerAuth := httpclient.WithHeader(config.HeaderAuthField, authorization)
	headerAuthDate := httpclient.WithHeader(config.HeaderAuthDateField, date)
	c, r, e := httpclient.Get(MailSearchHost+MailSearchUri, params, headerAuth, headerAuthDate)
	t.Log(c, string(r), e)
}
