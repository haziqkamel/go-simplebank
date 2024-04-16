package util

import "github.com/spf13/viper"

type Config struct {
	DBDriver      string `mapstructure:"DB_DRIVER" default:"postgres"`
	DBSource      string `mapstructure:"DB_SOURCE" default:"postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS" default:"0.0.0.0:8080"`
}

// LoadConfig loads the application configuration from the environment variables
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env") // JSON, XML, YAML, TOML, HCL, envfile

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
