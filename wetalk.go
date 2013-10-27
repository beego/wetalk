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

// An open source project for Gopher community.
package main

import (
	"fmt"

	"github.com/astaxie/beego"

	"github.com/beego/wetalk/routers"
	"github.com/beego/wetalk/utils"
)

// We have to call a initialize function manully
// because we use `bee bale` to pack static resources
// and we cannot make sure that which init() execute first.
func initialize() {
	utils.LoadConfig()
}

func main() {
	initialize()

	if utils.IsProMode {
		beego.Info("Product mode enabled")
	} else {
		beego.Info("Develment mode enabled")
	}
	beego.Info(beego.AppName, utils.APP_VER, utils.AppUrl)

	if !utils.IsProMode {
		beego.SetStaticPath("/static_source", "static_source")
	}

	// Add Filters
	beego.AddFilter("^/img/:", "AfterStatic", routers.ImageFilter)

	// Register routers.
	posts := new(routers.PostListRouter)
	beego.Router("/", posts, "get:Home")
	beego.Router("/:slug(recent|best|cold|favs|follow)", posts, "get:Navs")
	beego.Router("/category/:slug", posts, "get:Category")
	beego.Router("/topic/:slug", posts, "get:Topic;post:TopicSubmit")

	post := new(routers.PostRouter)
	beego.Router("/new", post, "get:New;post:NewSubmit")
	beego.Router("/p/:post([0-9]+)", post, "get:Single;post:SingleSubmit")
	beego.Router("/p/:post([0-9]+)/edit", post, "get:Edit;post:EditSubmit")

	user := new(routers.UserRouter)
	beego.Router("/u/:username", user, "get:Home")

	login := new(routers.LoginRouter)
	beego.Router("/login", login, "post:Login")
	beego.Router("/logout", login, "get:Logout")

	register := new(routers.RegisterRouter)
	beego.Router("/register", register, "post:Register")
	beego.Router("/active/success", register, "get:ActiveSuccess")
	beego.Router("/active/:code([0-9a-zA-Z]+)", register, "get:Active")

	settings := new(routers.SettingsRouter)
	beego.Router("/settings/profile", settings, "get:Profile;post:ProfileSave")

	forgot := new(routers.ForgotRouter)
	beego.Router("/forgot", forgot)
	beego.Router("/reset/:code([0-9a-zA-Z]+)", forgot, "get:Reset;post:ResetPost")

	upload := new(routers.UploadRouter)
	beego.Router("/upload", upload, "post:Post")

	adminDashboard := new(routers.AdminDashboardRouter)
	beego.Router("/admin", adminDashboard)

	routes := map[string]beego.ControllerInterface{
		"user":     new(routers.UserAdminRouter),
		"post":     new(routers.PostAdminRouter),
		"comment":  new(routers.CommentAdminRouter),
		"topic":    new(routers.TopicAdminRouter),
		"category": new(routers.CategoryAdminRouter),
	}
	for name, router := range routes {
		beego.Router(fmt.Sprintf("/admin/:model(%s)", name), router, "get:List")
		beego.Router(fmt.Sprintf("/admin/:model(%s)/:id(new)", name), router, "get:Create;post:Save")
		beego.Router(fmt.Sprintf("/admin/:model(%s)/:id([0-9]+)", name), router, "get:Edit;post:Update")
		beego.Router(fmt.Sprintf("/admin/:model(%s)/:id([0-9]+)/:action(delete)", name), router, "get:Confirm;post:Delete")
	}

	// "robot.txt"
	beego.Router("/robot.txt", &routers.RobotRouter{})

	// For all unknown pages.
	beego.Run()
}
