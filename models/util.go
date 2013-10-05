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
	"github.com/astaxie/beego/orm"
)

func CheckIsExist(qs orm.QuerySeter, field string, value interface{}, skipId int) bool {
	qs = qs.Filter(field, value)
	if skipId > 0 {
		qs = qs.Exclude("Id", skipId)
	}
	cnt, _ := qs.Count()
	return cnt > 0
}
