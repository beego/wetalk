package models

import (
	"testing"

	. "github.com/beego/beebbs/utils"
)

var encoded = "9e2a6b0a670d48bc9fae7f79503a0d7e888650ede413810ff762ae767181fb3e7cdec433d435ed2c671bf4b6ecc49ebae5c7"

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
	user.Id = 5
	user.UserName = "beebbs"
	user.Email = "service@beego.me"
	user.Password = encoded
	user.Rands = GetRandomString(10)
	SecretKey = encoded
	ActiveCodeLives = 1
	ResetPwdCodeLives = 1
	code := CreateUserActiveCode(user, nil)
	ThrowFail(t, AssertIs(VerifyUserActiveCode(user, code), true))
	user.Rands = GetRandomString(10)
	ThrowFail(t, AssertIs(VerifyUserActiveCode(user, code), false))

	code = CreateUserResetPwdCode(user, nil)
	ThrowFail(t, AssertIs(VerifyUserResetPwdCode(user, code), true))
	user.Rands = GetRandomString(10)
	ThrowFail(t, AssertIs(VerifyUserResetPwdCode(user, code), false))
}
