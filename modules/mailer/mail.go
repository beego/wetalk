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

package mailer

import (
	"github.com/beego/wetalk/modules/models"
	"github.com/beego/wetalk/setting"
)

// Create New mail message use MailFrom and MailUser
func NewMailMessage(To []string, subject, body string) Message {
	msg := NewHtmlMessage(To, setting.MailFrom, subject, body)
	msg.User = setting.MailUser
	return msg
}

func GetMailTmplData(lang string, user *models.User) map[interface{}]interface{} {
	data := make(map[interface{}]interface{}, 10)
	data["AppName"] = setting.AppName
	data["AppVer"] = setting.AppVer
	data["AppUrl"] = setting.AppUrl
	data["AppLogo"] = setting.AppLogo
	data["IsProMode"] = setting.IsProMode
	data["Lang"] = lang
	data["ActiveCodeLives"] = setting.ActiveCodeLives
	data["ResetPwdCodeLives"] = setting.ResetPwdCodeLives
	if user != nil {
		data["User"] = user
	}
	return data
}
