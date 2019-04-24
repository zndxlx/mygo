package control

import (
    "errors"

    "encoding/json"
    //"opencast/dao"
    "opencast/service"

    //log "github.com/thinkboy/log4go"
)

type castControlReq struct {
    Suid   int64  `json:"suid,string"` //发送请求的用户客户端uid
    Tuid   int64  `json:"tuid,string"` // TV端uid
    Sid    string //会话id,标识一个连接
    Uri    string //单次投屏唯一标示 （2018/07/30新增）
    St     int    //控制命令 1:play 2:pause 3:stop 4:seekto 5:volumeto 6:volumeAdd 7:volumeReduce
    Seekto int    //快进的period进度，单位秒
    Vt     int    //声音调整百分比，比如vt=50,表示音量调整到50%
    Appid  int    `json:"appid,string"` //应用id
    Token  string //初始化认证返回的token
}

func getCastControlPara(body []byte) (req castControlReq, err error) {
    if err = json.Unmarshal(body, &req); err != nil {
        return
    }

    if req.Suid == 0 {
        err = errors.New("no suid")
        return
    }
    if req.Tuid == 0 {
        err = errors.New("no tuid")
        return
    }
    if req.Appid == 0 {
        err = errors.New("no appid")
        return
    }
    if req.Sid == "" {
        err = errors.New("no sid")
        return
    }
    if req.Token == "" {
        err = errors.New("no token")
        return
    }
    if req.Uri == "" {
        err = errors.New("no uri")
        return
    }

    return
}

func CastControl(ctx *Context) {
    req, err := getCastControlPara(ctx.Body)
    ctx.Appid = req.Appid
    if err != nil {
        ctx.Rsp.SetStatus(RspParaErr)
        ctx.Rsp.SetMsg(err.Error())
        return
    }

    status, err := CheckToken(req.Appid, req.Suid, req.Token)
    if err != nil {
        ctx.Rsp.SetStatus(status)
        ctx.Rsp.SetMsg(err.Error())
        return
    }

    tvinfo, err := service.GetTvinfo(req.Tuid)
    if err != nil || tvinfo.Status == false {
        ctx.Rsp.SetStatus(RspTvOff)
        return
    }

    err = service.PushCastControl(req.Tuid, &service.IMCastControlReq{
        Sid:    req.Sid,
        St:     req.St,
        Seekto: req.Seekto,
        Vt:     req.Vt,
        Uri:    req.Uri})
    if err != nil {
        ctx.Rsp.SetStatus(RspPushErr)
        return
    }

    return
}
