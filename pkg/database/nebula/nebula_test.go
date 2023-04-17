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

import "testing"

func TestNebula(t *testing.T) {
	conn := NewConnector()
	// [scheme:][//[userinfo@]host][/]path[?query][#fragment]
	err := conn.Open("nebula", "nebula://root:test@192.168.1.16:9669")
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	ts, err := conn.QueryTableNames("entities")
	if err != nil {
		t.Fatal(err)
	}
	for _, tag := range ts {
		t.Log(tag)
	}

	tis, err := conn.QueryTableInfo("entities", "Entity")
	if err != nil {
		t.Fatal(err)
	}
	for _, ti := range tis {
		t.Log(ti)
	}
}
