package vdate

import (
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/samber/lo"
)

var (
	validate *validator.Validate
	trans    ut.Translator
)

func init() {
	validate = validator.New()

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		jsonTag := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if jsonTag == "" || jsonTag == "-" {
			return fld.Name
		}
		return jsonTag
	})

	enLocale := en.New()
	uni := ut.New(enLocale, enLocale)
	trans, _ = uni.GetTranslator("en")

	lo.Must0(en_translations.RegisterDefaultTranslations(validate, trans))
}

var (
	Struct     = validate.Struct
	StructCtx  = validate.StructCtx
	Translator = trans
)
