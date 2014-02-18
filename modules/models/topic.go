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

// post topic
type Topic struct {
	Id        int
	Name      string    `orm:"size(30);unique"`
	Intro     string    `orm:"type(text)"`
	NameZhCn  string    `orm:"size(30);unique"`
	IntroZhCn string    `orm:"type(text)"`
	Image     *Image    `orm:"rel(one);null"`
	Slug      string    `orm:"size(100);unique"`
	Followers int       `orm:"index"`
	Order     int       `orm:"index"`
	Created   time.Time `orm:"auto_now_add"`
	Updated   time.Time `orm:"auto_now;index"`
}

func (m *Topic) Insert() error {
	if _, err := orm.NewOrm().Insert(m); err != nil {
		return err
	}
	return nil
}

func (m *Topic) Read(fields ...string) error {
	if err := orm.NewOrm().Read(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *Topic) Update(fields ...string) error {
	fields = append(fields, "Updated")
	if _, err := orm.NewOrm().Update(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *Topic) Delete() error {
	if _, err := orm.NewOrm().Delete(m); err != nil {
		return err
	}
	return nil
}

func (m *Topic) RefreshFollowers() int {
	cnt, err := FollowTopics().Filter("Topic", m.Id).Count()
	if err == nil {
		m.Followers = int(cnt)
		m.Update("Followers")
	}
	return m.Followers
}

func (m *Topic) String() string {
	return utils.ToStr(m.Id)
}

func (m *Topic) Link() string {
	return fmt.Sprintf("%stopic/%s", setting.AppUrl, m.Slug)
}

func (m *Topic) GetName(lang string) string {
	var name string
	switch i18n.IndexLang(lang) {
	case setting.LangZhCN:
		name = m.NameZhCn
	default:
		name = m.Name
	}
	return name
}

func (m *Topic) GetIntro(lang string) string {
	var intro string
	switch i18n.IndexLang(lang) {
	case setting.LangZhCN:
		intro = m.IntroZhCn
	default:
		intro = m.Intro
	}
	return intro
}

func Topics() orm.QuerySeter {
	return orm.NewOrm().QueryTable("topic").OrderBy("-Id")
}

// topic category
type Category struct {
	Id    int
	Name  string `orm:"size(30);unique"`
	Slug  string `orm:"size(100);unique"`
	Order int    `orm:"index"`
}

func (m *Category) Insert() error {
	if _, err := orm.NewOrm().Insert(m); err != nil {
		return err
	}
	return nil
}

func (m *Category) Read(fields ...string) error {
	if err := orm.NewOrm().Read(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *Category) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *Category) Delete() error {
	if _, err := orm.NewOrm().Delete(m); err != nil {
		return err
	}
	return nil
}

func (m *Category) String() string {
	return utils.ToStr(m.Id)
}

func (m *Category) Link() string {
	return fmt.Sprintf("%scategory/%s", setting.AppUrl, m.Slug)
}

func Categories() orm.QuerySeter {
	return orm.NewOrm().QueryTable("category").OrderBy("-Id")
}

// user follow topics
type FollowTopic struct {
	Id      int
	User    *User     `orm:"rel(fk)"`
	Topic   *Topic    `orm:"rel(fk)"`
	Created time.Time `orm:"auto_now_add"`
}

func (*FollowTopic) TableUnique() [][]string {
	return [][]string{
		[]string{"User", "Topic"},
	}
}

func (m *FollowTopic) Insert() error {
	if _, err := orm.NewOrm().Insert(m); err != nil {
		return err
	}
	return nil
}

func (m *FollowTopic) Read(fields ...string) error {
	if err := orm.NewOrm().Read(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *FollowTopic) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *FollowTopic) Delete() error {
	if _, err := orm.NewOrm().Delete(m); err != nil {
		return err
	}
	return nil
}

func (m *FollowTopic) String() string {
	return utils.ToStr(m.Id)
}

func FollowTopics() orm.QuerySeter {
	return orm.NewOrm().QueryTable("follow_topic").OrderBy("-Id")
}

func init() {
	orm.RegisterModel(new(Topic), new(Category), new(FollowTopic))
}
