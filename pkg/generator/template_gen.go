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
	"path/filepath"
	"strings"
	"text/template"
)

type TemplGenerator struct {
	tmpl *template.Template
}

func NewTemplateGenerator(tmpl string, funcMaps ...map[string]interface{}) *TemplGenerator {
	fms := map[string]interface{}{
		// string functions
		"firstLower": stringfunc.FirstLower,
		"firstUpper": stringfunc.FirstUpper,

		// tag functions
		"withTag": WithTag,

		// result functions
		"resultPayloadName":   FindPayloadKeyName,
		"setResultTotal":      SetResultTotal,
		"setResultPagination": SetResultPagination,
		"resultDefined":       ResultDefined,
		"resultInfo":          ResultInfo,

		// io functions
		"dir": filepath.Dir,

		// db functions
		//"convertPrimaryKeyType": Convert2PrimaryKeyType,
		"primaryKeyName":           FindPrimaryKeyName,
		"setPrimaryKeyValue":       SetPrimaryKeyValue,
		"setPrimaryKeyValueImport": SetPrimaryKeyValueImport,
		"selectModuleKey":          SelectModuleKey,
		"selectModulePrimaryInfo":  SelectModulePrimaryInfo,
		"hasDB":                    HasDB,
	}
	for _, fm := range funcMaps {
		for k, v := range fm {
			fms[k] = v
		}
	}
	t, err := template.New("app").Funcs(fms).Option("missingkey=error").Parse(tmpl)
	if err != nil {
		panic(fmt.Errorf("Parse template failed: %v. ", err))
	}
	return &TemplGenerator{
		tmpl: t,
	}
}

func NewGeneratorWithTmplFile(tmplPath string, funcMaps ...map[string]interface{}) *TemplGenerator {
	d, err := ioutil.ReadFile(tmplPath)
	if err != nil {
		panic(fmt.Errorf("Cannot open template file: %s. ", tmplPath))
	}
	return NewTemplateGenerator(string(d), funcMaps...)
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

func WithTag(m model.ModelData, modelName, originTag string) string {
	originTag = withTag(modelName, originTag, m.Value.App.Web.Accept)
	originTag = withTag(modelName, originTag, m.Value.App.Web.Produce)
	return originTag
}

func withTag(modelName, originTag string, tag string) string {
	if tag == "" {
		tag = "json"
	}
	originTag = strings.TrimSpace(originTag)
	if originTag == "" {
		return fmt.Sprintf(`%s:"%s"`, tag, stringfunc.FirstLower(modelName))
	}
	tags := strings.Split(originTag, " ")
	have := false
	for _, s := range tags {
		if strings.Index(s, tag+":") != -1 {
			have = true
			break
		}
	}
	if !have {
		originTag += fmt.Sprintf(` %s:"%s"`, tag, stringfunc.FirstLower(modelName))
	}
	return originTag
}

func FindPayloadKeyName(m model.Result) string {
	i, have := m.FindKeyInfo("payload")
	if have {
		return i.Name
	}
	return ""
}

func SetResultTotal(m model.Result, param string) string {
	i, have := m.FindKeyInfo("count")
	if have {
		return fmt.Sprintf("%s = %s(%s)", stringfunc.FirstUpper(i.Name), i.DataType, param)
	}
	return ""
}

func SetResultPagination(m model.Result, param string) string {
	i, have := m.FindKeyInfo("pagination")
	if have {
		return fmt.Sprintf("%s = %s(%s)", stringfunc.FirstUpper(i.Name), i.DataType, param)
	}
	return ""
}

func ResultDefined(m model.ModelData) bool {
	return m.Value.App.Result.Defined()
}

func ResultInfo(m model.Result, key string) *model.Info {
	i, _ := m.FindKeyInfo(key)
	return i
}

func SelectModuleKey(m model.Module) string {
	i, have := m.FindPrimaryKeyInfo()
	if have {
		return i.Name
	}
	if len(m.Infos) > 0 {
		return m.Infos[0].Name
	}
	return ""
}

func SelectModulePrimaryInfo(m model.Module) *model.Info {
	i, have := m.FindPrimaryKeyInfo()
	if have {
		return i
	}
	return nil
}

func HasDB(c model.Value, dbType string) bool {
	for _, ds := range c.App.DataSources {
		if ds.DriverName == dbType {
			return true
		}
	}
	return false
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

func SetPrimaryKeyValueImport(m model.Module) string {
	info, have := m.FindPrimaryKeyInfo()
	if have {
		switch info.DataType {
		case "int", "int16", "int32", "int64", "uint", "uint16", "uint32", "uint64", "float", "float32", "float64":
			return "strconv"
		}
	}
	return ""
}

func SetPrimaryKeyValue(m model.Module, requestName, paramName string) string {
	info, have := m.FindPrimaryKeyInfo()
	if have {
		b := strings.Builder{}
		switch info.DataType {
		case "int", "int16", "int32", "int64":
			b.WriteString(fmt.Sprintf("if v, err := strconv.ParseInt(%s, 10, 64); err == nil { %s.%s = %s(v) }\n",
				paramName, requestName, stringfunc.FirstUpper(info.Name), info.DataType))
		case "uint", "uint16", "uint32", "uint64":
			b.WriteString(fmt.Sprintf("if v, err := strconv.ParseUint(%s, 10, 64); err == nil { %s.%s = %s(v) }\n",
				paramName, requestName, stringfunc.FirstUpper(info.Name), info.DataType))
		case "float", "float32", "float64":
			b.WriteString(fmt.Sprintf("if v, err := strconv.ParseFloat(%s, 10, 64); err == nil { %s.%s = %s(v) }\n",
				paramName, requestName, stringfunc.FirstUpper(info.Name), info.DataType))
		case "string":
			b.WriteString(fmt.Sprintf("%s.%s = %s\n",
				requestName, stringfunc.FirstUpper(info.Name), paramName))
		}
		return b.String()
	} else {
		return ""
	}
}
