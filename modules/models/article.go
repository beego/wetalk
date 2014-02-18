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

	"github.com/astaxie/beego/orm"
	"github.com/beego/i18n"

	"github.com/beego/wetalk/modules/utils"
	"github.com/beego/wetalk/setting"
)

type Article struct {
	Id               int
	User             *User     `orm:"rel(fk)"`
	Uri              string    `orm:"size(60);unqiue"`
	Title            string    `orm:"size(60)"`
	Content          string    `orm:"type(text)"`
	ContentCache     string    `orm:"type(text)"`
	TitleZhCn        string    `orm:"size(60)"`
	ContentZhCn      string    `orm:"type(text)"`
	ContentCacheZhCn string    `orm:"type(text)"`
	LastAuthor       *User     `orm:"rel(fk);null"`
	IsPublish        bool      `orm:"index"`
	Created          time.Time `orm:"auto_now_add"`
	Updated          time.Time `orm:"auto_now"`
}

func (m *Article) Insert() error {
	if _, err := orm.NewOrm().Insert(m); err != nil {
		return err
	}
	return nil
}

func (m *Article) Read(fields ...string) error {
	if err := orm.NewOrm().Read(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *Article) Update(fields ...string) error {
	fields = append(fields, "Updated")
	if _, err := orm.NewOrm().Update(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *Article) Delete() error {
	if _, err := orm.NewOrm().Delete(m); err != nil {
		return err
	}
	return nil
}

func (m *Article) String() string {
	return utils.ToStr(m.Id)
}

func (m *Article) Link() string {
	uri := m.Uri
	if len(uri) > 0 && uri[0] == '/' {
		uri = uri[1:]
	}
	return fmt.Sprintf("%s%s", setting.AppUrl, uri)
}

func (m *Article) GetTitle(lang string) string {
	var title string
	switch i18n.IndexLang(lang) {
	case setting.LangZhCN:
		title = m.TitleZhCn
	default:
		title = m.Title
	}
	return title
}

func (m *Article) GetContentCache(lang string) string {
	var content, contentCache string
	switch i18n.IndexLang(lang) {
	case setting.LangZhCN:
		content = m.ContentZhCn
		contentCache = m.ContentCacheZhCn
	default:
		content = m.Content
		contentCache = m.ContentCache
	}
	if setting.RealtimeRenderMD {
		return utils.RenderMarkdown(content)
	} else {
		return contentCache
	}
}

func Articles() orm.QuerySeter {
	return orm.NewOrm().QueryTable("article").OrderBy("-Id")
}

func init() {
	orm.RegisterModel(new(Article))
}
