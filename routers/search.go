package routers

import (
	"strings"

	"github.com/astaxie/beego"

	"github.com/beego/wetalk/models"
	"github.com/beego/wetalk/utils"
)

type SearchRouter struct {
	baseRouter
}

func (this *SearchRouter) Get() {
	this.TplNames = "search/posts.html"

	pers := 25

	q := strings.TrimSpace(this.GetString("q"))

	this.Data["Q"] = q

	if len(q) == 0 {
		return
	}

	page, _ := utils.StrTo(this.GetString("p")).Int()

	posts, meta, err := models.SearchPost(q, page)
	if err != nil {
		this.Data["SearchError"] = true
		beego.Error("SearchPosts: ", err)
		return
	}

	if len(posts) > 0 {
		this.SetPaginator(pers, meta.TotalFound)
	}

	this.Data["Posts"] = posts
	this.Data["Meta"] = meta
}
