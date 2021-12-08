// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package database

import "github.com/xfali/gobatis-cmd/pkg/common"

type DummyDriver struct {
}

func (d *DummyDriver) Open(driver, info string) error {
	return nil
}
func (d *DummyDriver) Close() error {
	return nil
}

func (d *DummyDriver) QueryTableInfo(dbName, tableName string) ([]common.ModelInfo, error) {
	return tableInfo[tableName], nil
}

func (d *DummyDriver) QueryTableNames(dbName string) ([]string, error) {
	var ret []string
	for k := range tableInfo {
		ret = append(ret, k)
	}
	return ret, nil
}

var tableInfo = map[string][]common.ModelInfo{
	"user": {
		{
			ColumnName: "id",
			DataType:   "bigint",
			Nullable:   "",
			ColumnKey:  "PRI",
			Comment:    "user id",
			Tag:        "",
		},
		{
			ColumnName: "name",
			DataType:   "varchar",
			Nullable:   "",
			ColumnKey:  "",
			Comment:    "user name",
			Tag:        "",
		},
	},
	"order": {
		{
			ColumnName: "id",
			DataType:   "bigint",
			Nullable:   "",
			ColumnKey:  "PRI",
			Comment:    "order id",
			Tag:        "",
		},
		{
			ColumnName: "productId",
			DataType:   "bigint",
			Nullable:   "",
			ColumnKey:  "",
			Comment:    "product id",
			Tag:        "",
		},
	},
	"product": {
		{
			ColumnName: "id",
			DataType:   "bigint",
			Nullable:   "",
			ColumnKey:  "PRI",
			Comment:    "product id",
			Tag:        "",
		},
		{
			ColumnName: "name",
			DataType:   "varchar",
			Nullable:   "",
			ColumnKey:  "",
			Comment:    "product name",
			Tag:        "",
		},
	},
}
