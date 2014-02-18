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
	"fmt"

	"github.com/astaxie/beego/orm"

	"github.com/beego/wetalk/modules/models"
	"github.com/beego/wetalk/modules/post"
	"github.com/beego/wetalk/modules/utils"
	"github.com/beego/wetalk/routers/base"
)

// HomeRouter serves home page.
type PostListRouter struct {
	base.BaseRouter
}

func (this *PostListRouter) setCategories(cats *[]models.Category) {
	post.ListCategories(cats)
	this.Data["Categories"] = *cats
}

func (this *PostListRouter) setTopicsOfCat(topics *[]models.Topic, cat *models.Category) {
	post.ListTopicsOfCat(topics, cat)
	this.Data["Topics"] = *topics
}

func (this *PostListRouter) postsFilter(qs orm.QuerySeter) orm.QuerySeter {
	if !this.IsLogin {
		return qs
	}
	args := []string{utils.ToStr(this.Locale.Index())}
	args = append(args, this.User.LangAdds...)
	args = append(args, utils.ToStr(this.User.Lang))
	qs = qs.Filter("Lang__in", args)
	return qs
}

// Get implemented Get method for HomeRouter.
func (this *PostListRouter) Home() {
	this.Data["IsHome"] = true
	this.TplNames = "post/home.html"

	var cats []models.Category
	this.setCategories(&cats)

	var posts []models.Post
	qs := models.Posts().OrderBy("-Created").Limit(25).RelatedSel()
	qs = this.postsFilter(qs)

	models.ListObjects(qs, &posts)
	this.Data["Posts"] = posts

	this.Data["CategorySlug"] = "hot"

	var topics []models.Topic
	post.ListTopics(&topics)
	this.Data["Topics"] = topics
}

// Get implemented Get method for HomeRouter.
func (this *PostListRouter) Category() {
	this.TplNames = "post/category.html"

	slug := this.GetString(":slug")
	cat := models.Category{Slug: slug}
	if err := cat.Read("Slug"); err != nil {
		this.Abort("404")
		return
	}

	pers := 25

	qs := models.Posts().Filter("Category", &cat)
	qs = this.postsFilter(qs)

	cnt, _ := models.CountObjects(qs)
	pager := this.SetPaginator(pers, cnt)

	qs = qs.OrderBy("-Created").Limit(pers, pager.Offset()).RelatedSel()

	var posts []models.Post
	models.ListObjects(qs, &posts)

	this.Data["Posts"] = posts
	this.Data["Category"] = &cat
	this.Data["CategorySlug"] = cat.Slug
	this.Data["IsCategory"] = true

	var cats []models.Category
	this.setCategories(&cats)

	var topics []models.Topic
	this.setTopicsOfCat(&topics, &cat)
}

// Get implemented Get method for HomeRouter.
func (this *PostListRouter) Navs() {
	slug := this.GetString(":slug")

	switch slug {
	case "favs", "follow":
		if this.CheckLoginRedirect() {
			return
		}
	}

	this.Data["CategorySlug"] = slug
	this.TplNames = fmt.Sprintf("post/navs/%s.html", slug)

	pers := 25

	var posts []models.Post

	switch slug {
	case "recent":
		qs := models.Posts()
		qs = this.postsFilter(qs)

		cnt, _ := models.CountObjects(qs)
		pager := this.SetPaginator(pers, cnt)

		qs = qs.OrderBy("-Updated").Limit(pers, pager.Offset()).RelatedSel()

		models.ListObjects(qs, &posts)

		var cats []models.Category
		this.setCategories(&cats)

	case "best":
		qs := models.Posts().Filter("IsBest", true)
		qs = this.postsFilter(qs)

		cnt, _ := models.CountObjects(qs)
		pager := this.SetPaginator(pers, cnt)

		qs = qs.OrderBy("-Created").Limit(pers, pager.Offset()).RelatedSel()

		models.ListObjects(qs, &posts)

		var cats []models.Category
		this.setCategories(&cats)

	case "cold":
		qs := models.Posts().Filter("Replys", 0)
		qs = this.postsFilter(qs)

		cnt, _ := models.CountObjects(qs)
		pager := this.SetPaginator(pers, cnt)

		qs = qs.OrderBy("-Created").Limit(pers, pager.Offset()).RelatedSel()

		models.ListObjects(qs, &posts)

		var cats []models.Category
		this.setCategories(&cats)

	case "favs":
		var topicIds orm.ParamsList
		nums, _ := models.FollowTopics().Filter("User", &this.User.Id).OrderBy("-Created").ValuesFlat(&topicIds, "Topic")
		if nums > 0 {
			qs := models.Posts().Filter("Topic__in", topicIds)
			qs = this.postsFilter(qs)

			cnt, _ := models.CountObjects(qs)
			pager := this.SetPaginator(pers, cnt)

			qs = qs.OrderBy("-Created").Limit(pers, pager.Offset()).RelatedSel()

			models.ListObjects(qs, &posts)

			var topics []models.Topic
			nums, _ = models.Topics().Filter("Id__in", topicIds).Limit(8).All(&topics)
			this.Data["Topics"] = topics
			this.Data["TopicsMore"] = nums >= 8
		}

	case "follow":
		var userIds orm.ParamsList
		nums, _ := this.User.FollowingUsers().OrderBy("-Created").ValuesFlat(&userIds, "FollowUser")
		if nums > 0 {
			qs := models.Posts().Filter("User__in", userIds)
			qs = this.postsFilter(qs)

			cnt, _ := models.CountObjects(qs)
			pager := this.SetPaginator(pers, cnt)

			qs = qs.OrderBy("-Created").Limit(pers, pager.Offset()).RelatedSel()

			models.ListObjects(qs, &posts)
		}
	}

	this.Data["Posts"] = posts
}

