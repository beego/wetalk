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
	"github.com/astaxie/beego/orm"
	"github.com/beego/wetalk/models"
	"net/url"
	"strings"
	"time"

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

func filter(content string) string {
	return content
}

func mainx() {
	initialize()

	// db, _ := sql.Open("mysql", "root:root@/mygolang?charset=utf8&loc="+url.QueryEscape("Asia/shanghai"))
	// db.Query("select * from gocnbbs_common_member", ...)
	orm.RegisterDataBase("bbs", "mysql", "root:root@/mygolang?charset=utf8&loc="+url.QueryEscape("Asia/shanghai"))
	o := orm.NewOrm()
	o.Using("bbs")
	type BUser struct {
		Uid      int
		Email    string
		Username string
		Regdate  int
		User     *models.User
	}
	var busers []*BUser
	o.Raw("select * from gocnbbs_common_member").QueryRows(&busers)
	type BPost struct {
		Authorid int
		First    bool
		Dateline int64
		Subject  string
		Message  string
		Fid      int
	}
	byUsers := make(map[int]*BUser)
	for _, u := range busers {
		byUsers[u.Uid] = u
	}

	for _, u := range busers {
		uu := models.User{Email: strings.ToLower(u.Email)}
		if err := orm.NewOrm().Read(&uu, "Email"); err == nil {
			u.User = &uu
		} else {
			var name []rune
			u.Username = strings.ToLower(u.Username)
			for _, r := range u.Username {
				if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' || r == '-' || r == '_' {
					name = append(name, r)
				}
			}
			uu.NickName = u.Username
			if len(name) >= 5 {
				uu.UserName = string(name)
			} else {
				uu.UserName = string(name) + strings.Split(uu.Email, "@")[0]
				if len(uu.UserName) < 5 {
					uu.UserName = uu.UserName + "2014"
				}
			}
			uu.GrEmail = utils.EncodeMd5(uu.Email)
			if uu.UserName == "golang" || uu.UserName == "gopher" {
				uu.UserName = uu.UserName + "2014"
			}
			if err := uu.Insert(); err != nil {
				uu.UserName = uu.UserName + "2014"
				if err := uu.Insert(); err != nil {
					fmt.Println(err)
				}
			}
			u.User = &uu
		}
	}

	mapper := map[int]int{47: 5, 51: 3, 46: 4, 43: 5, 45: 6, 44: 7, 2: 8, 37: 8, 49: 8}

	var bposts []*BPost
	o.Raw("select * from gocnbbs_forum_post where fid in (47,51,46,43,45,44,2,37,49) order by tid asc, first desc, dateline asc").QueryRows(&bposts)
	c := 0
	var lPost *models.Post
	var lComment *models.Comment
	var replys int
	for _, p := range bposts {
		if byUsers[p.Authorid] == nil {
			continue
		}
		u := byUsers[p.Authorid]
		if p.First {
			if lPost != nil && replys > 0 {
				lPost.LastReply = lComment.User
				lPost.Replys = replys
				if err := lPost.Update("LastReply", "Replys"); err != nil {
					fmt.Println(err)
				}
			}

			post := new(models.Post)
			post.Title = p.Subject
			post.User = u.User
			post.LastAuthor = u.User
			post.Content = filter(p.Message)
			post.Created = time.Unix(p.Dateline, 0)
			post.Updated = post.Created
			post.Topic = &models.Topic{Id: mapper[p.Fid]}
			post.Category = &models.Category{Id: 1}
			post.Lang = 1
			if post.Topic.Id == 0 {
				fmt.Println(post.Title)
				post.Topic.Id = 8
			}
			if err := post.Insert(); err != nil {
				fmt.Println(err)
			}

			lPost = post
			replys = 0
			c++
		} else {
			comment := new(models.Comment)
			comment.Post = lPost
			comment.Message = filter(p.Message)
			comment.Created = time.Unix(p.Dateline, 0)
			comment.Floor = replys + 1
			comment.User = u.User
			if err := comment.Insert(); err != nil {
				fmt.Println(err)
			}
			lComment = comment
			replys++
		}
	}
	fmt.Println(c)
}

