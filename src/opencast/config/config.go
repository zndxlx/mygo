package config

import (
    "flag"
    "runtime"
    "time"

    "github.com/Terry-Mao/goconf"
    log "github.com/thinkboy/log4go"
)

var (
    gconf    *goconf.Config
    Conf     *Config
    confFile string
)

func init() {
    flag.StringVar(&confFile, "c", "./opencast.conf", " set opencast config file path")
}

type Config struct {
    // base section
    PidFile     string `goconf:"base:pidfile"`
    Dir         string `goconf:"base:dir"`
    Log         string `goconf:"base:log"`
    MaxProc     int    `goconf:"base:maxproc"`
    PprofAddr   string `goconf:"base:pprof.addr:,"`
    MonitorAddr string `goconf:"base:monitor.addr:,"`

    HTTPAddrs        []string      `goconf:"base:http.addrs:,"`
    HTTPReadTimeout  time.Duration `goconf:"base:http.read.timeout:time"`
    HTTPWriteTimeout time.Duration `goconf:"base:http.write.timeout:time"`

    // yun db
    YunDbAddr     string `goconf:"yundb:addr"`
    YunDbUser     string `goconf:"yundb:user"`
    YunDbPass     string `goconf:"yundb:pass"`
    YunDbDataBase string `goconf:"yundb:database"`

    // redis
    RedisAddr string `goconf:"redis:addr"`
    RedisPwd  string `goconf:"redis:pwd"`
    RedisDB   int    `goconf:"redis:db"`
    RedisPool int    `goconf:"redis:pool"`

    // comet status redis
    CRedisAddr     string `goconf:"credis:addr"`
    CRedisPwd      string `goconf:"credis:pwd"`
    CRedisDB       int    `goconf:"credis:db"`
    CRedisPool     int    `goconf:"credis:pool"`
    CRedisCometKey string `goconf:"credis:comet.key"`
    //im server
    ImServerAddr string `goconf:"imserver:addr"`

    //report url
    ReportAuth     string `goconf:"report:auth"`
    ReportRelation string `goconf:"report:relation"`
    ReportConn     string `goconf:"report:conn"`
    ReportPush     string `goconf:"report:push"`
}

func NewConfig() *Config {
    return &Config{
        // base section
        PidFile:        "/tmp/opencast.pid",
        Dir:            "./",
        Log:            "./opencast-log.xml",
        MaxProc:        runtime.NumCPU(),
        PprofAddr:      "localhost:6971",
        MonitorAddr:    ":11001",
        CRedisCometKey: "open-comet-",
    }
}

// InitConfig init the global config.
func InitConfig() (err error) {
    Conf = NewConfig()
    gconf = goconf.New()
    if err = gconf.Parse(confFile); err != nil {
        return err
    }
    if err := gconf.Unmarshal(Conf); err != nil {
        return err
    }
    return nil
}

func ReloadConfig() (*Config, error) {
    conf := NewConfig()
    ngconf, err := gconf.Reload()
    if err != nil {
        return nil, err
    }
    if err := ngconf.Unmarshal(conf); err != nil {
        return nil, err
    }
    gconf = ngconf
    return conf, nil
}

func Reload() {
    newConf, err := ReloadConfig()
    if err != nil {
        log.Error("ReloadConfig() error(%v)", err)
        return
    }
    Conf = newConf
}
