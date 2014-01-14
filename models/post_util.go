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

package models

import (
	"fmt"
	"regexp"

	"github.com/slene/blackfriday"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"

	"github.com/beego/wetalk/utils"
)

func ListPostsOfCategory(cat *Category, posts *[]Post) (int64, error) {
	return Posts().Filter("Category", cat).RelatedSel().OrderBy("-Updated").All(posts)
}

func ListPostsOfTopic(topic *Topic, posts *[]Post) (int64, error) {
	return Posts().Filter("Topic", topic).RelatedSel().OrderBy("-Updated").All(posts)
}

func RenderPostContent(mdStr string) string {
	htmlFlags := 0
	htmlFlags |= blackfriday.HTML_USE_XHTML
	// htmlFlags |= blackfriday.HTML_USE_SMARTYPANTS
	// htmlFlags |= blackfriday.HTML_SMARTYPANTS_FRACTIONS
	// htmlFlags |= blackfriday.HTML_SMARTYPANTS_LATEX_DASHES
	htmlFlags |= blackfriday.HTML_SKIP_HTML
	htmlFlags |= blackfriday.HTML_SKIP_STYLE
	htmlFlags |= blackfriday.HTML_SKIP_SCRIPT
	htmlFlags |= blackfriday.HTML_GITHUB_BLOCKCODE
	htmlFlags |= blackfriday.HTML_OMIT_CONTENTS
	htmlFlags |= blackfriday.HTML_COMPLETE_PAGE
	renderer := blackfriday.HtmlRenderer(htmlFlags, "", "")

	// set up the parser
	extensions := 0
	extensions |= blackfriday.EXTENSION_NO_INTRA_EMPHASIS
	extensions |= blackfriday.EXTENSION_TABLES
	extensions |= blackfriday.EXTENSION_FENCED_CODE
	extensions |= blackfriday.EXTENSION_AUTOLINK
	extensions |= blackfriday.EXTENSION_STRIKETHROUGH
	extensions |= blackfriday.EXTENSION_HARD_LINE_BREAK
	extensions |= blackfriday.EXTENSION_SPACE_HEADERS
	extensions |= blackfriday.EXTENSION_NO_EMPTY_LINE_BEFORE_BLOCK

	body := blackfriday.Markdown([]byte(mdStr), renderer, extensions)

	return string(body)
}

var mentionRegexp = regexp.MustCompile(`\B@([\d\w-_]*)`)

func FilterMentions(user *User, content string) {
	matches := mentionRegexp.FindAllStringSubmatch(content, -1)
	mentions := make([]string, 0, len(matches))
	for _, m := range matches {
		if len(m) > 1 {
			mentions = append(mentions, m[1])
		}
	}
	// var users []*User
	// num, err := Users().Filter("UserName__in", mentions).Filter("Follow__User", user.Id).All(&users)
	// if err == nil && num > 0 {
	// TODO mention email to user
	// }
}

func PostBrowsersAdd(uid int, ip string, post *Post) {
	var key string
	if uid == 0 {
		key = ip
	} else {
		key = utils.ToStr(uid)
	}
	key = fmt.Sprintf("PCA.%d.%s", post.Id, key)
	if utils.Cache.Get(key) != nil {
		return
	}
	_, err := Posts().Filter("Id", post.Id).Update(orm.Params{
		"Browsers": orm.ColValue(orm.Col_Add, 1),
	})
	if err != nil {
		beego.Error("PostCounterAdd ", err)
	}
	utils.Cache.Put(key, true, 60)
}

func PostReplysCount(post *Post) {
	cnt, err := post.Comments().Count()
	if err == nil {
		post.Replys = int(cnt)
		err = post.Update("Replys")
	}
	if err != nil {
		beego.Error("PostReplysCount ", err)
	}
}
