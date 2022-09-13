package metric

import (
	"fmt"
	"gitee.com/phper95/pkg/httpclient"
	"gitee.com/phper95/pkg/sign"
	"mail-search/config"
	"net/http"
	"net/url"
	"strconv"
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

func TestMetric(t *testing.T) {
	params := url.Values{}
	params.Add("userid", "1")
	params.Add("keyword", "测试")
	params.Add("page_num", "1")
	params.Add("page_size", "10")

	for i := 0; i < 10000; i++ {
		params.Add("userid", strconv.Itoa(i))
		authorization, date, err := sign.New(ak, sk, ttl).Generate(MailSearchUri, http.MethodGet, params)
		if err != nil {
			fmt.Println(err)
			return
		}
		headerAuth := httpclient.WithHeader(config.HeaderAuthField, authorization)
		headerAuthDate := httpclient.WithHeader(config.HeaderAuthDateField, date)
		c, r, e := httpclient.Get(MailSearchHost+MailSearchUri, params, headerAuth, headerAuthDate)
		fmt.Println(c, string(r), e)
	}
}
