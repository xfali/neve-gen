// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package stage

import (
	"context"
	"fmt"
	"github.com/xfali/gobatis-cmd/pkg"
	"github.com/xfali/gobatis-cmd/pkg/common"
	"github.com/xfali/gobatis-cmd/pkg/config"
	"github.com/xfali/gobatis-cmd/pkg/db"
	"github.com/xfali/gobatis-cmd/pkg/generator"
	"github.com/xfali/gobatis-cmd/pkg/mapping"
	"github.com/xfali/neve-gen/pkg/model"
	"github.com/xfali/neve-gen/pkg/stringfunc"
	"github.com/xfali/neve-gen/pkg/utils"
	"github.com/xfali/stream"
	"github.com/xfali/xlog"
	"os"
	"path/filepath"
	"strings"
)

type GenGobatisStage struct {
	logger   xlog.Logger
	target   string
	tmplSpec model.TemplateSepc
	files    []string
}

func NeGenGobatisStage(target string, tmplSpec model.TemplateSepc) *GenGobatisStage {
	return &GenGobatisStage{
		logger:   xlog.GetLogger(),
		tmplSpec: tmplSpec,
		target:   target,
	}
}

func (s *GenGobatisStage) Name() string {
	return s.tmplSpec.Name
}

func (s *GenGobatisStage) Generate(ctx context.Context, model *model.ModelData) error {
	select {
	case <-ctx.Done():
		return context.Canceled
	default:
		for _, ds := range model.Config.DataSources {
			if ds.Scan.Enable {
				models, infos, err := readDbInfo(ds)
				if err != nil {
					return err
				}
				for _, m := range models {
					output := filepath.Join(s.target, s.tmplSpec.Target)
					output = strings.Replace(output, "${MODULE}", stringfunc.FirstLower(m.Name), -1)
					err := utils.Mkdir(output)
					if err != nil {
						s.logger.Errorln(err)
						return fmt.Errorf("Create Module dir : %s failed. ", output)
					}
					conf := config.Config{
						Driver:      ds.DriverName,
						Path:        output + "/",
						PackageName: m.Pkg,
						//ModelFile:   pkg.Camel2snake(m.Name),
						TagName:   "xfield,json,yaml,xml",
						Namespace: fmt.Sprintf("%s.%s", m.Pkg, pkg.Camel2snake(m.Name)),
					}
					s.files = append(s.files, filepath.Join(output, strings.ToLower(m.Name) + ".go"))
					generator.GenModel(conf, m.Name, infos[m.Name])
					s.files = append(s.files, filepath.Join(output, strings.ToLower(m.Name) + "_proxy.go"))
					generator.GenV2Proxy(conf, m.Name, infos[m.Name])
				}
			}
		}
	}
	return nil
}

func (s *GenGobatisStage) Rollback(ctx context.Context) error {
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
