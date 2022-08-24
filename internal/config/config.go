package config

import (
	"time"

	"go-wai-wong/internal/constant"

	"github.com/spf13/viper"
)

// load config, avoid using init()
func LoadConfig() {
	viper.SetDefault(constant.TokenSecret, "NXY4eS9CP0UoSCtLYlBlU2hWbVlxM3Q2dzl6JEMmRik=")
	viper.SetDefault(constant.TokenAudience, "local")
	viper.SetDefault(constant.TokenExpiresIn, constant.ExpiresInMinutes*time.Minute)
}
