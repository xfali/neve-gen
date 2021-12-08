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
	GenHandlerStageName = "Generate handlers"
)

type GenHandlerStage struct {
	logger xlog.Logger
	gen    generator.Generator
	target string
	files  []string
}

func NewGenHandlerStage(target string, tempPath string) *GenHandlerStage {
	return &GenHandlerStage{
		logger: xlog.GetLogger(),
		gen:    generator.NewGeneratorWithTmplFile(filepath.Join(tempPath, "handler.tmpl")),
		target: target,
	}
}

func (s *GenHandlerStage) Name() string {
	return GenHandlerStageName
}

func (s *GenHandlerStage) Generate(ctx context.Context, model *model.ModelData) error {
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
				err = generator.WriteGoHeader(f, "handler.tmpl")
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

func (s *GenHandlerStage) ShouldSkip(ctx context.Context, model *model.ModelData) bool {
	return false
}

func (s *GenHandlerStage) Rollback(ctx context.Context) error {
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
