package config

type UserSrvConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}
type UserWebConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}
type JWTConfig struct {
	SigningKey string `mapstructure:"key"`
}

type ServerConfig struct {
	UserWebInfo UserWebConfig `mapstructure:"user_web"`
	UserSrvInfo UserSrvConfig `mapstructure:"user_srv"`
	JWTInfo     JWTConfig     `mapstructure:"jwt"`
}
