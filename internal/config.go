package internal

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"sync"
)

type Config struct {
	Db DbConfig `yaml:"db"`
	Server ServerConfig `yaml:"server"`
}

type DbConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	SslMode  string `yaml:"ssl_mode"`
}

type ServerConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

func (s *ServerConfig ) Address() string {
	return fmt.Sprintf(":%s", s.Port)

}

func (d *DbConfig) DbUrlConnection() string {
	dbUrl := fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%s sslmode=%s",
		d.User, d.Database, d.Password, d.Host, d.Port, d.SslMode)
	return dbUrl
}

var once sync.Once

func GetConfig(path string) *Config {
	config := &Config{}

	once.Do(func() {
		err := cleanenv.ReadConfig(path, config)
		if err != nil {
			log.Println("parse config failed")
			log.Fatal(err)
		}
	})

	return config
}