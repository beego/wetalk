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

// Package models implemented database access funtions.
package models

import (
	"time"

	"github.com/astaxie/beego/orm"
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

func (u *User) Read(fields ...string) error {
	if err := orm.NewOrm().Read(u, fields...); err != nil {
		return err
	}
	return nil
}

func (u *User) Update(fields ...string) error {
	fields = append(fields, "Updated")
	if _, err := orm.NewOrm().Update(u, fields...); err != nil {
		return err
	}
	return nil
}

// NewUser saves 'User' into database.
func NewUser(u *User) error {
	u.Rands = GetUserSalt()
	_, err := orm.NewOrm().Insert(u)
	if err != nil {
		return err
	}
	return nil
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

// post content
type Post struct {
	Id           int
	User         *User  `orm:"rel(fk)"`
	Slug         string `orm:"size(100);unique"`
	Title        string `orm:"size(100)"`
	Content      string `orm:"type(text)"`
	ContentCache string `orm:"type(text)"`
	Browsers     int
	Replys       int
	Favorites    int
	LastReply    *User     `orm:"rel(fk)"`
	Created      time.Time `orm:"auto_now_add"`
	Updated      time.Time `orm:"auto_now;index"`
}

// post topic
type Topic struct {
	Id         int
	Name       string `orm:"size(30);unique"`
	Intro      string `orm:"type(text)"`
	IntroCache string `orm:"type(text)"`
	Slug       string `orm:"size(100);unique"`
	Followers  int
	Cat        *TopicCat `orm:"rel(one)"`
	Created    time.Time `orm:"auto_now_add"`
	Updated    time.Time `orm:"auto_now;index"`
}

// topic category
type TopicCat struct {
	Id    int
	Name  string
	Slug  string `orm:"size(100);unique"`
	Order int
}

// commnet content for post
type Comment struct {
	Id           int
	Post         *Post    `orm:"rel(fk)"`
	Parent       *Comment `orm:"rel(fk)"`
	Message      string   `orm:"type(text)"`
	MessageCache string   `orm:"type(text)"`
	Status       int
	Created      time.Time `orm:"auto_now_add;index"`
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
	orm.RegisterModel(new(Setting), new(User), new(Follow))
	orm.RegisterModel(new(Post), new(FavoritePost), new(Comment))
	orm.RegisterModel(new(Topic), new(TopicCat), new(FollowTopic))
}
