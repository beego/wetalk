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
	"github.com/beego/wetalk/utils"
)

func GetMailTmplData(lang string, user *User) map[interface{}]interface{} {
	data := make(map[interface{}]interface{}, 10)
	data["AppDescription"] = utils.AppDescription
	data["AppKeywords"] = utils.AppKeywords
	data["AppName"] = utils.AppName
	data["AppVer"] = utils.AppVer
	data["AppUrl"] = utils.AppUrl
	data["AppLogo"] = utils.AppLogo
	data["IsProMode"] = utils.IsProMode
	data["Lang"] = lang
	data["ActiveCodeLives"] = utils.ActiveCodeLives
	data["ResetPwdCodeLives"] = utils.ResetPwdCodeLives
	if user != nil {
		data["User"] = user
	}
	return data
}
