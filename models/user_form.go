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

package models

import (
	"strings"

	"github.com/astaxie/beego/validation"
	"github.com/beego/i18n"

	"github.com/beego/wetalk/utils"
)

// Register form
type RegisterForm struct {
	UserName   string      `valid:"Required;AlphaDash;MinSize(5);MaxSize(30)"`
	Email      string      `valid:"Required;Email;MaxSize(80)"`
	Password   string      `form:"type(password)" valid:"Required;MinSize(4);MaxSize(30)"`
	PasswordRe string      `form:"type(password)" valid:"Required;MinSize(4);MaxSize(30)"`
	Locale     i18n.Locale `form:"-"`
}

func (form *RegisterForm) Valid(v *validation.Validation) {

	// Check if passwords of two times are same.
	if form.Password != form.PasswordRe {
		v.SetError("PasswordRe", "Password not match first input")
		return
	}

	e1, e2, _ := CanRegistered(form.UserName, form.Email)

	if !e1 {
		v.SetError("UserName", "Username has been already taken")
	}

	if !e2 {
		v.SetError("Email", "Email has been already taken")
	}
}

func (form *RegisterForm) Labels() map[string]string {
	return map[string]string{
		"UserName":   "Username",
		"PasswordRe": "Retype Password",
	}
}

func (form *RegisterForm) Helps() map[string]string {
	return map[string]string{
		"UserName": form.Locale.Tr("Min-length is %d", 5) + ", " + form.Locale.Tr("only contains %s", "a-z 0-9 - _"),
	}
}

func (form *RegisterForm) Placeholders() map[string]string {
	return map[string]string{
		"UserName":   "Please enter your username",
		"Email":      "Please enter your e-mail address",
		"Password":   "Please enter your password",
		"PasswordRe": "Please reenter your password",
	}
}

// Login form
type LoginForm struct {
	UserName string `valid:"Required"`
	Password string `form:"type(password)" valid:"Required"`
	Remember bool
}

func (form *LoginForm) Labels() map[string]string {
	return map[string]string{
		"UserName": "Username or Email",
		"Remember": "Remember Me",
	}
}

// Forgot form
type ForgotForm struct {
	Email string `valid:"Required;Email;MaxSize(80)"`
	User  *User  `form:"-"`
}

func (form *ForgotForm) Helps() map[string]string {
	return map[string]string{
		"Email": "This operaion lead to send your an e-mail with a reset secure link",
	}
}

func (form *ForgotForm) Valid(v *validation.Validation) {
	if HasUser(form.User, form.Email) == false {
		v.SetError("Email", "Wong email address, please check your input.")
	}
}

// Reset password form
type ResetPwdForm struct {
	Password   string `form:"type(password)" valid:"Required;MinSize(4);MaxSize(30)"`
	PasswordRe string `form:"type(password)" valid:"Required;MinSize(4);MaxSize(30)"`
}

func (form *ResetPwdForm) Valid(v *validation.Validation) {
	// Check if passwords of two times are same.
	if form.Password != form.PasswordRe {
		v.SetError("PasswordRe", "Password not match first input")
		return
	}
}

func (form *ResetPwdForm) Labels() map[string]string {
	return map[string]string{
		"PasswordRe": "Retype Password",
	}
}

func (form *ResetPwdForm) Placeholders() map[string]string {
	return map[string]string{
		"Password":   "Please enter your password",
		"PasswordRe": "Please reenter your password",
	}
}

// Settings Profile form
type ProfileForm struct {
	NickName  string      `valid:"Required;MaxSize(30)"`
	Url       string      `valid:"MaxSize(100)"`
	Info      string      `form:"type(textarea)" valid:"MaxSize(255)"`
	Email     string      `valid:"Required;Email;MaxSize(100)"`
	HideEmail bool        `valid:""`
	GrEmail   string      `valid:"Required;MaxSize(80)"`
	Locale    i18n.Locale `form:"-"`
}

func (form *ProfileForm) SetFromUser(user *User) {
	utils.SetFormValues(user, form)
}

