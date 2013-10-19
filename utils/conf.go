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
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Unknwon/goconfig"
	"github.com/howeyc/fsnotify"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/orm"
	"github.com/beego/i18n"

	"github.com/beego/wetalk/mailer"
)

const (
	APP_VER = "0.0.2.0911"
)

var (
	AppName           string
	AppDescription    string
	AppKeywords       string
	AppVer            string
	AppHost           string
	AppUrl            string
	AppLogo           string
	AppJsVer          string
	AppCssVer         string
	AvatarURL         string
	SecretKey         string
	IsProMode         bool
	MailUser          string
	MailFrom          string
	ActiveCodeLives   int
	ResetPwdCodeLives int
	LoginRememberDays int
	DateFormat        string
	DateTimeFormat    string
	RealtimeRenderMD  bool
)

var (
	Cfg   *goconfig.ConfigFile
	Cache cache.Cache
)

// LoadConfig loads configuration file.
func LoadConfig() *goconfig.ConfigFile {
	var err error

	cfgPath := "conf/app.ini"

	if fh, _ := os.OpenFile(cfgPath, os.O_RDONLY|os.O_CREATE, 0600); fh != nil {
		fh.Close()
	}

	// Load configuration, set app version and log level.
	Cfg, err = goconfig.LoadConfigFile(cfgPath)
	Cfg.BlockMode = false
	if err != nil {
		panic("Fail to load configuration file: " + err.Error())
	}

	dirs, _ := ioutil.ReadDir("conf")
	for _, info := range dirs {
		if !info.IsDir() {
			name := info.Name()
			if filepath.HasPrefix(name, "locale_") {
				if filepath.Ext(name) == ".ini" {
					lang := name[7 : len(name)-4]
					if len(lang) > 0 {
						if err := i18n.SetMessage(lang, "conf/"+name); err != nil {
							panic("Fail to set message file: " + err.Error())
						}
						continue
					}
				}
				beego.Error("locale ", name, " not loaded")
			}
		}
	}

	// Trim 4th part.
	AppVer = strings.Join(strings.Split(APP_VER, ".")[:3], ".")

	beego.RunMode = Cfg.MustValue("beego", "run_mode")
	beego.HttpPort = Cfg.MustInt("beego", "http_port_"+beego.RunMode)

	ver := ToStr(time.Now().Unix())
	AppJsVer = ver
	AppCssVer = ver

	spawnWatcher()
	reloadConfig()

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

	MailUser = Cfg.MustValue("app", "mail_user")
	MailFrom = Cfg.MustValue("app", "mail_from")

	SecretKey = Cfg.MustValue("app", "secret_key")
	ActiveCodeLives = Cfg.MustInt("app", "acitve_code_live_days")
	ResetPwdCodeLives = Cfg.MustInt("app", "resetpwd_code_live_days")
	LoginRememberDays = Cfg.MustInt("app", "login_remember_days")
	RealtimeRenderMD = Cfg.MustBool("app", "realtime_render_markdown")

	// set mailer connect args
	mailer.MailHost = Cfg.MustValue("mailer", "host")
	mailer.AuthUser = Cfg.MustValue("mailer", "user")
	mailer.AuthPass = Cfg.MustValue("mailer", "pass")

	IsProMode = beego.RunMode == "pro"
	if IsProMode {
		beego.SetLevel(beego.LevelInfo)
		beego.Info("Product mode enabled")
		beego.Info(beego.AppName, APP_VER)
	}

	orm.Debug, _ = Cfg.Bool("orm", "debug_log")
}

func spawnWatcher() {
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

				case ".css":
					beego.Info(event)
					ver := ToStr(time.Now().Unix())
					AppCssVer = ver

				case ".js":
					beego.Info(event)
					ver := ToStr(time.Now().Unix())
					AppJsVer = ver
				}
			}
		}
	}()

	if err := watcher.WatchFlags("conf", fsnotify.FSN_MODIFY); err != nil {
		beego.Error(err)
	}

	if err := watcher.WatchFlags("static", fsnotify.FSN_MODIFY); err != nil {
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
