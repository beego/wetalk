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

package api

import (
	"github.com/astaxie/beego/orm"

	"github.com/beego/wetalk/modules/auth"
	"github.com/beego/wetalk/modules/models"
	"github.com/beego/wetalk/modules/utils"
	"github.com/beego/wetalk/routers/base"
)

type ApiRouter struct {
	base.BaseRouter
}

func (this *ApiRouter) Users() {
	result := map[string]interface{}{
		"success": false,
	}

	defer func() {
		this.Data["json"] = result
		this.ServeJson()
	}()

	if !this.IsAjax() {
		return
	}

	action := this.GetString("action")

	if this.IsLogin {

		switch action {
		case "get-follows":
			var data []orm.ParamsList
			this.User.FollowingUsers().ValuesList(&data, "FollowUser__NickName", "FollowUser__UserName")
			result["success"] = true
			result["data"] = data

		case "follow", "unfollow":
			id, err := utils.StrTo(this.GetString("user")).Int()
			if err == nil && id != this.User.Id {
				fuser := models.User{Id: id}
				if action == "follow" {
					auth.UserFollow(&this.User, &fuser)
				} else {
					auth.UserUnFollow(&this.User, &fuser)
				}
				result["success"] = true
			}
		}
	}
}
