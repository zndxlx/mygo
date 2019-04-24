package util

import (
    "github.com/astaxie/beego/logs"
    //"log"
)

var ULog *logs.BeeLogger

func init() {
    ULog = logs.NewLogger()
    ULog.SetLogger("file", `{"filename":"./logs/app.log","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":10,"color":true}`)
    ULog.EnableFuncCallDepth(true)
    ULog.SetLogFuncCallDepth(2)
}
