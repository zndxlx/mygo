package service

import (
    "encoding/json"
    log "github.com/thinkboy/log4go"
)

const CastControlPrefix = "020004ff,"

type IMCastControlReq struct {
    Sid    string `json:"sid"`    //会话id sessionid
    St     int    `json:"st"`     //play 2:pause 3:stop 4:seekto 5:volumeto 6:volumeAdd 7:volumeReduce
    Seekto int    `json:"seekto"` //快进的period进度，单位秒
    Vt     int    `json:"vt"`     //50 表示volumeto声音调整到50%
    Uri    string `json:"uri"`    //单次投屏唯一标示
}

func PushCastControl(uid int64, req *IMCastControlReq) (err error) {
    data, err := json.Marshal(req)
    if err != nil {
        log.Error("json.Marshal(\"%v\") error(%v)", req, err)
        return
    }

    dataStr := CastControlPrefix + string(data)

    //发送消息
    err = push2Tv(uid, dataStr)
    return
}
