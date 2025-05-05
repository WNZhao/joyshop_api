package config

type UserSrvConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type ServeConfig struct {
	Name        string        `mapstructure:"name"`
	UserSrvInfo UserSrvConfig `mapstructure:"user_srv"`
}
