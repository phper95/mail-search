package v1

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"go.uber.org/zap"
	"mail-search/global"
	"mail-search/internal/pkg/errcode"
	"mail-search/internal/repo/es/mail_repo"
	"mail-search/internal/server/api/api_response"
	"mail-search/internal/service/mail_service"
	"mail-search/metric"
	"time"
)

type serchResponse struct {
	Total int64                  `json:"total"`
	Hits  []*mail_repo.MailIndex `json:"hits"`
}

func MailSearch(c *gin.Context) {
	t := time.Now()
	cluster := "a"
	//监控上报
	defer func() {
		obs, err := metric.MailSearch.GetMetricWithLabelValues(cluster)
		if err != nil {
			global.LOG.Error("metric.MailSearch error", zap.Error(err))
		} else {
			obs.Observe(float64(time.Since(t).Milliseconds()))
		}
	}()
	appG := api_response.Gin{C: c}
	keyword := c.Query("keyword")
	if len(keyword) == 0 {
		appG.ResponseErr(errcode.ErrCodes.ErrParams)
		return
	}
	mailService := mail_service.Mail{
		Keyword:  keyword,
		PageNum:  com.StrTo(c.Query("page_num")).MustInt(),
		PageSize: com.StrTo(c.Query("page_size")).MustInt(),
		UserID:   com.StrTo(c.Query("userid")).MustInt64(),
	}

	//上报搜索日志
	mailService.CreateTime = time.Now().Unix()
	defer func() {
		mailService.LogReport()
	}()

	//模拟多集群上报
	if mailService.UserID%2 == 0 {
		cluster = "b"
	}

	res, err := mailService.SearchMail()
	global.LOG.Warn("resp", zap.Any("", res))
	if err != nil {
		global.LOG.Error("search error", zap.Error(err), zap.Any("param", mailService))
		appG.ResponseErr(errcode.ErrCodes.ErrSearch)
		return
	}
	resp := serchResponse{
		Total: 0,
		Hits:  make([]*mail_repo.MailIndex, 0),
	}
	if res == nil {
		appG.ResponseOk(errcode.ErrCodes.ErrNo, resp)
		return
	}
	resp.Total = res.Hits.TotalHits.Value
	for _, hit := range res.Hits.Hits {
		index := &mail_repo.MailIndex{}
		err = json.Unmarshal(hit.Source, index)
		if err != nil {
			global.LOG.Error("Unmarshal error", zap.Error(err))
			continue
		}
		//index.Id, err = strconv.ParseInt(hit.Id, 10, 64)
		//if err != nil {
		//	global.LOG.Error("strconv.ParseInt error", zap.Error(err), zap.String("id", hit.Id))
		//	continue
		//}
		resp.Hits = append(resp.Hits, index)
	}
	global.LOG.Warn("resp", zap.Any("resp", resp))
	appG.ResponseOk(errcode.ErrCodes.ErrNo, resp)
}
