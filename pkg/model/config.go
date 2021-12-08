// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package model

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type GoConfig struct {
	Version string `yaml:"version"`
}

type DataSource struct {
	Name       string `yaml:"name"`
	DriverName string `yaml:"driverName"`
	DriverInfo string `yaml:"driverInfo"`
}

type Web struct {
	Port int `yaml:"port"`
}

type GobatisConfig struct {
	Enable bool `yaml:"enable"`
}

type SwaggerConfig struct {
	Enable bool `yaml:"enable"`
}

type Config struct {
	Go          GoConfig      `yaml:"go"`
	Swagger     SwaggerConfig `yaml:"swagger"`
	Gobatis     GobatisConfig `yaml:"gobatis"`
	DataSources []DataSource  `yaml:"datasources"`
	Web         Web           `yaml:"web"`
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
