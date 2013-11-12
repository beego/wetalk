package routers

import (
	"github.com/astaxie/beego/orm"

	"github.com/beego/wetalk/models"
	"github.com/beego/wetalk/utils"
)

type ApiRouter struct {
	baseRouter
}

func (this *ApiRouter) User() {
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

	if this.isLogin {

		switch action {
		case "get-follows":
			var data []orm.ParamsList
			this.user.FollowingUsers().ValuesList(&data, "FollowUser__NickName", "FollowUser__UserName")
			result["success"] = true
			result["data"] = data

		case "follow", "unfollow":
			id, err := utils.StrTo(this.GetString("user")).Int()
			if err == nil && id != this.user.Id {
				fuser := models.User{Id: id}
				if action == "follow" {
					models.UserFollow(&this.user, &fuser)
				} else {
					models.UserUnFollow(&this.user, &fuser)
				}
				result["success"] = true
			}
		}
	}
}

func (this *ApiRouter) Post() {
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

	if this.isLogin && this.user.IsAdmin {

		switch action {
		case "toggle-best":
			id, _ := utils.StrTo(this.GetString("post")).Int()
			if id > 0 {
				post := models.Post{Id: id}
				if err := post.Read(); err == nil {
					post.IsBest = !post.IsBest
					post.Update("IsBest")
					result["success"] = true
					result["isBest"] = post.IsBest
				}
			}
		}
	}
}
