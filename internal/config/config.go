package config

import (
	"github.com/jinzhu/configor"
	"github.com/sirupsen/logrus"
	"sync"
)

type Config struct {
	DB struct {
		Host           string `default:"localhost" env:"MYSQL_HOST"`
		Port           string `default:"3306" env:"MYSQL_PORT"`
		Username       string `default:"" env:"MYSQL_USER"`
		Password       string `default:"" env:"MYSQL_PASSWORD"`
		Database       string `default:"" env:"MYSQL_DATABASE"`
		MaxConnections string `default:"50" env:"MYSQL_MAX_CONNECTION"`
	}
}

var Conf Config
var once sync.Once

func GetCfg() *Config {
	once.Do(func() {
		err := configor.Load(&Conf, "config.yml")
		if err != nil {
			logrus.WithError(err).Warn(err)
		}
	})
	return &Conf
}

func (c *Config) GetDbHost() string {
	return c.DB.Host
}

func (c *Config) GetDbPort() string {
	return c.DB.Port
}

func (c *Config) GetDbUsername() string {
	return c.DB.Username
}

func (c *Config) GetDbPassword() string {
	return c.DB.Password
}

func (c *Config) GetDbDatabase() string {
	return c.DB.Database
}

func (c *Config) GetMaxConnections() string {
	return c.DB.MaxConnections
}
