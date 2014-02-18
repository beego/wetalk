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
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"

	"github.com/beego/wetalk/modules/models"
)

func ListCategories(cats *[]models.Category) (int64, error) {
	return models.Categories().OrderBy("-order").All(cats)
}

func ListTopics(topics *[]models.Topic) (int64, error) {
	return models.Topics().OrderBy("-Followers").All(topics)
}

func ListTopicsOfCat(topics *[]models.Topic, cat *models.Category) (int64, error) {
	var list orm.ParamsList
	var where string
	if cat != nil {
		where = " WHERE category_id = ?"
	}

	sql := fmt.Sprintf(`SELECT topic_id
		FROM post%s
		GROUP BY topic_id
		ORDER BY COUNT(topic_id) DESC LIMIT 8`, where)

	rs := orm.NewOrm().Raw(sql)

	if cat != nil {
		rs = rs.SetArgs(cat.Id)
	}

	cnt, err := rs.ValuesFlat(&list)
	if err != nil {
		beego.Error("models.ListTopicsOfCat ", err)
		return 0, err
	}
	if cnt > 0 {
		nums, err := models.Topics().Filter("Id__in", list).All(topics)
		if err != nil {
			beego.Error("models.ListTopicsOfCat ", err)
			return 0, err
		}
		return nums, err
	}
	return 0, nil
}
