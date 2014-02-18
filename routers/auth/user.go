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

package auth

import (
	"github.com/astaxie/beego/orm"

	"github.com/beego/wetalk/modules/models"
	"github.com/beego/wetalk/modules/utils"
	"github.com/beego/wetalk/routers/base"
)

type UserRouter struct {
	base.BaseRouter
}

func (this *UserRouter) getUser(user *models.User) bool {
	username := this.GetString(":username")
	user.UserName = username

	err := user.Read("UserName")
	if err != nil {
		this.Abort("404")
		return true
	}

	IsFollowed := false

	if this.IsLogin {
		if this.User.Id != user.Id {
			IsFollowed = this.User.FollowingUsers().Filter("FollowUser", user.Id).Exist()
		}
	}

	this.Data["TheUser"] = &user
	this.Data["IsFollowed"] = IsFollowed

	return false
}

func (this *UserRouter) Home() {
	this.TplNames = "user/home.html"

	var user models.User
	if this.getUser(&user) {
		return
	}

	limit := 5

	var posts []*models.Post
	var comments []*models.Comment

	user.RecentPosts().Limit(limit).RelatedSel().All(&posts)
	user.RecentComments().Limit(limit).RelatedSel().All(&comments)

	var ftopics []*models.FollowTopic
	var topics []*models.Topic
	nums, _ := models.FollowTopics().Filter("User", &user.Id).Limit(8).OrderBy("-Created").RelatedSel("Topic").All(&ftopics, "Topic")
	if nums > 0 {
		topics = make([]*models.Topic, 0, nums)
		for _, ft := range ftopics {
			topics = append(topics, ft.Topic)
		}
	}
	this.Data["TheUserTopics"] = topics
	this.Data["TheUserTopicsMore"] = nums >= 8

	this.Data["TheUserPosts"] = posts
	this.Data["TheUserComments"] = comments
}

func (this *UserRouter) Posts() {
	this.TplNames = "user/posts.html"

	var user models.User
	if this.getUser(&user) {
		return
	}

	limit := 20

	qs := user.RecentPosts()
	nums, _ := qs.Count()

	pager := this.SetPaginator(limit, nums)

	var posts []*models.Post
	qs.Limit(limit, pager.Offset()).RelatedSel().All(&posts)

	this.Data["TheUserPosts"] = posts
}

func (this *UserRouter) Comments() {
	this.TplNames = "user/comments.html"

	var user models.User
	if this.getUser(&user) {
		return
	}

	limit := 20

	qs := user.RecentComments()
	nums, _ := qs.Count()

	pager := this.SetPaginator(limit, nums)

	var comments []*models.Comment
	qs.Limit(limit, pager.Offset()).RelatedSel().All(&comments)

	this.Data["TheUserComments"] = comments
}

func (this *UserRouter) getFollows(user *models.User, following bool) []map[string]interface{} {
	limit := 20

	var qs orm.QuerySeter

	if following {
		qs = user.FollowingUsers()
	} else {
		qs = user.FollowerUsers()
	}

	nums, _ := qs.Count()

	pager := this.SetPaginator(limit, nums)

	qs = qs.Limit(limit, pager.Offset())

	var follows []*models.Follow

	if following {
		qs.RelatedSel("FollowUser").All(&follows, "FollowUser")
	} else {
		qs.RelatedSel("User").All(&follows, "User")
	}

	if len(follows) == 0 {
		return nil
	}

	ids := make([]int, 0, len(follows))
	for _, follow := range follows {
		if following {
			ids = append(ids, follow.FollowUser.Id)
		} else {
			ids = append(ids, follow.User.Id)
		}
	}

	var eids orm.ParamsList
	this.User.FollowingUsers().Filter("FollowUser__in", ids).ValuesFlat(&eids, "FollowUser__Id")

	var fids map[int]bool
	if len(eids) > 0 {
		fids = make(map[int]bool)
		for _, id := range eids {
			tid, _ := utils.StrTo(utils.ToStr(id)).Int()
			if tid > 0 {
				fids[tid] = true
			}
		}
	}

	users := make([]map[string]interface{}, 0, len(follows))
	for _, follow := range follows {
		IsFollowed := false
		var u *models.User
		if following {
			u = follow.FollowUser
		} else {
			u = follow.User
		}
		if fids != nil {
			IsFollowed = fids[u.Id]
		}
		users = append(users, map[string]interface{}{
			"User":       u,
			"IsFollowed": IsFollowed,
		})
	}

	return users
}

func (this *UserRouter) Following() {
	this.TplNames = "user/following.html"

	var user models.User
	if this.getUser(&user) {
		return
	}

	users := this.getFollows(&user, true)

	this.Data["TheUserFollowing"] = users
}

func (this *UserRouter) Followers() {
	this.TplNames = "user/followers.html"

	var user models.User
	if this.getUser(&user) {
		return
	}

	users := this.getFollows(&user, false)

	this.Data["TheUserFollowers"] = users
}

func (this *UserRouter) Favs() {
	this.TplNames = "user/favs.html"

	var user models.User
	if this.getUser(&user) {
		return
	}

	var ftopics []*models.FollowTopic
	var topics []*models.Topic
	nums, _ := models.FollowTopics().Filter("User", &user.Id).OrderBy("-Created").RelatedSel("Topic").All(&ftopics, "Topic")
	if nums > 0 {
		topics = make([]*models.Topic, 0, nums)
		for _, ft := range ftopics {
			topics = append(topics, ft.Topic)
		}
	}
	this.Data["TheUserTopics"] = topics
}
