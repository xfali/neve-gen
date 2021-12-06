// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package model

type Author struct {
	Name  string
	Email string
}

type Module struct {
	Name string
	Pkg  string
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
