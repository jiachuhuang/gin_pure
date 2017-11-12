package module

import (
	"time"
	"strings"
	"errors"
	"strconv"
	"github.com/gin-gonic/gin"
	"sync"
)

const (
	LQOK = 0x001
	LQBUSY = 0x010
	LQWAIT = 0x100
)

type LimitReq struct {
	sync.Mutex
	rate int64
	capacity int64
	last time.Time
	excess int64
}

func NewLimitReq(r string, cap int64) (*LimitReq, error){
	rate, err := paserRate(r)
	if err != nil {
		return nil, err
	}

	return &LimitReq{rate:rate,capacity:cap * 1000,last:time.Now(),excess:0},nil
}

func (lq *LimitReq) Acquire() (int, time.Duration){
	lq.Lock()
	defer lq.Unlock()

	now := time.Now()
	excess := lq.excess - (lq.rate * now.Sub(lq.last).Nanoseconds() / 1000000000) + 1000
	if excess < 0 {
		excess = 0
	}

	if excess > lq.capacity {
		return LQBUSY, 0
	} else if excess == 0 {
		lq.last = now
		lq.excess = excess
		return LQOK, 0
	} else {
		waitTms := time.Duration(excess * 1000 / lq.rate)  * time.Millisecond
		lq.last = now
		lq.excess = excess
		return LQWAIT, waitTms
	}
}

func LimitReqAcquire(lq *LimitReq) gin.HandlerFunc {
	return func(context *gin.Context) {
		status, waitTms := lq.Acquire()

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