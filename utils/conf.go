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

// Package utils implemented some useful functions.
package utils

import (
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/Unknwon/goconfig"
	"github.com/howeyc/fsnotify"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/orm"
	"github.com/beego/compress"
	"github.com/beego/i18n"

	"github.com/beego/wetalk/mailer"
)

const (
	APP_VER = "0.0.9.1027"
)

var (
	AppName             string
	AppDescription      string
	AppKeywords         string
	AppVer              string
	AppHost             string
	AppUrl              string
	AppLogo             string
	AvatarURL           string
	SecretKey           string
	IsProMode           bool
	MailUser            string
	MailFrom            string
	ActiveCodeLives     int
	ResetPwdCodeLives   int
	LoginRememberDays   int
	DateFormat          string
	DateTimeFormat      string
	DateTimeShortFormat string
	RealtimeRenderMD    bool
	ImageSizeSmall      int
	ImageSizeMiddle     int
	ImageLinkAlphabets  []byte
	ImageXSend          bool
	ImageXSendHeader    string
	Langs               []string
)

const (
	LangEnUS = iota
	LangZhCN
)

var (
	Cfg   *goconfig.ConfigFile
	Cache cache.Cache
)

var (
	AppConfPath      = "conf/app.ini"
	CompressConfPath = "conf/compress.json"
)

// LoadConfig loads configuration file.
func LoadConfig() *goconfig.ConfigFile {
	var err error

	if fh, _ := os.OpenFile(AppConfPath, os.O_RDONLY|os.O_CREATE, 0600); fh != nil {
		fh.Close()
	}

	// Load configuration, set app version and log level.
	Cfg, err = goconfig.LoadConfigFile(AppConfPath)
	Cfg.BlockMode = false
	if err != nil {
		panic("Fail to load configuration file: " + err.Error())
	}

	// Trim 4th part.
	AppVer = strings.Join(strings.Split(APP_VER, ".")[:3], ".")

	beego.RunMode = Cfg.MustValue("beego", "run_mode")

	beego.RunMode = Cfg.MustValue("beego", "run_mode")
	beego.HttpPort = Cfg.MustInt("beego", "http_port_"+beego.RunMode)

	IsProMode = beego.RunMode == "pro"
	if IsProMode {
		beego.SetLevel(beego.LevelInfo)
	}

	// cache system
	Cache, err = cache.NewCache("memory", `{"interval":360}`)

	// session settings
	beego.SessionOn = true
	beego.SessionProvider = Cfg.MustValue("app", "session_provider")
	beego.SessionSavePath = Cfg.MustValue("app", "session_path")
	beego.SessionName = Cfg.MustValue("app", "session_name")

	beego.EnableXSRF = true
	// xsrf token expire time
	beego.XSRFExpire = 86400 * 365

	driverName := Cfg.MustValue("orm", "driver_name")
	dataSource := Cfg.MustValue("orm", "data_source")
	maxIdle := Cfg.MustInt("orm", "max_idle_conn")
	maxOpen := Cfg.MustInt("orm", "max_open_conn")

	// set default database
	orm.RegisterDataBase("default", driverName, dataSource, maxIdle, maxOpen)
	orm.RunCommand()

	err = orm.RunSyncdb("default", false, false)
	if err != nil {
		beego.Error(err)
	}

	configWatcher()
	reloadConfig()

	settingLocales()
	settingCompress()

	return Cfg
}

func reloadConfig() {
	AppName = Cfg.MustValue("app", "app_name")
	beego.AppName = AppName

	AppHost = Cfg.MustValue("app", "app_host")
	AppUrl = Cfg.MustValue("app", "app_url")
	AppLogo = Cfg.MustValue("app", "app_logo")
	AppDescription = Cfg.MustValue("app", "description")
	AppKeywords = Cfg.MustValue("app", "keywords")
	AvatarURL = Cfg.MustValue("app", "avatar_url")
	DateFormat = Cfg.MustValue("app", "date_format")
	DateTimeFormat = Cfg.MustValue("app", "datetime_format")
	DateTimeShortFormat = Cfg.MustValue("app", "datetime_short_format")

	MailUser = Cfg.MustValue("app", "mail_user")
	MailFrom = Cfg.MustValue("app", "mail_from")

	SecretKey = Cfg.MustValue("app", "secret_key")
	ActiveCodeLives = Cfg.MustInt("app", "acitve_code_live_days")
	ResetPwdCodeLives = Cfg.MustInt("app", "resetpwd_code_live_days")
	LoginRememberDays = Cfg.MustInt("app", "login_remember_days")
	RealtimeRenderMD = Cfg.MustBool("app", "realtime_render_markdown")

	ImageSizeSmall = Cfg.MustInt("app", "image_size_small")
	ImageSizeMiddle = Cfg.MustInt("app", "image_size_middle")

	if ImageSizeSmall <= 0 {
		ImageSizeSmall = 300
	}

	if ImageSizeMiddle <= ImageSizeSmall {
		ImageSizeMiddle = ImageSizeSmall + 400
	}

	str := Cfg.MustValue("app", "image_link_alphabets")
	if len(str) == 0 {
		str = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	}
	ImageLinkAlphabets = []byte(str)

	ImageXSend = Cfg.MustBool("app", "image_xsend")
	ImageXSendHeader = Cfg.MustValue("app", "image_xsend_header")

	// set mailer connect args
	mailer.MailHost = Cfg.MustValue("mailer", "host")
	mailer.AuthUser = Cfg.MustValue("mailer", "user")
	mailer.AuthPass = Cfg.MustValue("mailer", "pass")

	orm.Debug = Cfg.MustBool("orm", "debug_log")
}

func settingLocales() {
	// load locales with locale_LANG.ini files
	langs := "en-US|zh-CN"
	for _, lang := range strings.Split(langs, "|") {
		lang = strings.TrimSpace(lang)
		if err := i18n.SetMessage(lang, "conf/"+"locale_"+lang+".ini"); err != nil {
			beego.Error("Fail to set message file: " + err.Error())
			os.Exit(2)
		}
	}
	Langs = i18n.ListLangs()
}

func settingCompress() {
	setting, err := compress.LoadJsonConf(CompressConfPath, IsProMode, AppUrl)
	if err != nil {
		beego.Error(err)
		return
	}

	setting.RunCommand()

	if IsProMode {
		setting.RunCompress(true, false, true)
	}

	beego.AddFuncMap("compress_js", setting.Js.CompressJs)
	beego.AddFuncMap("compress_css", setting.Css.CompressCss)
}

func configWatcher() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic("Failed start app watcher: " + err.Error())
	}

	go func() {
		for {
			select {
			case event := <-watcher.Event:
				switch filepath.Ext(event.Name) {
				case ".ini":
					beego.Info(event)

					if err := Cfg.Reload(); err != nil {
						beego.Error("Conf Reload: ", err)
					}

					if err := i18n.ReloadLangs(); err != nil {
						beego.Error("Conf Reload: ", err)
					}

					reloadConfig()
					beego.Info("Config Reloaded")

				case ".json":
					if event.Name == CompressConfPath {
						settingCompress()
						beego.Info("Beego Compress Reloaded")
					}
				}
			}
		}
	}()

	if err := watcher.WatchFlags("conf", fsnotify.FSN_MODIFY); err != nil {
		beego.Error(err)
	}
}

func IsMatchHost(uri string) bool {
	if len(uri) == 0 {
		return false
	}

	u, err := url.ParseRequestURI(uri)
	if err != nil {
		return false
	}

	if u.Host != AppHost {
		return false
	}

	return true
}
