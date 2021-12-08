// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package layout

import (
	"context"
	"fmt"
	"github.com/xfali/gobatis-cmd/pkg/common"
	"github.com/xfali/neve-gen/pkg/model"
	"github.com/xfali/neve-gen/pkg/utils"
	"github.com/xfali/stream"
	"github.com/xfali/xlog"
)

func LoadDatabase(ctx context.Context, m *model.ModelData) (context.Context, error) {
	all := map[string][]common.ModelInfo{}
	for _, ds := range m.Config.DataSources {
		if ds.Scan.Enable {
			ms, infos, err := utils.ReadDbInfo(ds)
			if err != nil {
				return ctx, fmt.Errorf("Load modules from database %s %s failed: %v. ", ds.DriverName, ds.DriverInfo, err)
			}
			for _, dbMod := range ms {
				m.Value.App.Modules = stream.Slice(m.Value.App.Modules).Filter(func(om *model.Module) bool {
					if dbMod.Name == om.Name {
						xlog.Warnf("Module %v defined in value file have a same name with database, remove it.", om)
						return false
					}
					return true
				}).Collect(nil).([]*model.Module)
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

var tableInfoKey = struct{}{}

func WithTableInfo(ctx context.Context, v map[string][]common.ModelInfo) context.Context {
	return context.WithValue(ctx, tableInfoKey, v)
}

func GetTableInfo(ctx context.Context) (map[string][]common.ModelInfo, bool) {
	v := ctx.Value(tableInfoKey)
	if v == nil {
		return nil, false
	}
	return v.(map[string][]common.ModelInfo), true
}
