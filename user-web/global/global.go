package global

import (
	ut "github.com/go-playground/universal-translator"
	"joyshop_api/user-web/config"
)

var (
	Trans        ut.Translator
	ServerConfig *config.ServeConfig
)
