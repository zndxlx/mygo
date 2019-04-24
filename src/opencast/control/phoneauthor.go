package control

import (
    "common/util"
    "crypto/md5"
    "encoding/binary"
    "errors"
    "fmt"
    log "github.com/thinkboy/log4go"
    "net/http"
    "net/url"
    "opencast/config"
    "opencast/dao"
    "opencast/service"
    "strconv"
    "time"
)

type phoneAuthReq struct {
    Appid     int
    Timestamp int64
    AndroidId string
    IosId     string
    Mac       string
    Sign      string
    Tid       int
}

type PhoneAuthRsp struct {
    Token      string `json:"token"`
    Uid        string `json:"uid"`
    ExpireTime int    `json:"expire_time"`
    Tcpserver  string `json:"tcpserver"`
}

func getPhoneAuthPara(r *http.Request) (req phoneAuthReq, err error) {
    appidStr := r.URL.Query().Get("appid")
    if appidStr == "" {
        err = errors.New("no appid")
        return
    }
    if req.Appid, err = strconv.Atoi(appidStr); err != nil {
        return
    }

    req.AndroidId = r.URL.Query().Get("android_id")
    req.IosId = r.URL.Query().Get("ios_id")
    req.Mac = r.URL.Query().Get("mac")
    if req.AndroidId == "" && req.Mac == "" && req.IosId == "" {
        err = errors.New("no Imei|Mac|Idfa")
        return
    }

    if req.Sign = r.URL.Query().Get("sign"); req.Sign == "" {
        err = errors.New("no Sign")
        return
    }

    TimeStr := r.URL.Query().Get("timestamp")
    if TimeStr == "" {
        err = errors.New("no timestamp")
        return
    }
    if req.Timestamp, err = strconv.ParseInt(TimeStr, 10, 0); err != nil {
        err = errors.New("timestamp malformed")
        return
    }

    if util.Math.AbsInt64(time.Now().Unix()-req.Timestamp) > 60*60 {
        err = errors.New("timestamp timeout")
        return
    }

    return
}

//reportPhoneAuth信息上报给大数据
func reportPhoneAuth(req *phoneAuthReq, rsp *PhoneAuthRsp) {
    url := fmt.Sprintf("%s?tid=%d&cu=%s&as=%s&sc=%d&m=%s&im=%s&id=%s&rsv=A100001&v=1.0&a=6001&st=5&sn=1",
        config.Conf.ReportAuth, req.Tid, rsp.Uid, rsp.Token, req.Appid, url.QueryEscape(req.Mac),
        url.QueryEscape(req.AndroidId), url.QueryEscape(req.IosId))

    service.AddReportJob(url)
}

func PhoneAuth(ctx *Context) {
    var rsp PhoneAuthRsp
    req, err := getPhoneAuthPara(ctx.Request)
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

    md5string := fmt.Sprintf("%d%d%s%s%s%s", req.Appid, req.Timestamp,
        channel.Md5Key, req.AndroidId, req.IosId, req.Mac)

    needSign := fmt.Sprintf("%x", md5.Sum([]byte(md5string)))
    if needSign != req.Sign {
        log.Error("sign error %s, need %s, md5string %s", req.Sign,
            needSign, md5string)
        ctx.Rsp.SetStatus(RspAuthErr)
        ctx.Rsp.SetMsg("sign invalid")
        return
    }

    rsp.Tcpserver, err = service.GetComet()
    if err != nil {
        ctx.Rsp.SetStatus(RspServerErr)
        ctx.Rsp.SetMsg(err.Error())
        return
    }

    uid, err := calUid(req.Appid, req.AndroidId, req.IosId, req.Mac)
    if err != nil {
        ctx.Rsp.SetStatus(RspAuthErr)
        ctx.Rsp.SetMsg(err.Error())
        return
    }
    rsp.Uid = strconv.FormatInt(uid, 10)

    now := time.Now().Unix()
    rsp.Token = service.CreateToken(req.Appid, uid, now, channel.Md5Key)
    rsp.ExpireTime = service.AuthExpireTime

    ctx.Rsp.SetData(rsp)
    log.Info("req %+v, rsp %+v", req, rsp)
    reportPhoneAuth(&req, &rsp)
    return
}

func calUid(appid int, androidId string, iodId string, mac string) (uid int64, err error) {
    if androidId == "" && iodId == "" && mac == "" {
        err = errors.New("para invalid")
        return
    }

    var md5str string

    if androidId != "" {
        md5str = fmt.Sprintf("%d%s", appid, androidId)
    } else if iodId != "" {
        md5str = fmt.Sprintf("%d%s", appid, iodId)
    } else {
        md5str = fmt.Sprintf("%d%s", appid, mac)
    }
    md5value := md5.Sum([]byte(md5str))
    token := fmt.Sprintf("%x", md5value[4:12])
    uid = int64(binary.BigEndian.Uint64(md5value[4:12]))

    log.Debug("appid=%d, androidId=%s, iodId=%s, mac=%s, token=%s, uid=%d", appid,
        androidId, iodId, mac, token, uid)

    return
}
