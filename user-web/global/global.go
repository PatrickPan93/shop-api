package global

import (
	ut "github.com/go-playground/universal-translator"
	"shop-api/user-web/config"
)

var (
	ServerConfig = &config.ServerConfig{}
	Trans        ut.Translator
)
