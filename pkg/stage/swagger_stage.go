// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package stage

import (
	"context"
	"fmt"
	"github.com/swaggo/swag/gen"
	"github.com/xfali/neve-gen/pkg/model"
	"github.com/xfali/neve-gen/pkg/utils"
	"github.com/xfali/xlog"
	"os"
	"path/filepath"
)

type SwaggerStage struct {
	logger     xlog.Logger
	target     string
	tmpPath    string
	tmplSpec   model.TemplateSepc
	mainTarget string
}

func NewSwaggerStage(target string, tmplSpec model.TemplateSepc, allSpecs []model.TemplateSepc) *SwaggerStage {
	main := findMain(allSpecs)
	if main == "" {
		panic("layout.yaml template must include main.go with name [main]")
	}
	return &SwaggerStage{
		logger:     xlog.GetLogger(),
		target:     target,
		tmplSpec:   tmplSpec,
		mainTarget: main,
	}
}

func findMain(allSpecs []model.TemplateSepc) string {
	for _, spec := range allSpecs {
		if spec.Name == "main" {
			return spec.Target
		}
	}
	return ""
}

func (s *SwaggerStage) Name() string {
	return s.tmplSpec.Name
}

func (s *SwaggerStage) ShouldSkip(ctx context.Context, model *model.ModelData) bool {
	return !CheckCondition(ctx, s.tmplSpec.Condition, model)
}

func (s *SwaggerStage) Generate(ctx context.Context, m *model.ModelData) error {
	select {
	case <-ctx.Done():
		return context.Canceled
	default:
		dir := filepath.Join(s.target, m.Config.Swagger.DocPath)
		err := utils.Mkdir(dir)
		if err != nil {
			s.logger.Errorln(err)
			return fmt.Errorf("Create app dir : %s failed. ", dir)
		}
		return gen.New().Build(&gen.Config{
			SearchDir:   s.target,
			OutputDir:   dir,
			MainAPIFile: s.mainTarget,
		})
	}
}

func (s *SwaggerStage) Rollback(ctx context.Context) error {
	return os.RemoveAll(filepath.Join(s.target, s.tmplSpec.Target))
}
