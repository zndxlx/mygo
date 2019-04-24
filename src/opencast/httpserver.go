package main

import (
    "github.com/facebookgo/grace/gracehttp"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    log "github.com/thinkboy/log4go"
    "net/http"
    "net/http/pprof"
    "opencast/config"
    "opencast/control"
)

func StartHttpServer() {
    // http listen
    var serverList []*http.Server
    for i := 0; i < len(config.Conf.HTTPAddrs); i++ {
        log.Info("start http listen:\"%s\"", config.Conf.HTTPAddrs[i])
        serverList = append(serverList, &http.Server{Addr: config.Conf.HTTPAddrs[i], Handler: serviceHandler(),
            ReadTimeout: config.Conf.HTTPReadTimeout, WriteTimeout: config.Conf.HTTPWriteTimeout})
    }

    serverList = append(serverList, &http.Server{Addr: config.Conf.PprofAddr, Handler: pprofHandler()})
    serverList = append(serverList, &http.Server{Addr: config.Conf.MonitorAddr, Handler: monitorHandler()})

    if error := gracehttp.Serve(serverList...); error != nil {
        panic(error)
    }
}

func serviceHandler() http.Handler {
    mux := http.NewServeMux()
    mux.HandleFunc("/ping", ProcessPing)
    mux.Handle("/v1/author/phoneauthor", &(control.WrapHandler{control.PhoneAuth}))
    mux.Handle("/v1/device/report", &(control.WrapHandler{control.DeviceReport}))
    mux.Handle("/v1/device/tvlist", &(control.WrapHandler{control.GetTvList}))
    mux.Handle("/v1/cast/connect", &(control.WrapHandler{control.CastConnect}))
    mux.Handle("/v1/cast/pushurl", &(control.WrapHandler{control.CastPush}))
    mux.Handle("/v1/cast/playcontrol", &(control.WrapHandler{control.CastControl}))
    return mux
}

func pprofHandler() http.Handler {
    mux := http.NewServeMux()
    mux.HandleFunc("/debug/pprof/", pprof.Index)
    mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
    mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
    mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
    return mux
}

func monitorHandler() http.Handler {
    mux := http.NewServeMux()
    mux.Handle("/metrics", promhttp.Handler())
    return mux
}

func ProcessPing(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("pong"))
    return
}
