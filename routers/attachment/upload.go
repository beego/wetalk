// Copyright 2013 wetalk authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package attachment

import (
	"github.com/beego/wetalk/setting"
	"net/http"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"

	"github.com/beego/wetalk/modules/attachment"
	"github.com/beego/wetalk/modules/models"
	"github.com/beego/wetalk/routers/base"
)

type UploadRouter struct {
	base.BaseRouter
}

func (this *UploadRouter) Post() {
	result := map[string]interface{}{
		"success": false,
	}

	defer func() {
		this.Data["json"] = &result
		this.ServeJson()
	}()

	// check permition
	if !this.User.IsActive {
		return
	}

	// get file object
	file, handler, err := this.Ctx.Request.FormFile("image")
	if err != nil {
		return
	}
	defer file.Close()

	t := time.Now()

	image := models.Image{}
	image.User = &this.User

	// get mime type
	mime := handler.Header.Get("Content-Type")

	// save and resize image
	if err := attachment.SaveImage(&image, file, mime, handler.Filename, t); err != nil {
		beego.Error(err)
		return
	}

	result["link"] = image.LinkMiddle()
	result["success"] = true

}

func ImageFilter(ctx *context.Context) {
	token := path.Base(ctx.Request.RequestURI)

	// split token and file ext
	var filePath string
	if i := strings.IndexRune(token, '.'); i == -1 {
		return
	} else {
		filePath = token[i+1:]
		token = token[:i]
	}

	// decode token to file path
	var image models.Image
	if err := image.DecodeToken(token); err != nil {
		beego.Info(err)
		return
	}

	// file real path
	filePath = attachment.GenImagePath(&image) + filePath

	// if x-send on then set header and http status
	// fall back use proxy serve file
	if setting.ImageXSend {
		ext := filepath.Ext(filePath)
		ctx.Output.ContentType(ext)
		ctx.Output.Header(setting.ImageXSendHeader, "/"+filePath)
		ctx.Output.SetStatus(200)
	} else {
		// direct serve file use go
		http.ServeFile(ctx.ResponseWriter, ctx.Request, filePath)
	}
}
