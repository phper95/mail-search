package mail_service

import (
	"context"
	"fmt"
	"gitee.com/phper95/pkg/es"
	"gitee.com/phper95/pkg/strutil"
	"gitee.com/phper95/pkg/timeutil"
	"github.com/olivere/elastic/v7"
	"go.uber.org/zap"
	"mail-search/global"
	"sync"
	"time"
)

var (
	LogTableCreated sync.Map
)

type Mail struct {
	UserID     int64  `json:"userid" bson:"userid"`
	Keyword    string `json:"keyword" bson:"keyword"`
	PageNum    int    `json:"page_num" bson:"page_num"`
	PageSize   int    `json:"page_size" bson:"page_size"`
	CreateTime int64  `json:"create_time" bson:"create_time"`
}

func (m *Mail) SearchMail() (result *elastic.SearchResult, err error) {
	query := elastic.NewBoolQuery()
	from := m.PageNum * 20

	query.MinimumNumberShouldMatch(1)

	subjectMatchPhreaseQuery := elastic.NewMatchPhraseQuery("subject", m.Keyword).Boost(2).QueryName("subjectMatchPhraseQuery")
	contentMatchQuery := elastic.NewMatchPhraseQuery("content", m.Keyword).Boost(1).QueryName("contentMatchQuery")

	query.Must(elastic.NewTermQuery("to", m.UserID))

	query.Should(subjectMatchPhreaseQuery, contentMatchQuery)

	orders := make([]map[string]bool, 0)

	//默认按照相关度算分来排序
	orders = append(orders, map[string]bool{"_score": false})

	return global.ES.Query(context.Background(), global.MailIndexName,
		[]string{strutil.Int64ToString(m.UserID)}, query, from, m.PageSize, es.WithEnableDSL(true),
		es.WithPreference(strutil.Int64ToString(m.UserID)),
		es.WithFetchSource(true), es.WithOrders(orders))
}

func (m *Mail) LogReport() {
	if global.Mongo == nil {
		global.LOG.Error(" global.Mongo is nil", zap.Any("param", m))
		return
	}
	tablename := fmt.Sprintf(global.MailSearchLogCollectionNamePrefix, timeutil.YMDLayoutInt64(time.Now()))
	//本地缓存，避免每次写入都要创建索引
	if _, ok := LogTableCreated.Load(tablename); !ok {
		err := global.Mongo.CreateMultiIndex(global.SearchLogDbName, tablename, []string{"userid", "create_time"}, false)
		if err != nil {
			global.LOG.Error(" Mongo CreateMultiIndex error", zap.Error(err), zap.Any("param", m))
		}
		LogTableCreated.Store(tablename, true)
	}
	err := global.Mongo.InsertMany(global.SearchLogDbName, tablename, m)
	if err != nil {
		global.LOG.Error(" Mongo InsertMany error", zap.Error(err), zap.Any("param", m))
	}
}
