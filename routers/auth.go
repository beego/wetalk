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

import ()

// LoginRouter serves login page.
type LoginRouter struct {
	baseRouter
}

// Get implemented Get method for LoginRouter.
func (this *LoginRouter) Get() {
	this.TplNames = "login.html"
}

// RegisterRouter serves login page.
type RegisterRouter struct {
	baseRouter
}

// Get implemented Get method for RegisterRouter.
func (this *RegisterRouter) Get() {
	this.TplNames = "register.html"
}

// ForgotRouter serves login page.
type ForgotRouter struct {
	baseRouter
}

// Get implemented Get method for ForgotRouter.
func (this *ForgotRouter) Get() {
	this.TplNames = "forgot.html"
}

// ResetRouter serves login page.
type ResetRouter struct {
	baseRouter
}

// Get implemented Get method for ResetRouter.
func (this *ResetRouter) Get() {
	this.TplNames = "reset.html"
}
