package cache

import (
	"strings"
	"time"
)

type Cache interface {
	Init(config string) (error)
	Get(key string) (val interface{})
	Set(key string, data interface{}, expire time.Duration) (bool, error)
	Delete(key string) (bool, error)
	Flush() (bool, error)
}

type Driver func() Cache

var drivers = make(map[string]Driver)

type cache struct {
	driver Cache
	dns string
}

func NewCache(dns string) Cache {
	config := strings.SplitN(dns, "@", 2)
	if len(config) < 2 {
		panic("error param")
	}

	newDriver, ok := drivers[config[0]]
	if !ok {
		panic("driver not exists")
	}

	c := &cache{newDriver(), config[1]}
	c.Init(config[1])
	return c
}

func Register(name string, driver Driver) {
	if driver == nil {
		panic("driver can not be nil")
	}

	if _, ok := drivers[name]; ok {
		panic("driver had register")
	}

	drivers[name] = driver
}

func (this *cache) Init(config string) (error) {
	return this.driver.Init(config)
}

func (this *cache) Get(key string) (val interface{}) {
	return this.driver.Get(key)
}

func (this *cache) Set(key string, data interface{}, expire time.Duration) (bool, error) {
	return this.driver.Set(key, data, expire)
}

func (this *cache) Delete(key string) (bool, error) {
	return this.driver.Delete(key)
}

func (this *cache) Flush() (bool, error) {
	return this.driver.Flush()
}
