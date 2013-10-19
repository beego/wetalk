package utils

import (
	"fmt"
	"html/template"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
)

func init() {
	// a example to create a bootstrap style checkbox creater
	RegisterFieldCreater("checkbox", func(fSet *FieldSet) {
		value := false
		if b, ok := fSet.Value.(bool); ok {
			value = b
		}
		active := ""
		if value {
			active = " active"
		}
		fSet.Field = template.HTML(fmt.Sprintf(`<label>
            <input type="hidden" name="%s" value="%v">
            <button type="button" data-toggle="button" data-name="%s" class="btn btn-default btn-xs btn-checked%s">
            	<i class="icon icon-ok"></i>
        	</button>%s
        </label>`, fSet.Name, value, fSet.Name, active, fSet.LabelText))
	})

	// a example to create a select2 box
	RegisterFieldFilter("select", func(fSet *FieldSet) {
		if strings.Index(fSet.Attrs, `rel="select2"`) != -1 {
			field := string(fSet.Field)
			field = strings.Replace(field, "<option", "<option></option><option", 1)
			fSet.Field = template.HTML(field)
		}
	})
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

type FieldSet struct {
	Label       template.HTML
	Field       template.HTML
	Id          string
	Name        string
	LabelText   string
	Value       interface{}
	Help        string
	Error       string
	Type        string
	Placeholder string
	Attrs       string
}

type FormSets struct {
	FieldList []*FieldSet
	Fields    map[string]*FieldSet
	Locale    FormLocaler
}

func (this *FormSets) SetError(fieldName, errMsg string) {
	if fSet, ok := this.Fields[fieldName]; ok {
		fSet.Error = this.Locale.Tr(errMsg)
	}
}

// create formSets for generate label/field html code
func NewFormSets(form interface{}, errs map[string]*validation.ValidationError, locale FormLocaler) *FormSets {
	fSets := new(FormSets)
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
		fSet.Type = fTyp
		fSet.Attrs = attrs

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

		placeholders := ""
		if len(fSet.Placeholder) > 0 {
			placeholders = fmt.Sprintf(` placeholder="%s"`, fSet.Placeholder)
		}

		// create error string
		if errs != nil {
			if err, ok := errs[name]; ok {
				fSet.Error = locale.Tr(err.Tmpl, err.LimitValue)
			}
		}

		if creater, ok := customCreaters[fTyp]; ok {
			// use custome creater generate label/input html
			creater(&fSet)

		} else {
			// create field html
			switch fTyp {
			case "text":
				fSet.Field = template.HTML(fmt.Sprintf(
					`<input id="%s" name="%s" type="text" value="%v" class="form-control"%s%s>`, fId, fName, value, placeholders, attrs))

			case "textarea":
				fSet.Field = template.HTML(fmt.Sprintf(
					`<textarea id="%s" name="%s" rows="5" class="form-control"%s%s>%v</textarea>`, fId, fName, placeholders, attrs, value))

			case "password":
				fSet.Field = template.HTML(fmt.Sprintf(
					`<input id="%s" name="%s" type="password" value="%v" class="form-control"%s%s>`, fId, fName, value, placeholders, attrs))

			case "select":
				var options string
				str := fmt.Sprintf(`<select id="%s" name="%s" class="form-control"%s%s>%s</select>`, fId, fName, placeholders, attrs)
				fun := elm.Addr().MethodByName(name + "SelectData")

				if fun.IsValid() {
					results := fun.Call([]reflect.Value{})
					if len(results) > 0 {
						v := results[0]
						if v.CanInterface() {
							if vu, ok := v.Interface().([][]string); ok {
								strv := ToStr(value)
								seted := false
								for _, parts := range vu {
									var n, v string
									switch {
									case len(parts) > 1:
										n, v = locale.Tr(parts[0]), parts[1]
									case len(parts) == 1:
										n, v = locale.Tr(parts[0]), parts[0]
									}
									var selected string
									if !seted && strv == v {
										selected = ` selected="selected"`
										seted = true
									}
									options += fmt.Sprintf(`<option value="%s"%s>%s</option>`, v, selected, n)
								}
							}
						}
					}
				}

				if len(options) == 0 {
					options = fmt.Sprintf(`<option value="%v">%v</option>`, value, value)
				}

				str = fmt.Sprintf(str, options)
				fSet.Field = template.HTML(str)

			case "hidden":
				fSet.Field = template.HTML(fmt.Sprintf(
					`<input id="%s" name="%s" type="hidden" value="%v"%s>`, fId, fName, value, attrs))

			case "date", "datetime":
				t := value.(time.Time)
				tval := beego.Date(t, DateTimeFormat)
				if tval == "0001-01-01 00:00:00" {
					t = time.Now()
				}
				if fTyp == "date" {
					tval = beego.Date(t, DateFormat)
				}
				fSet.Field = template.HTML(fmt.Sprintf(
					`<input id="%s" name="%s" type="%s" value="%s" class="form-control"%s%s>`, fId, fName, fTyp, tval, placeholders, attrs))

			case "checkbox":
				var checked string
				if b, ok := value.(bool); ok && b {
					checked = "checked"
				}
				fSet.Field = template.HTML(fmt.Sprintf(
					`<label for="%s" class="checkbox">%s<input id="%s" name="%s" type="checkbox" %s></label>`,
					fId, fSet.LabelText, fId, fName, checked))
			}

			// create label html
			switch fTyp {
			case "checkbox", "hidden":
			default:
				fSet.Label = template.HTML(fmt.Sprintf(`
					<label class="control-label" for="%s">%s</label>`, fId, fSet.LabelText))
			}

			if filter, ok := customFilters[fTyp]; ok {
				// use custome filter replace label/input html
				filter(&fSet)
			}

		}

		fSets.FieldList = append(fSets.FieldList, &fSet)
		fSets.Fields[name] = &fSet
	}
	return fSets
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
		if v, ok := values[fName]; !ok {
			continue
		} else {
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
					t, err := beego.DateParse(value, DateTimeFormat)
					if err != nil {
						continue
					}
					f.Set(reflect.ValueOf(t))
				} else {
					t, err := beego.DateParse(value, DateFormat)
					if err != nil {
						continue
					}
					f.Set(reflect.ValueOf(t))
				}
			}
		case reflect.String:
			f.SetString(value)
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
