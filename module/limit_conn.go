package module

import (
	"errors"
	"github.com/gin-gonic/gin"
	"sync"
)

type LimitConn struct {
	sync.Mutex
	conn int32
	maxconn int32
	status int
}

func NewLimitConn(maxconn int32) (lc *LimitConn, err error) {
	if maxconn > 65536 {
		lc = nil
		err = errors.New("connection limit must be less 65536")
	}

	lc = &LimitConn{
		conn:0,
		maxconn:maxconn,
		status:503,
	}
	return
}

func LimitConnAcquire(lc *LimitConn) gin.HandlerFunc {
	return func(context *gin.Context) {
		lc.Lock()
		if lc.conn >= lc.maxconn {
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
