// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package layout

import (
	"context"
	"github.com/xfali/neve-gen/pkg/model"
	"github.com/xfali/neve-gen/pkg/stage"
	"github.com/xfali/neve-gen/pkg/utils"
	"github.com/xfali/xlog"
)

type ProjectGenerator struct {
	logger xlog.Logger
	stages []stage.Stage
}

func NewProjectGenerator(stages []stage.Stage) *ProjectGenerator {
	ret := &ProjectGenerator{
		logger: xlog.GetLogger(),
		stages: stages,
	}
	return ret
}

func (g *ProjectGenerator) Layout(targetDir string, model *model.ModelData) error {
	err := utils.Rmdir(targetDir)
	if err != nil {
		g.logger.Errorln(err)
		return err
	}
	err = utils.Mkdir(targetDir)
	if err != nil {
		g.logger.Errorln(err)
		return err
	}
	ctx := context.Background()
	doneStage := make([]stage.Stage, 0, len(g.stages))
	for _, s := range g.stages {
		g.logger.Infoln("Generate stage: ", s.Name())
		err = s.Generate(ctx, model)
		if err != nil {
			g.logger.Errorln(err)
			for i := len(doneStage) - 1; i >= 0; i-- {
				rerr := doneStage[i].Rollback(ctx)
				if rerr != nil {
					g.logger.Errorln(rerr)
				}
			}
			return err
		} else {
			doneStage = append(doneStage, s)
		}
	}
	return nil
}
