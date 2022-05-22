package initialize

import (
	"fmt"
	"log"
	"reflect"
	"shop-api/user-web/global"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	myValidator "shop-api/user-web/validator"
)

func InitTrans(locale string) {
	//修改gin框架中的validator引擎属性, 实现定制
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		//注册一个获取json的tag的自定义方法
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		zhT := zh.New() //中文翻译器
		enT := en.New() //英文翻译器
		//第一个参数是备用的语言环境，后面的参数是应该支持的语言环境
		uni := ut.New(enT, zhT, enT)
		global.Trans, ok = uni.GetTranslator(locale)
		if !ok {
			log.Fatalln(fmt.Errorf("uni.GetTranslator(%s)", locale))
		}

		switch locale {
		case "en":
			enTranslations.RegisterDefaultTranslations(v, global.Trans)
		case "zh":
			zhTranslations.RegisterDefaultTranslations(v, global.Trans)
		default:
			enTranslations.RegisterDefaultTranslations(v, global.Trans)
		}
		return
	}
	return
}

func InitCustomValidator() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("mobile", myValidator.ValidateMobile)
		// 注册自定义验证器翻译
		_ = v.RegisterTranslation("mobile", global.Trans, func(ut ut.Translator) error {
			return ut.Add("mobile", "非法的手机号码", true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("mobile", fe.Field())
			return t
		})
	}
}
