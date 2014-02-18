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

	"github.com/beego/wetalk/modules/models"
	"github.com/beego/wetalk/modules/utils"
)

type TopicAdminForm struct {
	Create    bool   `form:"-"`
	Id        int    `form:"-"`
	Name      string `valid:"Required;MaxSize(30)"`
	Intro     string `form:"type(textarea)" valid:"Required"`
	NameZhCn  string `valid:"Required;MaxSize(30)"`
	IntroZhCn string `form:"type(textarea)" valid:"Required"`
	Slug      string `valid:"Required;MaxSize(100)"`
	Followers int    ``
	Order     int    ``
	Image     string `valid:""`
}

func (form *TopicAdminForm) Valid(v *validation.Validation) {
	qs := models.Topics()

	if models.CheckIsExist(qs, "Name", form.Name, form.Id) {
		v.SetError("Name", "admin.field_need_unique")
	}

	if models.CheckIsExist(qs, "NameZhCn", form.NameZhCn, form.Id) {
		v.SetError("NameZhCn", "admin.field_need_unique")
	}

	if models.CheckIsExist(qs, "Slug", form.Slug, form.Id) {
		v.SetError("Slug", "admin.field_need_unique")
	}
}

func (form *TopicAdminForm) SetFromTopic(topic *models.Topic) {
	utils.SetFormValues(topic, form)
}

func (form *TopicAdminForm) SetToTopic(topic *models.Topic) {
	utils.SetFormValues(form, topic, "Id")
}

type CategoryAdminForm struct {
	Create bool   `form:"-"`
	Id     int    `form:"-"`
	Name   string `valid:"Required;MaxSize(30)"`
	Slug   string `valid:"Required;MaxSize(100)"`
	Order  int    ``
}

func (form *CategoryAdminForm) Valid(v *validation.Validation) {
	qs := models.Categories()

	if models.CheckIsExist(qs, "Name", form.Name, form.Id) {
		v.SetError("Name", "admin.field_need_unique")
	}

	if models.CheckIsExist(qs, "Slug", form.Slug, form.Id) {
		v.SetError("Slug", "admin.field_need_unique")
	}
}

func (form *CategoryAdminForm) SetFromCategory(cat *models.Category) {
	utils.SetFormValues(cat, form)
}

func (form *CategoryAdminForm) SetToCategory(cat *models.Category) {
	utils.SetFormValues(form, cat, "Id")
}
