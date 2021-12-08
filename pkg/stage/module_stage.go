// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package stage

import (
	"context"
	"fmt"
	"github.com/xfali/neve-gen/pkg/generator"
	"github.com/xfali/neve-gen/pkg/model"
	"github.com/xfali/neve-gen/pkg/stringfunc"
	"github.com/xfali/neve-gen/pkg/utils"
	"github.com/xfali/xlog"
	"os"
	"path/filepath"
	"strings"
)

type ModuleStage struct {
	logger   xlog.Logger
	gen      generator.Generator
	files    []string
	target   string
	tmpPath  string
	tmplSpec model.TemplateSepc
}

type ModuleModel struct {
	Config model.Config
	Value  model.Module
}

func NewModuleStage(target, tempPath string, tmplSpec model.TemplateSepc) *ModuleStage {
	t := filepath.Join(tempPath, tmplSpec.Template)
	return &ModuleStage{
		logger:   xlog.GetLogger(),
		gen:      generator.NewGeneratorWithTmplFile(t),
		target:   target,
		tmpPath:  t,
		tmplSpec: tmplSpec,
	}
}

func (s *ModuleStage) Name() string {
	return s.tmplSpec.Name
}

func (s *ModuleStage) ShouldSkip(ctx context.Context, model *model.ModelData) bool {
	return !CheckCondition(ctx, s.tmplSpec.Condition, model)
}

func (s *ModuleStage) Generate(ctx context.Context, m *model.ModelData) error {
	select {
	case <-ctx.Done():
		return context.Canceled
	default:
		for _, v := range m.Value.App.Modules {
			if strings.Index(s.tmplSpec.Condition, "${MODULE}") != -1 {
				cond := strings.Replace(s.tmplSpec.Condition, "${MODULE}", fmt.Sprintf(`"%s"`, v.Name), -1)
				if !CheckCondition(ctx, cond, m) {
					continue
				}
			}
			data := ModuleModel{
				Config: *m.Config,
				Value:  *v,
			}
			err := func() error {
				output := filepath.Join(s.target, s.tmplSpec.Target)
				output = strings.Replace(output, "${MODULE}", stringfunc.FirstLower(v.Name), -1)
				dir := filepath.Dir(output)
				err := utils.Mkdir(dir)
				if err != nil {
					s.logger.Errorln(err)
					return fmt.Errorf("Create Module dir : %s failed. ", dir)
				}
				f, err := os.Create(output)
				if err != nil {
					s.logger.Errorln(err)
					return fmt.Errorf("Create file: %s failed. ", output)
				}
				s.files = append(s.files, output)
				defer f.Close()
				if s.tmplSpec.Code == model.TemplateCodeGo {
					err = generator.WriteHeader(f, s.tmplSpec.Code, s.tmpPath)
					if err != nil {
						return err
					}
				}
				return s.gen.Generate(data, f)
			}()
			if err != nil {
				s.logger.Errorln(err)
				return err
			}
		}
		return nil
	}
}

func (s *ModuleStage) Rollback(ctx context.Context) error {
	var last error
	for _, v := range s.files {
		err := os.Remove(v)
		if err != nil {
			last = err
			s.logger.Errorln(err)
		}
	}
	return last
}
