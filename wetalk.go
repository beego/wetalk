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
	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/orm"

	"github.com/beego/wetalk/routers"
	"github.com/beego/wetalk/utils"
)

// We have to call a initialize function manully
// because we use `bee bale` to pack static resources
// and we cannot make sure that which init() execute first.
func initialize() {
	cfg := utils.LoadConfig()

	var err error

	// cache system
	utils.Cache, err = cache.NewCache("memory", `{"interval":360}`)

	// session settings
	beego.SessionOn = true
	beego.SessionProvider = cfg.MustValue("app", "session_provider")
	beego.SessionSavePath = cfg.MustValue("app", "session_path")
	beego.SessionName = cfg.MustValue("app", "session_name")

	beego.EnableXSRF = true
	// xsrf token expire time
	beego.XSRFExpire = 86400 * 365

	driverName := cfg.MustValue("orm", "driver_name")
	dataSource := cfg.MustValue("orm", "data_source")
	maxIdle := cfg.MustInt("orm", "max_idle_conn")
	maxOpen := cfg.MustInt("orm", "max_open_conn")

	// set default database
	orm.RegisterDataBase("default", driverName, dataSource, maxIdle, maxOpen)

	orm.RunCommand()

	err = orm.RunSyncdb("default", false, false)
	if err != nil {
		beego.Error(err)
	}
}

func main() {
	initialize()

	beego.Info(beego.AppName, utils.APP_VER)

	// Register routers.

	posts := new(routers.PostRouter)
	beego.Router("/", posts, "get:Home")
	beego.Router("/p/:post([0-9]+)", posts, "get:Single;post:SingleSubmit")
	beego.Router("/new", posts, "get:New;post:NewSubmit")
	beego.Router("/:slug(recent|best|cold|favs|follow)", posts, "get:Navs")
	beego.Router("/category/:slug", posts, "get:Category")
	beego.Router("/topic/:slug", posts, "get:Topic;post:TopicSubmit")

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
