package control

import (
    "encoding/json"
    "errors"
    "fmt"
    "opencast/config"
    "opencast/dao"
    "opencast/service"
    // log "github.com/thinkboy/log4go"
    "time"
)

type castConnectReq struct {
    Suid  int64  `json:"suid,string"`  //发送请求的用户客户端uid
    Sid   string `json:"sid"`          //会话id,标识一个连接
    Tuid  int64  `json:"tuid,string"`  // TV端uid
    Appid int    `json:"appid,string"` //应用id
    Token string `json:"token"`        //初始化认证返回的token
    Sname string `json:"sname"`        //发送端名称
    Tid   int    //租户id
}

func getCastConnectPara(body []byte) (req castConnectReq, err error) {
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

    return
}

func reportCastConn(req *castConnectReq) {
    url := fmt.Sprintf("%s?tid=%d&u=%d&cu=%d&as=%s&sc=%d&ls=%d&v=1.0&a=6001&st=3&sn=4&sta=1&rsv=A100001&lt=0",
        config.Conf.ReportConn, req.Tid, req.Tuid, req.Suid, req.Token, req.Appid, time.Now().Unix()*1000)
    service.AddReportJob(url)
}

func CastConnect(ctx *Context) {
    req, err := getCastConnectPara(ctx.Body)
    ctx.Appid = req.Appid
    if err != nil {
        ctx.Rsp.SetStatus(RspParaErr)
        ctx.Rsp.SetMsg(err.Error())
        return
    }

    channel, err := dao.GetChannelInfo(req.Appid)
    if err != nil {
        ctx.Rsp.SetStatus(RspAuthErr)
        ctx.Rsp.SetMsg(err.Error())
        return
    }

    req.Tid = channel.Tid

    status, err := CheckToken(req.Appid, req.Suid, req.Token)
    if err != nil {
        ctx.Rsp.SetStatus(status)
        ctx.Rsp.SetMsg(err.Error())
        return
    }

    hasTv, err := dao.PhoneHasTv(req.Suid, req.Tuid)
    if err != nil || !hasTv {
        ctx.Rsp.SetStatus(RspTvNotAllow)
        return
    }

    tvinfo, err := service.GetTvinfo(req.Tuid)
    if err != nil || tvinfo.Status == false {
        ctx.Rsp.SetStatus(RspTvOff)
        return
    }

    reportCastConn(&req)

    err = service.PushCastConnect(req.Tuid, &service.IMCastConnectReq{
        Sid:   req.Sid,
        Sname: req.Sname,
        Suid:  req.Suid,
        Sdkv:  "openapi/v1",
        Sc:    req.Appid})
    if err != nil {
        ctx.Rsp.SetStatus(RspPushErr)
        return
    }

    return
}
