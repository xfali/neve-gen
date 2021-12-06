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

const (
	ApplicationStageName = "Generate main"
)

type GenMainStage struct {
	logger xlog.Logger
	gen    generator.Generator
	target string
}

func NewGenMainStage(target string, tempPath string) *GenMainStage {
	return &GenMainStage{
		logger: xlog.GetLogger(),
		gen:    generator.NewGeneratorWithTmplFile(filepath.Join(tempPath, "application.tmpl")),
		target: target,
	}
}

func (s *GenMainStage) Name() string {
	return ApplicationStageName
}

func (s *GenMainStage) Generate(ctx context.Context, model *model.ModelData) error {
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
		return s.gen.Generate(model, f)
	}
}

func (s *GenMainStage) Rollback(ctx context.Context) error {
	return os.Remove(s.target)
}
