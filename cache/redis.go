package cache


type RedisCache struct {
	dns string
	password string
	DB int8
}

func init() {
	Register("redis", NewRedisCache)
}

func NewRedisCache() Cache {
	return new(RedisCache)
}

func (this *RedisCache) Init(config string) (error) {
	return nil
}

func (this *RedisCache) Get(key interface{}) (val interface{}) {
	return nil
}

func (this *RedisCache) Set(key interface{}, data interface{}) (bool, error) {
	return true, nil
}

func (this *RedisCache) Delete(key interface{}) (bool, error) {
	return true, nil
}

func (this *RedisCache) Flush() (bool, error) {
	return true, nil
}
