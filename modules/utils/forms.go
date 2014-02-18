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
	"fmt"
	"html/template"
	"math"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"

	"github.com/beego/wetalk/setting"
)

func init() {
	initCommonField()
	initExtraField()
}

type FormLocaler interface {
	Tr(string, ...interface{}) string
}

type FormHelper interface {
	Helps() map[string]string
}

type FormLabeler interface {
	Labels() map[string]string
}

type FormPlaceholder interface {
	Placeholders() map[string]string
}

type FieldCreater func(*FieldSet)

type FieldFilter func(*FieldSet)

var customCreaters = make(map[string]FieldCreater)

var customFilters = make(map[string]FieldFilter)

type fakeLocale struct{}

func (*fakeLocale) Tr(text string, args ...interface{}) string {
	return text
}

var fakeLocaler FormLocaler = new(fakeLocale)

// register a custom label/input creater
func RegisterFieldCreater(name string, field FieldCreater) {
	customCreaters[name] = field
}

// register a custom label/input creater
func RegisterFieldFilter(name string, field FieldFilter) {
	customFilters[name] = field
}

type HtmlLazyField func() template.HTML

func (f HtmlLazyField) String() string {
	return string(f())
}

type FieldSet struct {
	Label       template.HTML
	Field       HtmlLazyField
	Id          string
	Name        string
	LabelText   string
	Value       interface{}
	Help        string
	Error       string
	Type        string
	Kind        string
	Placeholder string
	Attrs       string
	FormElm     reflect.Value
	Locale      FormLocaler
}

type FormSets struct {
	FieldList []*FieldSet
	Fields    map[string]*FieldSet
	Locale    FormLocaler
	inited    bool
	form      interface{}
	errs      map[string]*validation.ValidationError
}

func (this *FormSets) SetError(fieldName, errMsg string) {
	if fSet, ok := this.Fields[fieldName]; ok {
		fSet.Error = this.Locale.Tr(errMsg)
	}
}

// create formSets for generate label/field html code
func NewFormSets(form interface{}, errs map[string]*validation.ValidationError, locale FormLocaler) *FormSets {
	fSets := new(FormSets)
	fSets.errs = errs
	fSets.Fields = make(map[string]*FieldSet)
	if locale != nil {
		fSets.Locale = locale
	} else {
		fSets.Locale = fakeLocaler
	}

	val := reflect.ValueOf(form)

	panicAssertStructPtr(val)

	elm := val.Elem()

	var helps map[string]string
	var labels map[string]string
	var places map[string]string

	// get custom field helo messages
	if f, ok := form.(FormHelper); ok {
		hlps := f.Helps()
		if hlps != nil {
			helps = hlps
		}
	}

	// ge custom field labels
	if f, ok := form.(FormLabeler); ok {
		lbls := f.Labels()
		if lbls != nil {
			labels = lbls
		}
	}

	// ge custom field placeholders
	if f, ok := form.(FormPlaceholder); ok {
		phs := f.Placeholders()
		if phs != nil {
			places = phs
		}
	}

outFor:
	for i := 0; i < elm.NumField(); i++ {
		f := elm.Field(i)
		fT := elm.Type().Field(i)

		name := fT.Name
		value := f.Interface()
		fTyp := "text"

		switch f.Kind() {
		case reflect.Bool:
			fTyp = "checkbox"
		default:
			switch value.(type) {
			case time.Time:
				fTyp = "datetime"
			}
		}

		fName := name

		var attrm map[string]string

		// parse struct tag settings
		for _, v := range strings.Split(fT.Tag.Get("form"), ";") {
			v = strings.TrimSpace(v)
			if v == "-" {
				continue outFor
			} else if i := strings.Index(v, "("); i > 0 && strings.Index(v, ")") == len(v)-1 {
				tN := v[:i]
				v = strings.TrimSpace(v[i+1 : len(v)-1])
				switch tN {
				case "type":
					fTyp = v
				case "name":
					fName = v
				case "attr":
					if attrm == nil {
						attrm = make(map[string]string)
					}
					parts := strings.SplitN(v, ",", 2)
					if len(parts) > 1 {
						attrm[parts[0]] = parts[1]
					} else {
						attrm[v] = v
					}
				}
			}
		}

		var attrs string
		if attrm != nil {
			for k, v := range attrm {
				attrs += fmt.Sprintf(` %s="%s"`, k, v)
			}
		}

		// set field id
		fId := elm.Type().Name() + "-" + fName

		var fSet FieldSet

		fSet.Id = fId
		fSet.Name = fName
		fSet.Value = value
		fSet.Attrs = attrs
		fSet.FormElm = elm
		fSet.Locale = locale

		if i := strings.IndexRune(fTyp, ','); i != -1 {
			fSet.Type = fTyp[:i]
			fSet.Kind = fTyp[i+1:]
			fTyp = fSet.Type
		} else {
			fSet.Type = fTyp
			fSet.Kind = fTyp
		}

		// get field label text
		fSet.LabelText = fName
		if labels != nil {
			if _, ok := labels[name]; ok {
				fSet.LabelText = labels[name]
			}
		}
		fSet.LabelText = locale.Tr(fSet.LabelText)

		// get field help
		if helps != nil {
			if _, ok := helps[name]; ok {
				fSet.Help = helps[name]
			}
		}
		fSet.Help = locale.Tr(helps[name])

		if places != nil {
			if _, ok := places[name]; ok {
				fSet.Placeholder = places[name]
			}
		}
		fSet.Placeholder = locale.Tr(fSet.Placeholder)

		if len(fSet.Placeholder) > 0 {
			fSet.Placeholder = fmt.Sprintf(` placeholder="%s"`, fSet.Placeholder)
		}

		// create error string
		if errs != nil {
			if err, ok := errs[name]; ok {
				fSet.Error = locale.Tr(err.Tmpl, err.LimitValue)
			}
		}

		// create label html
		switch fTyp {
		case "checkbox", "hidden":
		default:
			fSet.Label = template.HTML(fmt.Sprintf(`
					<label class="control-label" for="%s">%s</label>`, fSet.Id, fSet.LabelText))
		}

		if creater, ok := customCreaters[fTyp]; ok {
			// use custome creater generate label/input html
			creater(&fSet)

			if filter, ok := customFilters[fTyp]; ok {
				// use custome filter replace label/input html
				filter(&fSet)
			}
		}

		if fSet.Field == nil {
			fSet.Field = func() template.HTML { return "" }
		}

		fSets.FieldList = append(fSets.FieldList, &fSet)
		fSets.Fields[name] = &fSet
	}

	fSets.inited = true

	return fSets
}

