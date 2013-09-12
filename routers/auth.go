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
	"github.com/beego/wetalk/utils"
)

// LoginRouter serves login page.
type LoginRouter struct {
	baseRouter
}

// Get implemented login page.
func (this *LoginRouter) Get() {
	if this.isLogin {
		this.Redirect("/settings/profile", 302)
		return
	}

	this.Data["IsLoginPage"] = true
	this.TplNames = "auth/login.html"
}

// Login implemented user post login.
func (this *LoginRouter) Login() {
	if this.isLogin {
		this.Redirect("/settings/profile", 302)
		return
	}

}

// Logout implemented user logout page.
func (this *LoginRouter) Logout() {
	models.LogoutUser(this.CruSession)

	// set flash message
	flash := beego.NewFlash()
	flash.Data["HasLogout"] = "true"
	flash.Store(&this.Controller)

	this.Redirect("/login", 302)
}

// RegisterRouter serves register page.
type RegisterRouter struct {
	baseRouter
}

// Get implemented Get method for RegisterRouter.
func (this *RegisterRouter) Get() {
	if this.isLogin {
		this.Redirect("/settings/profile", 302)
		return
	}

	this.Data["IsRegister"] = true
	this.TplNames = "auth/register.html"
}

// Register implemented Post method for RegisterRouter.
func (this *RegisterRouter) Register() {
	if this.isLogin {
		this.Redirect("/settings/profile", 302)
		return
	}

	this.Data["IsRegister"] = true
	this.TplNames = "auth/register.html"

	// Get input form.
	form := models.RegisterForm{}
	this.ParseForm(&form)
	// Put data back in case users input invalid data for any section.
	this.Data["Form"] = form

	errs := make(map[string]validation.ValidationError)
	this.Data["FormError"] = errs

	// Verify basic input.
	valid := validation.Validation{}
	if ok, _ := valid.Valid(&form); !ok {
		utils.GetFirstValidErrors(valid.Errors, &errs)
		return
	}

	// Check if passwords of two times are same.
	if form.Password != form.PasswordRe {
		errs["PasswordRe"] = validation.ValidationError{
			Tmpl: this.Locale.Tr("Password not match first input"),
		}
		return
	}

	// Process register.
	e1, e2, err := models.CanRegistered(form.UserName, form.Email)
	if err != nil {
		beego.Error(err)
		return
	}

	if e1 && e2 {
		// Create new user.
		user := new(models.User)
		if err := models.RegisterUser(form, user); err == nil {
			models.SendRegisterMail(this.Locale, user)

			// login user
			models.LoginUser(this.CruSession, user)

			// set flash message
			flash := beego.NewFlash()
			flash.Data["RegSuccess"] = "true"
			flash.Store(&this.Controller)

			this.Redirect("/settings/profile", 302)

		} else {
			beego.Error(err)
		}

	} else {
		if !e1 {
			errs["UserName"] = validation.ValidationError{
				Tmpl: this.Locale.Tr("Username has been already taken"),
			}
		}

		if !e2 {
			errs["Email"] = validation.ValidationError{
				Tmpl: this.Locale.Tr("Email has been already taken"),
			}
		}
	}

}

// Active implemented check Email actice code.
func (this *RegisterRouter) Active() {
	code := this.Ctx.Input.Params(":code")

	if this.user.IsActive {
		this.Redirect("/settings/profile", 302)
		return
	}

	var user models.User

	beego.Info(models.VerifyUserActiveCode(&user, code))

	if models.VerifyUserActiveCode(&user, code) {
		user.IsActive = true
		user.Rands = utils.GetRandomString(10)
		if err := user.Update("IsActive", "Rands", "Updated"); err != nil {
			beego.Error(err)
		}
		if this.isLogin {
			this.user = user
		}
		this.Data["Success"] = true
	} else {
		this.Data["Success"] = false
	}

	this.TplNames = "auth/active.html"
}

// ForgotRouter serves login page.
type ForgotRouter struct {
	baseRouter
}

// Get implemented Get method for ForgotRouter.
func (this *ForgotRouter) Get() {
	this.TplNames = "auth/forgot.html"
}

// ResetRouter serves login page.
type ResetRouter struct {
	baseRouter
}

// Get implemented Get method for ResetRouter.
func (this *ResetRouter) Get() {
	this.TplNames = "auth/reset.html"
}

// SettingsRouter serves user settings.
type SettingsRouter struct {
	baseRouter
}

// Active implemented user account email active.
func (this *SettingsRouter) Active() {
	this.TplNames = "settings/profile.html"
}

// Profile implemented user profile settings page.
func (this *SettingsRouter) Profile() {
	if !this.isLogin {
		this.Redirect("/login", 302)
		return
	}

	this.TplNames = "settings/profile.html"
}

// ProfileSave implemented save user profile.
func (this *SettingsRouter) ProfileSave() {
	this.TplNames = "settings/profile.html"
}
