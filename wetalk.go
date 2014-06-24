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
	"github.com/beego/social-auth"

	"github.com/beego/wetalk/modules/utils"
	"github.com/beego/wetalk/routers/admin"
	"github.com/beego/wetalk/routers/api"
	"github.com/beego/wetalk/routers/article"
	"github.com/beego/wetalk/routers/attachment"
	"github.com/beego/wetalk/routers/auth"
	"github.com/beego/wetalk/routers/base"
	"github.com/beego/wetalk/routers/post"
	"github.com/beego/wetalk/setting"

	_ "github.com/go-sql-driver/mysql"
)

// We have to call a initialize function manully
// because we use `bee bale` to pack static resources
// and we cannot make sure that which init() execute first.
func initialize() {
	setting.LoadConfig()

	if err := utils.InitSphinxPools(); err != nil {
		beego.Error(fmt.Sprint("sphinx init pool", err))
	}

	setting.SocialAuth = social.NewSocial("/login/", auth.SocialAuther)
	setting.SocialAuth.ConnectSuccessURL = "/settings/profile"
	setting.SocialAuth.ConnectFailedURL = "/settings/profile"
	setting.SocialAuth.ConnectRegisterURL = "/register/connect"
	setting.SocialAuth.LoginURL = "/login"
}

func main() {
	beego.SetLogFuncCall(true)

	initialize()

	beego.Info("AppPath:", beego.AppPath)

	if setting.IsProMode {
		beego.Info("Product mode enabled")
	} else {
		beego.Info("Develment mode enabled")
	}
	beego.Info(beego.AppName, setting.APP_VER, setting.AppUrl)

	if !setting.IsProMode {
		beego.SetStaticPath("/static_source", "static_source")
		beego.DirectoryIndex = true
	}

	// Add Filters
	beego.InsertFilter("/img/*", beego.BeforeRouter, attachment.ImageFilter)

	beego.InsertFilter("/captcha/*", beego.BeforeRouter, setting.Captcha.Handler)

	// Register routers.
	posts := new(post.PostListRouter)
	beego.Router("/", posts, "get:Home")
	beego.Router("/:slug(recent|best|cold|favs|follow)", posts, "get:Navs")
	beego.Router("/category/:slug", posts, "get:Category")
	beego.Router("/topic/:slug", posts, "get:Topic;post:TopicSubmit")

	postR := new(post.PostRouter)
	beego.Router("/new", postR, "get:New;post:NewSubmit")
	beego.Router("/post/:post([0-9]+)", postR, "get:Single;post:SingleSubmit")
	beego.Router("/post/:post([0-9]+)/edit", postR, "get:Edit;post:EditSubmit")

	if setting.NativeSearch || setting.SphinxEnabled {
		searchR := new(post.SearchRouter)
		beego.Router("/search", searchR, "get:Get")
	}

	user := new(auth.UserRouter)
	beego.Router("/user/:username/comments", user, "get:Comments")
	beego.Router("/user/:username/posts", user, "get:Posts")
	beego.Router("/user/:username/following", user, "get:Following")
	beego.Router("/user/:username/followers", user, "get:Followers")
	beego.Router("/user/:username/favs", user, "get:Favs")
	beego.Router("/user/:username", user, "get:Home")

	login := new(auth.LoginRouter)
	beego.Router("/login", login, "get:Get;post:Login")
	beego.Router("/logout", login, "get:Logout")

	beego.InsertFilter("/login/*/access", beego.BeforeRouter, auth.OAuthAccess)
	beego.InsertFilter("/login/*", beego.BeforeRouter, auth.OAuthRedirect)

	socialR := new(auth.SocialAuthRouter)
	beego.Router("/register/connect", socialR, "get:Connect;post:ConnectPost")

	register := new(auth.RegisterRouter)
	beego.Router("/register", register, "get:Get;post:Register")
	beego.Router("/active/success", register, "get:ActiveSuccess")
	beego.Router("/active/:code([0-9a-zA-Z]+)", register, "get:Active")

	settings := new(auth.SettingsRouter)
	beego.Router("/settings/profile", settings, "get:Profile;post:ProfileSave")

	forgot := new(auth.ForgotRouter)
	beego.Router("/forgot", forgot)
	beego.Router("/reset/:code([0-9a-zA-Z]+)", forgot, "get:Reset;post:ResetPost")

	upload := new(attachment.UploadRouter)
	beego.Router("/upload", upload, "post:Post")

	apiR := new(api.ApiRouter)
	beego.Router("/api/user", apiR, "post:Users")
	beego.Router("/api/md", apiR, "post:Markdown")

	adminDashboard := new(admin.AdminDashboardRouter)
	beego.Router("/admin", adminDashboard)

	adminR := new(admin.AdminRouter)
	beego.Router("/admin/model/get", adminR, "post:ModelGet")
	beego.Router("/admin/model/select", adminR, "post:ModelSelect")

	routes := map[string]beego.ControllerInterface{
		"user":     new(admin.UserAdminRouter),
		"post":     new(admin.PostAdminRouter),
		"comment":  new(admin.CommentAdminRouter),
		"topic":    new(admin.TopicAdminRouter),
		"category": new(admin.CategoryAdminRouter),
		"article":  new(admin.ArticleAdminRouter),
	}
	for name, router := range routes {
		beego.Router(fmt.Sprintf("/admin/:model(%s)", name), router, "get:List")
		beego.Router(fmt.Sprintf("/admin/:model(%s)/:id(new)", name), router, "get:Create;post:Save")
		beego.Router(fmt.Sprintf("/admin/:model(%s)/:id([0-9]+)", name), router, "get:Edit;post:Update")
		beego.Router(fmt.Sprintf("/admin/:model(%s)/:id([0-9]+)/:action(delete)", name), router, "get:Confirm;post:Delete")
	}

	// "robot.txt"
	beego.Router("/robot.txt", &base.RobotRouter{})

	articleR := new(article.ArticleRouter)
	beego.Router("/:slug", articleR, "get:Show")

	if beego.RunMode == "dev" {
		beego.Router("/test/:tmpl(mail/.*)", new(base.TestRouter))
	}

	// For all unknown pages.
	beego.Run()
}
