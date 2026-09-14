package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port string `mapstructure:"port"`
		Env  string `mapstructure:"env"`
	} `mapstructure:"server"`
	Redis struct {
		URL string `mapstructure:"url"`
	} `mapstructure:"redis"`
	WebSocket struct {
		HeartbeatInterval int `mapstructure:"heartbeat_interval"`
		MaxMessageSize    int `mapstructure:"max_message_size"`
		RateLimit         struct {
			MessagesPerSecond int `mapstructure:"messages_per_second"`
			Burst             int `mapstructure:"burst"`
		} `mapstructure:"rate_limit"`
	} `mapstructure:"websocket"`
	Auth struct {
		Enabled   bool   `mapstructure:"enabled"`
		JWTSecret string `mapstructure:"jwt_secret"`
	} `mapstructure:"auth"`
}

var AppConfig *Config

func LoadConfig() error {
	viper.SetConfigName("default")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath("../../configs")
	viper.AddConfigPath(".")

	// Allow overriding via environment variables (e.g., REDIS_URL)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	AppConfig = &Config{}
	if err := viper.Unmarshal(AppConfig); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return nil
}
