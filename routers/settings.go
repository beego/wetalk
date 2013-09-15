package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"

	"github.com/beego/wetalk/models"
)

// SettingsRouter serves user settings.
type SettingsRouter struct {
	baseRouter
}

func (this *SettingsRouter) getProfileForm() {

}

// Profile implemented user profile settings page.
func (this *SettingsRouter) Profile() {
	this.TplNames = "settings/profile.html"

	// need login
	if this.CheckLoginRedirect() {
		return
	}

	form := models.ProfileForm{}
	form.SetFromUser(&this.user)

	this.Data["Form"] = form
}

// ProfileSave implemented save user profile.
func (this *SettingsRouter) ProfileSave() {
	this.TplNames = "settings/profile.html"

	// need login
	if this.CheckLoginRedirect() {
		return
	}

	action := this.GetString("action")

	if this.IsAjax() {
		switch action {
		case "send-verify-email":
			models.SendActiveMail(this.Locale, &this.user)

			this.Data["json"] = true
			this.ServeJson()
			return
		}
		return
	}

	form := models.ProfileForm{}
	form.SetFromUser(&this.user)

	this.Data["Form"] = form

	switch action {
	case "save-profile":
		if this.ValidForm(&form) {
			if err := form.SaveUserProfile(&this.user); err != nil {
				beego.Error("ProfileSave: save-profile", err)
			}
			this.FlashRedirect("/settings/profile", 302, "ProfileSave")
			return
		}

	case "change-password":
		form := models.PasswordForm{}
		if this.ValidForm(&form, "PwdForm") {
			if models.VerifyPassword(form.PasswordOld, this.user.Password) == false {
				this.SetFormError("PasswordOld", validation.ValidationError{
					Tmpl: this.Locale.Tr("Your old password not correct."),
				}, "PwdForm")
				return
			}

			if form.Password != form.PasswordRe {
				this.SetFormError("PasswordRe", validation.ValidationError{
					Tmpl: this.Locale.Tr("Password not match first input"),
				}, "PwdForm")
				return
			}

			if err := models.SaveNewPassword(&this.user, form.Password); err != nil {
				beego.Error("ProfileSave: change-password", err)
			}
			this.FlashRedirect("/settings/profile", 302, "PasswordSave")
		}

	default:
		this.Redirect("/settings/profile", 302)
	}

}
