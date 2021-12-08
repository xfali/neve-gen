// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package stage

import (
	"context"
	"fmt"
	"github.com/xfali/neve-gen/pkg/model"
	"github.com/xfali/neve-gen/pkg/utils"
	"github.com/xfali/xlog"
	"os"
	"path/filepath"
)

const (
	ProjectStageName = "Project layout"
)

type GenProjectStage struct {
	logger   xlog.Logger
	target   string
	tempPath string
}

func NewGenProjectStage(target string, tempPath string) *GenProjectStage {
	return &GenProjectStage{
		logger:   xlog.GetLogger(),
		target:   target,
		tempPath: filepath.Join(tempPath, "project"),
	}
}

func (s *GenProjectStage) Name() string {
	return ProjectStageName
}

func (s *GenProjectStage) Generate(ctx context.Context, model *model.ModelData) error {
	select {
	case <-ctx.Done():
		return context.Canceled
	default:
		err := utils.Mkdir(s.target)
		if err != nil {
			s.logger.Errorln(err)
			return err
		}
		info, err := os.Stat(s.tempPath)
		if err != nil {
			err = fmt.Errorf("Project template path %s is not exists ", s.tempPath)
			s.logger.Errorln(err)
			return err
		} else if !info.IsDir() {
			err = fmt.Errorf("Project template path %s is not a directory ", s.tempPath)
			s.logger.Errorln(err)
			return err
		}
		return utils.CopyDir(s.target, s.tempPath)
	}
}

func (s *GenProjectStage) ShouldSkip(ctx context.Context, model *model.ModelData) bool {
	return false
}

func (s *GenProjectStage) Rollback(ctx context.Context) error {
	return os.RemoveAll(s.target)
}
