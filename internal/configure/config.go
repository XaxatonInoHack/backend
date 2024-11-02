package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Fiber struct {
	SecretKey string `yaml:"secret_key"`
	Host      string `yaml:"host"`
	Port      int    `yaml:"port"`
}

type Config struct {
	Fiber    Fiber    `yaml:"fiber"`
	Postgres Postgres `yaml:"postgres"`
}

func New() *Config {
	return &Config{
		Fiber:    Fiber{},
		Postgres: Postgres{},
	}
}

func MustConfig(p *string) *Config {
	var path string
	if p == nil {
		path = fetchConfigPath()
	}

	if path == "" {
		path = "./config.yaml"
	}

	if _, ok := os.Stat(path); os.IsNotExist(ok) {
		panic("Config file does not exist: " + path)
	}

	cfg := New()

	if err := cleanenv.ReadConfig(path, cfg); err != nil {
		panic("failed to read config: " + err.Error())
	}

	return cfg
}

func fetchConfigPath() string {
	var res string

	// --config="path/to/config.yaml"
	flag.StringVar(&res, "config", "", "path to config")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}

func (f Fiber) String() string {
	return fmt.Sprintf("%s:%d", f.Host, f.Port)
}
