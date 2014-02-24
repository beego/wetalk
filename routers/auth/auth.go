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

package auth

import (
	"github.com/astaxie/beego"
	"strings"

	"github.com/beego/wetalk/modules/auth"
	"github.com/beego/wetalk/modules/models"
	"github.com/beego/wetalk/modules/utils"
	"github.com/beego/wetalk/routers/base"
	"github.com/beego/wetalk/setting"
)

// LoginRouter serves login page.
type LoginRouter struct {
	base.BaseRouter
}

// Get implemented login page.
func (this *LoginRouter) Get() {
	this.Data["IsLoginPage"] = true
	this.TplNames = "auth/login.html"

	loginRedirect := strings.TrimSpace(this.GetString("to"))
	if utils.IsMatchHost(loginRedirect) == false {
		loginRedirect = "/"
	}

	// no need login
	if this.CheckLoginRedirect(false, loginRedirect) {
		return
	}

	if len(loginRedirect) > 0 {
		this.Ctx.SetCookie("login_to", loginRedirect, 0, "/")
	}

	form := auth.LoginForm{}
	this.SetFormSets(&form)
}

// Login implemented user login.
func (this *LoginRouter) Login() {
	this.Data["IsLoginPage"] = true
	this.TplNames = "auth/login.html"

	// no need login
	if this.CheckLoginRedirect(false) {
		return
	}

	var user models.User
	var key string
	ajaxErrMsg := "auth.login_error_ajax"

	form := auth.LoginForm{}
	// valid form and put errors to template context
	if this.ValidFormSets(&form) == false {
		if this.IsAjax() {
			goto ajaxError
		}
		return
	}

	key = "auth.login." + form.UserName + this.Ctx.Input.IP()
	if times, ok := utils.TimesReachedTest(key, setting.LoginMaxRetries); ok {
		if this.IsAjax() {
			ajaxErrMsg = "auth.login_error_times_reached"
			goto ajaxError
		}
		this.Data["ErrorReached"] = true

	} else if auth.VerifyUser(&user, form.UserName, form.Password) {
		loginRedirect := this.LoginUser(&user, form.Remember)

		if this.IsAjax() {
			this.Data["json"] = map[string]interface{}{
				"success":  true,
				"message":  this.Tr("auth.login_success_ajax"),
				"redirect": loginRedirect,
			}
			this.ServeJson()
			return
		}

		this.Redirect(loginRedirect, 302)
		return
	} else {
		utils.TimesReachedSet(key, times, setting.LoginFailedBlocks)
		if this.IsAjax() {
			goto ajaxError
		}
	}
	this.Data["Error"] = true
	return

ajaxError:
	this.Data["json"] = map[string]interface{}{
		"success": false,
		"message": this.Tr(ajaxErrMsg),
		"once":    this.Data["once_token"],
	}
	this.ServeJson()
}

// Logout implemented user logout page.
func (this *LoginRouter) Logout() {
	auth.LogoutUser(this.Ctx)

	// write flash message
	this.FlashWrite("HasLogout", "true")

	this.Redirect("/login", 302)
}

// RegisterRouter serves register page.
type RegisterRouter struct {
	base.BaseRouter
}

// Get implemented Get method for RegisterRouter.
func (this *RegisterRouter) Get() {
	// no need login
	if this.CheckLoginRedirect(false) {
		return
	}

	this.Data["IsRegister"] = true
	this.TplNames = "auth/register.html"

	form := auth.RegisterForm{Locale: this.Locale}
	this.SetFormSets(&form)
}

// Register implemented Post method for RegisterRouter.
func (this *RegisterRouter) Register() {
	this.Data["IsRegister"] = true
	this.TplNames = "auth/register.html"

	// no need login
	if this.CheckLoginRedirect(false) {
		return
	}

	form := auth.RegisterForm{Locale: this.Locale}
	// valid form and put errors to template context
	if this.ValidFormSets(&form) == false {
		return
	}

	// Create new user.
	user := new(models.User)

	if err := auth.RegisterUser(user, form.UserName, form.Email, form.Password); err == nil {
		auth.SendRegisterMail(this.Locale, user)

		loginRedirect := this.LoginUser(user, false)
		if loginRedirect == "/" {
			this.FlashRedirect("/settings/profile", 302, "RegSuccess")
		} else {
			this.Redirect(loginRedirect, 302)
		}

		return

	} else {
		beego.Error("Register: Failed ", err)
	}
}

// Active implemented check Email actice code.
func (this *RegisterRouter) Active() {
	this.TplNames = "auth/active.html"

	// no need active
	if this.CheckActiveRedirect(false) {
		return
	}

	code := this.GetString(":code")

	var user models.User

	if auth.VerifyUserActiveCode(&user, code) {
		user.IsActive = true
		user.Rands = models.GetUserSalt()
		if err := user.Update("IsActive", "Rands", "Updated"); err != nil {
			beego.Error("Active: user Update ", err)
		}
		if this.IsLogin {
			this.User = user
		}

		this.Redirect("/active/success", 302)

	} else {
		this.Data["Success"] = false
	}
}

// ActiveSuccess implemented success page when email active code verified.
func (this *RegisterRouter) ActiveSuccess() {
	this.TplNames = "auth/active.html"

	this.Data["Success"] = true
}

// ForgotRouter serves login page.
type ForgotRouter struct {
	base.BaseRouter
}

// Get implemented Get method for ForgotRouter.
func (this *ForgotRouter) Get() {
	this.TplNames = "auth/forgot.html"

	// no need login
	if this.CheckLoginRedirect(false) {
		return
	}

	form := auth.ForgotForm{}
	this.SetFormSets(&form)
}

// Get implemented Post method for ForgotRouter.
func (this *ForgotRouter) Post() {
	this.TplNames = "auth/forgot.html"

	// no need login
	if this.CheckLoginRedirect(false) {
		return
	}

	var user models.User
	form := auth.ForgotForm{User: &user}
	// valid form and put errors to template context
	if this.ValidFormSets(&form) == false {
		return
	}

	// send reset password email
	auth.SendResetPwdMail(this.Locale, &user)

	this.FlashRedirect("/forgot", 302, "SuccessSend")
}

// Reset implemented user password reset.
func (this *ForgotRouter) Reset() {
	this.TplNames = "auth/reset.html"

	code := this.GetString(":code")
	this.Data["Code"] = code

	var user models.User

	if auth.VerifyUserResetPwdCode(&user, code) {
		this.Data["Success"] = true
		form := auth.ResetPwdForm{}
		this.SetFormSets(&form)
	} else {
		this.Data["Success"] = false
	}
}

// Reset implemented user password reset.
func (this *ForgotRouter) ResetPost() {
	this.TplNames = "auth/reset.html"

	code := this.GetString(":code")
	this.Data["Code"] = code

	var user models.User

	if auth.VerifyUserResetPwdCode(&user, code) {
		this.Data["Success"] = true

		form := auth.ResetPwdForm{}
		if this.ValidFormSets(&form) == false {
			return
		}

		user.IsActive = true
		user.Rands = models.GetUserSalt()

		if err := auth.SaveNewPassword(&user, form.Password); err != nil {
			beego.Error("ResetPost Save New Password: ", err)
		}

		if this.IsLogin {
			auth.LogoutUser(this.Ctx)
		}

		this.FlashRedirect("/login", 302, "ResetSuccess")

	} else {
		this.Data["Success"] = false
	}
}
