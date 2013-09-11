package models

import (
	"fmt"

	"github.com/beego/beebbs/mailer"
	"github.com/beego/beebbs/utils"
	"github.com/beego/i18n"
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

	subject := locale.Tr("Verify you email address.")
	body := locale.Tr("code: %s", code)

	msg := NewMailMessage([]string{user.Email}, subject, body)
	msg.Info = fmt.Sprintf("UID: %d, send register mail", user.Id)

	// async send mail
	mailer.SendAsync(msg)
}
