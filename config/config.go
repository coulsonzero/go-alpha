package config

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	MySQL MysqlConfig `mapstructure:"mysql"`
	Redis RedisConfig `mapstructure:"redis"`
}

type MysqlConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

func (r RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}

var Conf Config

func init() {
	viper.SetConfigType("yaml")
	viper.SetConfigFile("config.yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	if err := viper.Unmarshal(&Conf); err != nil {
		log.Fatalf("Failed to parse config: %v", err)
	}

	// 环境变量覆盖（用于 Docker 容器）
	if v := os.Getenv("MYSQL_HOST"); v != "" {
		Conf.MySQL.Host = v
	}
	if v := os.Getenv("REDIS_HOST"); v != "" {
		Conf.Redis.Host = v
	}
}

func GetYamlDsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		Conf.MySQL.User,
		Conf.MySQL.Password,
		Conf.MySQL.Host,
		Conf.MySQL.Port,
		Conf.MySQL.DBName,
	)
}
