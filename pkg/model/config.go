// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package model

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Scan struct {
	Databases []Database `yaml:"databases"`
}

type DB struct {
	DBName string   `yaml:"dbname"`
	Format string   `yaml:"format"`
	Tables []string `yaml:"tables"`
}

type Database struct {
	Enable   bool   `yaml:"enable"`
	Driver   string `yaml:"driver"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DBs      []DB   `yaml:"dbs"`
}

type GobatisConfig struct {
	Enable bool `yaml:"enable"`
}

type SwaggerConfig struct {
	Enable  bool   `yaml:"enable"`
	DocPath string `yaml:"docPath"`
}

type RestClientConfig struct {
	Enable bool `yaml:"enable"`
}

type Config struct {
	Swagger    SwaggerConfig    `yaml:"swagger"`
	Gobatis    GobatisConfig    `yaml:"gobatis"`
	RestClient RestClientConfig `yaml:"restclient"`
	Scan       Scan             `yaml:"scan"`
}

func LoadConfig(path string) (*Config, error) {
	d, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	ret := &Config{}
	err = yaml.Unmarshal(d, ret)
	return ret, err
}
