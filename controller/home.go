package controller

import (
	"github.com/gin-gonic/gin"
	"pure/utils"
	"pure/cache"
	"fmt"
	"time"
)

type HomeController struct {

}

func (this *HomeController) Index(context *gin.Context) {
	c := utils.GetVar("cache").(cache.Cache)
	c.Set("abc","aaa", 3*time.Second)
	v := c.Get("abc")
	fmt.Println(v)
	context.String(200, "ok")
}
