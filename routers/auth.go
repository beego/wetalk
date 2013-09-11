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
)

// LoginRouter serves login page.
type LoginRouter struct {
	baseRouter
}

// Get implemented Get method for LoginRouter.
func (this *LoginRouter) Get() {
	this.Data["IsLoginPage"] = true
	this.TplNames = "login.html"
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
	this.TplNames = "register.html"
}

// Post implemented Post method for RegisterRouter.
func (this *RegisterRouter) Post() {
	this.Data["IsRegister"] = true
	this.TplNames = "register.html"

	// Get input form.
	form := RegisterForm{}
	this.ParseForm(&form)
	// Put data back in case users input invalid data for any section.
	this.Data["Form"] = form

	errs := make(map[string]validation.ValidationError)
	this.Data["FormError"] = errs

	// Verify basic input.
	valid := validation.Validation{}
	if ok, _ := valid.Valid(&form); !ok {
		getFirstValidError(valid.Errors, &errs)
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
	if e1, e2, err := models.CanRegistered(form.UserName, form.Email); err != nil {
		beego.Error(err)
	} else {

		if e1 && e2 {
			if err := registerUser(form); err == nil {

				// TODO
				// forbid re submit
				// need send verify email
				// and redirect to /register/success

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
}

// ForgotRouter serves login page.
type ForgotRouter struct {
	baseRouter
}

// Get implemented Get method for ForgotRouter.
func (this *ForgotRouter) Get() {
	this.TplNames = "forgot.html"
}

// ResetRouter serves login page.
type ResetRouter struct {
	baseRouter
}

// Get implemented Get method for ResetRouter.
func (this *ResetRouter) Get() {
	this.TplNames = "reset.html"
}
