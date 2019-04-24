package main

import (
    "flag"
    //"fmt"
    "opencast/config"
    "opencast/dao"
    "opencast/monitor"
    "opencast/service"
    "runtime"

    _ "github.com/go-sql-driver/mysql"
    log "github.com/thinkboy/log4go"
)

const (
    VERSION = "0.1"
)

func main() {
    flag.Parse()
    if err := config.InitConfig(); err != nil {
        panic(err)
    }
    runtime.GOMAXPROCS(config.Conf.MaxProc)
    log.LoadConfiguration(config.Conf.Log)
    defer log.Close()

    log.Info("opencast[%s] start", VERSION)
    log.Info("conf %#v", config.Conf)

    monitor.Init()
    dao.Init()
    service.Init()
    StartHttpServer()
}