func (form *ProfileForm) SaveUserProfile(user *User) error {
	// set md5 value if the value is an email
	if strings.IndexRune(form.GrEmail, '@') != -1 {
		form.GrEmail = utils.EncodeMd5(form.GrEmail)
	}

	changes := utils.FormChanges(user, form)
	if len(changes) > 0 {
		// if email changed then need re-active
		if user.Email != form.Email {
			user.IsActive = false
			changes = append(changes, "IsActive")
		}

		utils.SetFormValues(form, user)
		return user.Update(changes...)
	}
	return nil
}

func (form *ProfileForm) Labels() map[string]string {
	return map[string]string{
		"NickName":  "Nickname",
		"HideEmail": "Private your email",
		"GrEmail":   "Gravatar Token",
		"Url":       "Website",
	}
}

func (form *ProfileForm) Helps() map[string]string {
	return map[string]string{
		"GrEmail": "Enter an email will convert to token, direct input token is supported",
		"Info":    form.Locale.Tr("Max-length is %d", 255),
	}
}

func (form *ProfileForm) Placeholders() map[string]string {
	return map[string]string{
		"NickName": "Please enter your nickname",
		"GrEmail":  "You can enter an another gravatar token",
		"Url":      "Please enter your site url",
		"Info":     "Please say something introduce yourself",
	}
}

// Change password form
type PasswordForm struct {
	PasswordOld string `form:"type(password)" valid:"Required"`
	Password    string `form:"type(password)" valid:"Required;MinSize(4);MaxSize(30)"`
	PasswordRe  string `form:"type(password)" valid:"Required;MinSize(4);MaxSize(30)"`
}

func (form *PasswordForm) Valid(v *validation.Validation) {
	// Check if passwords of two times are same.
	if form.Password != form.PasswordRe {
		v.SetError("PasswordRe", "Password not match first input")
		return
	}

	// if models.VerifyPassword(form.PasswordOld, this.user.Password) == false {
	// 	this.SetFormError(&form, fieldName, errMsg, ...)
	// 	return
	// }
}

func (form *PasswordForm) Labels() map[string]string {
	return map[string]string{
		"PasswordOld": "Old Password",
		"Password":    "New Password",
		"PasswordRe":  "Retype Password",
	}
}

func (form *PasswordForm) Placeholders() map[string]string {
	return map[string]string{
		"PasswordOld": "Please enter your old password",
		"Password":    "Please enter your new password",
		"PasswordRe":  "Please reenter your password",
	}
}

type UserAdminForm struct {
	Create    bool   `form:"-"`
	Id        int    `form:"-"`
	UserName  string `valid:"Required;AlphaDash;MinSize(5);MaxSize(30)"`
	Email     string `valid:"Required;Email;MaxSize(100)"`
	HideEmail bool   ``
	NickName  string `valid:"Required;MaxSize(30)"`
	Url       string `valid:"MaxSize(100)"`
	Info      string `form:"type(textarea)" valid:"MaxSize(255)"`
	GrEmail   string `valid:"Required;MaxSize(80)"`
	Followers int    ``
	Following int    ``
	IsAdmin   bool   ``
	IsActive  bool   ``
	IsForbid  bool   ``
}

func (form *UserAdminForm) Valid(v *validation.Validation) {
	qs := Users()

	if CheckIsExist(qs, "UserName", form.UserName, form.Id) {
		v.SetError("UserName", "Username has been already taken")
	}

	if CheckIsExist(qs, "Email", form.Email, form.Id) {
		v.SetError("Email", "Email has been already taken")
	}
}

func (form *UserAdminForm) Helps() map[string]string {
	return nil
}

func (form *UserAdminForm) Labels() map[string]string {
	return nil
}

func (form *UserAdminForm) SetFromUser(user *User) {
	utils.SetFormValues(user, form)
}

func (form *UserAdminForm) SetToUser(user *User) {
	// set md5 value if the value is an email
	if strings.IndexRune(form.GrEmail, '@') != -1 {
		form.GrEmail = utils.EncodeMd5(form.GrEmail)
	}

	utils.SetFormValues(form, user)
}
