package control

import (
    "net/http"

    //log "github.com/thinkboy/log4go"
)

type Context struct {
    Writer  http.ResponseWriter // 响应
    Request *http.Request       // 请求
    Rsp     rsp                 //响应
    Body    []byte              //请求包体
    Appid   int                 //方便监控时候区别应用
}

func NewContext(w http.ResponseWriter, req *http.Request) *Context {
    return &Context{
        Writer:  w,
        Request: req,
        Rsp:     NewRsp(200),
        Body:    []byte{},
        Appid:   0,
    }
}
