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
	all := map[string]TableInfo{}
	for _, ds := range m.Config.DataSources {
		if ds.Scan.Enable {
			ms, infos, err := ReadDbInfo(ds)
			if err != nil {
				return ctx, fmt.Errorf("Load modules from database %s %s failed: %v. ", ds.DriverName, ds.DriverInfo, err)
			}
			m.Value.App.Modules = stream.Slice(m.Value.App.Modules).Filter(func(om *model.Module) bool {
				return !stream.Slice(ms).AnyMatch(func(nm model.Module) bool {
					if nm.Name == om.Name {
						xlog.Warnf("Duplicate definition Module: [%s] found both in value file and database, keep database one.", om.Name)
						return true
					}
					return false
				})
			}).Collect(nil).([]*model.Module)
			for _, dbMod := range ms {
				m.Value.App.Modules = append(m.Value.App.Modules, &dbMod)
			}
			for k, v := range infos {
				if _, ok := all[k]; ok {
					return ctx, fmt.Errorf("Load modules from database %s %s found same tablename: %s. ", ds.DriverName, ds.DriverInfo, k)
				}
				all[k] = v
			}
		}
	}
	return WithTableInfo(ctx, all), nil
}

