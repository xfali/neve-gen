// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package model

type Author struct {
	Name  string
	Email string
}

type Info struct {
	Name     string
	DataType string
	Nullable string
	Key      string
	Comment  string
	Tag      string
}

type Module struct {
	Name string
	Pkg  string
	Info []Info
}

type App struct {
	Name        string
	Version     string
	Description string
	ModName     string
	Modules     []Module
}

type Value struct {
	Author Author
	App    App
}

func (m Module) FindPrimaryKeyInfo() (pki Info, have bool) {
	for _, v := range m.Info {
		if v.Key == "PRI" {
			return v, true
		}
	}
	return Info{}, false
}
