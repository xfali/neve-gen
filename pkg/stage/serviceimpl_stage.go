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
	GenServiceImplStageName = "Generate service implements"
)

type GenServiceImplStage struct {
	logger   xlog.Logger
	tempPath string
	target   string
	files    []string
}

func NeGenServiceImplStage(target string, tempPath string) *GenServiceImplStage {
	return &GenServiceImplStage{
		logger:   xlog.GetLogger(),
		tempPath: tempPath,
		target:   target,
	}
}

func (s *GenServiceImplStage) Name() string {
	return GenServiceImplStageName
}

func (s *GenServiceImplStage) Generate(ctx context.Context, model *model.ModelData) error {
	name := "service_dummy_impl.tmpl"
	if model.Config.Gobatis.Enable {
		name = "service_gobatis_impl.tmpl"
	}
	gen := generator.NewGeneratorWithTmplFile(filepath.Join(s.tempPath, name))
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
				err = generator.WriteGoHeader(f, name)
				if err != nil {
					return err
				}
				return gen.Generate(v, f)
			}()
			if err != nil {
				s.logger.Errorln(err)
				return err
			}
		}

		return nil
	}
}

func (s *GenServiceImplStage) Rollback(ctx context.Context) error {
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
