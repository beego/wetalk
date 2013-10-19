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

	"github.com/beego/wetalk/utils"
)

// post content
type Post struct {
	Id           int
	User         *User  `orm:"rel(fk)"`
	Title        string `orm:"size(60)"`
	Content      string `orm:"type(text)"`
	ContentCache string `orm:"type(text)"`
	Browsers     int
	Replys       int
	Favorites    int
	LastReply    *User  `orm:"rel(fk)"`
	Topic        *Topic `orm:"rel(fk)"`
	IsBest       bool
	Category     *Category `orm:"rel(fk)"`
	Created      time.Time `orm:"auto_now_add"`
	Updated      time.Time `orm:"auto_now;index"`
}

func (m *Post) Insert() error {
	if _, err := orm.NewOrm().Insert(m); err != nil {
		return err
	}
	return nil
}

func (m *Post) Read(fields ...string) error {
	if err := orm.NewOrm().Read(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *Post) Update(fields ...string) error {
	fields = append(fields, "Updated")
	if _, err := orm.NewOrm().Update(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *Post) Delete() error {
	if _, err := orm.NewOrm().Delete(m); err != nil {
		return err
	}
	return nil
}

func (m *Post) String() string {
	return utils.ToStr(m.Id)
}

func (m *Post) Link() string {
	return fmt.Sprintf("%sp/%d", utils.AppUrl, m.Id)
}

func (m *Post) GetContentCache() string {
	if utils.RealtimeRenderMD {
		return RenderPostContent(m.Content)
	} else {
		return m.ContentCache
	}
}

func Posts() orm.QuerySeter {
	return orm.NewOrm().QueryTable("post").OrderBy("-Id")
}

// commnet content for post
type Comment struct {
	Id           int
	User         *User  `orm:"rel(fk)"`
	Post         *Post  `orm:"rel(fk)"`
	Message      string `orm:"type(text)"`
	MessageCache string `orm:"type(text)"`
	Status       int
	Created      time.Time `orm:"auto_now_add;index"`
}

func (m *Comment) Insert() error {
	if _, err := orm.NewOrm().Insert(m); err != nil {
		return err
	}
	return nil
}

func (m *Comment) Read(fields ...string) error {
	if err := orm.NewOrm().Read(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *Comment) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *Comment) Delete() error {
	if _, err := orm.NewOrm().Delete(m); err != nil {
		return err
	}
	return nil
}

func (m *Comment) String() string {
	return utils.ToStr(m.Id)
}

func Comments() orm.QuerySeter {
	return orm.NewOrm().QueryTable("comment").OrderBy("-Id")
}

// user favorite posts
type FavoritePost struct {
	Id      int
	User    *User     `orm:"rel(fk)"`
	Post    *Post     `orm:"rel(fk)"`
	Created time.Time `orm:"auto_now_add"`
}

func (*FavoritePost) TableUnique() [][]string {
	return [][]string{
		[]string{"User", "Post"},
	}
}

func init() {
	orm.RegisterModel(new(Post), new(FavoritePost), new(Comment))
}
