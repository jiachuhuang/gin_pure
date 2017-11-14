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
	set map[string]*MemoryCacheNode
	sum int
	num int
	pool *sync.Pool
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
	return &MemoryCache{lruQueue: utils.NewQueue(), set: make(map[string]*MemoryCacheNode), num:0}
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

	this.pool = &sync.Pool{
		New: func() interface{} {
			return new(MemoryCacheNode)
		},
	}
	return nil
}

func (this *MemoryCache) NewMemoryCacheNode (data interface{}, expire time.Duration, n *utils.QNode) *MemoryCacheNode {
	mcn := this.pool.Get().(*MemoryCacheNode)
	mcn.reset()

	mcn.data = data
	mcn.expire = expire
	mcn.n = n
	return mcn
}

func (this *MemoryCache) Recycle(mcn *MemoryCacheNode) {
	if mcn != nil {
		this.pool.Put(mcn)
	}
}

func (mcn *MemoryCacheNode) reset () {
	mcn.createTime = time.Now()
	mcn.data = nil
	mcn.expire = 0
	mcn.n = nil
}

func (this *MemoryCache) Get(key string) (val interface{}) {
	this.RLock()
	defer this.RUnlock()

	if mcn, ok := this.set[key]; ok {
		if mcn.isExpire() {
			this.lruQueue.RemoveNode(mcn.n)
			this.lruQueue.Recycle(mcn.n)
			delete(this.set, key)
			this.Recycle(mcn)
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
		tn := this.lruQueue.Pop()
		if tn != nil {
			tmcn := this.set[tn.Value.(string)]
			this.Recycle(tmcn)
			delete(this.set, tn.Value.(string))
			this.lruQueue.Recycle(tn)
			this.num--
		}
	}

	n := this.lruQueue.NewQNode(key)
	this.lruQueue.InsertHeader(n)
	mcn := this.NewMemoryCacheNode(data,expire,n)
	this.set[key] = mcn
	this.num++
	return true, nil
}

func (this *MemoryCache) Delete(key string) (bool, error) {
	this.Lock()
	defer this.Unlock()

	if mcn, ok := this.set[key]; ok {
		this.lruQueue.RemoveNode(mcn.n)
		this.lruQueue.Recycle(mcn.n)
		delete(this.set, key)
		this.Recycle(mcn)
		this.num--
	}
	return true, nil
}

func (this *MemoryCache) Flush() (bool, error) {
	this.Lock()
	this.Unlock()

	this.lruQueue.Clear()
	this.num = 0
	this.set = make(map[string]*MemoryCacheNode)
	return true, nil
}

func (mcn MemoryCacheNode) isExpire() bool {
	if mcn.expire == 0  {
		return false
	}

	return time.Now().Sub(mcn.createTime) > mcn.expire
}


