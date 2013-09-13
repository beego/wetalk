package utils

import (
	"html/template"

	"github.com/astaxie/beego"
	"github.com/beego/i18n"
)

// get HTML i18n string
func i18nHTML(locale, format string, args ...interface{}) template.HTML {
	return template.HTML(i18n.Tr(locale, format, args...))
}

func init() {
	// Register template functions.
	beego.AddFuncMap("i18n", i18nHTML)
}
