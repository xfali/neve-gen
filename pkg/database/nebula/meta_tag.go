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

type TagMeta struct {
	TagName string
	Fields  []FieldMeta
}

type FieldMeta struct {
	Field   string `column:"Field"`
	Type    string `column:"Type"`
	Null    string `column:"Null"`
	Default string `column:"Default"`
	Comment string `column:"Comment"`
}
