package main

import (
	"github.com/gin-gonic/gin"
	"pure/controller"
	"pure/module"
)

var (
	maxconn int32 = 1
)

func main() {
	app := gin.Default()
	lc, err := module.NewLimitConn(maxconn)

	if err != nil {
		//log.Printf("[APP-error] "+"APP start error: %s", err)
		return
	}

	v1 := app.Group("v1", module.LimitConnAcquire(lc))
	{
		home := new(controller.HomeController)
		v1.GET("/", home.Index)
	}

	app.Run(":8080")
}
