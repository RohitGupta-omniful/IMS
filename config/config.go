package config

import (
	"github.com/spf13/viper"
	"log"
	"time"
)

type Config struct {
	Server   ServerConfig
	Database DBConfig
	AWS      AWSConfig
}

type ServerConfig struct {
	Name         string        `mapstructure:"name"`
	Port         string        `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
}

type DBConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
}

type AWSConfig struct {
	Region        string `mapstructure:"region"`
	PublicBucket  string `mapstructure:"public_bucket"`
	PrivateBucket string `mapstructure:"private_bucket"`
}

var AppConfig *Config

func LoadConfig(path string) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Failed to unmarshal config: %v", err)
	}

	AppConfig = &config
}
