package controller

import (
	"github.com/gin-gonic/gin"
	"pure/utils"
	"pure/cache"
)

type HomeController struct {

}

func (this *HomeController) Index(context *gin.Context) {
	c := utils.GetVar("cache").(cache.Cache)
	c.Get("abc")
	context.String(200, "ok")
}
