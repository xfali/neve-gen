// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package database

import (
	"context"
	"fmt"
	"github.com/xfali/neve-gen/pkg/model"
	"github.com/xfali/stream"
	"github.com/xfali/xlog"
)

func LoadDatabase(ctx context.Context, m *model.ModelData) (context.Context, error) {
	allMs := make([]model.Module, 0, 64)
	allInfo := make(map[string]TableInfo, 64)
	for _, ds := range m.Config.Scan.Databases {
		if ds.Enable {
			ms, infos, err := ReadAllDbInfo(ds)
			if err != nil {
				return ctx, fmt.Errorf("Load modules from database %s %v failed: %v. ", ds.Driver, ds, err)
			}
			allMs = append(allMs, ms...)
			for k, v := range infos {
				if info, ok := allInfo[k]; ok {
					return ctx, fmt.Errorf("Load modules from database %s %s found same tablename: %s. ", ds.Driver, info.TableName, k)
				}
				allInfo[k] = v
			}
		}
	}
	m.Value.App.Modules = stream.Slice(m.Value.App.Modules).Filter(func(om *model.Module) bool {
		return !stream.Slice(allMs).AnyMatch(func(nm model.Module) bool {
			if nm.Name == om.Name {
				xlog.Warnf("Duplicate definition Module: [%s] found both in value file and database, keep database one.", om.Name)
				return true
			}
			return false
		})
	}).Collect(nil).([]*model.Module)
	for _, dbMod := range allMs {
		nm := dbMod
		m.Value.App.Modules = append(m.Value.App.Modules, &nm)
	}
	return WithTableInfo(ctx, allInfo), nil
}

