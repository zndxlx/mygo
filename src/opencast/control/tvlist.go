package control

import (
    "errors"
    "net/http"
    "opencast/dao"
    "opencast/service"
    "strconv"
    "strings"

    log "github.com/thinkboy/log4go"
)

type getTvlistReq struct {
    Uid   int64   //发送端的uid(初始化认证时候返回)
    Appid int     //应用id
    Token string  //初始化认证返回的token
    Tvs   []int64 //Tv端uid列表，请求中是逗号分割,比如 123，123，125
}

type TvStatus struct {
    Uid int64 `json:"uid"` //tv的uid
    //Name   string `json:"name"`   //设备名
    Online bool `json:"online"` //是否在线
}

type GetTvlistRsp struct {
    Tvs []TvStatus `json:"tvs"` //tv信息列表
}

func getTvlistPara(r *http.Request) (req getTvlistReq, err error) {
    uidStr := r.URL.Query().Get("uid")
    if uidStr == "" {
        err = errors.New("no uid")
        return
    }
    if req.Uid, err = strconv.ParseInt(uidStr, 10, 0); err != nil {
        err = errors.New("uid malform")
        return
    }

    appidStr := r.URL.Query().Get("appid")
    if appidStr == "" {
        err = errors.New("no appid")
        return
    }
    if req.Appid, err = strconv.Atoi(appidStr); err != nil {
        err = errors.New("appid malform")
        return
    }

    if req.Token = r.URL.Query().Get("token"); req.Token == "" {
        err = errors.New("no token")
        return
    }

    tvsStr := r.URL.Query().Get("tvs")
    if tvsStr == "" {
        err = errors.New("no tvs")
        return
    }
    tvs := strings.Split(tvsStr, ",")
    for _, tv := range tvs {
        if tvUid, tmpErr := strconv.ParseInt(tv, 10, 0); tmpErr != nil {
            err = tmpErr
            return
        } else {
            req.Tvs = append(req.Tvs, tvUid)
        }
    }
    return
}

func GetTvList(ctx *Context) {
    var rsp GetTvlistRsp
    req, err := getTvlistPara(ctx.Request)
    ctx.Appid = req.Appid
    if err != nil {
        ctx.Rsp.SetStatus(RspParaErr)
        ctx.Rsp.SetMsg(err.Error())
        return
    }

    status, err := CheckToken(req.Appid, req.Uid, req.Token)
    if err != nil {
        ctx.Rsp.SetStatus(status)
        ctx.Rsp.SetMsg(err.Error())
    }

    tvSet, err := dao.GetPhoneNearTv(req.Uid)
    if err != nil {
        ctx.Rsp.SetStatus(RspServerErr)
        ctx.Rsp.SetMsg(err.Error())
        return
    }

    for _, tv := range req.Tvs {
        if _, ok := tvSet[strconv.FormatInt(tv, 10)]; ok {
            if tvInfo, err := service.GetTvinfo(tv); err != nil {
                ctx.Rsp.SetStatus(RspServerErr)
                ctx.Rsp.SetMsg(err.Error())
                return
            } else {
                rsp.Tvs = append(rsp.Tvs, TvStatus{Uid: tv, Online: tvInfo.Status})
            }
        } else {
            log.Error("not allow to query %v", tv)
        }
    }
    ctx.Rsp.SetData(rsp)
    return
}
