package dao

import (
    //"fmt"
    "strconv"

    log "github.com/thinkboy/log4go"
)

func AddPhoneNearTv(uid int64, nearTvs []int64) (err error) {
    rediskey := "neartv:" + strconv.FormatInt(uid, 10)

    var tvSlice []interface{} = make([]interface{}, len(nearTvs))
    for i, d := range nearTvs {
        tvSlice[i] = d
    }

    err = gRedisClient.SAdd(rediskey, tvSlice...).Err()
    if err != nil {
        log.Error("redis SAdd err, uid=%d, nearTvs=%v, err = %s \n",
            uid, nearTvs, err.Error())
        return
    }

    return
}

func GetPhoneNearTv(uid int64) (nearTvs map[string]struct{}, err error) {
    rediskey := "neartv:" + strconv.FormatInt(uid, 10)

    nearTvs, err = gRedisClient.SMembersMap(rediskey).Result()
    if err != nil {
        log.Error("redis get nearTv err, uid=%d, err = %s \n",
            uid, err.Error())
        return
    }

    return
}

func PhoneHasTv(phone int64, tv int64) (has bool, err error) {
    rediskey := "neartv:" + strconv.FormatInt(phone, 10)

    has, err = gRedisClient.SIsMember(rediskey, tv).Result()
    if err != nil {
        log.Error("redis ismember tv err, phone=%d, tv=%d, err = %s \n",
            phone, tv, err.Error())
        return
    }

    return
}
