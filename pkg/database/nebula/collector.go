/*
 * Copyright (C) 2023, Xiongfa Li.
 * All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package nebuladrv

import (
	"context"
	nebula_go "github.com/vesoft-inc/nebula-go/v3"
	"github.com/xfali/gobatis-cmd/pkg/common"
	"github.com/xfali/lean/connection"
	"github.com/xfali/lean/drivers/nebuladrv"
	"github.com/xfali/lean/mapping"
	"github.com/xfali/xlog"
	"net/url"
	"strconv"
)

type Connector struct {
	conn connection.Connection
}

func NewConnector() *Connector {
	return &Connector{}
}

func (c *Connector) QueryTableInfo(dbName, tableName string) ([]common.ModelInfo, error) {
	sess, err := c.conn.GetSession()
	if err != nil {
		return nil, err
	}
	defer sess.Close()
	ctx := context.Background()
	var modules []common.ModelInfo
	_, err = sess.Query(ctx, "USE "+dbName)
	if err != nil {
		return nil, err
	}
	r, err := sess.Query(ctx, "DESC SPACE "+dbName)
	if err != nil {
		return nil, err
	}
	spc := struct {
		Type string `column:"Vid Type"`
	}{}
	_, err = mapping.ScanRows(&spc, r)
	if err != nil {
		return nil, err
	}
	vidType := "int64"
	if spc.Type != "INT64" {
		vidType = "string"
	}
	t := tableName
	r, err = sess.Query(ctx, "DESC TAG "+t)
	if err != nil {
		return modules, err
	}
	m := TagMeta{}
	_, err = mapping.ScanRows(&m.Fields, r)
	if err != nil {
		return modules, err
	}

	m.TagName = t
	modules = append(modules, common.ModelInfo{
		ColumnName: "vid",
		DataType:   vidType,
		ColumnKey:  "PRI",
		Nullable:   "",
		Comment:    "vertex id",
		Tag:        "",
	})
	for _, f := range m.Fields {
		modules = append(modules, common.ModelInfo{
			ColumnName: f.Field,
			DataType:   f.Type,
			Nullable:   f.Null,
			Comment:    f.Comment,
			Tag:        "",
		})
	}
	return modules, nil
}

func (c *Connector) QueryTableNames(dbName string) ([]string, error) {
	sess, err := c.conn.GetSession()
	if err != nil {
		return nil, err
	}
	defer sess.Close()
	ctx := context.Background()
	_, err = sess.Query(ctx, "USE "+dbName)
	if err != nil {
		return nil, err
	}
	r, err := sess.Query(ctx, "SHOW TAGS")
	if err != nil {
		return nil, err
	}
	var tns []string
	_, err = mapping.ScanRows(&tns, r)
	if err != nil {
		return nil, err
	}
	return tns, nil
}

//
//func (c *Connector) ReadDbInfo(ctx context.Context, ds model.Database) ([]model.Module, map[string]database.TableInfo, error) {
//	sess, err := c.conn.GetSession()
//	if err != nil {
//		return nil, nil, err
//	}
//	defer sess.Close()
//
//	var modules []model.Module
//	tables := map[string]database.TableInfo{}
//	for _, namespace := range ds.DBs {
//		_, err = sess.Query(ctx, "USE "+namespace.DBName)
//		if err != nil {
//			return nil, nil, err
//		}
//		r, err := sess.Query(ctx, "DESC SPACE "+namespace.DBName)
//		if err != nil {
//			return nil, nil, err
//		}
//		spc := struct {
//			Type string `column:"Vid Type"`
//		}{}
//		_, err = mapping.ScanRows(&spc, r)
//		if err != nil {
//			return nil, nil, err
//		}
//		vidType := "int64"
//		if spc.Type != "INT64" {
//			vidType = "string"
//		}
//		r, err = sess.Query(ctx, "SHOW TAGS")
//		if err != nil {
//			return nil, nil, err
//		}
//		tns := struct {
//			Tags []string `column:"Name"`
//		}{}
//		_, err = mapping.ScanRows(&tns, r)
//		if err != nil {
//			return nil, nil, err
//		}
//		for _, t := range tns.Tags {
//			r, err := sess.Query(ctx, "DESC TAG "+t)
//			if err != nil {
//				return modules, tables, err
//			}
//			m := TagMeta{}
//			_, err = mapping.ScanRows(&m.Fields, r)
//			if err != nil {
//				return modules, tables, err
//			}
//
//			m.TagName = t
//			ti := database.TableInfo{
//				DriverName: "nebula",
//				DbName:     namespace.DBName,
//				Format:     "nebula",
//				TableName:  t,
//			}
//			ti.Info = append(ti.Info, common.ModelInfo{
//				ColumnName: "vid",
//				DataType:   vidType,
//				ColumnKey:  "PRI",
//				Nullable:   "",
//				Comment:    "vertex id",
//				Tag:        "",
//			})
//			for _, f := range m.Fields {
//				ti.Info = append(ti.Info, common.ModelInfo{
//					ColumnName: f.Field,
//					DataType:   f.Type,
//					Nullable:   f.Null,
//					Comment:    f.Comment,
//					Tag:        "",
//				})
//			}
//			tables[t] = ti
//		}
//	}
//	return modules, tables, err
//}

func (c *Connector) Open(driver, info string) error {
	u, err := url.Parse(info)
	p, _ := strconv.Atoi(u.Port())
	pool, err := nebuladrv.NebulaConnPoolCreator([]nebula_go.HostAddress{
		{
			Host: u.Hostname(),
			Port: p,
		},
	}, nebula_go.GetDefaultConf())(xlog.GetLogger())
	if err != nil {
		return err
	}
	pw, _ := u.User.Password()
	if pw == "" {
		pw = "test"
	}
	conn := nebuladrv.NewNebulaConnection(
		nebuladrv.ConnOpts.WithUserInfo(u.User.Username(), pw),
		nebuladrv.ConnOpts.WithConnectionPool(pool),
	)
	err = conn.Open()
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

func (c *Connector) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

var SqlType2GoMap = map[string]string{
	"bool":      "bool",
	"int":       "int",
	"int8":      "int8",
	"int16":     "int16",
	"int32":     "int32",
	"int64":     "int64",
	"float":     "float64",
	"double":    "float64",
	"string":    "string",
	"date":      "time.Time",
	"time":      "time.Time",
	"datetime":  "time.Time",
	"timestamp": "time.Time",
}

func NebulaTypeMapping(sqlType string) string {
	return SqlType2GoMap[sqlType]
}