func initCommonField() {
	RegisterFieldCreater("text", func(fSet *FieldSet) {
		fSet.Field = func() template.HTML {
			return template.HTML(fmt.Sprintf(
				`<input id="%s" name="%s" type="text" value="%v" class="form-control"%s%s>`,
				fSet.Id, fSet.Name, fSet.Value, fSet.Placeholder, fSet.Attrs))
		}
	})

	RegisterFieldCreater("textarea", func(fSet *FieldSet) {
		fSet.Field = func() template.HTML {
			return template.HTML(fmt.Sprintf(
				`<textarea id="%s" name="%s" rows="5" class="form-control"%s%s>%v</textarea>`,
				fSet.Id, fSet.Name, fSet.Placeholder, fSet.Attrs, fSet.Value))
		}
	})

	RegisterFieldCreater("password", func(fSet *FieldSet) {
		fSet.Field = func() template.HTML {
			return template.HTML(fmt.Sprintf(
				`<input id="%s" name="%s" type="password" value="%v" class="form-control"%s%s>`,
				fSet.Id, fSet.Name, fSet.Value, fSet.Placeholder, fSet.Attrs))
		}
	})

	RegisterFieldCreater("hidden", func(fSet *FieldSet) {
		fSet.Field = func() template.HTML {
			return template.HTML(fmt.Sprintf(
				`<input id="%s" name="%s" type="hidden" value="%v"%s>`, fSet.Id, fSet.Name, fSet.Value, fSet.Attrs))
		}
	})

	datetimeFunc := func(fSet *FieldSet) {
		fSet.Field = func() template.HTML {
			t := fSet.Value.(time.Time)
			tval := beego.Date(t, setting.DateTimeFormat)
			if tval == "0001-01-01 00:00:00" {
				t = time.Now()
			}
			if fSet.Type == "date" {
				tval = beego.Date(t, setting.DateFormat)
			}
			return template.HTML(fmt.Sprintf(
				`<input id="%s" name="%s" type="%s" value="%s" class="form-control"%s%s>`,
				fSet.Id, fSet.Name, fSet.Type, tval, fSet.Placeholder, fSet.Attrs))
		}
	}

	RegisterFieldCreater("date", datetimeFunc)
	RegisterFieldCreater("datetime", datetimeFunc)

	RegisterFieldCreater("checkbox", func(fSet *FieldSet) {
		fSet.Field = func() template.HTML {
			var checked string
			if b, ok := fSet.Value.(bool); ok && b {
				checked = "checked"
			}
			return template.HTML(fmt.Sprintf(
				`<label for="%s" class="checkbox">%s<input id="%s" name="%s" type="checkbox" %s></label>`,
				fSet.Id, fSet.LabelText, fSet.Id, fSet.Name, checked))
		}
	})

	RegisterFieldCreater("select", func(fSet *FieldSet) {
		fSet.Field = func() template.HTML {
			var options string
			str := fmt.Sprintf(`<select id="%s" name="%s" class="form-control"%s%s>%s</select>`,
				fSet.Id, fSet.Name, fSet.Placeholder, fSet.Attrs)

			fun := fSet.FormElm.Addr().MethodByName(fSet.Name + "SelectData")

			if fun.IsValid() {
				results := fun.Call([]reflect.Value{})
				if len(results) > 0 {
					v := results[0]
					if v.CanInterface() {
						if vu, ok := v.Interface().([][]string); ok {

							var vs []string
							val := reflect.ValueOf(fSet.Value)
							if val.Kind() == reflect.Slice {
								vs = make([]string, 0, val.Len())
								for i := 0; i < val.Len(); i++ {
									vs = append(vs, ToStr(val.Index(i).Interface()))
								}
							}

							isMulti := len(vs) > 0
							for _, parts := range vu {
								var n, v string
								switch {
								case len(parts) > 1:
									n, v = fSet.Locale.Tr(parts[0]), parts[1]
								case len(parts) == 1:
									n, v = fSet.Locale.Tr(parts[0]), parts[0]
								}
								var selected string
								if isMulti {
									for _, e := range vs {
										if e == v {
											selected = ` selected="selected"`
											break
										}
									}
								} else if ToStr(fSet.Value) == v {
									selected = ` selected="selected"`
								}
								options += fmt.Sprintf(`<option value="%s"%s>%s</option>`, v, selected, n)
							}
						}
					}
				}
			}

			if len(options) == 0 {
				options = fmt.Sprintf(`<option value="%v">%v</option>`, fSet.Value, fSet.Value)
			}

			return template.HTML(fmt.Sprintf(str, options))
		}
	})
}

