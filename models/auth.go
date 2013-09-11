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

package models

import (
	"fmt"
	"strings"

	"github.com/astaxie/beego/orm"
	"github.com/beego/beebbs/utils"
)

// CanRegistered checks if the username or e-mail is available.
func CanRegistered(userName string, email string) (bool, bool, error) {
	cond := orm.NewCondition()
	cond = cond.Or("UserName", userName).Or("Email", email)

	var maps []orm.Params
	o := orm.NewOrm()
	n, err := o.QueryTable("user").SetCond(cond).Values(&maps, "UserName", "Email")
	if err != nil {
		return false, false, err
	}

	e1 := true
	e2 := true

	if n > 0 {
		for _, m := range maps {
			if e1 && orm.ToStr(m["UserName"]) == userName {
				e1 = false
			}
			if e2 && orm.ToStr(m["Email"]) == email {
				e2 = false
			}
		}
	}

	return e1, e2, nil
}

// register create user
func RegisterUser(form RegisterForm, user *User) error {
	// use random salt encode password
	salt := utils.GetRandomString(10)
	pwd := utils.EncodePassword(form.Password, salt)

	user.UserName = form.UserName
	user.Email = form.Email

	// save salt and encode password, use $ as split char
	user.Password = fmt.Sprintf("%s$%s", salt, pwd)

	// save md5 email value for gravatar
	user.GrEmail = utils.EncodeMd5(form.Email)

	return NewUser(user)
}

// verify username/email and password
func VerifyUser(username, password string, user *User) bool {
	// search user by username or email
	qs := orm.NewOrm().QueryTable("user")
	if strings.Index(username, "@") == -1 {
		qs = qs.Filter("UserName", username)
	} else {
		qs = qs.Filter("Email", username)
	}
	err := qs.One(user)
	if err != nil {
		// user not found
		return false
	}

	// split
	var salt, encoded string
	if len(user.Password) > 11 {
		salt = user.Password[:10]
		encoded = user.Password[:11]
	}

	if VerifyPassword(password, salt, encoded) {
		// success
		return true
	}
	return false
}

// compare raw password and encoded password
func VerifyPassword(rawPwd, salt, encodedPwd string) bool {
	return utils.EncodePassword(rawPwd, salt) == encodedPwd
}

// verify active code when active account
func VerifyUserActiveCode(user *User, code string) bool {
	days := utils.ActiveCodeLives
	data := utils.ToStr(user.Id) + user.Email + user.UserName + user.Password + user.Rands
	return utils.VerifyTimeLimitCode(data, days, code)
}

// create a time limit code for user active
func CreateUserActiveCode(user *User, startInf interface{}) string {
	days := utils.ActiveCodeLives
	data := utils.ToStr(user.Id) + user.Email + user.UserName + user.Password + user.Rands
	return utils.CreateTimeLimitCode(data, days, startInf)
}

// verify code when reset password
func VerifyUserResetPwdCode(user *User, code string) bool {
	days := utils.ResetPwdCodeLives
	data := utils.ToStr(user.Id) + user.Email + user.UserName + user.Password + user.Rands + user.Updated.String()
	return utils.VerifyTimeLimitCode(data, days, code)
}

// create a time limit code for user reset password
func CreateUserResetPwdCode(user *User, startInf interface{}) string {
	days := utils.ResetPwdCodeLives
	data := utils.ToStr(user.Id) + user.Email + user.UserName + user.Password + user.Rands + user.Updated.String()
	return utils.CreateTimeLimitCode(data, days, startInf)
}
