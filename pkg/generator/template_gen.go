// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package generator

import (
	"fmt"
	"github.com/xfali/neve-gen/pkg/model"
	"github.com/xfali/neve-gen/pkg/stringfunc"
	"io"
	"io/ioutil"
	"strings"
	"text/template"
)

type TemplGenerator struct {
	tmpl *template.Template
}

func NewTemplateGenerator(tmpl string) *TemplGenerator {
	t, err := template.New("app").Funcs(map[string]interface{}{
		"firstLower":     stringfunc.FirstLower,
		"firstUpper":     stringfunc.FirstUpper,
		"primaryKeyName": FindPrimaryKeyName,
		//"convertPrimaryKeyType": Convert2PrimaryKeyType,
		"setPrimaryKeyValue": SetPrimaryKeyValue,
		"selectModuleKey":    SelectModuleKey,
	}).Option("missingkey=error").Parse(tmpl)
	if err != nil {
		panic(fmt.Errorf("Parse template failed: %v. ", err))
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
	return NewTemplateGenerator(string(d))
}

func (g *TemplGenerator) Generate(model interface{}, w io.Writer) error {
	return g.tmpl.Execute(w, model)
}

func FindPrimaryKeyName(m model.Module) string {
	i, have := m.FindPrimaryKeyInfo()
	if have {
		return i.Name
	}
	return ""
}

func SelectModuleKey(m model.Module) string {
	i, have := m.FindPrimaryKeyInfo()
	if have {
		return i.Name
	}
	if len(m.Info) > 0 {
		return m.Info[0].Name
	}
	return ""
}

func Convert2PrimaryKeyType(m model.Module, paramName string) string {
	_, have := m.FindPrimaryKeyInfo()
	if have {
		b := strings.Builder{}
		b.WriteString(fmt.Sprintf("v, _ := strconv.ParseInt(%s, 10, 64)\n", paramName))
		return b.String()
	} else {
		return paramName
	}
}

func SetPrimaryKeyValue(m model.Module, requestName, paramName string) string {
	info, have := m.FindPrimaryKeyInfo()
	if have {
		b := strings.Builder{}
		b.WriteString(fmt.Sprintf("v, _ := strconv.ParseInt(%s, 10, 64)\n", paramName))
		b.WriteString(fmt.Sprintf("%s.%s = v", requestName, info.Name))
		return b.String()
	} else {
		return ""
	}
}
