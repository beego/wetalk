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

	"github.com/beego/wetalk/modules/models"
	"github.com/beego/wetalk/modules/post"
	"github.com/beego/wetalk/modules/utils"
)

type CategoryAdminRouter struct {
	ModelAdminRouter
	object models.Category
}

func (this *CategoryAdminRouter) Object() interface{} {
	return &this.object
}

func (this *CategoryAdminRouter) ObjectQs() orm.QuerySeter {
	return models.Categories().RelatedSel()
}

// view for list model data
func (this *CategoryAdminRouter) List() {
	var cats []models.Category
	qs := models.Categories().RelatedSel()
	if err := this.SetObjects(qs, &cats); err != nil {
		this.Data["Error"] = err
		beego.Error(err)
	}
}

// view for create object
func (this *CategoryAdminRouter) Create() {
	form := post.CategoryAdminForm{Create: true}
	this.SetFormSets(&form)
}

// view for new object save
func (this *CategoryAdminRouter) Save() {
	form := post.CategoryAdminForm{Create: true}
	if this.ValidFormSets(&form) == false {
		return
	}

	var cat models.Category
	form.SetToCategory(&cat)
	if err := cat.Insert(); err == nil {
		this.FlashRedirect(fmt.Sprintf("/admin/category/%d", cat.Id), 302, "CreateSuccess")
		return
	} else {
		beego.Error(err)
		this.Data["Error"] = err
	}
}

// view for edit object
func (this *CategoryAdminRouter) Edit() {
	form := post.CategoryAdminForm{}
	form.SetFromCategory(&this.object)
	this.SetFormSets(&form)
}

// view for update object
func (this *CategoryAdminRouter) Update() {
	form := post.CategoryAdminForm{Id: this.object.Id}
	if this.ValidFormSets(&form) == false {
		return
	}

	// get changed field names
	changes := utils.FormChanges(&this.object, &form)

	url := fmt.Sprintf("/admin/category/%d", this.object.Id)

	// update changed fields only
	if len(changes) > 0 {
		form.SetToCategory(&this.object)
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
func (this *CategoryAdminRouter) Confirm() {
}

// view for delete object
func (this *CategoryAdminRouter) Delete() {
	if this.FormOnceNotMatch() {
		return
	}

	// delete object
	if err := this.object.Delete(); err == nil {
		this.FlashRedirect("/admin/category", 302, "DeleteSuccess")
		return
	} else {
		beego.Error(err)
		this.Data["Error"] = err
	}
}
