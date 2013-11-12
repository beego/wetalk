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
	"strings"

	"github.com/astaxie/beego/orm"
)

// A true/false field.
type SliceStringField []string

func (e SliceStringField) Value() []string {
	return []string(e)
}

func (e *SliceStringField) Set(d []string) {
	*e = SliceStringField(d)
}

func (e *SliceStringField) Add(v string) {
	*e = append(*e, v)
}

func (e *SliceStringField) String() string {
	return strings.Join(e.Value(), ",")
}

func (e *SliceStringField) FieldType() int {
	return orm.TypeCharField
}

func (e *SliceStringField) SetRaw(value interface{}) error {
	switch d := value.(type) {
	case []string:
		e.Set(d)
	case string:
		if len(d) > 0 {
			parts := strings.Split(d, ",")
			v := make([]string, 0, len(parts))
			for _, p := range parts {
				v = append(v, strings.TrimSpace(p))
			}
			e.Set(v)
		}
	default:
		return fmt.Errorf("<SliceStringField.SetRaw> unknown value `%v`", value)
	}
	return nil
}

func (e *SliceStringField) RawValue() interface{} {
	return e.String()
}

func (e *SliceStringField) Clean() error {
	return nil
}

var _ orm.Fielder = new(SliceStringField)
