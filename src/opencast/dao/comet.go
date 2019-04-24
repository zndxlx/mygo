package dao

import (
    "encoding/json"
    log "github.com/thinkboy/log4go"
    "opencast/config"
    "strings"
)

type CometStatus struct {
    ServerId   int
    IP         string
    ExtranetIp string
    CPU        int
    CPUUsed    float32
    Memory     int
    MemoryUsed int
    ConCount   int
    Timestamp  int64
}

func GetComets() (comets []CometStatus) {
    rediskey := config.Conf.CRedisCometKey + "*"
    keys, err := gCRedisClient.Keys(rediskey).Result()
    if err != nil {
        log.Error("redis cmd Keys  err, err = %s \n", err.Error())
        return
    }

    for _, key := range keys {
        val, err := gCRedisClient.Get(key).Bytes()
        if err != nil {
            log.Error("redis cmd Get err, key=%s, err=%s", key, err.Error())
            continue
        }

        var cometStatus CometStatus
        if err = json.Unmarshal(val, &cometStatus); err != nil {
            log.Error("json.Unmarshal failed, val=%s err=%s ", string(val), err.Error())
            continue
        }
        cometStatus.ExtranetIp = strings.TrimSpace(cometStatus.ExtranetIp)
        comets = append(comets, cometStatus)
    }

    log.Info("comets=%+v", comets)
    return
}
