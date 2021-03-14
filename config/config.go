package config

import (
	"github.com/jackc/pgx"
	"github.com/spf13/viper"
	"log"
)

var (
	Cfg *Config
	Db   *pgx.Conn
)

type Config struct {
	DB     ConfDB     `mapstructure:"database"`
	Server ConfServer `mapstructure:"server"`
}

type ConfDB struct {
	Driver   string `mapstructure:"db_driver"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DbName   string `mapstructure:"db_name"`
	SslMode  string `mapstructure:"ssl_mode"`
	Host     string `mapstructure:"host"`
	MaxConn  string `mapstructure:"max_conn"`
}

type ConfServer struct {
	Address  string `mapstructure:"address"`
	Port     string `mapstructure:"port"`
	Host     string `mapstructure:"host"`
	Protocol string `mapstructure:"protocol"`
}

func Init() *Config {
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatal("cannot open config:", err)
	}
	return &cfg
}

func LoadConfig() (config Config, err error) {
	viper.AddConfigPath("./config")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}