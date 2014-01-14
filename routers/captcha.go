package routers

import (
	"github.com/astaxie/beego/context"

	"github.com/beego/wetalk/models"
)

func CaptchaFilter(ctx *context.Context) {
	models.CaptchaHandler.ServeHTTP(ctx.ResponseWriter, ctx.Request)
}
