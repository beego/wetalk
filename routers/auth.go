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

package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"

	"github.com/beego/wetalk/models"
)

// LoginRouter serves login page.
type LoginRouter struct {
	baseRouter
}

// Get implemented login page.
func (this *LoginRouter) Get() {
	this.Data["IsLoginPage"] = true
	this.TplNames = "auth/login.html"

	// no need login
	if this.CheckLoginRedirect(false) {
		return
	}
}

// Login implemented user login.
func (this *LoginRouter) Login() {
	this.Data["IsLoginPage"] = true
	this.TplNames = "auth/login.html"

	// no need login
	if this.CheckLoginRedirect(false) {
		return
	}

	form := models.LoginForm{}
	// valid form and put errors to template context
	if this.ValidForm(&form) == false {
		return
	}

	if models.VerifyUser(&this.user, form.UserName, form.Password) {
		// login user
		models.LoginUser(&this.user, &this.Controller)

		this.Redirect("/", 302)
		return
	}

	this.Data["Error"] = true
}

// Logout implemented user logout page.
func (this *LoginRouter) Logout() {
	models.LogoutUser(&this.Controller)

	// write flash message
	this.FlashWrite("HasLogout", "true")

	this.Redirect("/login", 302)
}

// RegisterRouter serves register page.
type RegisterRouter struct {
	baseRouter
}

// Get implemented Get method for RegisterRouter.
func (this *RegisterRouter) Get() {
	// no need login
	if this.CheckLoginRedirect(false) {
		return
	}

	this.Data["IsRegister"] = true
	this.TplNames = "auth/register.html"
}

// Register implemented Post method for RegisterRouter.
func (this *RegisterRouter) Register() {
	this.Data["IsRegister"] = true
	this.TplNames = "auth/register.html"

	flashKey := "RegSuccess"

	// no need login
	if this.CheckLoginRedirect(false) {
		return
	}

	form := models.RegisterForm{}
	// valid form and put errors to template context
	if this.ValidForm(&form) == false {
		return
	}

	// Check if passwords of two times are same.
	if form.Password != form.PasswordRe {
		this.SetFormError("PasswordRe", validation.ValidationError{
			Tmpl: this.Locale.Tr("Password not match first input"),
		})
		return
	}

	// Process register.
	e1, e2, err := models.CanRegistered(form.UserName, form.Email)
	if err != nil {
		beego.Error("Register: CanRegistered", err)
		return
	}

	if e1 && e2 {
		// Create new user.
		user := new(models.User)
		if err := models.RegisterUser(user, form); err == nil {
			models.SendRegisterMail(this.Locale, user)

			// login user
			models.LoginUser(user, &this.Controller)

			this.FlashRedirect("/settings/profile", 302, flashKey)

			return

		} else {
			beego.Error("Register: RegisterUser", err)
		}

	} else {
		if !e1 {
			this.SetFormError("UserName", validation.ValidationError{
				Tmpl: this.Locale.Tr("Username has been already taken"),
			})
		}

		if !e2 {
			this.SetFormError("Email", validation.ValidationError{
				Tmpl: this.Locale.Tr("Email has been already taken"),
			})
		}
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

	if models.VerifyUserActiveCode(&user, code) {
		user.IsActive = true
		user.Rands = models.GetUserSalt()
		if err := user.Update("IsActive", "Rands", "Updated"); err != nil {
			beego.Error("Active: user Update ", err)
		}
		if this.isLogin {
			this.user = user
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
	baseRouter
}

// Get implemented Get method for ForgotRouter.
func (this *ForgotRouter) Get() {
	this.TplNames = "auth/forgot.html"

	// no need login
	if this.CheckLoginRedirect(false) {
		return
	}
}

// Get implemented Post method for ForgotRouter.
func (this *ForgotRouter) Post() {
	this.TplNames = "auth/forgot.html"

	// no need login
	if this.CheckLoginRedirect(false) {
		return
	}

	flashKey := "SuccessSend"

	form := models.FogotForm{}
	// valid form and put errors to template context
	if this.ValidForm(&form) == false {
		return
	}

	var user models.User
	if models.HasUser(&user, form.Email) {
		models.SendResetPwdMail(this.Locale, &user)

		this.FlashRedirect("/forgot", 302, flashKey)
		return

	} else {
		this.SetFormError("Email", validation.ValidationError{
			Tmpl: this.Locale.Tr("Wong email address, please check your input."),
		})
	}
}

// Reset implemented user password reset.
func (this *ForgotRouter) Reset() {
	this.TplNames = "auth/reset.html"

	code := this.GetString(":code")

	var user models.User

	if models.VerifyUserResetPwdCode(&user, code) {
		this.Data["Success"] = true
	} else {
		this.Data["Success"] = false
	}
}

// Reset implemented user password reset.
func (this *ForgotRouter) ResetPost() {
	this.TplNames = "auth/reset.html"

	code := this.GetString(":code")

	var user models.User

	if models.VerifyUserResetPwdCode(&user, code) {
		this.Data["Success"] = true

		form := models.RestPwdForm{}
		if this.ValidForm(&form) == false {
			return
		}

		// Check if passwords of two times are same.
		if form.Password != form.PasswordRe {
			this.SetFormError("PasswordRe", validation.ValidationError{
				Tmpl: this.Locale.Tr("Password not match first input"),
			})
			return
		}

		user.IsActive = true
		if err := models.SaveNewPassword(&user, form.Password); err != nil {
			beego.Error("ResetPost Save New Password: ", err)
		}

		if this.isLogin {
			models.LogoutUser(&this.Controller)
		}

		this.FlashRedirect("/login", 302, "ResetSuccess")

	} else {
		this.Data["Success"] = false
	}
}
