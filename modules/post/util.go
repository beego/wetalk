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
	"regexp"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"

	"github.com/beego/wetalk/modules/models"
	"github.com/beego/wetalk/modules/utils"
	"github.com/beego/wetalk/setting"
)

func ListPostsOfCategory(cat *models.Category, posts *[]models.Post) (int64, error) {
	return models.Posts().Filter("Category", cat).RelatedSel().OrderBy("-Updated").All(posts)
}

func ListPostsOfTopic(topic *models.Topic, posts *[]models.Post) (int64, error) {
	return models.Posts().Filter("Topic", topic).RelatedSel().OrderBy("-Updated").All(posts)
}

var mentionRegexp = regexp.MustCompile(`\B@([\d\w-_]*)`)

func FilterMentions(user *models.User, content string) {
	matches := mentionRegexp.FindAllStringSubmatch(content, -1)
	mentions := make([]string, 0, len(matches))
	for _, m := range matches {
		if len(m) > 1 {
			mentions = append(mentions, m[1])
		}
	}
	// var users []*User
	// num, err := Users().Filter("UserName__in", mentions).Filter("Follow__User", user.Id).All(&users)
	// if err == nil && num > 0 {
	// TODO mention email to user
	// }
}

func PostBrowsersAdd(uid int, ip string, post *models.Post) {
	var key string
	if uid == 0 {
		key = ip
	} else {
		key = utils.ToStr(uid)
	}
	key = fmt.Sprintf("PCA.%d.%s", post.Id, key)
	if setting.Cache.Get(key) != nil {
		return
	}
	_, err := models.Posts().Filter("Id", post.Id).Update(orm.Params{
		"Browsers": orm.ColValue(orm.Col_Add, 1),
	})
	if err != nil {
		beego.Error("PostCounterAdd ", err)
	}
	setting.Cache.Put(key, true, 60)
}

func PostReplysCount(post *models.Post) {
	cnt, err := post.Comments().Count()
	if err == nil {
		post.Replys = int(cnt)
		err = post.Update("Replys")
	}
	if err != nil {
		beego.Error("PostReplysCount ", err)
	}
}