func initExtraField() {
	// a example to create a bootstrap style checkbox creater
	RegisterFieldCreater("checkbox", func(fSet *FieldSet) {
		fSet.Field = func() template.HTML {
			value := false
			if b, ok := fSet.Value.(bool); ok {
				value = b
			}
			active := ""
			if value {
				active = " active"
			}
			return template.HTML(fmt.Sprintf(`<label>
            <input type="hidden" name="%s" value="%v">
            <button type="button" data-toggle="button" data-name="%s" class="btn btn-default btn-xs btn-checked%s">
            	<i class="icon icon-ok"></i>
        	</button>%s
        </label>`, fSet.Name, value, fSet.Name, active, fSet.LabelText))
		}
	})

	// a example to create a select2 box
	RegisterFieldFilter("select", func(fSet *FieldSet) {
		if strings.Index(fSet.Attrs, `rel="select2"`) != -1 {
			field := fSet.Field.String()
			fSet.Field = func() template.HTML {
				field = strings.Replace(field, "<option", "<option></option><option", 1)
				return template.HTML(field)
			}
		}
	})

	// a example to create a captcha field
	RegisterFieldCreater("captcha", func(fSet *FieldSet) {
		fSet.Label = template.HTML(fmt.Sprintf(`
					<label class="control-label" for="%s">%s</label>`, fSet.Id, fSet.LabelText))

		fSet.Field = func() template.HTML {
			return template.HTML(fmt.Sprintf(`%v
			<input id="%s" name="%s" type="text" value="" class="form-control" autocomplete="off"%s%s>`,
				setting.Captcha.CreateCaptchaHtml(), fSet.Id, fSet.Name, fSet.Placeholder, fSet.Attrs))
		}
	})

	RegisterFieldCreater("empty", func(fSet *FieldSet) {
		fSet.Label = ""
		fSet.Field = func() template.HTML { return "" }
	})
}

