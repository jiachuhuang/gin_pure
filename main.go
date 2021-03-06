package main

import (
	"github.com/gin-gonic/gin"
	"pure/controller"
	"pure/cache"
	"pure/utils"
	"pure/module"
)


func main() {
	app := gin.Default()
	lq, err := module.NewLimitReq("100r/s",5,2,1024)

	if err != nil {
	//	//log.Printf("[APP-error] "+"APP start error: %s", err)
		return
	}

	c := cache.NewCache("memory@256")
	utils.SetVar("cache", c, false)

	v1 := app.Group("v1")
	v1.Use(module.LimitReqAcquire(lq))
	{
		home := new(controller.HomeController)
		v1.GET("/", home.Index)
	}

	app.Run(":8080")
}
