// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package model

type GoConfig struct {
	Version string `yaml:"version"`
}

type GobatisConfig struct {
	Enable bool `yaml:"enable"`
}

type SwaggerConfig struct {
	Enable bool `yaml:"enable"`
}

type Config struct {
	Go      GoConfig      `yaml:"go"`
	Swagger SwaggerConfig `yaml:"swagger"`
	Gobatis GobatisConfig `yaml:"gobatis"`
}
