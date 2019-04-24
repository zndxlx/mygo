package control

import (
    "encoding/json"
    "errors"
    "fmt"
    // log "github.com/thinkboy/log4go"
    "net/url"
    "opencast/config"
    "opencast/dao"
    "opencast/service"
    "strconv"
    "strings"
)

type deviceReportReq struct {
    Uid    int64   `json:"uid,string"`   //发送端的uid(初始化认证时候返回)
    Appid  int     `json:"appid,string"` //应用id
    Token  string  `json:"token"`        //初始化认证返回的token
    Tvs    string  `json:"tvs"`          //Tv端uid列表，请求中是逗号分割,比如 123，123，125
    TvList []int64 //解析后的Tv端uid列表
    Tid    int     //租户id
}

func getDeviceReportPara(body []byte) (req deviceReportReq, err error) {
    if err = json.Unmarshal(body, &req); err != nil {
        return
    }

    if req.Uid == 0 {
        err = errors.New("no uid")
        return
    }

    if req.Appid == 0 {
        err = errors.New("no appid")
        return
    }

    if req.Token == "" {
        err = errors.New("no token")
        return
    }

    if req.Tvs == "" {
        err = errors.New("no tvs")
        return
    }

    tvs := strings.Split(req.Tvs, ",")
    for _, tv := range tvs {
        if tvUid, tmpErr := strconv.ParseInt(tv, 10, 0); tmpErr != nil {
            err = tmpErr
            return
        } else {
            req.TvList = append(req.TvList, tvUid)
        }
    }

    return
}

func reportDeviceReportReq(req *deviceReportReq) {
    url := fmt.Sprintf("%s?tid=%d&cu=%d&cut=6&appid=%d&ulist=%s", config.Conf.ReportRelation,
        req.Tid, req.Uid, req.Appid, url.QueryEscape(req.Tvs))
    service.AddReportJob(url)
}

func DeviceReport(ctx *Context) {
    req, err := getDeviceReportPara(ctx.Body)
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

    status, err := CheckToken(req.Appid, req.Uid, req.Token)
    if err != nil {
        ctx.Rsp.SetStatus(status)
        ctx.Rsp.SetMsg(err.Error())
        return
    }

    dao.AddPhoneNearTv(req.Uid, req.TvList)

    reportDeviceReportReq(&req)
    return
}
