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

package attachment

import (
	"fmt"
	"github.com/beego/wetalk/setting"
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

	"github.com/beego/wetalk/modules/models"
	"github.com/beego/wetalk/modules/utils"
)

func SaveImage(m *models.Image, r io.ReadSeeker, mime string, filename string, created time.Time) error {
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

		if m.Width > setting.ImageSizeSmall {
			if err := ImageResize(m, img, setting.ImageSizeSmall); err != nil {
				os.RemoveAll(fullPath)
				return err
			}
		}

		if m.Width > setting.ImageSizeMiddle {
			if err := ImageResize(m, img, setting.ImageSizeMiddle); err != nil {
				os.RemoveAll(fullPath)
				return err
			}
		}

	}

	return nil
}

func ImageResize(img *models.Image, im image.Image, width int) error {
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

func GenImagePath(img *models.Image) string {
	return "upload/img/" + beego.Date(img.Created, "y/m/d/s/") + utils.ToStr(img.Id) + "/"
}

func GenImageFilePath(img *models.Image, width int) string {
	var size string
	if width == 0 {
		size = "full"
	} else {
		size = utils.ToStr(width)
	}
	return GenImagePath(img) + size + img.GetExt()
}
