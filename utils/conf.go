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

	"github.com/Unknwon/com"
	"github.com/Unknwon/goconfig"
	"github.com/astaxie/beego/cache"
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
	IsBeta            bool
	MailUser          string
	MailFrom          string
	ActiveCodeLives   int
	ResetPwdCodeLives int
	LoginRememberDays int
	DateFormat        string
	DateTimeFormat    string
)

var (
	Cfg   *goconfig.ConfigFile
	Cache cache.Cache
)

// LoadConfig loads configuration file.
func LoadConfig(cfgPath string) (*goconfig.ConfigFile, error) {
	if !com.IsExist(cfgPath) {
		os.Create(cfgPath)
	}

	return goconfig.LoadConfigFile(cfgPath)
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
