package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	LogLevel string
}

type ServerConfig struct {
	Port int
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type JWTConfig struct {
	Secret     string
	Expiration int
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath("/etc/order")

	viper.SetDefault("server.port", 8080)
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.user", "order")
	viper.SetDefault("database.password", "order123")
	viper.SetDefault("database.dbname", "order")
	viper.SetDefault("database.sslmode", "disable")
	viper.SetDefault("jwt.secret", "your-secret-key")
	viper.SetDefault("jwt.expiration", 24)
	viper.SetDefault("loglevel", "info")

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	viper.BindEnv("database.host", "DB_HOST")
	viper.BindEnv("database.port", "DB_PORT")
	viper.BindEnv("database.user", "DB_USER")
	viper.BindEnv("database.password", "DB_PASSWORD")
	viper.BindEnv("database.dbname", "DB_NAME")
	viper.BindEnv("jwt.secret", "JWT_SECRET")
	viper.BindEnv("server.port", "SERVER_PORT")
	viper.BindEnv("loglevel", "LOG_LEVEL")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
