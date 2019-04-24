package service

import (
    "encoding/json"
    "errors"
    "fmt"
    "opencast/config"

    log "github.com/thinkboy/log4go"
)

type TvInfoEx struct {
    Fe string `json:"fe"` //支持特性
}

type TvInfo struct {
    Status bool   `json:"status"`
    Appid  int    `json:"appid"`
    Info   string `json:"info"`
    Err    string `json:"err"`
    Safe   bool   //支持加密
}

func GetTvinfo(uid int64) (tvinfo TvInfo, err error) {
    url := fmt.Sprintf("http://%s/1/client/info?uid=%d", config.Conf.ImServerAddr, uid)
    code, body, getErr := HttpGet(url)
    if getErr != nil {
        err = getErr
        return
    } else if code != 200 {
        err = errors.New("req error")
        return
    } else {
        if err = json.Unmarshal(body, &tvinfo); err != nil {
            log.Error("err=%v", err)
            return
        }

        if tvinfo.Err != "" {
            err = errors.New(tvinfo.Err)
            log.Error("err msg: %s", tvinfo.Err)
            return
        }

        if tvinfo.Status == true && len(tvinfo.Info) != 0 {
            var tvInfoEx TvInfoEx
            if tmpErr := json.Unmarshal([]byte(tvinfo.Info), &tvInfoEx); tmpErr != nil {
                log.Error("tmpErr=%v", tmpErr)
            } else {
                if len(tvInfoEx.Fe) > 0 && tvInfoEx.Fe[0] == '1' {
                    tvinfo.Safe = true
                }
            }
        }
    }

    return
}
