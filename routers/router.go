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
	"github.com/astaxie/beego/validation"
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
	// check flash redirect, if match url then end, else for redirect return
	if match, redir := this.CheckFlashRedirect(this.Ctx.Request.RequestURI); redir {
		return
	} else if match {
		this.EndFlashRedirect()
	}

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
	this.Data["AvatarURL"] = utils.AvatarURL
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
	if models.GetUserFromSession(&this.user, sess) {
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

	// if method is GET then auto create a form once token
	if this.Ctx.Request.Method == "GET" {
		this.FormOnceCreate()
	}
}

// check if user not active then redirect
func (this *baseRouter) CheckActiveRedirect(args ...interface{}) bool {
	var url string
	needActive := true
	for _, arg := range args {
		switch v := arg.(type) {
		case bool:
			needActive = v
		case string:
			// custom redirect url
			url = v
		}
	}
	if needActive {
		// if need active and no login then redirect to login
		if this.CheckLoginRedirect() {
			return true
		}
		// redirect to active page
		if !this.user.IsActive {
			this.FlashRedirect("/settings/profile", 302, "NeedActive")
			return true
		}
	} else {
		// no need active
		if this.user.IsActive {
			if url == "" {
				url = "/"
			}
			this.Redirect(url, 302)
			return true
		}
	}
	return false

}

// check if not login then redirect
func (this *baseRouter) CheckLoginRedirect(args ...interface{}) bool {
	var url string
	needLogin := true
	for _, arg := range args {
		switch v := arg.(type) {
		case bool:
			needLogin = v
		case string:
			// custom redirect url
			url = v
		}
	}

	// if need login then redirect to /login
	if needLogin && !this.isLogin {
		this.Redirect("/login", 302)
		return true
	}

	// if not need login then redirect to /
	if !needLogin && this.isLogin {
		if url == "" {
			url = "/"
		}
		this.Redirect(url, 302)
		return true
	}
	return false
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

// check flash redirect, ensure browser redirect to uri and display flash message.
func (this *baseRouter) CheckFlashRedirect(value string) (match bool, redirect bool) {
	v := this.GetSession("on_redirect")
	if params, ok := v.([]interface{}); ok {
		if len(params) != 5 {
			this.EndFlashRedirect()
			goto end
		}
		uri := utils.ToStr(params[0])
		code := 302
		if c, ok := params[1].(int); ok {
			if c/100 == 3 {
				code = c
			}
		}
		flag := utils.ToStr(params[2])
		flagVal := utils.ToStr(params[3])
		times := 0
		if v, ok := params[4].(int); ok {
			times = v
		}

		times += 1
		if times > 3 {
			// if max retry times reached then end
			this.EndFlashRedirect()
			goto end
		}

		// match uri or flash flag
		if uri == value || flag == value {
			match = true
		} else {
			// if no match then continue redirect
			this.FlashRedirect(uri, code, flag, flagVal, times)
			redirect = true
		}
	}
end:
	return match, redirect
}

// set flash redirect
func (this *baseRouter) FlashRedirect(uri string, code int, flag string, args ...interface{}) {
	flagVal := "true"
	times := 0
	for _, arg := range args {
		switch v := arg.(type) {
		case string:
			flagVal = v
		case int:
			times = v
		}
	}

	if len(uri) == 0 || uri[0] != '/' {
		panic("flash reirect only support same host redirect")
	}

	params := []interface{}{uri, code, flag, flagVal, times}
	this.SetSession("on_redirect", params)

	this.FlashWrite(flag, flagVal)
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

// valid form and put errors to tempalte context
func (this *baseRouter) ValidForm(form interface{}, names ...string) bool {
	// parse request params to form ptr struct
	this.ParseForm(form)

	// Put data back in case users input invalid data for any section.
	name := "Form"
	if len(names) > 0 {
		name = names[0]
	}
	this.Data[name] = form

	errName := "FormError"
	if len(names) > 1 {
		errName = names[1]
	}

	// check form once
	if this.FormOnceNotMatch() {
		return false
	}

	// Verify basic input.
	valid := validation.Validation{}
	if ok, _ := valid.Valid(form); !ok {
		errs := make(map[string]validation.ValidationError)
		utils.GetFirstValidErrors(valid.Errors, &errs)
		this.Data[errName] = errs
		return false
	}
	return true
}

// add valid error to FormError
func (this *baseRouter) SetFormError(field string, err validation.ValidationError, names ...string) {
	errName := "FormError"
	if len(names) > 0 {
		errName = names[0]
	}

	var errs map[string]validation.ValidationError
	if er, ok := this.Data[errName].(map[string]validation.ValidationError); ok {
		errs = er
	} else {
		errs = make(map[string]validation.ValidationError)
		this.Data[errName] = errs
	}
	errs[field] = err
}

// check xsrf and show a friendly page
func (this *baseRouter) CheckXsrfCookie() {
	token := this.GetString("_xsrf")
	if token == "" {
		token = this.Ctx.Request.Header.Get("X-Xsrftoken")
	}
	if token == "" {
		token = this.Ctx.Request.Header.Get("X-Csrftoken")
	}
	if token == "" {
		this.Ctx.Abort(403, "'_xsrf' argument missing from POST")
	} else if this.XsrfToken() != token {
		this.Ctx.Abort(403, "XSRF cookie does not match POST argument")
	}
}

func (this *baseRouter) SystemException() {

}

func (this *baseRouter) IsAjax() bool {
	return this.Ctx.Input.Header("X-Requested-With") == "XMLHttpRequest"
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
