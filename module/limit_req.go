package module

import (
	"time"
	"strings"
	"errors"
	"strconv"
	"github.com/gin-gonic/gin"
	"sync"
	"pure/utils"
)

const (
	LQOK = 0x001
	LQBUSY = 0x010
	LQWAIT = 0x100
)

const (
	LQMDEFAULT uint8 = 1 << iota
	LQMIP
	LQMEND
)

type LimitReq struct {
	sync.Mutex
	rate int64
	capacity int64
	mode uint8
	set map[string]*LimitReqNode
	lruQueue *utils.Queue
	qsize int
	qlen int
}

type LimitReqNode struct {
	last time.Time
	excess int64
	n *utils.QNode
}

func NewLimitReq(r string, cap int64, m uint8, qs int) (*LimitReq, error){
	if m >= LQMEND {
		return nil, errors.New("unknown model")
	}

	if qs < 128 || qs > 65535 {
		return nil, errors.New("lru queue size[128,65535] is not valid")
	}

	rate, err := paserRate(r)
	if err != nil {
		return nil, err
	}

	lq := &LimitReq{
		rate:rate,
		capacity:cap * 1000,
		set:make(map[string]*LimitReqNode),
		lruQueue:utils.NewQueue(),
		mode:m,
		qsize:qs,
		}
	return lq,nil
}

func (lq *LimitReq) GetSetKey(context *gin.Context) (string, error) {
	switch lq.mode {
	case LQMIP:
		return context.ClientIP(), nil
	case LQMDEFAULT:
		return "default", nil
	default:
		return "", errors.New("unknown limit req mode")
	}
}

func (lq *LimitReq) Acquire(key string) (int, time.Duration) {
	lq.Lock()
	defer lq.Unlock()

	now := time.Now()

	lrn, ok := lq.set[key]
	if ok {
		lq.lruQueue.RemoveNode(lrn.n)
		lq.lruQueue.InsertHeader(lrn.n)
	} else {
		if lq.qlen >= lq.qsize {
			tn := lq.lruQueue.GetTailNode()
			_, tnok := tn.Value.(string)
			if tn != nil && tnok{
				delete(lq.set, tn.Value.(string))
				lq.lruQueue.RemoveNode(tn)
				lq.qlen--
			} else {
				lq.qlen = 0
				lq.set = make(map[string]*LimitReqNode)
				lq.lruQueue = utils.NewQueue()
			}
		}
		n := &utils.QNode{Value: key}
		lrn = &LimitReqNode{now,0,n}
		lq.set[key] = lrn
		lq.lruQueue.InsertHeader(n)
		lq.qlen++
	}

	excess := lrn.excess - (lq.rate * now.Sub(lrn.last).Nanoseconds() / 1000000000) + 1000
	if excess < 0 {
		excess = 0
	}

	if excess > lq.capacity {
		return LQBUSY, 0
	} else if excess == 0 {
		lrn.last = now
		lrn.excess = excess
		return LQOK, 0
	} else {
		waitTms := time.Duration(excess * 1000 / lq.rate)  * time.Millisecond
		lrn.last = now
		lrn.excess = excess
		return LQWAIT, waitTms
	}
}

func LimitReqAcquire(lq *LimitReq) gin.HandlerFunc {
	return func(context *gin.Context) {
		key, err := lq.GetSetKey(context)
		if err != nil {
			context.AbortWithStatusJSON(503, "system error")
			return
		}

		status, waitTms := lq.Acquire(key)

		switch status {
		case LQBUSY:
			context.AbortWithStatusJSON(503, "busy")
			return
		case LQWAIT:
			time.Sleep(waitTms)
		case LQOK:
			context.Next()
			break
		}
	}
}

func paserRate(r string) (int64, error) {
	l := len(r)
	if l < 4 {
		return 0, errors.New("limit req rate param error")
	}

	slash := strings.LastIndexByte(r, '/')
	u := r[slash + 1:l]

	scale := 0
	switch u {
	case "s":
		scale = 1
	case "m":
		scale = 60
	default:
		scale = 0
	}

	rs := strings.TrimRight(r[:slash-1], " ")
	rate, err := strconv.Atoi(rs)

	if err != nil || scale == 0 {
		return 0, errors.New("limit req rate param error")
	}

	return int64(rate * 1000 / scale), nil
}