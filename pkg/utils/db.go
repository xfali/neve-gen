// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package utils

import (
	"fmt"
	"github.com/xfali/gobatis-cmd/pkg"
	"github.com/xfali/gobatis-cmd/pkg/common"
	"github.com/xfali/gobatis-cmd/pkg/db"
	"github.com/xfali/gobatis-cmd/pkg/mapping"
	"github.com/xfali/neve-gen/pkg/model"
	"github.com/xfali/neve-gen/pkg/stringfunc"
	"github.com/xfali/stream"
	"strings"
)

func ReadDbInfo(ds model.DataSource) ([]model.Module, map[string][]common.ModelInfo, error) {
	dbDriver := db.GetDriver(ds.DriverName)
	if dbDriver == nil {
		return nil, nil, fmt.Errorf("DB type %s not support. ", ds.DriverName)
	}
	err := dbDriver.Open(ds.DriverName, ds.DriverInfo)
	if err != nil {
		return nil, nil, fmt.Errorf("DB Open %s info: %s failed. ", ds.DriverName, ds.DriverInfo)
	}
	defer dbDriver.Close()

	tables, err := dbDriver.QueryTableNames(ds.Scan.DBName)
	if err != nil {
		return nil, nil, err
	}
	if len(tables) == 0 {
		return nil, nil, nil
	}
	if len(ds.Scan.Tables) > 0 {
		tables = stream.Slice(tables).Filter(func(s string) bool {
			return stream.Slice(ds.Scan.Tables).AnyMatch(func(s1 string) bool {
				return s == s1
			})
		}).Collect(nil).([]string)
	}

	ret := make([]model.Module, 0, len(tables))
	ti := make(map[string][]common.ModelInfo, len(tables))
	for _, t := range tables {
		mds, err := dbDriver.QueryTableInfo(ds.Scan.DBName, t)
		if err != nil {
			return nil, nil, err
		}
		infos := make([]*model.Info, len(mds))
		for i, md := range mds {
			info := GobatisInfo2ModuleInfo(md)
			infos[i] = &info
		}
		n := Snake2Camel(t)
		m := model.Module{
			Name:  n,
			Pkg:   stringfunc.FirstLower(n),
			Infos: infos,
		}
		ti[n] = mds
		ret = append(ret, m)
	}
	return ret, ti, nil
}

func Snake2Camel(table string) string {
	if strings.Index(table, "_") != -1 {
		return pkg.Snake2camel(strings.ToLower(table))
	}
	return table
}

func GobatisInfo2ModuleInfo(info common.ModelInfo) model.Info {
	return model.Info{
		Name:     Snake2Camel(info.ColumnName),
		DataType: mapping.SqlType2GoMap[info.DataType],
		Nullable: info.Nullable,
		Key:      info.ColumnKey,
		Comment:  info.Comment,
		Tag:      info.Tag,
	}
}
