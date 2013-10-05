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

package routers

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"

	"github.com/beego/wetalk/models"
	"github.com/beego/wetalk/utils"
)

type PostAdminRouter struct {
	ModelAdminRouter
	object models.Post
}

func (this *PostAdminRouter) Object() interface{} {
	return &this.object
}

func (this *PostAdminRouter) ObjectQs() orm.QuerySeter {
	return models.Posts().RelatedSel()
}

// view for list model data
func (this *PostAdminRouter) List() {
	var posts []models.Post
	qs := models.Posts().RelatedSel()
	if err := this.SetObjects(qs, &posts); err != nil {
		this.Data["Error"] = err
		beego.Error(err)
	}
}

// view for create object
func (this *PostAdminRouter) Create() {
	form := models.PostAdminForm{Create: true}
	this.SetFormSets(&form)
}

// view for new object save
func (this *PostAdminRouter) Save() {
	form := models.PostAdminForm{Create: true}
	if this.ValidFormSets(&form) == false {
		return
	}

	var post models.Post
	form.SetToPost(&post)
	if err := post.Insert(); err == nil {
		this.FlashRedirect(fmt.Sprintf("/admin/post/%d", post.Id), 302, "CreateSuccess")
		return
	} else {
		beego.Error(err)
		this.Data["Error"] = err
	}
}

// view for edit object
func (this *PostAdminRouter) Edit() {
	form := models.PostAdminForm{}
	form.SetFromPost(&this.object)
	this.SetFormSets(&form)
}

// view for update object
func (this *PostAdminRouter) Update() {
	form := models.PostAdminForm{}
	if this.ValidFormSets(&form) == false {
		return
	}

	// get changed field names
	changes := utils.FormChanges(&this.object, &form)

	url := fmt.Sprintf("/admin/post/%d", this.object.Id)

	// update changed fields only
	if len(changes) > 0 {
		form.SetToPost(&this.object)
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
func (this *PostAdminRouter) Confirm() {
}

// view for delete object
func (this *PostAdminRouter) Delete() {
	if this.FormOnceNotMatch() {
		return
	}

	// delete object
	if err := this.object.Delete(); err == nil {
		this.FlashRedirect("/admin/post", 302, "DeleteSuccess")
		return
	} else {
		beego.Error(err)
		this.Data["Error"] = err
	}
}
