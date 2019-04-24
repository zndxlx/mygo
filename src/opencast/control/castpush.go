package control

import (
    "errors"
    //"net/http"
    "encoding/json"
    "fmt"
    log "github.com/thinkboy/log4go"
    "net/url"
    "opencast/config"
    "opencast/dao"
    "opencast/service"
)

type castPushReq struct {
    Suid  int64  `json:"suid,string"`  //发送请求的用户客户端uid
    Sid   string `json:"sid"`          //会话id,标识一个连接
    Tuid  int64  `json:"tuid,string"`  // TV端uid
    Appid int    `json:"appid,string"` //应用id
    Pos   int    `json:"pos"`          //启播位置
    Url   string `json:"url"`          //播放地址
    Uri   string `json:"uri"`          //单次投屏唯一标示
    Token string `json:"token"`        //初始化认证返回的token
    Surl  bool   `json:"surl"`         //播放地址是否进行了加密
    Tid   int    //租户id
}

func getCastPushPara(body []byte) (req castPushReq, err error) {
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

func reportCastPush(req *castPushReq) {
    url := fmt.Sprintf("%s?sta=1&rsv=A100001&v=1.0&a=6001&st=1&p=4&mt=102&dt=100&bid=none&tid=%d&u=%d&cu=%d&as=%s&sc=%d&s=%s&uri=%s",
        config.Conf.ReportPush, req.Tid, req.Tuid, req.Suid, req.Token, req.Appid,
        url.QueryEscape(req.Sid), url.QueryEscape(req.Uri))

    service.AddReportJob(url + "&sn=1")
    service.AddReportJob(url + "&sn=2")
}

func CastPush(ctx *Context) {
    req, err := getCastPushPara(ctx.Body)
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

    tvinfo, err := service.GetTvinfo(req.Tuid)
    if err != nil || tvinfo.Status == false {
        ctx.Rsp.SetStatus(RspTvOff)
        return
    }

    if req.Surl == true {
        orgUrl, err := service.AesCBCDecrypt(req.Url, channel.Md5Key)
        if err != nil {
            ctx.Rsp.SetStatus(RspServerErr)
            ctx.Rsp.SetMsg(err.Error())
            return
        }
        req.Url = string(orgUrl)
    }

    var safe bool = false
    var key string = ""
    if req.Surl && tvinfo.Safe {
        if tvChannel, tmpErr := dao.GetChannelInfo(tvinfo.Appid); tmpErr != nil {
            log.Error("get Appid channel info failed appid=%d,err=%s",
                tvinfo.Appid, tmpErr.Error())
        } else {
            safe = true
            key = tvChannel.Md5Key
        }
    }

    reportCastPush(&req)

    err = service.PushCastPush(req.Tuid,
        &service.IMCastPushReq{
            Url:     req.Url,
            Suid:    req.Suid,
            Uri:     req.Uri,
            Timeout: "5",
            Sdkv:    "openapi/v1",
            Appid:   "",
            Sid:     req.Sid,
            Pos:     req.Pos,
            Mt:      102},
        &service.IMCastPushSafePara{
            Safe: safe,
            Key:  key,
        })
    if err != nil {
        ctx.Rsp.SetStatus(RspPushErr)
        return
    }

    return
}
