package cache

import (
	"sync"
	"pure/utils"
	"strconv"
	"errors"
	"time"
)

type MemoryCache struct{
	sync.RWMutex
	lruQueue *utils.Queue
	set map[string]MemoryCacheNode
	sum int
	num int
}

type MemoryCacheNode struct {
	data interface{}
	n *utils.QNode
	expire time.Duration
	createTime time.Time
}

func init() {
	Register("memory", NewMemoryCache)
}

func NewMemoryCache() Cache{
	return &MemoryCache{lruQueue: utils.NewQueue(), set: make(map[string]MemoryCacheNode), num:0}
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

func (this *MemoryCache) Get(key string) (val interface{}) {
	this.RLock()
	defer this.RUnlock()

	if mcn, ok := this.set[key]; ok {
		if mcn.isExpire() {
			this.lruQueue.RemoveNode(mcn.n)
			delete(this.set, key)
			this.num--
			return nil
		}
		return mcn.data
	}
	return nil
}

func (this *MemoryCache) Set(key string, data interface{}, expire time.Duration) (bool, error) {
	this.Lock()
	defer this.Unlock()

	if this.num >= this.sum {
		tn := this.lruQueue.GetTailNode()
		if tn != nil{
			delete(this.set, tn.Value.(string))
			this.lruQueue.RemoveNode(tn)
			this.num--
		}
	}

	n := &utils.QNode{Value: key}
	this.lruQueue.InsertHeader(n)
	mcn := MemoryCacheNode{data, n, expire, time.Now()}
	this.set[key] = mcn
	this.num++
	return true, nil
}

func (this *MemoryCache) Delete(key string) (bool, error) {
	return true, nil
}

func (this *MemoryCache) Flush() (bool, error) {
	return true, nil
}

func (mcn MemoryCacheNode) isExpire() bool {
	if mcn.expire == 0  {
		return false
	}

	return time.Now().Sub(mcn.createTime) > mcn.expire
}


