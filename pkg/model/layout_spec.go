// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package model

const (
	TemplateLayoutSpecFile = "layout.yaml"

	TemplateTypeApp    = "app"
	TemplateTypeModule = "module"

	TemplateCodeGoMode = "go.mod"
	TemplateCodeGo     = "go"
)

type TemplateLayoutConfig struct {
	Name    string             `yaml:"name"`
	Version string             `yaml:"version"`
	Sepc    TemplateLayoutSepc `yaml:"spec"`
}

type TemplateLayoutSepc struct {
	TemplateSepcs []TemplateSepc `yaml:"templates"`
}

type TemplateSepc struct {
	Name     string `yaml:"name"`
	Template string `yaml:"template"`
	Type     string `yaml:"type"`
	Code     string `yaml:"code"`
	Target   string `yaml:"target"`
}
