// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package stage

import (
	"context"
	"fmt"
	"github.com/xfali/neve-gen/pkg/database"
	"github.com/xfali/neve-gen/pkg/generator"
	"github.com/xfali/neve-gen/pkg/model"
	"github.com/xfali/neve-gen/pkg/stringfunc"
	"github.com/xfali/neve-gen/pkg/utils"
	"github.com/xfali/xlog"
	"os"
	"path/filepath"
	"strings"
)

type ModuleStage struct {
	logger   xlog.Logger
	files    []string
	target   string
	factory  generator.Factory
	tmplSpec model.TemplateSepc
}

type ModuleModel struct {
	Module    *model.Module
	TableInfo *database.TableInfo
}

func (m ModuleModel) getModule() *model.Module {
	return m.Module
}

func (m ModuleModel) getTableInfo() *database.TableInfo {
	return m.TableInfo
}

func NewModuleStage(target string, factory generator.Factory, tmplSpec model.TemplateSepc) *ModuleStage {
	//t := filepath.Join(tempPath, tmplSpec.Template)
	return &ModuleStage{
		logger:   xlog.GetLogger(),
		target:   target,
		factory:  factory,
		tmplSpec: tmplSpec,
	}
}

func (s *ModuleStage) Name() string {
	return s.tmplSpec.Name
}

func (s *ModuleStage) ShouldSkip(ctx context.Context, model *model.ModelData) bool {
	return !CheckCondition(ctx, s.tmplSpec.Condition, model)
}

func (s *ModuleStage) Generate(ctx context.Context, m *model.ModelData) error {
	select {
	case <-ctx.Done():
		return context.Canceled
	default:
		for _, v := range m.Value.App.Modules {
			if strings.Index(s.tmplSpec.Condition, "${MODULE}") != -1 {
				cond := strings.Replace(s.tmplSpec.Condition, "${MODULE}", fmt.Sprintf(`"%s"`, v.Name), -1)
				if !CheckCondition(ctx, cond, m) {
					continue
				}
			}
			data := ModuleModel{
				Module: v,
			}
			infos, load := database.GetTableInfo(ctx)
			if load {
				ds := infos[v.Name]
				data.TableInfo = &ds
			}

			err := func() error {
				output := filepath.Join(s.target, s.tmplSpec.Target)
				output = strings.Replace(output, "${MODULE}", stringfunc.FirstLower(v.Name), -1)
				dir := filepath.Dir(output)
				err := utils.Mkdir(dir)
				if err != nil {
					s.logger.Errorln(err)
					return fmt.Errorf("Create Module dir : %s failed. ", dir)
				}
				f, err := os.Create(output)
				if err != nil {
					s.logger.Errorln(err)
					return fmt.Errorf("Create file: %s failed. ", output)
				}
				s.files = append(s.files, output)
				defer f.Close()
				if s.tmplSpec.Code == model.TemplateCodeGo {
					err = generator.WriteHeader(f, s.tmplSpec.Code, s.tmplSpec.Template)
					if err != nil {
						return err
					}
				}
				gen := s.factory.CreateGenerator(s.tmplSpec.Template, map[string]interface{}{
					"currentModule":          data.getModule,
					"currentModuleTableInfo": data.getTableInfo,
					"currentTemplateSpec": func() model.TemplateSepc {
						return s.tmplSpec
					},
				})
				return gen.Generate(m, f)
			}()
			if err != nil {
				s.logger.Errorln(err)
				return err
			}
		}
		return nil
	}
}

func (s *ModuleStage) Rollback(ctx context.Context) error {
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
