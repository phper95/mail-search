package metric

import (
	"gitee.com/phper95/pkg/prome"
	"github.com/prometheus/client_golang/prometheus"
	"mail-search/config"
)

var MailSearch = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:        "mail_search",
		Help:        "histogram for mail search",
		Buckets:     prome.DefaultBuckets,
		ConstLabels: prometheus.Labels{"machine": prome.GetHostName(), "app": config.AppName},
	},
	[]string{"cluster"},
)
