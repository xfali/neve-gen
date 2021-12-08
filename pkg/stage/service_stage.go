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
	GenServiceStageName = "Generate services"
)

type GenServiceStage struct {
	logger xlog.Logger
	gen    generator.Generator
	target string
	files  []string
}

func NewGenServiceStage(target string, tempPath string) *GenServiceStage {
	return &GenServiceStage{
		logger: xlog.GetLogger(),
		gen:    generator.NewGeneratorWithTmplFile(filepath.Join(tempPath, "service.tmpl")),
		target: target,
	}
}

func (s *GenServiceStage) Name() string {
	return GenServiceStageName
}

func (s *GenServiceStage) Generate(ctx context.Context, model *model.ModelData) error {
	select {
	case <-ctx.Done():
		return context.Canceled
	default:
		for _, v := range model.Value.App.Modules {
			err := func() error {
				output := filepath.Join(s.target, )
				f, err := os.Create(output)
				if err != nil {
					s.logger.Errorln(err)
					return fmt.Errorf("Create file: %s failed. ", s.target)
				}
				s.files = append(s.files, output)
				defer f.Close()
				err = generator.WriteGoHeader(f, "service.tmpl")
				if err != nil {
					return err
				}
				return s.gen.Generate(v, f)
			}()
			if err != nil {
				s.logger.Errorln(err)
				return err
			}
		}

		return nil
	}
}

func (s *GenServiceStage) Rollback(ctx context.Context) error {
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
