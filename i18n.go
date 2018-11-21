package middleware

import (
	"github.com/forease/i18n/v2/i18n"
	"github.com/gin-gonic/gin"
)

var (
	DefaultLang = "zh_CN"
	lang        = DefaultLang
)

func I18N(args ...string) gin.HandlerFunc {
	if len(args) > 0 {
		lang = args[0]
	}

	return func(c *gin.Context) {

	}
}

func LoadLocales(dir string) error {
	return i18n.LoadLocales(dir)
}

func Tr(format string, args ...interface{}) string {
	return i18n.Tr(lang, format, args...)
}

func getLang() string {
	return ""
}
