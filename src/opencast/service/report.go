package service

import (
    "errors"
    // "fmt"
    log "github.com/thinkboy/log4go"
    // "strconv"
    //"time"
    "opencast/monitor"
)

//数据定义channel
var (
    reportChan chan string
)

const (
    chanSize       = 10240 * 4
    reportPoolSize = 100
)

//AddReportJob 增加一个上报任务, if chan full discard it
func AddReportJob(url string) (err error) {
    select {
    case reportChan <- url:
    default:
        err = errors.New("channel full")
        monitor.ReportCount.WithLabelValues("full").Inc()
    }
    return
}

func doReport() {
    for {
        reportURL := <-reportChan

        code, _, getErr := HttpGet(reportURL)
        if code != 200 || getErr != nil {
            log.Error("do report failed, code=%d, err=%v, url=%s", code, getErr.Error(), reportURL)
            monitor.ReportCount.WithLabelValues("err").Inc()
        } else {
            monitor.ReportCount.WithLabelValues("ok").Inc()
        }
    }
}

//InitReportTask 多个协程处理
func InitReportTask() {
    reportChan = make(chan string, chanSize)
    for i := 0; i < reportPoolSize; i++ {
        go doReport()
    }
    return
}
