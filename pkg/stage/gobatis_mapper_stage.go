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

type GenGobatisMapperStage struct {
	logger   xlog.Logger
	target   string
	tmplSpec model.TemplateSepc
	files    []string
}

func NeGenGobatisMapperStage(target string, tmplSpec model.TemplateSepc) *GenGobatisMapperStage {
	return &GenGobatisMapperStage{
		logger:   xlog.GetLogger(),
		tmplSpec: tmplSpec,
		target:   target,
	}
}

func (s *GenGobatisMapperStage) Name() string {
	return s.tmplSpec.Name
}

func (s *GenGobatisMapperStage) ShouldSkip(ctx context.Context, model *model.ModelData) bool {
	return !CheckCondition(ctx, s.tmplSpec.Condition, model)
}

func (s *GenGobatisMapperStage) Generate(ctx context.Context, model *model.ModelData) error {
	select {
	case <-ctx.Done():
		return context.Canceled
	default:
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
						ds := info.DataSource
						conf := config.Config{
							Driver:      ds.DriverName,
							Path:        output + "/",
							PackageName: m.Pkg,
							ModelFile:   pkg.Camel2snake(m.Name),
							TagName:     "xfield,json,yaml,xml",
							Keyword:     false,
							Namespace:   fmt.Sprintf("%s.%s", m.Pkg, pkg.Camel2snake(m.Name)),
						}
						conf.MapperFile = ds.Scan.Format
						if ds.Scan.Format == "xml" {
							s.files = append(s.files, filepath.Join(output, "xml", strings.ToLower(m.Name)+"_mapper.xml"))
							generator.GenXml(conf, info.TableName, info.Info)
						} else if ds.Scan.Format == "template" {
							s.files = append(s.files, filepath.Join(output, "template", strings.ToLower(m.Name)+"_mapper.tmpl"))
							generator.GenTemplate(conf, info.TableName, info.Info)
						}
					}
				}
			}
		}
	}
	return nil
}

func (s *GenGobatisMapperStage) Rollback(ctx context.Context) error {
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
