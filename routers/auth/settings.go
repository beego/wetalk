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

	"github.com/beego/wetalk/modules/auth"
	"github.com/beego/wetalk/routers/base"
)

// SettingsRouter serves user settings.
type SettingsRouter struct {
	base.BaseRouter
}

// Profile implemented user profile settings page.
func (this *SettingsRouter) Profile() {
	this.TplNames = "settings/profile.html"

	// need login
	if this.CheckLoginRedirect() {
		return
	}

	form := auth.ProfileForm{Locale: this.Locale}
	form.SetFromUser(&this.User)
	this.SetFormSets(&form)

	formPwd := auth.PasswordForm{}
	this.SetFormSets(&formPwd)
}

// ProfileSave implemented save user profile.
func (this *SettingsRouter) ProfileSave() {
	this.TplNames = "settings/profile.html"

	if this.CheckLoginRedirect() {
		return
	}

	action := this.GetString("action")

	if this.IsAjax() {
		switch action {
		case "send-verify-email":
			if this.User.IsActive {
				this.Data["json"] = false
			} else {
				auth.SendActiveMail(this.Locale, &this.User)
				this.Data["json"] = true
			}

			this.ServeJson()
			return
		}
		return
	}

	profileForm := auth.ProfileForm{Locale: this.Locale}
	profileForm.SetFromUser(&this.User)

	pwdForm := auth.PasswordForm{User: &this.User}

	this.Data["Form"] = profileForm

	switch action {
	case "save-profile":
		if this.ValidFormSets(&profileForm) {
			if err := profileForm.SaveUserProfile(&this.User); err != nil {
				beego.Error("ProfileSave: save-profile", err)
			}
			this.FlashRedirect("/settings/profile", 302, "ProfileSave")
			return
		}

	case "change-password":
		if this.ValidFormSets(&pwdForm) {
			// verify success and save new password
			if err := auth.SaveNewPassword(&this.User, pwdForm.Password); err == nil {
				this.FlashRedirect("/settings/profile", 302, "PasswordSave")
				return
			} else {
				beego.Error("ProfileSave: change-password", err)
			}
		}

	default:
		this.Redirect("/settings/profile", 302)
		return
	}

	if action != "save-profile" {
		this.SetFormSets(&profileForm)
	}
	if action != "change-password" {
		this.SetFormSets(&pwdForm)
	}
}
