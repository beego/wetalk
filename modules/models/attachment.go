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
	"fmt"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"

	"github.com/beego/wetalk/modules/utils"
	"github.com/beego/wetalk/setting"
)

type Image struct {
	Id      int
	User    *User `orm:"rel(fk)"`
	Width   int
	Height  int
	Ext     int `orm:"index"`
	Created time.Time
}

func (m *Image) Insert() error {
	if _, err := orm.NewOrm().Insert(m); err != nil {
		return err
	}
	return nil
}

func (m *Image) Read(fields ...string) error {
	if err := orm.NewOrm().Read(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *Image) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *Image) Delete() error {
	if _, err := orm.NewOrm().Delete(m); err != nil {
		return err
	}
	return nil
}

func (m *Image) LinkFull() string {
	return m.LinkSize(0)
}

func (m *Image) LinkSmall() string {
	var width int
	switch {
	case m.Width > setting.ImageSizeSmall:
		width = setting.ImageSizeSmall
	}
	return m.LinkSize(width)
}

func (m *Image) LinkMiddle() string {
	var width int
	switch {
	case m.Width > setting.ImageSizeMiddle:
		width = setting.ImageSizeMiddle
	}
	return m.LinkSize(width)
}

func (m *Image) LinkSize(width int) string {
	if m.Ext == 3 {
		// if image is gif then return full size
		width = 0
	}
	var size string
	switch width {
	case setting.ImageSizeSmall, setting.ImageSizeMiddle:
		size = utils.ToStr(width)
	default:
		size = "full"
	}
	return "/img/" + m.GetToken() + "." + size + m.GetExt()
}

func (m *Image) GetExt() string {
	var ext string
	switch m.Ext {
	case 1:
		ext = ".jpg"
	case 2:
		ext = ".png"
	case 3:
		ext = ".gif"
	}
	return ext
}

func (m *Image) GetToken() string {
	number := beego.Date(m.Created, "ymds") + utils.ToStr(m.Id)
	return utils.NumberEncode(number, setting.ImageLinkAlphabets)
}

func (m *Image) DecodeToken(token string) error {
	number := utils.NumberDecode(token, setting.ImageLinkAlphabets)
	if len(number) < 9 {
		return fmt.Errorf("token `%s` too short <- `%s`", token, number)
	}

	if t, err := beego.DateParse(number[:8], "ymds"); err != nil {
		return fmt.Errorf("token `%s` date parse error <- `%s`", token, number)
	} else {
		m.Created = t
	}

	var err error
	m.Id, err = utils.StrTo(number[8:]).Int()
	if err != nil {
		return fmt.Errorf("token `%s` id parse error <- `%s`", token, err)
	}

	return nil
}

func init() {
	orm.RegisterModel(new(Image))
}
