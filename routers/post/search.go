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
	"strings"

	"github.com/astaxie/beego"

	"github.com/beego/wetalk/modules/post"
	"github.com/beego/wetalk/modules/utils"
	"github.com/beego/wetalk/routers/base"
)

type SearchRouter struct {
	base.BaseRouter
}

func (this *SearchRouter) Get() {
	this.TplNames = "search/posts.html"

	pers := 25

	q := strings.TrimSpace(this.GetString("q"))

	this.Data["Q"] = q

	if len(q) == 0 {
		return
	}

	page, _ := utils.StrTo(this.GetString("p")).Int()

	posts, meta, err := post.SearchPost(q, page)
	if err != nil {
		this.Data["SearchError"] = true
		beego.Error("SearchPosts: ", err)
		return
	}

	if len(posts) > 0 {
		this.SetPaginator(pers, meta.TotalFound)
	}

	this.Data["Posts"] = posts
	this.Data["Meta"] = meta
}
