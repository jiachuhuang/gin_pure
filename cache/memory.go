package cache

import (
	"sync"
	"pure/utils"
	"strconv"
	"errors"
)

type MemoryCache struct{
	sync.Mutex
	lruQueue *utils.Queue
	set map[string]MemoryCacheNode
	sum int
	num int
}

type MemoryCacheNode struct {
	data interface{}
	n *utils.QNode
}

func init() {
	Register("memory", NewMemoryCache)
}

func NewMemoryCache() Cache{
	return &MemoryCache{lruQueue: utils.NewQueue(), num:0}
}

func (this *MemoryCache) Init(config string) (error) {
	sum, err := strconv.Atoi(config)
	if err != nil {
		return err
	}

	if sum < 128 || sum > 65535 {
		return errors.New("cache nodes num less than 128 or large than 65535")
	}

	this.sum = sum
	return nil
}

func (this *MemoryCache) Get(key interface{}) (val interface{}) {
	return nil
}

func (this *MemoryCache) Set(key interface{}, data interface{}) (bool, error) {

	return true, nil
}

func (this *MemoryCache) Delete(key interface{}) (bool, error) {
	return true, nil
}

func (this *MemoryCache) Flush() (bool, error) {
	return true, nil
}


