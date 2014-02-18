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

package post

import (
	"github.com/astaxie/beego/validation"
	"github.com/beego/i18n"

	"github.com/beego/wetalk/modules/models"
	"github.com/beego/wetalk/modules/utils"
	"github.com/beego/wetalk/setting"
)

type PostForm struct {
	Lang       int               `form:"type(select);attr(rel,select2)"`
	Category   int               `form:"type(select);attr(rel,select2)" valid:"Required"`
	Topic      int               `form:"type(select);attr(rel,select2)" valid:"Required"`
	Title      string            `form:"attr(autocomplete,off)" valid:"Required;MinSize(5);MaxSize(60)"`
	Content    string            `form:"type(textarea)" valid:"Required;MinSize(10)"`
	Categories []models.Category `form:"-"`
	Topics     []models.Topic    `form:"-"`
	Locale     i18n.Locale       `form:"-"`
}

func (form *PostForm) LangSelectData() [][]string {
	langs := setting.Langs
	data := make([][]string, 0, len(langs))
	for i, lang := range langs {
		data = append(data, []string{lang, utils.ToStr(i)})
	}
	return data
}

func (form *PostForm) CategorySelectData() [][]string {
	data := make([][]string, 0, len(form.Categories))
	for _, cat := range form.Categories {
		data = append(data, []string{"category." + cat.Name, utils.ToStr(cat.Id)})
	}
	return data
}

