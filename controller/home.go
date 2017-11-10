package controller

import "github.com/gin-gonic/gin"

type HomeController struct {

}

func (this *HomeController) Index(context *gin.Context) {
	context.String(200, "OK")
}
