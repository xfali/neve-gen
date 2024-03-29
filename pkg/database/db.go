// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package database

import (
	"context"
	"fmt"
	"github.com/xfali/gobatis-cmd/pkg/common"
	"github.com/xfali/gobatis-cmd/pkg/mapping"
	"github.com/xfali/gobatis-cmd/pkg/stringutils"
	"github.com/xfali/neve-gen/pkg/model"
	"github.com/xfali/neve-gen/pkg/stringfunc"
	"github.com/xfali/stream"
	"strings"
)

type TableInfo struct {
	DriverName string
	DbName     string
	Format     string
	TableName  string
	Info       []common.ModelInfo
}

func ReadAllDbInfo(ds model.Database) ([]model.Module, map[string]TableInfo, error) {
	ret := make([]model.Module, 0, 64)
	ti := make(map[string]TableInfo, 64)
	for _, db := range ds.DBs {
		m, i, err := ReadDbInfo(ds, db)
		if err != nil {
			return nil, nil, err
		}
		ret = append(ret, m...)
		for k, v := range i {
			if info, ok := ti[k]; ok {
				return nil, nil, fmt.Errorf("Load modules from database %s %s found same tablename: %s. ", ds.Driver, info.TableName, k)
			}
			ti[k] = v
		}
	}
	return ret, ti, nil
}

type TypeMapping func(sqlType string) string

func defaultSqlTypeMapping(sqlType string) string {
	return mapping.SqlType2GoMap[sqlType]
}

func ReadDbInfo(ds model.Database, info model.DB) ([]model.Module, map[string]TableInfo, error) {
	m := getMetaReader(ds.Driver)
	if m == nil {
		return nil, nil, fmt.Errorf("DB type %s not support ", ds.Driver)
	}
	dbDriver := m.driver
	typeMapping := m.mapping
	di := m.formatter(&ds, info.DBName)
	err := dbDriver.Open(ds.Driver, di)
	if err != nil {
		return nil, nil, fmt.Errorf("DB Open %s info: %s failed: %v ", ds.Driver, di, err)
	}

	defer dbDriver.Close()

	tables, err := dbDriver.QueryTableNames(info.DBName)
	if err != nil {
		return nil, nil, err
	}
	if len(tables) == 0 {
		return nil, nil, nil
	}
	if len(info.Tables) > 0 {
		tables = stream.Slice(tables).Filter(func(s string) bool {
			return stream.Slice(info.Tables).AnyMatch(func(s1 string) bool {
				return s == s1
			})
		}).Collect(nil).([]string)
	}

	ret := make([]model.Module, 0, len(tables))
	ti := make(map[string]TableInfo, len(tables))
	for _, t := range tables {
		mds, err := dbDriver.QueryTableInfo(info.DBName, t)
		if err != nil {
			return nil, nil, err
		}
		infos := make([]*model.Info, len(mds))
		pri := false
		for i, md := range mds {
			info := GobatisInfo2ModuleInfo(md)
			info.DataType = typeMapping(md.DataType)
			infos[i] = &info
			if info.Key == "PRI" {
				pri = true
			}
		}
		// No Primary key found, set first one.
		if !pri && len(infos) > 0 {
			infos[0].Key = "PRI"
		}
		n := Snake2Camel(t)
		m := model.Module{
			Name:  n,
			Pkg:   stringfunc.FirstLower(n),
			Infos: infos,
		}
		ti[n] = TableInfo{
			DriverName: ds.Driver,
			DbName:     info.DBName,
			Format:     info.Format,
			Info:       mds,
			TableName:  t,
		}
		ret = append(ret, m)
	}
	return ret, ti, nil
}

func Snake2Camel(s string) string {
	return stringutils.Snake2camel(s)
}

func GobatisInfo2ModuleInfo(info common.ModelInfo) model.Info {
	name := Snake2Camel(info.ColumnName)
	tag := info.Tag
	if strings.Index(tag, "xfield:") == -1 {
		tag = fmt.Sprintf(`%s:"%s"`, "xfield", info.ColumnName)
	}
	return model.Info{
		Name:     name,
		DataType: mapping.SqlType2GoMap[info.DataType],
		Nullable: info.Nullable,
		Key:      info.ColumnKey,
		Comment:  info.Comment,
		Tag:      tag,
	}
}

var tableInfoKey = struct{}{}

func WithTableInfo(ctx context.Context, v map[string]TableInfo) context.Context {
	return context.WithValue(ctx, tableInfoKey, v)
}

func IsTable(ctx context.Context, name string) bool {
	v, ok := GetTableInfo(ctx)
	if !ok {
		return false
	}
	_, ok = v[name]
	return ok
}

func GetTableInfo(ctx context.Context) (map[string]TableInfo, bool) {
	v := ctx.Value(tableInfoKey)
	if v == nil {
		return nil, false
	}
	return v.(map[string]TableInfo), true
}