func (form *PostForm) TopicSelectData() [][]string {
	data := make([][]string, 0, len(form.Topics))
	for _, topic := range form.Topics {
		data = append(data, []string{topic.GetName(form.Locale.Lang), utils.ToStr(topic.Id)})
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

	if len(i18n.GetLangByIndex(form.Lang)) == 0 {
		v.SetError("Lang", "error")
	}
}

func (form *PostForm) SavePost(post *models.Post, user *models.User) error {
	utils.SetFormValues(form, post)
	post.Category = &models.Category{Id: form.Category}
	post.Topic = &models.Topic{Id: form.Topic}
	post.User = user
	post.LastReply = user
	post.LastAuthor = user
	post.ContentCache = utils.RenderMarkdown(form.Content)

	// mentioned follow users
	FilterMentions(user, post.ContentCache)

	return post.Insert()
}

func (form *PostForm) SetFromPost(post *models.Post) {
	utils.SetFormValues(post, form)
	form.Category = post.Category.Id
	form.Topic = post.Topic.Id
}

func (form *PostForm) UpdatePost(post *models.Post, user *models.User) error {
	changes := utils.FormChanges(post, form)
	if len(changes) == 0 {
		return nil
	}
	utils.SetFormValues(form, post)
	post.Category.Id = form.Category
	post.Topic.Id = form.Topic
	for _, c := range changes {
		if c == "Content" {
			post.ContentCache = utils.RenderMarkdown(form.Content)
			changes = append(changes, "ContentCache")
		}
	}

	// update last edit author
	if post.LastAuthor != nil && post.LastAuthor.Id != user.Id {
		post.LastAuthor = user
		changes = append(changes, "LastAuthor")
	}

	changes = append(changes, "Updated")

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
	PostForm   `form:"-"`
	Create     bool   `form:"-"`
	User       int    `form:"attr(rel,select2-admin-model);attr(data-model,User)" valid:"Required"`
	Title      string `valid:"Required;MaxSize(60)"`
	Content    string `form:"type(textarea,markdown)" valid:"Required"`
	Browsers   int    ``
	Replys     int    ``
	Favorites  int    ``
	LastReply  int    `form:"attr(rel,select2-admin-model);attr(data-model,User)" valid:""`
	LastAuthor int    `form:"attr(rel,select2-admin-model);attr(data-model,User)" valid:""`
	Topic      int    `form:"type(select);attr(rel,select2)" valid:"Required"`
	Category   int    `form:"type(select);attr(rel,select2)" valid:"Required"`
	Lang       int    `form:"type(select);attr(rel,select2)"`
	IsBest     bool   ``
}

func (form *PostAdminForm) Valid(v *validation.Validation) {
	user := models.User{Id: form.User}
	if user.Read() != nil {
		v.SetError("User", "admin.not_found_by_id")
	}

	user.Id = form.LastReply
	if user.Read() != nil {
		v.SetError("LastReply", "admin.not_found_by_id")
	}

	user.Id = form.LastAuthor
	if user.Read() != nil {
		v.SetError("LastReply", "admin.not_found_by_id")
	}

	topic := models.Topic{Id: form.Topic}
	if topic.Read() != nil {
		v.SetError("Topic", "admin.not_found_by_id")
	}

	cat := models.Category{Id: form.Category}
	if cat.Read() != nil {
		v.SetError("Category", "admin.not_found_by_id")
	}

	if len(i18n.GetLangByIndex(form.Lang)) == 0 {
		v.SetError("Lang", "Not Found")
	}
}

func (form *PostAdminForm) SetFromPost(post *models.Post) {
	utils.SetFormValues(post, form)

	if post.User != nil {
		form.User = post.User.Id
	}

	if post.LastReply != nil {
		form.LastReply = post.LastReply.Id
	}

	if post.LastAuthor != nil {
		form.LastAuthor = post.LastAuthor.Id
	}

	if post.Topic != nil {
		form.Topic = post.Topic.Id
	}

	if post.Category != nil {
		form.Category = post.Category.Id
	}
}

func (form *PostAdminForm) SetToPost(post *models.Post) {
	utils.SetFormValues(form, post)

	if post.User == nil {
		post.User = &models.User{}
	}
	post.User.Id = form.User

	if post.LastReply == nil {
		post.LastReply = &models.User{}
	}
	post.LastReply.Id = form.LastReply

	if post.LastAuthor == nil {
		post.LastAuthor = &models.User{}
	}
	post.LastAuthor.Id = form.LastAuthor

	if post.Topic == nil {
		post.Topic = &models.Topic{}
	}
	post.Topic.Id = form.Topic

	if post.Category == nil {
		post.Category = &models.Category{}
	}
	post.Category.Id = form.Category

	post.ContentCache = utils.RenderMarkdown(post.Content)
}

type CommentForm struct {
	Message string `form:"type(textarea,markdown)" valid:"Required;MinSize(5)"`
}

func (form *CommentForm) SaveComment(comment *models.Comment, user *models.User, post *models.Post) error {
	comment.Message = form.Message
	comment.MessageCache = utils.RenderMarkdown(form.Message)
	comment.User = user
	comment.Post = post
	if err := comment.Insert(); err == nil {
		post.LastReply = user
		post.Update("LastReply", "Updated")

		cnt, _ := post.Comments().Filter("Id__lte", comment.Id).Count()
		comment.Floor = int(cnt)
		return comment.Update("Floor")
	} else {
		return err
	}
}

type CommentAdminForm struct {
	Create  bool   `form:"-"`
	User    int    `form:"attr(rel,select2-admin-model);attr(data-model,User)" valid:"Required"`
	Post    int    `valid:"Required"`
	Message string `form:"type(textarea)" valid:"Required"`
	Floor   int    `valid:"Required"`
	Status  int    `valid:""`
}

func (form *CommentAdminForm) Valid(v *validation.Validation) {
	user := models.User{Id: form.User}
	if user.Read() != nil {
		v.SetError("User", "admin.not_found_by_id")
	}

	post := models.Post{Id: form.Post}
	if post.Read() != nil {
		v.SetError("Post", "admin.not_found_by_id")
	}
}

func (form *CommentAdminForm) SetFromComment(comment *models.Comment) {
	utils.SetFormValues(comment, form)

	if comment.User != nil {
		form.User = comment.User.Id
	}

	if comment.Post != nil {
		form.Post = comment.Post.Id
	}
}

func (form *CommentAdminForm) SetToComment(comment *models.Comment) {
	utils.SetFormValues(form, comment)

	if comment.User == nil {
		comment.User = &models.User{}
	}
	comment.User.Id = form.User

	if comment.Post == nil {
		comment.Post = &models.Post{}
	}
	comment.Post.Id = form.Post

	comment.MessageCache = utils.RenderMarkdown(comment.Message)
}
