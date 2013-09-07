// Copyright 2013 beebbs authors
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

package routers

import (
	"github.com/beego/beebbs/utils"
)

// HomeRouter serves home page.
type HomeRouter struct {
	baseRouter
}

// Get implemented Get method for HomeRouter.
func (this *HomeRouter) Get() {
	this.Data["IsHome"] = true
	this.TplNames = "home.html"

	// Get language.
	this.Data["Hello"] = utils.I18n(this.langVer, "hello, %s", "joe")
	this.Data["Hey"] = "hello, %s"
	this.Data["Name"] = "jim"
}
