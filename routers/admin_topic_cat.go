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

type TopicCatAdminRouter struct {
	ModelAdminRouter
	object models.TopicCat
}

func (this *TopicCatAdminRouter) Object() interface{} {
	return &this.object
}

func (this *TopicCatAdminRouter) ObjectQs() orm.QuerySeter {
	return models.TopicCats().RelatedSel()
}

// view for list model data
func (this *TopicCatAdminRouter) List() {
	var topicCats []models.TopicCat
	qs := models.TopicCats().RelatedSel()
	if err := this.SetObjects(qs, &topicCats); err != nil {
		this.Data["Error"] = err
		beego.Error(err)
	}
}

// view for create object
func (this *TopicCatAdminRouter) Create() {
	form := models.TopicCatAdminForm{Create: true}
	this.SetFormSets(&form)
}

// view for new object save
func (this *TopicCatAdminRouter) Save() {
	form := models.TopicCatAdminForm{Create: true}
	if this.ValidFormSets(&form) == false {
		return
	}

	var topicCat models.TopicCat
	form.SetToTopicCat(&topicCat)
	if err := topicCat.Insert(); err == nil {
		this.FlashRedirect(fmt.Sprintf("/admin/topicCat/%d", topicCat.Id), 302, "CreateSuccess")
		return
	} else {
		beego.Error(err)
		this.Data["Error"] = err
	}
}

// view for edit object
func (this *TopicCatAdminRouter) Edit() {
	form := models.TopicCatAdminForm{}
	form.SetFromTopicCat(&this.object)
	this.SetFormSets(&form)
}

// view for update object
func (this *TopicCatAdminRouter) Update() {
	form := models.TopicCatAdminForm{Id: this.object.Id}
	if this.ValidFormSets(&form) == false {
		return
	}

	// get changed field names
	changes := utils.FormChanges(&this.object, &form)

	url := fmt.Sprintf("/admin/topicCat/%d", this.object.Id)

	// update changed fields only
	if len(changes) > 0 {
		form.SetToTopicCat(&this.object)
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
func (this *TopicCatAdminRouter) Confirm() {
}

// view for delete object
func (this *TopicCatAdminRouter) Delete() {
	if this.FormOnceNotMatch() {
		return
	}

	// delete object
	if err := this.object.Delete(); err == nil {
		this.FlashRedirect("/admin/topicCat", 302, "DeleteSuccess")
		return
	} else {
		beego.Error(err)
		this.Data["Error"] = err
	}
}
