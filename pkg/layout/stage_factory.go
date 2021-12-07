// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package layout

import (
	"fmt"
	"github.com/xfali/fig"
	"github.com/xfali/neve-gen/pkg/model"
	"github.com/xfali/neve-gen/pkg/stage"
	"path/filepath"
)

func CreateStages(target, tmplRoot string) []stage.Stage {
	var ret []stage.Stage
	ret = append(ret, stage.NewGenProjectStage(target, tmplRoot))
	ret = append(ret, stage.NewGenGoModStage(target, tmplRoot))
	ret = append(ret, stage.NewGenMainStage(target, tmplRoot))
	ret = append(ret, stage.NewGenHandlerStage(target, tmplRoot))
	ret = append(ret, stage.NewGenServiceStage(target, tmplRoot))
	ret = append(ret, stage.NeGenServiceImplStage(target, tmplRoot))
	return ret
}

func ParseStages(target, tmplRoot string) ([]stage.Stage, error) {
	m, err := LoadLayoutSpec(tmplRoot)
	if err != nil {
		return nil, err
	}
	var ret []stage.Stage
	ret = append(ret, stage.NewGenProjectStage(target, tmplRoot))
	for _, spec := range m.Sepc.TemplateSepcs {
		s, err := CreateStagesByTemplateSpec(target, tmplRoot, spec)
		if err != nil {
			return nil, err
		}
		ret = append(ret, s)
	}
	return ret, nil
}

func CreateStagesByTemplateSpec(target, tmplRoot string, spec model.TemplateSepc) (stage.Stage, error) {
	switch spec.Type {
	case model.TemplateTypeApp:
		return stage.NewAppStage(target, tmplRoot, spec), nil
	case model.TemplateTypeModule:
		return stage.NewModuleStage(target, tmplRoot, spec), nil
	default:
		return nil, fmt.Errorf("Type: %s not support\nSpec: %v .", spec.Type, spec)
	}
}

func LoadLayoutSpec(path string) (*model.TemplateLayoutConfig, error) {
	f, err := fig.LoadYamlFile(filepath.Join(path, model.TemplateLayoutSpecFile))
	if err != nil {
		return nil, err
	}
	m := &model.TemplateLayoutConfig{
	}
	err = f.GetValue("", m)
	if err != nil {
		return nil, err
	}
	return m, err
}
