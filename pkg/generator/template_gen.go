// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package generator

import (
	"fmt"
	"github.com/xfali/neve-gen/pkg/model"
	"io"
	"io/ioutil"
	"text/template"
)

type TemplGenerator struct {
	tmpl *template.Template
}

func NewTemplateGenerator(tmpl string) *TemplGenerator {
	t, err := template.New("app").Option("missingkey=error").Parse(tmpl)
	if err != nil {
		panic(err)
	}
	return &TemplGenerator{
		tmpl: t,
	}
}

func NewGeneratorWithTmplFile(tmplPath string) *TemplGenerator {
	d, err := ioutil.ReadFile(tmplPath)
	if err != nil {
		panic(fmt.Errorf("Cannot open template file: %s. ", tmplPath))
	}
	t, err := template.New("app").Option("missingkey=error").Parse(string(d))
	if err != nil {
		panic(fmt.Errorf("Parse template failed: %v. ", err))
	}
	return &TemplGenerator{
		tmpl: t,
	}
}

func (g *TemplGenerator) Generate(model *model.TemplateModel, w io.Writer) error {
	return g.tmpl.Execute(w, model)
}
