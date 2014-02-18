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

package article

import (
	"github.com/beego/wetalk/modules/models"
	"github.com/beego/wetalk/routers/base"
)

type ArticleRouter struct {
	base.BaseRouter
}

func (this *ArticleRouter) loadArticle(article *models.Article) bool {
	uri := this.Ctx.Request.RequestURI
	err := models.Articles().RelatedSel("User").Filter("IsPublish", true).Filter("Uri", uri).One(article)
	if err == nil {
		this.Data["Article"] = article
	} else {
		this.Abort("404")
	}
	return err != nil
}

func (this *ArticleRouter) Show() {
	this.TplNames = "article/show.html"
	article := models.Article{}
	if this.loadArticle(&article) {
		return
	}
}
