// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package model

import (
	"github.com/xfali/neve-gen/pkg/stringfunc"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
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
	Name  string  `yaml:"name"`
	Pkg   string  `yaml:"pkg"`
	Infos []*Info `yaml:"infos"`
}

type DataSource struct {
	Name       string `yaml:"name"`
	DriverName string `yaml:"driverName"`
	DriverInfo string `yaml:"driverInfo"`
}

type Web struct {
	Port int `yaml:"port"`
}

type GoConfig struct {
	Version string `yaml:"version"`
}

type App struct {
	Go          GoConfig     `yaml:"go"`
	Name        string       `yaml:"name"`
	Version     string       `yaml:"version"`
	Description string       `yaml:"description"`
	ModName     string       `yaml:"modName"`
	DataSources []DataSource `yaml:"datasources"`
	Web         Web          `yaml:"web"`
	Modules     []*Module    `yaml:"modules"`
	Result      Result       `yaml:"result"`
}

type Result struct {
	Name  string  `yaml:"name"`
	Pkg   string  `yaml:"pkg"`
	Infos []*Info `yaml:"infos"`
}

func (m Result) FindKeyInfo(key string) (pki *Info, have bool) {
	for _, v := range m.Infos {
		if strings.ToLower(v.Key) == key {
			return v, true
		}
	}
	return nil, false
}

func (m Result) Defined() bool {
	if _, ok := m.FindKeyInfo("payload"); ok {
		return true
	}
	return false
}

type Value struct {
	Author *Author `yaml:"author"`
	App    *App    `yaml:"app"`
}

func (m Module) FindPrimaryKeyInfo() (pki *Info, have bool) {
	for _, v := range m.Infos {
		if v.Key == "PRI" {
			return v, true
		}
	}
	return nil, false
}

func LoadValue(path string) (*Value, error) {
	d, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	ret := &Value{}
	err = yaml.Unmarshal(d, ret)
	return ret, err
}

func (v *Value) Normalize() {
	for _, m := range v.App.Modules {
		m.Name = stringfunc.FirstUpper(m.Name)
		for _, info := range m.Infos {
			info.Name = stringfunc.FirstUpper(info.Name)
		}
	}
}
