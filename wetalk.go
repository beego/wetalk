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

// An open source project for Gopher community.
package main

import (
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/orm"
	"github.com/beego/i18n"

	"github.com/beego/wetalk/mailer"
	"github.com/beego/wetalk/routers"
	"github.com/beego/wetalk/utils"
)

const (
	APP_VER = "0.0.3.0913"
)

// We have to call a initialize function manully
// because we use `bee bale` to pack static resources
// and we cannot make sure that which init() execute first.
func initialize() {
	var err error
	// Load configuration, set app version and log level.
	utils.Cfg, err = utils.LoadConfig("conf/app.ini")
	if err != nil {
		panic("Fail to load configuration file: " + err.Error())
	}
	err = i18n.SetMessage("conf/message.ini")
	if err != nil {
		panic("Fail to set message file: " + err.Error())
	}

	// Trim 4th part.
	utils.AppVer = strings.Join(strings.Split(APP_VER, ".")[:3], ".")

	beego.AppName = utils.Cfg.MustValue("beego", "app_name")
	beego.RunMode = utils.Cfg.MustValue("beego", "run_mode")
	beego.HttpPort = utils.Cfg.MustInt("beego", "http_port_"+beego.RunMode)

	utils.AppName = beego.AppName
	utils.AppUrl = utils.Cfg.MustValue("app", "app_url")
	utils.AppDescription = utils.Cfg.MustValue("app", "description")
	utils.AppKeywords = utils.Cfg.MustValue("app", "keywords")
	utils.AppJsVer = utils.Cfg.MustValue("app", "js_ver")
	utils.AppCssVer = utils.Cfg.MustValue("app", "css_ver")
	utils.AvatarURL = utils.Cfg.MustValue("app", "avatar_url")

	utils.MailUser = utils.Cfg.MustValue("app", "mail_user")
	utils.MailFrom = utils.Cfg.MustValue("app", "mail_from")

	utils.SecretKey = utils.Cfg.MustValue("app", "secret_key")
	utils.ActiveCodeLives = utils.Cfg.MustInt("app", "acitve_code_live_days")
	utils.ResetPwdCodeLives = utils.Cfg.MustInt("app", "resetpwd_code_live_days")

	utils.IsBeta = utils.Cfg.MustBool("server", "beta")
	utils.IsProMode = beego.RunMode == "pro"
	if utils.IsProMode {
		beego.SetLevel(beego.LevelInfo)
		beego.Info("Product mode enabled")
		beego.Info(beego.AppName, APP_VER)
	}

	orm.Debug, _ = utils.Cfg.Bool("orm", "debug_log")

	driverName, _ := utils.Cfg.GetValue("orm", "driver_name")
	dataSource, _ := utils.Cfg.GetValue("orm", "data_source")
	maxIdle, _ := utils.Cfg.Int("orm", "max_idle_conn")

	// cache system
	utils.Cache, err = cache.NewCache("memory", `{"interval":360}`)

	// session settings
	beego.SessionOn = true
	beego.SessionProvider = utils.Cfg.MustValue("app", "session_provider")
	beego.SessionSavePath = utils.Cfg.MustValue("app", "session_path")
	beego.SessionName = utils.Cfg.MustValue("app", "session_name")

	beego.EnableXSRF = true
	// xsrf token expire time
	beego.XSRFExpire = 86400 * 365

	// set mailer connect args
	mailer.MailHost = utils.Cfg.MustValue("mailer", "host")
	mailer.AuthUser = utils.Cfg.MustValue("mailer", "user")
	mailer.AuthPass = utils.Cfg.MustValue("mailer", "pass")

	// set default database
	orm.RegisterDataBase("default", driverName, dataSource, maxIdle)

	orm.RunCommand()

	err = orm.RunSyncdb("default", false, false)
	if err != nil {
		beego.Error(err)
	}
}

func main() {
	initialize()

	beego.Info(beego.AppName, APP_VER)

	// Register routers.
	beego.Router("/", &routers.HomeRouter{})

	login := &routers.LoginRouter{}
	beego.Router("/login", login, "post:Login")
	beego.Router("/logout", login, "get:Logout")

	register := &routers.RegisterRouter{}
	beego.Router("/register", register, "post:Register")
	beego.Router("/active/success", register, "get:ActiveSuccess")
	beego.Router("/active/:code([0-9a-zA-Z]+)", register, "get:Active")

	settings := new(routers.SettingsRouter)
	beego.Router("/settings/profile", settings, "get:Profile;post:ProfileSave")

	fogot := &routers.ForgotRouter{}
	beego.Router("/forgot", fogot)
	beego.Router("/reset/:code([0-9a-zA-Z]+)", fogot, "get:Reset;post:ResetPost")

	// "robot.txt"
	beego.Router("/robot.txt", &routers.RobotRouter{})

	// For all unknown pages.
	beego.Run()
}
