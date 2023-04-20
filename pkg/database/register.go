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

package database

import (
	"github.com/xfali/gobatis-cmd/pkg/common"
	"github.com/xfali/gobatis-cmd/pkg/db"
	nebuladrv "github.com/xfali/neve-gen/pkg/database/nebula"
	"github.com/xfali/neve-gen/pkg/model"
	"sync"
)

type InfoFormatter func(db *model.Database, dbname string) string

var (
	drivers = map[string]*metaReader{
		"mysql": {
			driver:    &db.MysqlDB{},
			mapping:   defaultSqlTypeMapping,
			formatter: defaultFormatter,
		},
		"postgres": {
			driver:    &db.PostgresDB{},
			mapping:   defaultSqlTypeMapping,
			formatter: defaultFormatter,
		},
		"sqlite": {
			driver:    &db.SqliteDB{},
			mapping:   defaultSqlTypeMapping,
			formatter: defaultFormatter,
		},
		"neve_dummydb": {
			driver:    &DummyDriver{},
			mapping:   defaultSqlTypeMapping,
			formatter: defaultFormatter,
		},
		"nebula": {
			driver:    nebuladrv.NewConnector(),
			mapping:   nebuladrv.NebulaTypeMapping,
			formatter: nebuladrv.NebulaInfoFormatter,
		},
	}
	driverLocker sync.Mutex
)

type metaReader struct {
	driver    common.DBDriver
	mapping   TypeMapping
	formatter InfoFormatter
}

func RegisterMetaReader(driver string, dbDriver common.DBDriver, mapping TypeMapping) {
	driverLocker.Lock()
	defer driverLocker.Unlock()
	drivers[driver] = &metaReader{
		driver:  dbDriver,
		mapping: mapping,
	}
}

func getMetaReader(driver string) *metaReader {
	driverLocker.Lock()
	defer driverLocker.Unlock()

	return drivers[driver]
}

func defaultFormatter(ds *model.Database, dbname string) string {
	return db.GenDBInfo(ds.Driver, dbname, ds.Username, ds.Password, ds.Host, ds.Port)
}
