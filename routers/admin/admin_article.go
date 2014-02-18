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

package admin

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"

	"github.com/beego/wetalk/modules/article"
	"github.com/beego/wetalk/modules/models"
	"github.com/beego/wetalk/modules/utils"
)

type ArticleAdminRouter struct {
	ModelAdminRouter
	object models.Article
}

func (this *ArticleAdminRouter) Object() interface{} {
	return &this.object
}

func (this *ArticleAdminRouter) ObjectQs() orm.QuerySeter {
	return models.Articles().RelatedSel()
}

// view for list model data
func (this *ArticleAdminRouter) List() {
	var articles []models.Article
	qs := models.Articles().RelatedSel()
	if err := this.SetObjects(qs, &articles); err != nil {
		this.Data["Error"] = err
		beego.Error(err)
	}
}

// view for create object
func (this *ArticleAdminRouter) Create() {
	form := article.ArticleAdminForm{Create: true}
	this.SetFormSets(&form)
}

// view for new object save
func (this *ArticleAdminRouter) Save() {
	form := article.ArticleAdminForm{Create: true}
	if !this.ValidFormSets(&form) {
		return
	}

	var a models.Article
	form.SetToArticle(&a)
	if err := a.Insert(); err == nil {
		this.FlashRedirect(fmt.Sprintf("/admin/article/%d", a.Id), 302, "CreateSuccess")
		return
	} else {
		beego.Error(err)
		this.Data["Error"] = err
	}
}

// view for edit object
func (this *ArticleAdminRouter) Edit() {
	form := article.ArticleAdminForm{}
	form.SetFromArticle(&this.object)
	this.SetFormSets(&form)
}

// view for update object
func (this *ArticleAdminRouter) Update() {
	form := article.ArticleAdminForm{}
	if this.ValidFormSets(&form) == false {
		return
	}

	// get changed field names
	changes := utils.FormChanges(&this.object, &form)

	url := fmt.Sprintf("/admin/article/%d", this.object.Id)

	// update changed fields only
	if len(changes) > 0 {
		form.SetToArticle(&this.object)
		if err := this.object.Update(changes...); err == nil {
			this.FlashRedirect(url, 302, "UpdateSuccess")
			return
		} else {
			beego.Error(err)
			this.Data["Error"] = err
		}
	} else {
		this.Redirect(url, 302)
	}
}

// view for confirm delete object
func (this *ArticleAdminRouter) Confirm() {
}

// view for delete object
func (this *ArticleAdminRouter) Delete() {
	if this.FormOnceNotMatch() {
		return
	}

	// delete object
	if err := this.object.Delete(); err == nil {
		this.FlashRedirect("/admin/article", 302, "DeleteSuccess")
		return
	} else {
		beego.Error(err)
		this.Data["Error"] = err
	}
}
