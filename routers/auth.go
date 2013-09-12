// Copyright 2013 beebbs authors
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

	"github.com/beego/beebbs/models"
	"github.com/beego/beebbs/utils"
)

// LoginRouter serves login page.
type LoginRouter struct {
	baseRouter
}

// Get implemented Get method for LoginRouter.
func (this *LoginRouter) Get() {
	this.Data["IsLoginPage"] = true
	this.TplNames = "auth/login.html"
}

// Post implemented Post method for LoginRouter.
func (this *LoginRouter) Post() {

}

// RegisterRouter serves login page.
type RegisterRouter struct {
	baseRouter
}

// Get implemented Get method for RegisterRouter.
func (this *RegisterRouter) Get() {
	this.Data["IsRegister"] = true
	this.TplNames = "auth/register.html"
}

// Post implemented Post method for RegisterRouter.
func (this *RegisterRouter) Post() {
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
			this.Data["IsSuccess"] = true

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

// Success implemented Register Success Page.
func (this *RegisterRouter) Success() {

}

// Resend implemented post resend active code.
func (this *RegisterRouter) Resend() {

}

// Active implemented check Email actice code.
func (this *RegisterRouter) Active() {

}

// ActiveSuccess implemented Email active success page .
func (this *RegisterRouter) ActiveSuccess() {

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
