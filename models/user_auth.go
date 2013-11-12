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
	"encoding/hex"
	"fmt"
	"strings"
	// "time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/session"

	"github.com/beego/wetalk/utils"
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

// check if exist user by username or email
func HasUser(user *User, username string) bool {
	var err error
	qs := orm.NewOrm()
	if strings.IndexRune(username, '@') == -1 {
		user.UserName = username
		err = qs.Read(user, "UserName")
	} else {
		user.Email = username
		err = qs.Read(user, "Email")
	}
	if err == nil {
		return true
	}
	return false
}

// return a user salt token
func GetUserSalt() string {
	return utils.GetRandomString(10)
}

// register create user
func RegisterUser(user *User, form RegisterForm) error {
	// use random salt encode password
	salt := GetUserSalt()
	pwd := utils.EncodePassword(form.Password, salt)

	user.UserName = strings.ToLower(form.UserName)
	user.Email = strings.ToLower(form.Email)

	// save salt and encode password, use $ as split char
	user.Password = fmt.Sprintf("%s$%s", salt, pwd)

	// save md5 email value for gravatar
	user.GrEmail = utils.EncodeMd5(form.Email)

	// Use username as default nickname.
	user.NickName = user.UserName

	return user.Insert()
}

// set a new password to user
func SaveNewPassword(user *User, password string) error {
	salt := GetUserSalt()
	user.Password = fmt.Sprintf("%s$%s", salt, utils.EncodePassword(password, salt))
	user.Rands = GetUserSalt()
	return user.Update("Password", "Rands")
}

// login user
func LoginUser(user *User, c *beego.Controller, remember bool) {
	// werid way of beego session regenerate id...
	c.SessionRegenerateID()
	c.CruSession.Set("auth_user_id", user.Id)
}

// logout user
func LogoutUser(c *beego.Controller) {
	c.CruSession.Delete("auth_user_id")
	c.DestroySession()
}

// get user if key exist in session
func GetUserFromSession(user *User, sess session.SessionStore) bool {
	if id, ok := sess.Get("auth_user_id").(int); ok && id > 0 {
		*user = User{Id: id}
		if user.Read() == nil {
			return true
		}
	}

	return false
}

// verify username/email and password
func VerifyUser(user *User, username, password string) (success bool) {
	// search user by username or email
	if HasUser(user, username) == false {
		return
	}

	if VerifyPassword(password, user.Password) {
		// success
		success = true
	}
	return
}

// compare raw password and encoded password
func VerifyPassword(rawPwd, encodedPwd string) bool {

	// split
	var salt, encoded string
	if len(encodedPwd) > 11 {
		salt = encodedPwd[:10]
		encoded = encodedPwd[11:]
	}

	return utils.EncodePassword(rawPwd, salt) == encoded
}

// get user by erify code
func getVerifyUser(user *User, code string) bool {
	if len(code) <= utils.TimeLimitCodeLength {
		return false
	}

	// use tail hex username query user
	hexStr := code[utils.TimeLimitCodeLength:]
	if b, err := hex.DecodeString(hexStr); err == nil {
		user.UserName = string(b)
		if user.Read("UserName") == nil {
			return true
		}
	}

	return false
}

// verify active code when active account
func VerifyUserActiveCode(user *User, code string) bool {
	hours := utils.ActiveCodeLives

	if getVerifyUser(user, code) {
		// time limit code
		prefix := code[:utils.TimeLimitCodeLength]
		data := utils.ToStr(user.Id) + user.Email + user.UserName + user.Password + user.Rands

		return utils.VerifyTimeLimitCode(data, hours, prefix)
	}

	return false
}

// create a time limit code for user active
func CreateUserActiveCode(user *User, startInf interface{}) string {
	hours := utils.ActiveCodeLives
	data := utils.ToStr(user.Id) + user.Email + user.UserName + user.Password + user.Rands
	code := utils.CreateTimeLimitCode(data, hours, startInf)

	// add tail hex username
	code += hex.EncodeToString([]byte(user.UserName))
	return code
}

// verify code when reset password
func VerifyUserResetPwdCode(user *User, code string) bool {
	hours := utils.ResetPwdCodeLives

	if getVerifyUser(user, code) {
		// time limit code
		prefix := code[:utils.TimeLimitCodeLength]
		data := utils.ToStr(user.Id) + user.Email + user.UserName + user.Password + user.Rands + user.Updated.String()

		return utils.VerifyTimeLimitCode(data, hours, prefix)
	}

	return false
}

// create a time limit code for user reset password
func CreateUserResetPwdCode(user *User, startInf interface{}) string {
	hours := utils.ResetPwdCodeLives
	data := utils.ToStr(user.Id) + user.Email + user.UserName + user.Password + user.Rands + user.Updated.String()
	code := utils.CreateTimeLimitCode(data, hours, startInf)

	// add tail hex username
	code += hex.EncodeToString([]byte(user.UserName))
	return code
}
