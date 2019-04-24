package control

import (
    "encoding/json"
    log "github.com/thinkboy/log4go"
    "io/ioutil"
    "net/http"
    "opencast/monitor"
    "strconv"
    "time"
)

type WrapHandler struct {
    Handle func(ctx *Context)
}

func (self *WrapHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    start := time.Now()
    ctx := NewContext(w, r)

    if r.Method == "POST" {
        var err error
        if ctx.Body, err = ioutil.ReadAll(r.Body); err != nil {
            log.Error("ioutil.ReadAll() failed (%s)", err)
            ctx.Rsp.SetStatus(RspServerErr)
            ctx.Rsp.SetMsg(err.Error())
            self.retWrite(ctx, start)
            return
        }
    }

    self.Handle(ctx)
    self.retWrite(ctx, start)
}

func (self *WrapHandler) retWrite(ctx *Context, start time.Time) {
    data, err := json.Marshal(ctx.Rsp)
    if err != nil {
        log.Error("json.Marshal(\"%v\") error(%v)", ctx.Rsp, err)
        return
    }
    dataStr := string(data)
    if _, err := ctx.Writer.Write([]byte(dataStr)); err != nil {
        log.Error("w.Write(\"%s\") error(%v)", dataStr, err)
    }
    elapsed := time.Since(start)
    monitor.HttpRequestCount.WithLabelValues(ctx.Request.URL.Path,
                            strconv.Itoa(ctx.Rsp.GetStatus()),
                            strconv.Itoa(ctx.Appid)).Inc()
    monitor.HttpRequestDuration.WithLabelValues(ctx.Request.URL.Path).Observe((float64)(elapsed / time.Microsecond))
    log.Info("reqUrl: %s, reqBody:%s, rsp:%s, ip:%s, time:%fs", ctx.Request.URL.String(),
        string(ctx.Body), dataStr, ctx.Request.RemoteAddr, elapsed.Seconds())
}
