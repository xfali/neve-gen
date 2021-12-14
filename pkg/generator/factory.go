// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package generator

import (
	"github.com/xfali/neve-gen/pkg/buildin"
	"path/filepath"
)

type Factory interface {
	CreateGenerator(tmplPath string, funcMaps ...map[string]interface{}) Generator
}

type FileTemplateFactory struct {
	root string
}

func NewFileTemplateFactory(templateRoot string) *FileTemplateFactory {
	return &FileTemplateFactory{
		root: templateRoot,
	}
}

func (f *FileTemplateFactory) CreateGenerator(tmplPath string, funcMaps ...map[string]interface{}) Generator {
	return NewGeneratorWithTmplFile(filepath.Join(f.root, tmplPath), funcMaps...)
}

type BuildinTemplateFactory struct {
}

func NewBuildinTemplateFactory() *BuildinTemplateFactory {
	return &BuildinTemplateFactory{
	}
}

func (f *BuildinTemplateFactory) CreateGenerator(tmplPath string, funcMaps ...map[string]interface{}) Generator {
	return NewTemplateGenerator(buildin.GetBuildTemplate(tmplPath), funcMaps...)
}
