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
	"github.com/beego/wetalk/utils"
	"time"

	"github.com/astaxie/beego/orm"
)

// post topic
type Topic struct {
	Id         int
	Name       string `orm:"size(30);unique"`
	Intro      string `orm:"type(text)"`
	IntroCache string `orm:"type(text)"`
	Slug       string `orm:"size(100);unique"`
	Followers  int
	Cat        *TopicCat `orm:"rel(one)"`
	Order      int
	Created    time.Time `orm:"auto_now_add"`
	Updated    time.Time `orm:"auto_now;index"`
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

func (m *Topic) String() string {
	return utils.ToStr(m.Id)
}

func Topics() orm.QuerySeter {
	return orm.NewOrm().QueryTable("topic").OrderBy("-Id")
}

// topic category
type TopicCat struct {
	Id    int
	Name  string `orm:"size(30);unique"`
	Slug  string `orm:"size(100);unique"`
	Order int
}

func (m *TopicCat) Insert() error {
	if _, err := orm.NewOrm().Insert(m); err != nil {
		return err
	}
	return nil
}

func (m *TopicCat) Read(fields ...string) error {
	if err := orm.NewOrm().Read(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *TopicCat) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *TopicCat) Delete() error {
	if _, err := orm.NewOrm().Delete(m); err != nil {
		return err
	}
	return nil
}

func (m *TopicCat) String() string {
	return utils.ToStr(m.Id)
}

func TopicCats() orm.QuerySeter {
	return orm.NewOrm().QueryTable("topic_cat").OrderBy("-Id")
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

func init() {
	orm.RegisterModel(new(Topic), new(TopicCat), new(FollowTopic))
}
