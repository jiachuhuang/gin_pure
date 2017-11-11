package cache

import (
	"github.com/go-redis/redis"
	"strings"
	"errors"
	"strconv"
	"time"
)

type RedisCache struct {
	dns string
	password string
	DB int8
	client *redis.Client
}

func init() {
	Register("redis", NewRedisCache)
}

func NewRedisCache() Cache {
	return new(RedisCache)
}

func (this *RedisCache) Init(config string) (error) {
	options, err := this.paserConfig(config)
	if err != nil {
		return err
	}

	this.dns = options["addr"]
	this.password = options["password"]
	DB, _ := strconv.Atoi(options["DB"])
	this.DB = int8(DB)

	this.client = redis.NewClient(&redis.Options{
		Addr: this.dns,
		Password: this.password,
		DB: DB,
		PoolSize: 1000,
	})
	return nil
}

func (this *RedisCache) paserConfig(config string) (map[string]string, error) {
	configs := strings.Split(config, "@")
	if len(configs) < 1 {
		return nil, errors.New("error redis config")
	}

	options := make(map[string]string)
	switch len(configs) {
	case 1:
		options["addr"] = configs[0]
		options["password"] = ""
		options["DB"] = "0"
		break
	case 2:
		options["addr"] = configs[0]
		options["password"] = configs[1]
		options["DB"] = "0"
		break
	case 3:
		options["addr"] = configs[0]
		options["password"] = configs[1]
		options["DB"] = configs[2]
		break
	}
	return options, nil
}

func (this *RedisCache) Get(key string) (val interface{}) {
	return nil
}

func (this *RedisCache) Set(key string, data interface{}, expire time.Duration) (bool, error) {
	return true, nil
}

func (this *RedisCache) Delete(key string) (bool, error) {
	return true, nil
}

func (this *RedisCache) Flush() (bool, error) {
	return true, nil
}
