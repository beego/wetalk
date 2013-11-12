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

func UserFollow(user *User, theUser *User) {
	if theUser.Read() == nil {
		var mutual bool
		tFollow := Follow{User: theUser, FollowUser: user}
		if err := tFollow.Read("User", "FollowUser"); err == nil {
			mutual = true
		}

		follow := Follow{User: user, FollowUser: theUser, Mutual: mutual}
		if err := follow.Insert(); err == nil && mutual {
			tFollow.Mutual = mutual
			tFollow.Update("Mutual")
		}

		if nums, err := user.FollowingUsers().Count(); err == nil {
			user.Following = int(nums)
			user.Update("Following")
		}

		if nums, err := theUser.FollowerUsers().Count(); err == nil {
			theUser.Followers = int(nums)
			theUser.Update("Followers")
		}
	}
}

func UserUnFollow(user *User, theUser *User) {
	num, _ := user.FollowingUsers().Filter("FollowUser", theUser.Id).Delete()
	if num > 0 {
		theUser.FollowingUsers().Filter("FollowUser", user.Id).Update(orm.Params{
			"Mutual": false,
		})

		if nums, err := user.FollowingUsers().Count(); err == nil {
			user.Following = int(nums)
			user.Update("Following")
		}

		if nums, err := theUser.FollowerUsers().Count(); err == nil {
			theUser.Followers = int(nums)
			theUser.Update("Followers")
		}
	}
}
