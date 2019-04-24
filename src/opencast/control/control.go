package control

import (
    "opencast/dao"
    "opencast/service"

    log "github.com/thinkboy/log4go"
)

const (
    RspSucess      = 200
    RspParaErr     = 401 //参数错误
    RspAuthErr     = 402 //认证失败
    RspTvOff       = 403 //tv不在线
    RspServerErr   = 405 //服务器错误
    RspPushErr     = 406 //消息推送失败
    RspTokenExpire = 410 //token过期
    RspTvNotAllow  = 411 //没有推送该tv的权限
)

func CheckToken(appid int, uid int64, token string) (status int, err error) {
    status = RspSucess
    channel, err := dao.GetChannelInfo(appid)
    if err != nil {
        status = RspAuthErr //认证失败
        return
    }

    err = service.CheckToken(appid, uid, channel.Md5Key, token)
    if err != nil {
        if err == service.ErrTokenTimeOut {
            status = RspTokenExpire //token过期
        } else {
            status = RspAuthErr //认证失败
        }
    }
    return
}

func Init() {
    log.Info("control init")
    return
}
