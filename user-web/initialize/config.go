package initialize

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
	"shop-api/user-web/config"
	"shop-api/user-web/global"
)

func GetEnvInfo(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
}
func InitConfig() {

	var configFileName string
	debug := GetEnvInfo("SHOP_DEBUG")
	configFilePrefix := "config"

	configFileName = fmt.Sprintf("user-web/%s-pro.yaml", configFilePrefix)
	if debug {
		configFileName = fmt.Sprintf("user-web/%s-debug.yaml", configFilePrefix)
	}
	v := viper.New()
	v.SetConfigFile(configFileName)

	if err := v.ReadInConfig(); err != nil {
		log.Panicln(err)
	}
	global.ServerConfig = &config.ServerConfig{}
	if err := v.Unmarshal(global.ServerConfig); err != nil {
		log.Panicln(err)
	}

	zap.S().Infof("配置信息: %v", global.ServerConfig)

	// viper watch配置变化
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed: ", e.Name)
		_ = v.ReadInConfig()
		_ = v.Unmarshal(global.ServerConfig)
		zap.S().Infof("配置信息: %v", global.ServerConfig)
	})
}
