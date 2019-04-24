package monitor

import (
    "github.com/prometheus/client_golang/prometheus"
)

var (
    HttpRequestCount = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "opencast_http_request_count",
            Help: "opencast http request count",
        },
        []string{"endpoint", "status", "appid"},
    )

    HttpRequestDuration = prometheus.NewSummaryVec(
        prometheus.SummaryOpts{
            Name: "opencast_http_request_duration",
            Help: "opencast http request duration",
        },
        []string{"endpoint"},
    )

    ReportCount = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "opencast_report_count",
            Help: "opencast report count",
        },
        []string{"status"},
    )
)

func Init() {
    prometheus.MustRegister(HttpRequestCount)
    prometheus.MustRegister(HttpRequestDuration)
    prometheus.MustRegister(ReportCount)
}
