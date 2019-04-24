package main

import (
    "github.com/astaxie/beego"
    "github.com/astaxie/beego/orm"
    _ "github.com/go-sql-driver/mysql"
    _ "myearAdmin/routers"
    //"myearAdmin/util"
    "fmt"
)

func main() {
    //设置beego内置的日志打印
    beego.SetLogger("file", `{"filename":"./logs/beego.log","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":10,"color":true}`)
    beego.BeeLogger.DelLogger("console")
    beego.SetLogFuncCall(true)
    beego.BConfig.Log.AccessLogs = true

    //连接数据库
    maxIdleConn := 30
    maxOpenConn := 30
    dbUser := "lebocloud"
    dbPass := ""
    dbAddr := "192.168.8.237:3306"
    dbName := "test"
    dbLink := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", dbUser, dbPass, dbAddr, dbName) + "&loc=Asia%2FShanghai"

    orm.RegisterDriver("mysql", orm.DRMySQL)
    orm.RegisterDataBase("default", "mysql", dbLink, maxIdleConn, maxOpenConn)
    orm.Debug = true

    beego.BConfig.Listen.Graceful = true

    beego.Run()
}
