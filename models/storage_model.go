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
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/nfnt/resize"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"

	"github.com/beego/wetalk/utils"
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
	case m.Width > utils.ImageSizeSmall:
		width = utils.ImageSizeSmall
	}
	return m.LinkSize(width)
}

func (m *Image) LinkMiddle() string {
	var width int
	switch {
	case m.Width > utils.ImageSizeMiddle:
		width = utils.ImageSizeMiddle
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
	case utils.ImageSizeSmall, utils.ImageSizeMiddle:
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
	return utils.NumberEncode(number, utils.ImageLinkAlphabets)
}

func (m *Image) DecodeToken(token string) error {
	number := utils.NumberDecode(token, utils.ImageLinkAlphabets)
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

func SaveImage(m *Image, r io.ReadSeeker, mime string, filename string, created time.Time) error {
	var ext string

	// test image mime type
	switch mime {
	case "image/jpeg":
		ext = ".jpg"

	case "image/png":
		ext = ".png"

	case "image/gif":
		ext = ".gif"

	default:
		ext = filepath.Ext(filename)
		switch ext {
		case ".jpg", ".png", ".gif":
		default:
			return fmt.Errorf("unsupport image format `%s`", filename)
		}
	}

	// decode image
	var img image.Image
	var err error
	switch ext {
	case ".jpg":
		m.Ext = 1
		img, err = jpeg.Decode(r)
	case ".png":
		m.Ext = 2
		img, err = png.Decode(r)
	case ".gif":
		m.Ext = 3
		img, err = gif.Decode(r)
	}

	if err != nil {
		return err
	}

	m.Width = img.Bounds().Dx()
	m.Height = img.Bounds().Dy()
	m.Created = created

	if err := m.Insert(); err != nil || m.Id <= 0 {
		return err
	}

	path := GenImagePath(m)
	os.MkdirAll(path, 0755)

	fullPath := GenImageFilePath(m, 0)
	if _, err := r.Seek(0, 0); err != nil {
		return err
	}

	var file *os.File
	if f, err := os.OpenFile(fullPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644); err != nil {
		return err
	} else {
		file = f
	}

	if _, err := io.Copy(file, r); err != nil {
		os.RemoveAll(fullPath)
		return err
	}

	if ext != ".gif" {

		if m.Width > utils.ImageSizeSmall {
			if err := ImageResize(m, img, utils.ImageSizeSmall); err != nil {
				os.RemoveAll(fullPath)
				return err
			}
		}

		if m.Width > utils.ImageSizeMiddle {
			if err := ImageResize(m, img, utils.ImageSizeMiddle); err != nil {
				os.RemoveAll(fullPath)
				return err
			}
		}

	}

	return nil
}

func ImageResize(img *Image, im image.Image, width int) error {
	savePath := GenImageFilePath(img, width)
	im = resize.Resize(uint(width), 0, im, resize.Bilinear)

	var file *os.File
	if f, err := os.OpenFile(savePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644); err != nil {
		return err
	} else {
		file = f
	}
	defer file.Close()

	var err error
	switch img.Ext {
	case 1:
		err = jpeg.Encode(file, im, &jpeg.Options{90})
	case 2:
		err = png.Encode(file, im)
	default:
		return fmt.Errorf("<ImageResize> unsupport image format")
	}

	return err
}

func GenImagePath(img *Image) string {
	return "upload/img/" + beego.Date(img.Created, "y/m/d/s/") + utils.ToStr(img.Id) + "/"
}

func GenImageFilePath(img *Image, width int) string {
	var size string
	if width == 0 {
		size = "full"
	} else {
		size = utils.ToStr(width)
	}
	return GenImagePath(img) + size + img.GetExt()
}
