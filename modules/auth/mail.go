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
	"fmt"

	"github.com/beego/i18n"

	"github.com/beego/wetalk/modules/mailer"
	"github.com/beego/wetalk/modules/models"
	"github.com/beego/wetalk/modules/utils"
)

// Send user register mail with active code
func SendRegisterMail(locale i18n.Locale, user *models.User) {
	code := CreateUserActiveCode(user, nil)

	subject := locale.Tr("mail.register_success_subject")

	data := mailer.GetMailTmplData(locale.Lang, user)
	data["Code"] = code
	body := utils.RenderTemplate("mail/auth/register_success.html", data)

	msg := mailer.NewMailMessage([]string{user.Email}, subject, body)
	msg.Info = fmt.Sprintf("UID: %d, send register mail", user.Id)

	// async send mail
	mailer.SendAsync(msg)
}

// Send user reset password mail with verify code
func SendResetPwdMail(locale i18n.Locale, user *models.User) {
	code := CreateUserResetPwdCode(user, nil)

	subject := locale.Tr("mail.reset_password_subject")

	data := mailer.GetMailTmplData(locale.Lang, user)
	data["Code"] = code
	body := utils.RenderTemplate("mail/auth/reset_password.html", data)

	msg := mailer.NewMailMessage([]string{user.Email}, subject, body)
	msg.Info = fmt.Sprintf("UID: %d, send reset password mail", user.Id)

	// async send mail
	mailer.SendAsync(msg)
}

// Send email verify active email.
func SendActiveMail(locale i18n.Locale, user *models.User) {
	code := CreateUserActiveCode(user, nil)

	subject := locale.Tr("mail.verify_your_email_subject")

	data := mailer.GetMailTmplData(locale.Lang, user)
	data["Code"] = code
	body := utils.RenderTemplate("mail/auth/active_email.html", data)

	msg := mailer.NewMailMessage([]string{user.Email}, subject, body)
	msg.Info = fmt.Sprintf("UID: %d, send email verify mail", user.Id)

	// async send mail
	mailer.SendAsync(msg)
}
