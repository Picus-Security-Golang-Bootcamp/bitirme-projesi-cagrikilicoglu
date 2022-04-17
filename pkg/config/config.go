package config

import (
	"errors"

	"github.com/spf13/viper"
)

// Config
type Config struct {
	ServerConfig ServerConfig `yaml:"ServerConfig"`
	JWTConfig    JWTConfig    `yaml:"JWTConfig"`
	DBConfig     DBConfig     `yaml:"DBConfig"`
	Logger       Logger       `yaml:"Logger"`
}

// ServerConfig
type ServerConfig struct {
	Port                int    `yaml:"Port"`
	TimeoutSecs         int    `yaml:"TimeoutSecs"`
	ReadTimeoutSecs     int    `yaml:"ReadTimeoutSecs"`
	WriteTimeoutSecs    int    `yaml:"WriteTimeoutSecs"`
	AppVersion          string `yaml:"AppVersion"`
	Mode                string `yaml:"Mode"`
	RoutePrefix         string `yaml:"RoutePrefix"`
	Debug               bool   `yaml:"Debug"`
	ShutdownTimeoutSecs int    `yaml:"ShutdownTimeoutSecs"`
}

// JWTConfig
type JWTConfig struct {
	SessionTime               int    `yaml:"SessionTime"`
	SecretKey                 string `yaml:"SecretKey"`
	RefreshSecretKey          string `yaml:"RefreshSecretKey"`
	AccessTokenDurationMins   int    `yaml:"AccessTokenDurationMins"`
	RefreshTokenDurationHours int    `yaml:"RefreshTokenDurationHours"`
}

// DBConfig
type DBConfig struct {
	MigrationFolder string `yaml:"MigrationFolder"`
	DataSourceName  string `yaml:"DataSourceName"`
	Name            string `yaml:"Name"`
	MaxOpen         int    `yaml:"MaxOpen"`
	MaxIdle         int    `yaml:"MaxIdle"`
	MaxLifetime     int    `yaml:"MaxLifetime"`
}

// Logger
type Logger struct {
	Development bool   `yaml:"Development"`
	Encoding    string `yaml:"Encoding"`
	Level       string `yaml:"Level"`
}

// LoadConfig reads configuration from a file
func LoadConfig(fileName string) (*Config, error) {
	v := viper.New()

	v.SetConfigName(fileName)
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}
	var c Config
	err := v.Unmarshal(&c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}
