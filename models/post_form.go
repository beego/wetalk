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
	"github.com/astaxie/beego/validation"

	"github.com/beego/wetalk/utils"
)

type PostForm struct {
	Category   int        `form:"type(select);attr(rel,select2)" valid:"Required"`
	Topic      int        `form:"type(select);attr(rel,select2)" valid:"Required"`
	Title      string     `form:"attr(autocomplete,off)" valid:"Required;MinSize(5);MaxSize(60)"`
	Content    string     `form:"type(textarea)" valid:"Required;MinSize(10)"`
	Categories []Category `form:"-"`
	Topics     []Topic    `form:"-"`
}

func (form *PostForm) CategorySelectData() [][]string {
	data := make([][]string, 0, len(form.Categories)+1)
	for _, cat := range form.Categories {
		data = append(data, []string{"category." + cat.Name, utils.ToStr(cat.Id)})
	}
	return data
}

func (form *PostForm) TopicSelectData() [][]string {
	data := make([][]string, 0, len(form.Topics)+1)
	for _, topic := range form.Topics {
		data = append(data, []string{"topic." + topic.Name, utils.ToStr(topic.Id)})
	}
	return data
}

func (form *PostForm) Valid(v *validation.Validation) {
	valid := false
	for _, topic := range form.Topics {
		if topic.Id == form.Topic {
			valid = true
		}
	}

	if !valid {
		v.SetError("Topic", "error")
	}

	valid = false
	for _, cat := range form.Categories {
		if cat.Id == form.Category {
			valid = true
		}
	}

	if !valid {
		v.SetError("Category", "error")
	}
}

func (form *PostForm) SavePost(post *Post, user *User) error {
	post.Title = form.Title
	post.Content = form.Content
	post.Category = &Category{Id: form.Category}
	post.Topic = &Topic{Id: form.Topic}
	post.User = user
	post.LastReply = user
	post.ContentCache = RenderPostContent(form.Content)
	return post.Insert()
}

func (form *PostForm) SetFromPost(post *Post) {
	form.Title = post.Title
	form.Content = post.Content
	form.Category = post.Category.Id
	form.Topic = post.Topic.Id
}

func (form *PostForm) UpdatePost(post *Post, user *User) error {
	changes := utils.FormChanges(post, form)
	if len(changes) == 0 {
		return nil
	}
	post.Title = form.Title
	post.Content = form.Content
	post.Category = &Category{Id: form.Category}
	post.Topic = &Topic{Id: form.Topic}
	post.ContentCache = RenderPostContent(form.Content)
	for _, c := range changes {
		if c == "Content" {
			changes = append(changes, "ContentCache")
		}
	}
	fmt.Println(changes)
	return post.Update(changes...)
}

func (form *PostForm) Placeholders() map[string]string {
	return map[string]string{
		"Category": "model.category_choose_dot",
		"Topic":    "model.topic_choose_dot",
		"Title":    "post.plz_enter_title",
	}
}

type PostAdminForm struct {
	Create    bool   `form:"-"`
	User      int    `valid:"Required"`
	Title     string `valid:"Required;MaxSize(60)"`
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
		v.SetError("User", "admin.not_found_by_id")
	}

	user.Id = form.LastReply
	if user.Read() != nil {
		v.SetError("LastReply", "admin.not_found_by_id")
	}

	topic := Topic{Id: form.Topic}
	if topic.Read() != nil {
		v.SetError("Topic", "admin.not_found_by_id")
	}

	cat := Category{Id: form.Category}
	if cat.Read() != nil {
		v.SetError("Category", "admin.not_found_by_id")
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

	post.ContentCache = RenderPostContent(post.Content)
}

type CommentForm struct {
	Message string `form:"type(textarea)" valid:"Required;MinSize(5)"`
}

func (form *CommentForm) SaveComment(comment *Comment, user *User, post *Post) error {
	comment.Message = form.Message
	comment.MessageCache = RenderPostContent(form.Message)
	comment.User = user
	comment.Post = post
	return comment.Insert()
}

type CommentAdminForm struct {
	Create  bool   `form:"-"`
	User    int    `valid:"Required"`
	Post    int    `valid:"Required"`
	Message string `form:"type(textarea)" valid:"Required"`
	Status  int8   `valid:"Required"`
}

func (form *CommentAdminForm) Valid(v *validation.Validation) {
	user := User{Id: form.User}
	if user.Read() != nil {
		v.SetError("User", "admin.not_found_by_id")
	}

	post := Post{Id: form.Post}
	if post.Read() != nil {
		v.SetError("Post", "admin.not_found_by_id")
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

	comment.MessageCache = RenderPostContent(comment.Message)
}