// Get implemented Get method for HomeRouter.
func (this *PostListRouter) Topic() {
	slug := this.GetString(":slug")

	switch slug {
	default: // View topic.
		this.TplNames = "post/topic.html"
		topic := models.Topic{Slug: slug}
		if err := topic.Read("Slug"); err != nil {
			this.Abort("404")
			return
		}

		pers := 25

		qs := models.Posts().Filter("Topic", &topic)
		qs = this.postsFilter(qs)

		cnt, _ := models.CountObjects(qs)
		pager := this.SetPaginator(pers, cnt)

		qs = qs.OrderBy("-Created").Limit(pers, pager.Offset()).RelatedSel()

		var posts []models.Post
		models.ListObjects(qs, &posts)

		this.Data["Posts"] = posts
		this.Data["Topic"] = &topic
		this.Data["IsTopic"] = true

		HasFavorite := false
		if this.IsLogin {
			HasFavorite = models.FollowTopics().Filter("User", &this.User).Filter("Topic", &topic).Exist()
		}
		this.Data["HasFavorite"] = HasFavorite
	}
}

// Get implemented Get method for HomeRouter.
func (this *PostListRouter) TopicSubmit() {
	slug := this.GetString(":slug")

	topic := models.Topic{Slug: slug}
	if err := topic.Read("Slug"); err != nil {
		this.Abort("404")
		return
	}

	result := map[string]interface{}{
		"success": false,
	}

	if this.IsAjax() {
		action := this.GetString("action")
		switch action {
		case "favorite":
			if this.IsLogin {
				qs := models.FollowTopics().Filter("User", &this.User).Filter("Topic", &topic)
				if qs.Exist() {
					qs.Delete()
				} else {
					fav := models.FollowTopic{User: &this.User, Topic: &topic}
					fav.Insert()
				}
				topic.RefreshFollowers()
				this.User.RefreshFavTopics()
				result["success"] = true
			}
		}
	}

	this.Data["json"] = result
	this.ServeJson()
}

type PostRouter struct {
	base.BaseRouter
}

func (this *PostRouter) New() {
	this.TplNames = "post/new.html"

	if this.CheckActiveRedirect() {
		return
	}

	form := post.PostForm{Locale: this.Locale}

	if v := this.Ctx.GetCookie("post_topic"); len(v) > 0 {
		form.Topic, _ = utils.StrTo(v).Int()
	}

	if v := this.Ctx.GetCookie("post_cat"); len(v) > 0 {
		form.Category, _ = utils.StrTo(v).Int()
	}

	if v := this.Ctx.GetCookie("post_lang"); len(v) > 0 {
		form.Lang, _ = utils.StrTo(v).Int()
	} else {
		form.Lang = this.Locale.Index()
	}

	slug := this.GetString("topic")
	if len(slug) > 0 {
		topic := models.Topic{Slug: slug}
		topic.Read("Slug")
		form.Topic = topic.Id
		this.Data["Topic"] = &topic
	}

	post.ListCategories(&form.Categories)
	post.ListTopics(&form.Topics)
	this.SetFormSets(&form)
}

