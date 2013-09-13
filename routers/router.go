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

// Package routers implemented controller methods of beego.
package routers

import (
	"html/template"
	"net/url"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/beego/i18n"

	"github.com/beego/wetalk/models"
	"github.com/beego/wetalk/utils"
)

var langTypes []*langType // Languages are supported.

// langType represents a language type.
type langType struct {
	Lang, Name string
}

// baseRouter implemented global settings for all other routers.
type baseRouter struct {
	beego.Controller
	i18n.Locale
	user    models.User
	isLogin bool
}

// Prepare implemented Prepare method for baseRouter.
func (this *baseRouter) Prepare() {
	if utils.IsProMode {
	} else {
		utils.AppJsVer = beego.Date(time.Now(), "YmdHis")
		utils.AppCssVer = beego.Date(time.Now(), "YmdHis")
	}

	// Setting properties.
	this.Data["AppDescription"] = utils.AppDescription
	this.Data["AppKeywords"] = utils.AppKeywords
	this.Data["AppName"] = utils.AppName
	this.Data["AppVer"] = utils.AppVer
	this.Data["AppUrl"] = utils.AppUrl
	this.Data["AppJsVer"] = utils.AppJsVer
	this.Data["AppCssVer"] = utils.AppCssVer
	this.Data["IsProMode"] = utils.IsProMode
	this.Data["IsBeta"] = utils.IsBeta

	// Setting language version.
	if len(langTypes) == 0 {
		// Initialize languages.
		langs := strings.Split(utils.Cfg.MustValue("lang", "types"), "|")
		names := strings.Split(utils.Cfg.MustValue("lang", "names"), "|")
		langTypes = make([]*langType, 0, len(langs))
		for i, v := range langs {
			langTypes = append(langTypes, &langType{
				Lang: v,
				Name: names[i],
			})
		}
	}

	isNeedRedir, langVer := setLangVer(this.Ctx, this.Input(), this.Data)
	this.Locale.CurrentLocale = langVer
	// Redirect to make URL clean.
	if isNeedRedir {
		i := strings.Index(this.Ctx.Request.RequestURI, "?")
		this.Redirect(this.Ctx.Request.RequestURI[:i], 302)
	}

	// read flash message
	beego.ReadFromRequest(&this.Controller)

	// start session
	sess := this.StartSession()

	// save logined user if exist in session
	if models.GetUserFromSession(sess, &this.user) {
		this.isLogin = true
		this.Data["User"] = this.user
		this.Data["IsLogin"] = this.isLogin
	} else {
		this.isLogin = false
	}

	// pass xsrf helper to template context
	xsrfToken := this.Controller.XsrfToken()
	this.Data["xsrf_token"] = xsrfToken
	this.Data["xsrf_html"] = template.HTML(this.Controller.XsrfFormHtml())

	if this.NeedFlashRedirect(this.Ctx.Request.RequestURI) {
		this.EndFlashRedirect()
	}
}

// read beego flash message
func (this *baseRouter) FlashRead(key string) (string, bool) {
	if data, ok := this.Data["flash"].(map[string]string); ok {
		value, ok := data[key]
		return value, ok
	}
	return "", false
}

// write beego flash message
func (this *baseRouter) FlashWrite(key string, value string) {
	flash := beego.NewFlash()
	flash.Data[key] = value
	flash.Store(&this.Controller)
}

// check flash redirect
func (this *baseRouter) NeedFlashRedirect(anys ...string) bool {
	v := this.GetSession("on_redirect")
	if s, ok := v.(string); ok {
		parts := strings.Split(s, "\r\n")
		flag := parts[0]
		value := ""
		if len(parts) > 1 {
			value = parts[1]
		}
		// if match any then return true
		for _, s := range anys {
			if flag == s || value == s {
				return true
			}
		}
	}
	return false
}

// set flash redirect
func (this *baseRouter) FlashRedirect(flag string, uri string, code int) {
	var value string
	if uri != "" {
		value = flag + "\r\n" + uri
	} else {
		value = flag
	}
	this.SetSession("on_redirect", value)
	this.Redirect(uri, code)
}

// clear flash redirect
func (this *baseRouter) EndFlashRedirect() {
	this.DelSession("on_redirect")
}

// check form once, void re-submit
func (this *baseRouter) FormOnceNotMatch() bool {
	notMatch := false
	recreat := false
	// exist in request value
	if value, ok := this.Input()["_once"]; ok && len(value) > 0 {
		// exist in session
		if v, ok := this.GetSession("form_once").(string); ok && v != "" {
			// not match
			if value[0] != v {
				notMatch = true
			} else {
				// if matched then re-creat once
				recreat = true
			}
		}
	}
	this.FormOnceCreate(recreat)
	return notMatch
}

// create form once html
func (this *baseRouter) FormOnceCreate(args ...bool) {
	var value string
	var creat bool
	creat = len(args) > 0 && args[0]
	if !creat {
		if v, ok := this.GetSession("form_once").(string); ok && v != "" {
			value = v
		} else {
			creat = true
		}
	}
	if creat {
		value = utils.GetRandomString(10)
		this.SetSession("form_once", value)
	}
	this.Data["once_html"] = template.HTML(`<input type="hidden" name="_once" value="` + value + `">`)
}

// setLangVer sets site language version.
func setLangVer(ctx *context.Context, input url.Values, data map[interface{}]interface{}) (bool, string) {
	isNeedRedir := false

	// 1. Check URL arguments.
	lang := input.Get("lang")

	// 2. Get language information from cookies.
	if len(lang) == 0 {
		ck, err := ctx.Request.Cookie("lang")
		if err == nil {
			lang = ck.Value
		}
	} else {
		isNeedRedir = true
	}

	// Check again in case someone modify by purpose.
	isValid := false
	for _, v := range langTypes {
		if lang == v.Lang {
			isValid = true
			break
		}
	}
	if !isValid {
		lang = ""
		isNeedRedir = false
	}

	// 3. Get language information from 'Accept-Language'.
	if len(lang) == 0 {
		al := ctx.Request.Header.Get("Accept-Language")
		if len(al) > 4 {
			al = al[:5] // Only compare first 5 letters.
			for _, v := range langTypes {
				if al == v.Lang {
					lang = al
					break
				}
			}
		}
	}

	// 4. DefaucurLang language is English.
	if len(lang) == 0 {
		lang = "en-US"
		isNeedRedir = false
	}

	curLang := langType{
		Lang: lang,
	}

	// Save language information in cookies.
	ctx.SetCookie("lang", curLang.Lang, 1<<31-1, "/")

	restLangs := make([]*langType, 0, len(langTypes)-1)
	for _, v := range langTypes {
		if lang != v.Lang {
			restLangs = append(restLangs, v)
		} else {
			curLang.Name = v.Name
		}
	}

	// Set language properties.
	data["Lang"] = curLang.Lang
	data["CurLang"] = curLang.Name
	data["RestLangs"] = restLangs

	return isNeedRedir, curLang.Lang
}
