package models

import (
	"testing"

	"github.com/astaxie/beego/orm"
	. "github.com/beego/wetalk/utils"
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
