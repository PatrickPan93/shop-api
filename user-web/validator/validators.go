package validator

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

// ValidateMobile 自定义mobile字段验证器
func ValidateMobile(fl validator.FieldLevel) bool {
	mobile := fl.Field().String()
	// 使用正则表达式判断是否合法手机号
	ok, _ := regexp.MatchString(`^1([38][0-9]|14[579]|5[^4]|16[6]|7[1-35-8]|9[189])\d{8}$`, mobile)
	return ok
}
