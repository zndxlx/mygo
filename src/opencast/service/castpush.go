package service

import (
    "encoding/json"

    log "github.com/thinkboy/log4go"
)

const CastPushPrefix = "020001ff,"

type IMCastPushSafePara struct {
    Safe bool   //是否加密
    Key  string //加密密钥
}

type IMCastPushReq struct {
    Url     string `json:"url"`         //播放地址
    Suid    int64  `json:"suid,string"` //发送端uid
    Uri     string `json:"uri"`         //单次投屏唯一标示
    Timeout string `json:"timeout"`     //处理超时时间
    Sdkv    string `json:"sdkv"`        //发送端SDK版本，例如3.0.0
    Appid   string `json:"app_id"`      //发送端类型
    Sid     string `json:"sid" `        //会话id sessionid
    Pos     int    `json:"pos" `        //起播位置
    Mt      int    `json:"mt" `         // "mt":"url媒体类型，101、音频；102、视频；103、图片；104、幻灯片"（2018/07/30 新增，满足 公网投屏协议增加支持投在线音乐和在线图片 的需求）
}

type IMCastPushSafeReq struct {
    Sbody string `json:"sbody"` //加密信息
}

func PushCastPush(uid int64, req *IMCastPushReq, safe *IMCastPushSafePara) (err error) {
    data, err := json.Marshal(req)
    if err != nil {
        log.Error("json.Marshal(\"%v\") error(%v)", req, err)
        return
    }
    if safe.Safe {
        var safeReq IMCastPushSafeReq
        safeReq.Sbody, err = AesCBCEncrypt(data, safe.Key)
        if err != nil {
            return
        }
        safeData, tmpErr := json.Marshal(safeReq)
        if err != nil {
            log.Error("json.Marshal(\"%v\") error(%v)", safeReq, err)
            err = tmpErr
            return
        }
        err = push2Tv(uid, CastPushPrefix+string(safeData))
    } else {
        err = push2Tv(uid, CastPushPrefix+string(data))
    }

    return
}