func (this *PostRouter) NewSubmit() {
	this.TplNames = "post/new.html"

	if this.CheckActiveRedirect() {
		return
	}

	form := post.PostForm{Locale: this.Locale}
	slug := this.GetString("topic")
	if len(slug) > 0 {
		topic := models.Topic{Slug: slug}
		topic.Read("Slug")
		form.Topic = topic.Id
		this.Data["Topic"] = &topic
	}

	post.ListCategories(&form.Categories)
	post.ListTopics(&form.Topics)
	if !this.ValidFormSets(&form) {
		return
	}

	var post models.Post
	if err := form.SavePost(&post, &this.User); err == nil {

		this.Ctx.SetCookie("post_topic", utils.ToStr(form.Topic), 1<<31-1, "/")
		this.Ctx.SetCookie("post_cat", utils.ToStr(form.Category), 1<<31-1, "/")
		this.Ctx.SetCookie("post_lang", utils.ToStr(form.Lang), 1<<31-1, "/")

		this.JsStorage("deleteKey", "post/new")
		this.Redirect(post.Link(), 302)
	}
}

func (this *PostRouter) loadPost(post *models.Post, user *models.User) bool {
	id, _ := this.GetInt(":post")
	if id > 0 {
		qs := models.Posts().Filter("Id", id)
		if user != nil {
			qs = qs.Filter("User", user.Id)
		}
		qs.RelatedSel(1).One(post)
	}

	if post.Id == 0 {
		this.Abort("404")
		return true
	}

	this.Data["Post"] = post

	return false
}

func (this *PostRouter) loadComments(post *models.Post, comments *[]*models.Comment) {
	qs := post.Comments()
	if num, err := qs.RelatedSel("User").OrderBy("Id").All(comments); err == nil {
		this.Data["Comments"] = *comments
		this.Data["CommentsNum"] = num
	}
}

func (this *PostRouter) Single() {
	this.TplNames = "post/post.html"

	var postMd models.Post
	if this.loadPost(&postMd, nil) {
		return
	}

	var comments []*models.Comment
	this.loadComments(&postMd, &comments)

	form := post.CommentForm{}
	this.SetFormSets(&form)

	post.PostBrowsersAdd(this.User.Id, this.Ctx.Input.IP(), &postMd)
}

func (this *PostRouter) SingleSubmit() {
	this.TplNames = "post/post.html"

	if this.CheckActiveRedirect() {
		return
	}

	var postMd models.Post
	if this.loadPost(&postMd, nil) {
		return
	}

	var redir bool

	defer func() {
		if !redir {
			var comments []*models.Comment
			this.loadComments(&postMd, &comments)
		}
	}()

	form := post.CommentForm{}
	if !this.ValidFormSets(&form) {
		return
	}

	comment := models.Comment{}
	if err := form.SaveComment(&comment, &this.User, &postMd); err == nil {
		this.JsStorage("deleteKey", "post/comment")
		this.Redirect(postMd.Link(), 302)
		redir = true

		post.PostReplysCount(&postMd)
	}
}

func (this *PostRouter) Edit() {
	this.TplNames = "post/edit.html"

	if this.CheckActiveRedirect() {
		return
	}

	var postMd models.Post
	if this.loadPost(&postMd, &this.User) {
		return
	}

	form := post.PostForm{}
	form.SetFromPost(&postMd)
	post.ListCategories(&form.Categories)
	post.ListTopics(&form.Topics)
	this.SetFormSets(&form)
}

func (this *PostRouter) EditSubmit() {
	this.TplNames = "post/edit.html"

	if this.CheckActiveRedirect() {
		return
	}

	var postMd models.Post
	if this.loadPost(&postMd, &this.User) {
		return
	}

	form := post.PostForm{}
	form.SetFromPost(&postMd)
	post.ListCategories(&form.Categories)
	post.ListTopics(&form.Topics)
	if !this.ValidFormSets(&form) {
		return
	}

	if err := form.UpdatePost(&postMd, &this.User); err == nil {
		this.JsStorage("deleteKey", "post/edit")
		this.Redirect(postMd.Link(), 302)
	}
}
