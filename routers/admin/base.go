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
	"github.com/astaxie/beego/orm"

	"github.com/beego/wetalk/modules/models"
)

type AdminRouter struct {
	BaseAdminRouter
}

func (this *AdminRouter) ModelGet() {
	id := this.GetString("id")
	model := this.GetString("model")
	result := map[string]interface{}{
		"success": false,
	}

	var data []orm.ParamsList

	defer func() {
		if len(data) > 0 {
			result["success"] = true
			result["data"] = data[0]
		}
		this.Data["json"] = result
		this.ServeJson()
	}()

	var qs orm.QuerySeter

	switch model {
	case "User":
		qs = models.Users()
	}

	qs = qs.Filter("Id", id).Limit(1)

	switch model {
	case "User":
		qs.ValuesList(&data, "Id", "UserName")
	}
}

func (this *AdminRouter) ModelSelect() {
	search := this.GetString("search")
	model := this.GetString("model")
	result := map[string]interface{}{
		"success": false,
	}

	var data []orm.ParamsList

	defer func() {
		if len(data) > 0 {
			result["success"] = true
			result["data"] = data
		}
		this.Data["json"] = result
		this.ServeJson()
	}()

	if len(search) < 3 {
		return
	}

	switch model {
	case "User":
		models.Users().Filter("UserName__icontains", search).Limit(10).ValuesList(&data, "Id", "UserName")
	}
}
