// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package model

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Author struct {
	Name  string `yaml:"name"`
	Email string `yaml:"email"`
}

type Info struct {
	Name     string `yaml:"name"`
	DataType string `yaml:"dataType"`
	Nullable string `yaml:"nullable"`
	Key      string `yaml:"key"`
	Comment  string `yaml:"comment"`
	Tag      string `yaml:"tag"`
}

type Module struct {
	Name  string `yaml:"name"`
	Pkg   string `yaml:"pkg"`
	Infos []Info `yaml:"infos"`
}

type App struct {
	Name        string   `yaml:"name"`
	Version     string   `yaml:"version"`
	Description string   `yaml:"description"`
	ModName     string   `yaml:"modName"`
	Modules     []Module `yaml:"modules"`
}

type Value struct {
	Author Author `yaml:"author"`
	App    App    `yaml:"app"`
}

func (m Module) FindPrimaryKeyInfo() (pki Info, have bool) {
	for _, v := range m.Infos {
		if v.Key == "PRI" {
			return v, true
		}
	}
	return Info{}, false
}

func LoadValue(path string) (*Value, error){
	d, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	ret := &Value{}
	err = yaml.Unmarshal(d, ret)
	return ret, err
}
