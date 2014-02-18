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
	"fmt"
	"strings"

	"github.com/astaxie/beego/orm"

	"github.com/beego/wetalk/modules/models"
	"github.com/beego/wetalk/modules/utils"
	"github.com/beego/wetalk/setting"
)

var searchEscapePattern = []string{
	`\\`, `(`, `)`, `|`, `-`, `!`, `@`, `~`, `'`, `&`, `/`, `^`, `$`, `=`,
	`\\\\`, `\(`, `\)`, `\|`, `\-`, `\!`, `\@`, `\~`, `\'`, `\&`, `\/`, `\^`, `\$`, `\=`,
}

func filterSearchQ(q string) string {
	q = strings.TrimSpace(q)
	replacer := strings.NewReplacer(searchEscapePattern...)
	return replacer.Replace(q)
}

func SearchPost(q string, page int) ([]*models.Post, *utils.SphinxMeta, error) {
	q = filterSearchQ(q)
	if len(q) == 0 {
		return nil, nil, fmt.Errorf("empty query")
	}

	sdb, err := utils.SphinxPools.GetConn()
	if err != nil {
		return nil, nil, err
	}
	defer sdb.Close()

	pers := 20
	if page == 0 {
		page = 1
	}
	offset := (page - 1) * pers

	var pids orm.ParamsList
	query := fmt.Sprintf(`SELECT @id AS pid, updated
		FROM `+setting.SphinxIndex+`
		WHERE MATCH('`+q+`')
		ORDER BY @weight DESC, updated DESC
		LIMIT %d, %d OPTION ranker=bm25`, offset, pers)

	if _, err = sdb.RawValuesFlat(query, &pids, "pid"); err != nil {
		return nil, nil, err
	}

	var meta *utils.SphinxMeta
	if meta, err = sdb.ShowMeta(); err != nil {
		return nil, nil, err
	}
	sdb.Close()

	if len(pids) == 0 {
		return nil, meta, nil
	}

	var posts []*models.Post
	_, err = models.Posts().Filter("Id__in", pids).RelatedSel().All(&posts)
	if err != nil {
		return nil, nil, err
	}

	return posts, meta, nil
}
