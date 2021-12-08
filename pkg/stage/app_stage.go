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
	"github.com/xfali/xlog"
	"os"
	"path/filepath"
)

type AppStage struct {
	logger   xlog.Logger
	gen      generator.Generator
	target   string
	tmpPath  string
	tmplSpec model.TemplateSepc
}

func NewAppStage(target, tempPath string, tmplSpec model.TemplateSepc) *AppStage {
	t := filepath.Join(tempPath, tmplSpec.Template)
	return &AppStage{
		logger:   xlog.GetLogger(),
		gen:      generator.NewGeneratorWithTmplFile(t),
		target:   filepath.Join(target, tmplSpec.Target),
		tmpPath:  t,
		tmplSpec: tmplSpec,
	}
}

func (s *AppStage) Name() string {
	return s.tmplSpec.Name
}

func (s *AppStage) ShouldSkip(ctx context.Context, model *model.ModelData) bool {
	return !CheckCondition(ctx, s.tmplSpec.Condition, model)
}

func (s *AppStage) Generate(ctx context.Context, m *model.ModelData) error {
	select {
	case <-ctx.Done():
		return context.Canceled
	default:
		f, err := os.Create(s.target)
		if err != nil {
			s.logger.Errorln(err)
			return fmt.Errorf("Create file: %s failed. ", s.target)
		}
		defer f.Close()
		err = generator.WriteHeader(f, s.tmplSpec.Code, s.tmpPath)
		if err != nil {
			return err
		}
		return s.gen.Generate(m, f)
	}
}

func (s *AppStage) Rollback(ctx context.Context) error {
	return os.Remove(s.target)
}
