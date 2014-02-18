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
	"testing"

	"github.com/astaxie/beego/orm"
	. "github.com/beego/wetalk/modules/utils"
)

var encoded = "9e2a6b0a670d48bc9fae7f79503a0d7e888650ede413810ff762ae767181fb3e7cdec433d435ed2c671bf4b6ecc49ebae5c7"

func init() {
	orm.RegisterDataBase("default", "mysql", "root:root@/wetalk?charset=utf8", 30)
	orm.RunSyncdb("default", true, false)
}

func TestPasswordVerify(t *testing.T) {
	pwd := "111111"
	salt := "010101"

	ThrowFailNow(t, AssertIs(EncodePassword(pwd, salt), encoded))

	ThrowFail(t, AssertIs(VerifyPassword(pwd, salt, encoded), true))
	ThrowFail(t, AssertIs(VerifyPassword(pwd, "fake", encoded), false))
	ThrowFail(t, AssertIs(VerifyPassword("fake", salt, encoded), false))
}

func TestUserVerifyCode(t *testing.T) {
	user := new(User)
	user.UserName = "wetalk"
	user.Email = "service@beego.me"
	user.Password = encoded
	user.Rands = GetRandomString(10)
	SecretKey = encoded
	ActiveCodeLives = 1
	ResetPwdCodeLives = 1

	ThrowFail(t, NewUser(user))

	code := CreateUserActiveCode(user, nil)
	ThrowFail(t, AssertIs(VerifyUserActiveCode(user, code), true))
	ThrowFail(t, AssertIs(VerifyUserActiveCode(user, code+"1"), false))

	code = CreateUserResetPwdCode(user, nil)
	ThrowFail(t, AssertIs(VerifyUserResetPwdCode(user, code), true))
	ThrowFail(t, AssertIs(VerifyUserResetPwdCode(user, code+"1"), false))
}
