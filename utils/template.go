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

package utils

import (
	"errors"
	"html/template"
	"reflect"
	"time"

	"github.com/astaxie/beego"
	"github.com/beego/i18n"
)

// get HTML i18n string
func i18nHTML(lang, format string, args ...interface{}) template.HTML {
	return template.HTML(i18n.Tr(lang, format, args...))
}

// get HTML i18n string with specify section
func i18nsHTML(lang, section, format string, args ...interface{}) template.HTML {
	return template.HTML(i18n.Trs(lang, section, format, args...))
}

func boolicon(b bool) (s template.HTML) {
	if b {
		s = `<i style="color:green;" class="icon-check""></i>`
	} else {
		s = `<i class="icon-check-empty""></i>`
	}
	return
}

func date(t time.Time) string {
	return beego.Date(t, DateFormat)
}

func datetime(t time.Time) string {
	return beego.Date(t, DateTimeFormat)
}

func loadtimes(t time.Time) int {
	return int(time.Now().Sub(t).Nanoseconds() / 1e6)
}

func dict(values ...interface{}) (map[string]interface{}, error) {
	if len(values)%2 != 0 {
		return nil, errors.New("invalid dict call")
	}
	dict := make(map[string]interface{}, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, errors.New("dict keys must be strings")
		}
		dict[key] = values[i+1]
	}
	return dict, nil
}

func timesince(lang string, t time.Time) string {
	now := time.Now()
	seonds := int(now.Sub(t).Seconds())
	if seonds < 60 {
		return i18n.Tr(lang, "%d seconds ago", seonds)
	} else if seonds < 60*60 {
		return i18n.Tr(lang, "%d minutes ago", seonds/60)
	} else if seonds < 60*60*24 {
		return i18n.Tr(lang, "%d hours ago", seonds/(60*60))
	} else {
		return i18n.Tr(lang, "%d days ago", seonds/(60*60*24))
	}
}

func init() {
	// Register template functions.
	beego.AddFuncMap("i18n", i18nHTML)
	beego.AddFuncMap("i18ns", i18nsHTML)
	beego.AddFuncMap("boolicon", boolicon)
	beego.AddFuncMap("date", date)
	beego.AddFuncMap("datetime", datetime)
	beego.AddFuncMap("dict", dict)
	beego.AddFuncMap("timesince", timesince)

	// move go1.2 template funcs to this
	// Comparisons
	beego.AddFuncMap("eq", eq) // ==
	beego.AddFuncMap("ge", ge) // >=
	beego.AddFuncMap("gt", gt) // >
	beego.AddFuncMap("le", le) // <=
	beego.AddFuncMap("lt", lt) // <
	beego.AddFuncMap("ne", ne) // !=

	beego.AddFuncMap("loadtimes", loadtimes)
}

var (
	errBadComparisonType = errors.New("invalid type for comparison")
	errBadComparison     = errors.New("incompatible types for comparison")
	errNoComparison      = errors.New("missing argument for comparison")
)

type kind int

const (
	invalidKind kind = iota
	boolKind
	complexKind
	intKind
	floatKind
	integerKind
	stringKind
	uintKind
)

func basicKind(v reflect.Value) (kind, error) {
	switch v.Kind() {
	case reflect.Bool:
		return boolKind, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return intKind, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return uintKind, nil
	case reflect.Float32, reflect.Float64:
		return floatKind, nil
	case reflect.Complex64, reflect.Complex128:
		return complexKind, nil
	case reflect.String:
		return stringKind, nil
	}
	return invalidKind, errBadComparisonType
}

// eq evaluates the comparison a == b || a == c || ...
func eq(arg1 interface{}, arg2 ...interface{}) (bool, error) {
	v1 := reflect.ValueOf(arg1)
	k1, err := basicKind(v1)
	if err != nil {
		return false, err
	}
	if len(arg2) == 0 {
		return false, errNoComparison
	}
	for _, arg := range arg2 {
		v2 := reflect.ValueOf(arg)
		k2, err := basicKind(v2)
		if err != nil {
			return false, err
		}
		if k1 != k2 {
			return false, errBadComparison
		}
		truth := false
		switch k1 {
		case boolKind:
			truth = v1.Bool() == v2.Bool()
		case complexKind:
			truth = v1.Complex() == v2.Complex()
		case floatKind:
			truth = v1.Float() == v2.Float()
		case intKind:
			truth = v1.Int() == v2.Int()
		case stringKind:
			truth = v1.String() == v2.String()
		case uintKind:
			truth = v1.Uint() == v2.Uint()
		default:
			panic("invalid kind")
		}
		if truth {
			return true, nil
		}
	}
	return false, nil
}

// ne evaluates the comparison a != b.
func ne(arg1, arg2 interface{}) (bool, error) {
	// != is the inverse of ==.
	equal, err := eq(arg1, arg2)
	return !equal, err
}

// lt evaluates the comparison a < b.
func lt(arg1, arg2 interface{}) (bool, error) {
	v1 := reflect.ValueOf(arg1)
	k1, err := basicKind(v1)
	if err != nil {
		return false, err
	}
	v2 := reflect.ValueOf(arg2)
	k2, err := basicKind(v2)
	if err != nil {
		return false, err
	}
	if k1 != k2 {
		return false, errBadComparison
	}
	truth := false
	switch k1 {
	case boolKind, complexKind:
		return false, errBadComparisonType
	case floatKind:
		truth = v1.Float() < v2.Float()
	case intKind:
		truth = v1.Int() < v2.Int()
	case stringKind:
		truth = v1.String() < v2.String()
	case uintKind:
		truth = v1.Uint() < v2.Uint()
	default:
		panic("invalid kind")
	}
	return truth, nil
}

// le evaluates the comparison <= b.
func le(arg1, arg2 interface{}) (bool, error) {
	// <= is < or ==.
	lessThan, err := lt(arg1, arg2)
	if lessThan || err != nil {
		return lessThan, err
	}
	return eq(arg1, arg2)
}

// gt evaluates the comparison a > b.
func gt(arg1, arg2 interface{}) (bool, error) {
	// > is the inverse of <=.
	lessOrEqual, err := le(arg1, arg2)
	if err != nil {
		return false, err
	}
	return !lessOrEqual, nil
}

// ge evaluates the comparison a >= b.
func ge(arg1, arg2 interface{}) (bool, error) {
	// >= is the inverse of <.
	lessThan, err := lt(arg1, arg2)
	if err != nil {
		return false, err
	}
	return !lessThan, nil
}
