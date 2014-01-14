package models

import (
	"github.com/dchest/captcha"
)

var CaptchaHandler = captcha.Server(captcha.StdWidth, captcha.StdHeight)