func main() {
	initialize()

	beego.Info("AppPath:", beego.AppPath)

	if utils.IsProMode {
		beego.Info("Product mode enabled")
	} else {
		beego.Info("Develment mode enabled")
	}
	beego.Info(beego.AppName, utils.APP_VER, utils.AppUrl)

	if !utils.IsProMode {
		beego.SetStaticPath("/static_source", "static_source")
		beego.DirectoryIndex = true
	}

	// Add Filters
	beego.AddFilter("^/img/:", "BeforRouter", routers.ImageFilter)
	beego.AddFilter("^/captcha/:", "BeforeRouter", routers.CaptchaFilter)

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
	beego.Router("/u/:username/comments", user, "get:Comments")
	beego.Router("/u/:username/posts", user, "get:Posts")
	beego.Router("/u/:username/following", user, "get:Following")
	beego.Router("/u/:username/followers", user, "get:Followers")
	beego.Router("/u/:username/favs", user, "get:Favs")
	beego.Router("/u/:username", user, "get:Home")

	login := new(routers.LoginRouter)
	beego.Router("/login", login, "get:Get;post:Login")
	beego.Router("/logout", login, "get:Logout")

	register := new(routers.RegisterRouter)
	beego.Router("/register", register, "get:Get;post:Register")
	beego.Router("/active/success", register, "get:ActiveSuccess")
	beego.Router("/active/:code([0-9a-zA-Z]+)", register, "get:Active")

	settings := new(routers.SettingsRouter)
	beego.Router("/settings/profile", settings, "get:Profile;post:ProfileSave")

	forgot := new(routers.ForgotRouter)
	beego.Router("/forgot", forgot)
	beego.Router("/reset/:code([0-9a-zA-Z]+)", forgot, "get:Reset;post:ResetPost")

	upload := new(routers.UploadRouter)
	beego.Router("/upload", upload, "post:Post")

	api := new(routers.ApiRouter)
	beego.Router("/api/user", api, "post:User")
	beego.Router("/api/post", api, "post:Post")

	adminDashboard := new(routers.AdminDashboardRouter)
	beego.Router("/admin", adminDashboard)

	admin := new(routers.AdminRouter)
	beego.Router("/admin/model/get", admin, "post:ModelGet")
	beego.Router("/admin/model/select", admin, "post:ModelSelect")

	routes := map[string]beego.ControllerInterface{
		"user":     new(routers.UserAdminRouter),
		"post":     new(routers.PostAdminRouter),
		"comment":  new(routers.CommentAdminRouter),
		"topic":    new(routers.TopicAdminRouter),
		"category": new(routers.CategoryAdminRouter),
		"article":  new(routers.ArticleAdminRouter),
	}
	for name, router := range routes {
		beego.Router(fmt.Sprintf("/admin/:model(%s)", name), router, "get:List")
		beego.Router(fmt.Sprintf("/admin/:model(%s)/:id(new)", name), router, "get:Create;post:Save")
		beego.Router(fmt.Sprintf("/admin/:model(%s)/:id([0-9]+)", name), router, "get:Edit;post:Update")
		beego.Router(fmt.Sprintf("/admin/:model(%s)/:id([0-9]+)/:action(delete)", name), router, "get:Confirm;post:Delete")
	}

	// "robot.txt"
	beego.Router("/robot.txt", &routers.RobotRouter{})

	article := new(routers.ArticleRouter)
	beego.Router("/:slug([0-9a-z-./]+)", article, "get:Show")

	if beego.RunMode == "dev" {
		beego.Router("/test/:tmpl(mail/.*)", new(routers.TestRouter))
	}

	// For all unknown pages.
	beego.Run()
}
