package service

import (
    log "github.com/thinkboy/log4go"
)

func Init() {
    if err := initHttpClient(); err != nil {
        panic(err)
    }
    log.Info("service init")
    initCometBl()
    InitReportTask()
}
