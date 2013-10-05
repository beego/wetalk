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
	"time"

	"github.com/astaxie/beego/orm"

	"github.com/beego/wetalk/utils"
)

// global settings name -> value
type Setting struct {
	Id      int
	Name    string `orm:"unique"`
	Value   string `orm:"type(text)"`
	Updated string `orm:"auto_now"`
}

// main user table
// IsAdmin: user is admininstator
// IsActive: set active when email is verified
// IsForbid: forbid user login
type User struct {
	Id        int
	UserName  string `orm:"size(30);unique"`
	NickName  string `orm:"size(30)"`
	Password  string `orm:"size(128)"`
	Url       string `orm:"size(100)"`
	Email     string `orm:"size(80);unique"`
	GrEmail   string `orm:"size(32)"`
	Info      string
	HideEmail bool
	Followers int
	Following int
	IsAdmin   bool      `orm:"index"`
	IsActive  bool      `orm:"index"`
	IsForbid  bool      `orm:"index"`
	Rands     string    `orm:"size(10)"`
	Created   time.Time `orm:"auto_now_add"`
	Updated   time.Time `orm:"auto_now"`
}

func (m *User) Insert() error {
	m.Rands = GetUserSalt()
	if _, err := orm.NewOrm().Insert(m); err != nil {
		return err
	}
	return nil
}

func (m *User) Read(fields ...string) error {
	if err := orm.NewOrm().Read(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *User) Update(fields ...string) error {
	fields = append(fields, "Updated")
	if _, err := orm.NewOrm().Update(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *User) Delete() error {
	if _, err := orm.NewOrm().Delete(m); err != nil {
		return err
	}
	return nil
}

func (m *User) String() string {
	return utils.ToStr(m.Id)
}

func Users() orm.QuerySeter {
	return orm.NewOrm().QueryTable("user").OrderBy("-Id")
}

// user follow
type Follow struct {
	Id         int
	User       *User `orm:"rel(fk)"`
	FollowUser *User `orm:"rel(fk)"`
	Mutual     bool
	Created    time.Time `orm:"auto_now_add"`
}

func (*Follow) TableUnique() [][]string {
	return [][]string{
		[]string{"User", "FollowUser"},
	}
}

func init() {
	orm.RegisterModel(new(Setting), new(User), new(Follow))
}
