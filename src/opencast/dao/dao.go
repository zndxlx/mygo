package dao

import (
    "database/sql"
    "fmt"
    "opencast/config"
    "time"

    "github.com/go-redis/redis"
    _ "github.com/go-sql-driver/mysql"
    log "github.com/thinkboy/log4go"
)

var (
    gRedisClient  *redis.Client
    gCRedisClient *redis.Client //comet status redis client
    gYunDb        *sql.DB
)

func initDb() {
    var err error

    yunDbName := fmt.Sprintf(`%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true`,
        config.Conf.YunDbUser, config.Conf.YunDbPass, config.Conf.YunDbAddr, config.Conf.YunDbDataBase)
    gYunDb, err = sql.Open("mysql", yunDbName)
    if err != nil {
        panic(err)
    }
    gYunDb.SetMaxIdleConns(5)
    gYunDb.SetMaxOpenConns(200)
}

func initRedis() {
    gRedisClient = redis.NewClient(&redis.Options{
        Addr:     config.Conf.RedisAddr,
        Password: config.Conf.RedisPwd,
        DB:       config.Conf.RedisDB,
        PoolSize: config.Conf.RedisPool,

        IdleTimeout: 1 * time.Minute,
    })

    if _, err := gRedisClient.Ping().Result(); err != nil {
        panic(err)
    }

    gCRedisClient = redis.NewClient(&redis.Options{
        Addr:     config.Conf.CRedisAddr,
        Password: config.Conf.CRedisPwd,
        DB:       config.Conf.CRedisDB,
        PoolSize: config.Conf.CRedisPool,

        IdleTimeout: 1 * time.Minute,
    })

    if _, err := gCRedisClient.Ping().Result(); err != nil {
        panic(err)
    }
}

func Init() {
    initDb()
    initRedis()
    initChannelCache()
    log.Info("dao init")
}
