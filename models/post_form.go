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
	"github.com/astaxie/beego/validation"

	"github.com/beego/wetalk/utils"
)

type PostAdminForm struct {
	Create    bool   `form:"-"`
	User      int    `valid:"Required"`
	Title     string `valid:"Required;MaxSize(100)"`
	Content   string `form:"type(textarea)" valid:"Required"`
	Browsers  int    ``
	Replys    int    ``
	Favorites int    ``
	LastReply int    `valid:"Required"`
	Topic     int    `valid:"Required"`
	Category  int    `valid:"Required"`
	IsBest    bool
}

func (form *PostAdminForm) Valid(v *validation.Validation) {
	user := User{Id: form.User}
	if user.Read() != nil {
		v.SetError("User", "Not found by this id")
	}

	user.Id = form.LastReply
	if user.Read() != nil {
		v.SetError("LastReply", "Not found by this id")
	}

	topic := Topic{Id: form.Topic}
	if topic.Read() != nil {
		v.SetError("Topic", "Not found by this id")
	}

	cat := Category{Id: form.Category}
	if cat.Read() != nil {
		v.SetError("Category", "Not found by this id")
	}
}

func (form *PostAdminForm) SetFromPost(post *Post) {
	utils.SetFormValues(post, form)

	if post.User != nil {
		form.User = post.User.Id
	}

	if post.LastReply != nil {
		form.LastReply = post.LastReply.Id
	}

	if post.Topic != nil {
		form.Topic = post.Topic.Id
	}

	if post.Category != nil {
		form.Category = post.Category.Id
	}
}

func (form *PostAdminForm) SetToPost(post *Post) {
	utils.SetFormValues(form, post)

	if post.User == nil {
		post.User = &User{}
	}
	post.User.Id = form.User

	if post.LastReply == nil {
		post.LastReply = &User{}
	}
	post.LastReply.Id = form.LastReply

	if post.Topic == nil {
		post.Topic = &Topic{}
	}
	post.Topic.Id = form.Topic

	if post.Category == nil {
		post.Category = &Category{}
	}
	post.Category.Id = form.Category

	// TODO make ContentCache
}

type CommentAdminForm struct {
	Create  bool   `form:"-"`
	User    int    `valid:"Required"`
	Post    int    `valid:"Required"`
	Message string `form:"type(textarea)" valid:"Required"`
	Status  int    `valid:"Required"`
}

func (form *CommentAdminForm) Valid(v *validation.Validation) {
	user := User{Id: form.User}
	if user.Read() != nil {
		v.SetError("User", "Not found by this id")
	}

	post := Post{Id: form.Post}
	if post.Read() != nil {
		v.SetError("Post", "Not found by this id")
	}
}

func (form *CommentAdminForm) SetFromComment(comment *Comment) {
	utils.SetFormValues(comment, form)

	if comment.User != nil {
		form.User = comment.User.Id
	}

	if comment.Post != nil {
		form.Post = comment.Post.Id
	}
}

func (form *CommentAdminForm) SetToComment(comment *Comment) {
	utils.SetFormValues(form, comment)

	if comment.User == nil {
		comment.User = &User{}
	}
	comment.User.Id = form.User

	if comment.Post == nil {
		comment.Post = &Post{}
	}
	comment.Post.Id = form.Post

	// TODO make MessageCache
}
