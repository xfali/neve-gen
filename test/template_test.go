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

var output = os.Stdout

func TestGenApp(t *testing.T) {
	m := getTestModel(t)
	app := generator.NewGeneratorWithTmplFile("../templates/application.tmpl")
	err := generator.WriteHeader(os.Stdout, "../templates/application.tmpl")
	if err != nil {
		t.Fatal(err)
	}
	err = app.Generate(m, output)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGenGomod(t *testing.T) {
	m := getTestModel(t)
	app := generator.NewGeneratorWithTmplFile("../templates/gomod.tmpl")
	err := app.Generate(m, output)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGenService(t *testing.T) {
	m := getTestModel(t)
	app := generator.NewGeneratorWithTmplFile("../templates/service.tmpl")
	for _, v := range m.Value.App.Modules {
		err := generator.WriteHeader(output, "../templates/service.tmpl")
		if err != nil {
			t.Fatal(err)
		}
		err = app.Generate(v, output)
		if err != nil {
			t.Fatal(err)
		}
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestGenHandler(t *testing.T) {
	m := getTestModel(t)
	app := generator.NewGeneratorWithTmplFile("../templates/handler.tmpl")
	for _, v := range m.Value.App.Modules {
		err := generator.WriteHeader(output, "../templates/handler.tmpl")
		if err != nil {
			t.Fatal(err)
		}
		err = app.Generate(v, output)
		if err != nil {
			t.Fatal(err)
		}
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestGenServiceImpl(t *testing.T) {
	t.Run("dummy", func(t *testing.T) {
		m := getTestModel(t)
		app := generator.NewGeneratorWithTmplFile("../templates/service_dummy_impl.tmpl")
		for _, v := range m.Value.App.Modules {
			err := generator.WriteHeader(output, "../templates/service_dummy_impl.tmpl")
			if err != nil {
				t.Fatal(err)
			}
			err = app.Generate(v, output)
			if err != nil {
				t.Fatal(err)
			}
			if err != nil {
				t.Fatal(err)
			}
		}
	})

	t.Run("map", func(t *testing.T) {
		m := getTestModel(t)
		app := generator.NewGeneratorWithTmplFile("../templates/service_impl.tmpl")
		for _, v := range m.Value.App.Modules {
			err := generator.WriteHeader(output, "../templates/service_impl.tmpl")
			if err != nil {
				t.Fatal(err)
			}
			err = app.Generate(v, output)
			if err != nil {
				t.Fatal(err)
			}
			if err != nil {
				t.Fatal(err)
			}
		}
	})
}

func getTestModel(t *testing.T) *model.ModelData {
	f, err := fig.LoadYamlFile("./test-conf.yaml")
	if err != nil {
		t.Fatal(err)
	}
	m := model.ModelData{
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
			ModName:     "My/App/test",
			Modules: []model.Module{
				{
					Name: "User",
					Pkg:  "user",
					Infos: []model.Info{
						{
							Name:"Id",
							DataType:"int",
							Nullable:"false",
							Key:"PRI",
							Comment:"",
							Tag:"",
						},
					},
				},
				{
					Name: "Order",
					Pkg:  "order",
					Infos: []model.Info{
						{
							Name:"Id",
							DataType:"int",
							Nullable:"false",
							Key:"PRI",
							Comment:"",
							Tag:"",
						},
					},
				},
			},
		},
	}
}
