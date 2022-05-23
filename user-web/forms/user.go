package forms

type PassWordLoginForm struct {
	// required必填字段, mobile为自定义验证器注册名称.
	Mobile string `json:"mobile" form:"mobile" binding:"required,mobile"`
	// 限制长度
	PassWord  string `json:"password" form:"password" binding:"required,min=3,max=20"`
	Captcha   string `json:"captcha" form:"captcha" binding:"required,min=5,max=5"`
	CaptchaId string `json:"captcha_id" form:"captcha_id" binding:"required"`
}
