package service

import (
    "crypto/md5"
    "errors"
    "fmt"
    "strconv"
    "time"
    log "github.com/thinkboy/log4go"
)

const (
    AuthExpireTime = 60 * 60 * 6 //6小时
)

var (
    ErrTokenTimeOut = errors.New("token time out")
    ErrTokenInvalid = errors.New("token Invalid")
)

func CreateToken(appid int, uid int64, time int64, key string) (token string) {
    md5str := fmt.Sprintf("%d%d%d%s", uid, appid, time, key)
    md5value := md5.Sum([]byte(md5str))
    token = fmt.Sprintf("%x%d", md5value[4:12], time)

    log.Debug("appid=%d, uid=%d, time=%d, token=%s", appid,
        uid, time, token)

    return
}

func CheckToken(appid int, uid int64, key string, token string) (err error) {
    if len(token) <= 16 {
        err = ErrTokenInvalid
        return
    }

    tokenData := []byte(token)

    tokenTime, err := strconv.ParseInt(string(tokenData[16:]), 10, 0)
    if err != nil {
        log.Error("%s pares to int failed ", string(tokenData[16:]))
        return
    }

    if time.Now().Unix()-tokenTime > AuthExpireTime {
        log.Info("token time out now %d, token_time %d ", time.Now().Unix(), tokenTime)
        err = ErrTokenTimeOut
        return
    }

    needToken := CreateToken(appid, uid, tokenTime, key)
    if token != needToken {
        log.Error("token error %s, need %s", token, needToken)
        err = ErrTokenInvalid
        return
    }

    return
}
