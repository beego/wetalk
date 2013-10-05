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
	"fmt"

	"github.com/beego/i18n"
	"github.com/beego/wetalk/mailer"
	"github.com/beego/wetalk/utils"
)

// Create New mail message use MailFrom and MailUser
func NewMailMessage(To []string, subject, body string) mailer.Message {
	msg := mailer.NewHtmlMessage(To, utils.MailFrom, subject, body)
	msg.User = utils.MailUser
	return msg
}

// Send user register mail with active code
func SendRegisterMail(locale i18n.Locale, user *User) {
	code := CreateUserActiveCode(user, nil)

	subject := locale.Tr("Register success, Welcome")
	body := locale.Tr("code: %s", code)

	msg := NewMailMessage([]string{user.Email}, subject, body)
	msg.Info = fmt.Sprintf("UID: %d, send register mail", user.Id)

	// async send mail
	mailer.SendAsync(msg)
}

// Send user reset password mail with verify code
func SendResetPwdMail(locale i18n.Locale, user *User) {
	code := CreateUserResetPwdCode(user, nil)

	subject := locale.Tr("Fogot password verify email")
	body := locale.Tr("code: %s", code)

	msg := NewMailMessage([]string{user.Email}, subject, body)
	msg.Info = fmt.Sprintf("UID: %d, send reset password mail", user.Id)

	// async send mail
	mailer.SendAsync(msg)
}

// Send email verify active email.
func SendActiveMail(locale i18n.Locale, user *User) {
	code := CreateUserActiveCode(user, nil)

	subject := locale.Tr("Verify your email address")
	body := locale.Tr("code: %s", code)

	msg := NewMailMessage([]string{user.Email}, subject, body)
	msg.Info = fmt.Sprintf("UID: %d, send email verify mail", user.Id)

	// async send mail
	mailer.SendAsync(msg)
}