// parse request.Form values to form
func ParseForm(form interface{}, values url.Values) {
	val := reflect.ValueOf(form)
	elm := reflect.Indirect(val)

	panicAssertStructPtr(val)

outFor:
	for i := 0; i < elm.NumField(); i++ {
		f := elm.Field(i)
		fT := elm.Type().Field(i)

		fName := fT.Name

		for _, v := range strings.Split(fT.Tag.Get("form"), ";") {
			v = strings.TrimSpace(v)
			if v == "-" {
				continue outFor
			} else if i := strings.Index(v, "("); i > 0 && strings.Index(v, ")") == len(v)-1 {
				tN := v[:i]
				v = strings.TrimSpace(v[i+1 : len(v)-1])
				switch tN {
				case "name":
					fName = v
				}
			}
		}

		value := ""
		var vs []string
		if v, ok := values[fName]; ok {
			vs = v
			if len(v) > 0 {
				value = v[0]
			}
		}

		switch fT.Type.Kind() {
		case reflect.Bool:
			b, _ := StrTo(value).Bool()
			f.SetBool(b)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			x, _ := StrTo(value).Int64()
			f.SetInt(x)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			x, _ := StrTo(value).Uint64()
			f.SetUint(x)
		case reflect.Float32, reflect.Float64:
			x, _ := StrTo(value).Float64()
			f.SetFloat(x)
		case reflect.Struct:
			if fT.Type.String() == "time.Time" {
				if len(value) > 10 {
					t, err := beego.DateParse(value, setting.DateTimeFormat)
					if err != nil {
						continue
					}
					f.Set(reflect.ValueOf(t))
				} else {
					t, err := beego.DateParse(value, setting.DateFormat)
					if err != nil {
						continue
					}
					f.Set(reflect.ValueOf(t))
				}
			}
		case reflect.String:
			f.SetString(value)
		case reflect.Slice:
			f.Set(reflect.ValueOf(vs))
		}
	}
}

// assert an object must be a struct pointer
func panicAssertStructPtr(val reflect.Value) {
	if val.Kind() == reflect.Ptr && val.Elem().Kind() == reflect.Struct {
		return
	}
	panic(fmt.Errorf("%s must be a struct pointer", val.Type().Name()))
}

// set values from one struct to other struct
// both need ptr struct
func SetFormValues(from interface{}, to interface{}, skips ...string) {
	val := reflect.ValueOf(from)
	elm := reflect.Indirect(val)

	valTo := reflect.ValueOf(to)
	elmTo := reflect.Indirect(valTo)

	panicAssertStructPtr(val)
	panicAssertStructPtr(valTo)

outFor:
	for i := 0; i < elmTo.NumField(); i++ {
		toF := elmTo.Field(i)
		name := elmTo.Type().Field(i).Name

		// skip specify field
		for _, skip := range skips {
			if skip == name {
				continue outFor
			}
		}
		f := elm.FieldByName(name)
		if f.Kind() != reflect.Invalid {
			// set value if type matched
			if f.Type().String() == toF.Type().String() {
				toF.Set(f)
			} else {
				fInt := false
				switch f.Interface().(type) {
				case int, int8, int16, int32, int64:
					fInt = true
				case uint, uint8, uint16, uint32, uint64:
				default:
					continue outFor
				}
				switch toF.Interface().(type) {
				case int, int8, int16, int32, int64:
					var v int64
					if fInt {
						v = f.Int()
					} else {
						vu := f.Uint()
						if vu > math.MaxInt64 {
							continue outFor
						}
						v = int64(vu)
					}
					if toF.OverflowInt(v) {
						continue outFor
					}
					toF.SetInt(v)
				case uint, uint8, uint16, uint32, uint64:
					var v uint64
					if fInt {
						vu := f.Int()
						if vu < 0 {
							continue outFor
						}
						v = uint64(vu)
					} else {
						v = f.Uint()
					}
					if toF.OverflowUint(v) {
						continue outFor
					}
					toF.SetUint(v)
				}
			}
		}
	}
}

// compare field values between two struct pointer
// return changed field names
func FormChanges(base interface{}, modified interface{}, skips ...string) (fields []string) {
	val := reflect.ValueOf(base)
	elm := reflect.Indirect(val)

	valMod := reflect.ValueOf(modified)
	elmMod := reflect.Indirect(valMod)

	panicAssertStructPtr(val)
	panicAssertStructPtr(valMod)

outFor:
	for i := 0; i < elmMod.NumField(); i++ {
		modF := elmMod.Field(i)
		name := elmMod.Type().Field(i).Name

		fT := elmMod.Type().Field(i)

		for _, v := range strings.Split(fT.Tag.Get("form"), ";") {
			v = strings.TrimSpace(v)
			if v == "-" {
				continue outFor
			}
		}

		// skip specify field
		for _, skip := range skips {
			if skip == name {
				continue outFor
			}
		}
		f := elm.FieldByName(name)
		if f.Kind() == reflect.Invalid {
			continue
		}

		// compare two values use string
		if ToStr(modF.Interface()) != ToStr(f.Interface()) {
			fields = append(fields, name)
		}
	}

	return
}
