// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package test

import (
	"github.com/xfali/fig"
	"github.com/xfali/neve-gen/pkg/generator"
	"github.com/xfali/neve-gen/pkg/model"
	"os"
	"testing"
)

func TestGenApp(t *testing.T) {
	m := getTestModel(t)
	app := generator.NewGeneratorWithTmplFile("../templates/application.tmpl")
	err := generator.WriteHeader(os.Stdout, "../templates/application.tmpl")
	if err != nil {
		t.Fatal(err)
	}
	err = app.Generate(m, os.Stdout)
	if err != nil {
		t.Fatal(err)
	}
}

func getTestModel(t *testing.T) *model.TemplateModel {
	f, err := fig.LoadYamlFile("./test-conf.yaml")
	if err != nil {
		t.Fatal(err)
	}
	m := model.TemplateModel{
	}
	err = f.GetValue("", &m.Config)
	if err != nil {
		t.Fatal(err)
	}
	m.Value = testModules()
	return &m
}

func testModules() model.Value {
	return model.Value{
		Author: model.Author{
			Name:  "testUser",
			Email: "testUser@test.org",
		},
		App: model.App{
			Name:        "testApp",
			Version:     "v0.0.1",
			Description: "auth generator test",
			ModName:     "gihub.com/xfali/neve-generator/test",
			Modules: []model.Module{
				{
					Name: "User",
					Pkg:  "user",
				},
				{
					Name: "Order",
					Pkg:  "order",
				},
			},
		},
	}
}
