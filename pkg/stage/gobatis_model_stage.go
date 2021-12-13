// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package stage

import (
	"context"
	"fmt"
	"github.com/xfali/gobatis-cmd/pkg"
	"github.com/xfali/gobatis-cmd/pkg/config"
	"github.com/xfali/gobatis-cmd/pkg/generator"
	"github.com/xfali/neve-gen/pkg/database"
	"github.com/xfali/neve-gen/pkg/model"
	"github.com/xfali/neve-gen/pkg/stringfunc"
	"github.com/xfali/neve-gen/pkg/utils"
	"github.com/xfali/xlog"
	"os"
	"path/filepath"
	"strings"
)

type GenGobatisModelStage struct {
	logger   xlog.Logger
	target   string
	tmplSpec model.TemplateSepc
	files    []string
}

func NewGenGobatisModelStage(target string, tmplSpec model.TemplateSepc) *GenGobatisModelStage {
	return &GenGobatisModelStage{
		logger:   xlog.GetLogger(),
		tmplSpec: tmplSpec,
		target:   target,
	}
}

func (s *GenGobatisModelStage) Name() string {
	return s.tmplSpec.Name
}

func (s *GenGobatisModelStage) ShouldSkip(ctx context.Context, model *model.ModelData) bool {
	return !CheckCondition(ctx, s.tmplSpec.Condition, model)
}

func (s *GenGobatisModelStage) Generate(ctx context.Context, model *model.ModelData) error {
	select {
	case <-ctx.Done():
		return context.Canceled
	default:
		infos, have := database.GetTableInfo(ctx)
		if have {
			for _, m := range model.Value.App.Modules {
				info, ok := infos[m.Name]
				if ok {
					output := filepath.Join(s.target, s.tmplSpec.Target)
					output = strings.Replace(output, "${MODULE}", stringfunc.FirstLower(m.Name), -1)
					err := utils.Mkdir(output)
					if err != nil {
						s.logger.Errorln(err)
						return fmt.Errorf("Create Module dir : %s failed. ", output)
					}
					conf := config.Config{
						Driver:      info.DriverName,
						Path:        output + "/",
						PackageName: m.Pkg,
						//ModelFile:   pkg.Camel2snake(m.Name),
						TagName:   "xfield,json,yaml,xml",
						Namespace: fmt.Sprintf("%s.%s", m.Pkg, pkg.Camel2snake(m.Name)),
					}
					s.files = append(s.files, filepath.Join(output, strings.ToLower(m.Name)+".go"))
					generator.GenModel(conf, m.Name, info.Info)
				}
			}
		}
	}
	return nil
}

func (s *GenGobatisModelStage) Rollback(ctx context.Context) error {
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
