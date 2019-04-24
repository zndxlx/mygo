package service

import (
    "encoding/json"
    log "github.com/thinkboy/log4go"
)

const CastConnectPrefix = "020005ff,"

type IMCastConnectReq struct {
    Sid   string `json:"sid"`         //会话id sessionid
    Suid  int64  `json:"suid,string"` //发送端uid
    Sname string `json:"sname"`       //发送端用户名称
    Sicon string `json:"sicon"`       //发送端用户icon
    Sdkv  string `json:"sdkv"`        //发送端SDK版本，例如3.0.0
    Appid string `json:"app_id"`      //发送端类型(不知道填什么)
    Mac   string `json:"mac"`         //发送端mac地址
    Sdid  string `json:"sdid"`        //发送端唯一ID
    Sc    int    `json:"sc,string"`   //发送端渠道,就是开发者平台上申请的appid
}

func PushCastConnect(uid int64, req *IMCastConnectReq) (err error) {
    data, err := json.Marshal(req)
    if err != nil {
        log.Error("json.Marshal(\"%v\") error(%v)", req, err)
        return
    }

    dataStr := CastConnectPrefix + string(data)

    err = push2Tv(uid, dataStr)

    return
}
