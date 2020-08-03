package db

import (
	"encoding/json"
	"fmt"
	"nms-controller/model"
	"os"
	"strconv"
	"strings"

	"github.com/golang/glog"
	"gopkg.in/redis.v4"
)

// # REDIS 環境變數
// # REDIS_ADDR="61.219.26.45:32041"
// # REDIS_PASSWORD=""
// # REDIS_DB="1"

type theRedis struct {
	RedisClient *redis.Client
}

func newTheRedis() *theRedis {
	return &theRedis{
		RedisClient: createRedisClient(),
	}
}

// 建立 redis 客戶端
func createRedisClient() *redis.Client {
	redisAddr := os.Getenv("REDIS_ADDR")
	redisIndex, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		glog.Error("env string to int:", err)
		panic(err)
	}
	redisPwd := os.Getenv("REDIS_PASSWORD")

	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPwd,
		DB:       redisIndex,
		PoolSize: 5,
	})

	_, err = client.Ping().Result()
	if err != nil {
		glog.Error("ping redis err:", err)
		panic(err)
	}

	fmt.Println("Successfully connected redis!", "db=", redisIndex)

	return client
}

func (r *theRedis) FlushRedis() {
	r.RedisClient.FlushDb()
}

func (r *theRedis) InsertRedis(cfg model.CustomConfig, metric model.Metric, rMetric model.RedisMetric) {
	//存入Redis
	queryName := strings.Replace(cfg.QueryName, "-", "_", -1)
	redisKey := queryName + ":" + metric.Header + ":" + metric.Type
	// fmt.Println("Redis key:" + redisKey)
	redisData, _ := json.Marshal(rMetric)
	r.RedisClient.RPush(redisKey, redisData)
}

func (r *theRedis) DeleteRedis(cfg model.CustomConfig, metric model.Metric) {
	//存入Redis
	queryName := strings.Replace(cfg.QueryName, "-", "_", -1)
	redisKey := queryName + ":" + metric.Header + ":" + metric.Type
	r.RedisClient.Del(redisKey)
}

//redis list 給exporter用
func (r *theRedis) InsertQueryNameKeys(cfg model.CustomConfig) {
	metrics := cfg.Metrics
	for _, metric := range metrics {
		queryName := strings.Replace(cfg.QueryName, "-", "_", -1)
		redisKey := queryName + ":" + metric.Header + ":" + metric.Type
		r.RedisClient.SAdd("querys", redisKey)
	}
}
