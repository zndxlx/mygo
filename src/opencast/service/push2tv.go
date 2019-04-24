package service

import (
    "encoding/json"
    "errors"
    "fmt"
    "opencast/config"

    //"strconv"
    //"time"

    log "github.com/thinkboy/log4go"
)

type PushRsp struct {
    Ret int `json:"ret"`
}

func push2Tv(uid int64, msg string) (err error) {
    url := fmt.Sprintf("http://%s/1/push?uid=%d", config.Conf.ImServerAddr, uid)
    code, body, err := HttpPost(url, msg)

    if err != nil {
        return
    } else if code != 200 {
        err = errors.New("req error")
        return
    } else {
        var rsp PushRsp
        if err = json.Unmarshal(body, &rsp); err != nil {
            log.Error("json decode failed, err=%v", err)
            return
        }

        if rsp.Ret != 1 {
            err = errors.New("im server error")
            log.Error("im server rsp error, rsp string %s", body)
            return
        }
        log.Info("im server  rsp string %s", body)
    }
    return
}
