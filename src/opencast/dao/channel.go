package dao

import (
    "fmt"
    "time"

    "github.com/bluele/gcache"
    log "github.com/thinkboy/log4go"
)

const (
    syncInterval       = time.Duration(60*5) * time.Second
    channelCacheMax    = 100000
    channelCacheExpire = time.Duration(10*24*3600) * time.Second
)

type ChannelInfo struct {
    Id     int
    Tid    int
    Secret string
    Md5Key string
}

var channelCache gcache.Cache

func GetChannelInfo(id int) (info ChannelInfo, err error) {
    var value interface{}
    value, err = channelCache.Get(id)
    if err != nil {
        log.Error("get %d failed, error(%v)", id, err)
        return
    }

    var ok bool
    info, ok = value.(ChannelInfo)
    if !ok {
        log.Error("value type error")
        return
    }

    log.Debug("ChannelInfo=%+v", info)

    return
}

func getAllChannelInfo() (channels []ChannelInfo, err error) {
    rows, err := gYunDb.Query(`SELECT id,tid,secret,md5_key FROM tb_channel`)
    if err != nil {
        log.Error("query failed, err = %s \n", err.Error())
        return
    }

    defer rows.Close()
    for rows.Next() {
        var info ChannelInfo
        err = rows.Scan(&info.Id, &info.Tid, &info.Secret, &info.Md5Key)
        if err != nil {
            log.Error("err=%v", err)
            return
        }
        //log.Debug("ChannelInfo=%+v", info)
        channels = append(channels, info)
    }
    err = rows.Err()
    if err != nil {
        log.Error("query failed, err = %s \n", err.Error())
        return
    }
    return
}

func getUpdateChannelInfo() (channels []ChannelInfo, err error) {
    rows, err := gYunDb.Query(`SELECT id,tid,package,secret,md5_key FROM tb_channel
        where UNIX_TIMESTAMP(last_time)>= UNIX_TIMESTAMP()-660
        and UNIX_TIMESTAMP(last_time) <= UNIX_TIMESTAMP()`)
    if err != nil {
        log.Error("query failed, err = %s \n", err.Error())
        return
    }

    defer rows.Close()
    for rows.Next() {
        var info ChannelInfo
        err = rows.Scan(&info.Id, &info.Tid, &info.Secret, &info.Md5Key)
        if err != nil {
            log.Error("err=%v", err)
            return
        }
        log.Debug("ChannelInfo=%+v", info)
        channels = append(channels, info)
    }
    err = rows.Err()
    if err != nil {
        log.Error("query failed, err = %s \n", err.Error())
        return
    }
    return
}

func syncUpdateChanel() {
    channels, err := getUpdateChannelInfo()
    if err != nil {
        return
    }

    for _, channel := range channels {
        channelCache.Set(channel.Id, channel)
    }
}

func syncAllChanel() {
    channels, err := getAllChannelInfo()
    if err != nil {
        return
    }

    for _, channel := range channels {
        channelCache.Set(channel.Id, channel)
    }
    return
}

func syncAllTask() {
    c := time.Tick(syncInterval)
    for now := range c {
        fmt.Printf("now %v \n", now)
        go syncAllChanel()
    }
}

func syncUpdateTask() {
    c := time.Tick(channelCacheExpire / 2)
    for now := range c {
        fmt.Printf("now %v \n", now)
        go syncUpdateChanel()
    }
}

func initChannelCache() {
    channelCache = gcache.New(channelCacheMax).LRU().
        Expiration(channelCacheExpire).
        Build()

    syncAllChanel()
    go syncUpdateTask()
    go syncAllTask()
}
