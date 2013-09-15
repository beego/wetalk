package models

import (
	"github.com/beego/wetalk/utils"
	"strings"
)

// Register form
type RegisterForm struct {
	UserName   string `form:"username" valid:"Required;AlphaDash;MinSize(5);MaxSize(30)"`
	Email      string `form:"email" valid:"Required;Email;MaxSize(80)"`
	Password   string `form:"password" valid:"Required;MinSize(4);MaxSize(30)"`
	PasswordRe string `form:"passwordre" valid:"Required;MinSize(4);MaxSize(30)"`
}

// Login form
type LoginForm struct {
	UserName string `form:"username" valid:"Required"`
	Password string `form:"password" valid:"Required"`
}

// Forgot form
type FogotForm struct {
	Email string `form:"email" valid:"Required;Email;MaxSize(80)"`
}

// Reset password form
type RestPwdForm struct {
	Password   string `form:"password" valid:"Required;MinSize(4);MaxSize(30)"`
	PasswordRe string `form:"passwordre" valid:"Required;MinSize(4);MaxSize(30)"`
}

// Settings Profile form
type ProfileForm struct {
	NickName  string `form:"nickname" valid:"Required;MaxSize(30)"`
	Url       string `form:"url" valid:"MaxSize(100)"`
	Info      string `form:"info" valid:"MaxSize(255)"`
	Email     string `form:"email" valid:"Required;Email;MaxSize(100)"`
	HideEmail bool   `form:"hideemail" valid:""`
	GrEmail   string `form:"gremail" valid:"Required;MaxSize(80)"`
}

func (form *ProfileForm) SetFromUser(user *User) {
	form.NickName = user.NickName
	form.Url = user.Url
	form.Info = user.Info
	form.Email = user.Email
	form.HideEmail = user.HideEmail
	form.GrEmail = user.GrEmail
}

func (form *ProfileForm) SaveUserProfile(user *User) error {
	user.NickName = form.NickName
	user.Url = form.Url
	user.Info = form.Info
	user.HideEmail = form.HideEmail

	// if email changed then need re-active
	if user.Email != form.Email {
		user.IsActive = false
		user.Email = form.Email
	}

	// set md5 value if the value is an email
	user.GrEmail = form.GrEmail
	if strings.IndexRune(form.GrEmail, '@') != -1 {
		user.GrEmail = utils.EncodeMd5(form.GrEmail)
	}

	return user.Update("NickName", "Url", "Info", "Email", "HideEmail", "GrEmail", "IsActive")
}

// Change password form
type PasswordForm struct {
	PasswordOld string `form:"passwordold" valid:"Required;MinSize(4);MaxSize(30)"`
	Password    string `form:"password" valid:"Required;MinSize(4);MaxSize(30)"`
	PasswordRe  string `form:"passwordre" valid:"Required;MinSize(4);MaxSize(30)"`
}
