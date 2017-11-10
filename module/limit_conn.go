package module

import (
	"errors"
	"github.com/gin-gonic/gin"
	"sync"
)

type LimitConn struct {
	sync.Mutex
	conn int32
	maxConn int32
	status int
}

func NewLimitConn(maxConn int32) (lc *LimitConn, err error) {
	if maxConn > 65536 {
		lc = nil
		err = errors.New("connection limit must be less 65536")
	}

	lc = &LimitConn{
		conn:0,
		maxConn:maxConn,
		status:503,
	}
	return
}

func LimitConnAcquire(lc *LimitConn) gin.HandlerFunc {
	return func(context *gin.Context) {
		lc.Lock()
		if lc.conn >= lc.maxConn {
			context.AbortWithStatusJSON(lc.status, "too many connection")
			lc.Unlock()
			return
		} else {
			lc.conn++
			lc.Unlock()
		}

		context.Next()

		lc.Lock()
		lc.conn--
		lc.Unlock()
	}
}
