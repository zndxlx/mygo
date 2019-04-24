package service

import (
    //"strings"
    "common/util"
    "errors"
    // "fmt"
    log "github.com/thinkboy/log4go"
    "opencast/dao"
    "sync"
    "time"
)

type CometBalancer struct {
    mtx   *sync.Mutex
    round *util.SmoothRoundRobinAlg
}

const (
    syncInterval = time.Duration(5) * time.Second
    maxCpuLoad   = 95
)

var cometBl *CometBalancer

func GetComet() (addr string, err error) {
    cometBl.mtx.Lock()

    addrInterface := cometBl.round.Next()
    if addrInterface == nil {
        log.Error("cometbl failed, no comet")
        err = errors.New("cometbl failed, no comet")
    } else {
        var ok bool
        if addr, ok = addrInterface.(string); !ok {
            log.Error("cometbl failed, type error addrInterface = %v", addrInterface)
            err = errors.New("cometbl failed, type error")
        }
    }

    cometBl.mtx.Unlock()
    return
}

func syncCometBl() {
    comets := dao.GetComets()

    if len(comets) == 0 {
        log.Error("have no comets")
        return
    }
    cometBl.mtx.Lock()
    cometBl.round = nil
    cometBl.round = &util.SmoothRoundRobinAlg{}
    for _, comet := range comets {
        weight := 1
        cpuUsed := int(comet.CPUUsed)
        if cpuUsed < maxCpuLoad {
            weight = maxCpuLoad - cpuUsed
        }
        log.Debug("comet.ExtranetIp %v, weight %d", comet.ExtranetIp, weight)
        cometBl.round.Add(comet.ExtranetIp+":8080", weight)
    }
    cometBl.mtx.Unlock()
}

func syncCometBlTask() {
    c := time.Tick(syncInterval)
    for now := range c {
        now = now
        go syncCometBl()
    }
}

func initCometBl() {
    cometBl = &CometBalancer{
        mtx:   &sync.Mutex{},
        round: &util.SmoothRoundRobinAlg{},
    }
    syncCometBl()
    go syncCometBlTask()
}
