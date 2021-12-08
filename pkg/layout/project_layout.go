// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package layout

import (
	"context"
	"github.com/xfali/neve-gen/pkg/database"
	"github.com/xfali/neve-gen/pkg/model"
	"github.com/xfali/neve-gen/pkg/stage"
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

func (g *ProjectGenerator) Layout(model *model.ModelData) error {
	ctx, err := database.LoadDatabase(context.Background(), model)
	if err != nil {
		return err
	}
	doneStage := make([]stage.Stage, 0, len(g.stages))
	for _, s := range g.stages {
		g.logger.Infof("Generating %s...\n", s.Name())
		err := s.Generate(ctx, model)
		if err != nil {
			g.logger.Infof("Generate %s failed. %v \n", s.Name(), err)
			for i := len(doneStage) - 1; i >= 0; i-- {
				rerr := doneStage[i].Rollback(ctx)
				if rerr != nil {
					g.logger.Errorln(rerr)
				}
			}
			return err
		} else {
			doneStage = append(doneStage, s)
			g.logger.Infof("Generate %s success. \n", s.Name())
		}
	}
	g.logger.Infof("All code generated. \n")
	return nil
}
