// Copyright 2014 beego Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package orm

// mysql operators.
var mysqlOperators = map[string]string{
	"exact":     "= ?",
	"iexact":    "LIKE ?",
	"contains":  "LIKE BINARY ?",
	"icontains": "LIKE ?",
	// "regex":       "REGEXP BINARY ?",
	// "iregex":      "REGEXP ?",
	"gt":          "> ?",
	"gte":         ">= ?",
	"lt":          "< ?",
	"lte":         "<= ?",
	"eq":          "= ?",
	"ne":          "!= ?",
	"startswith":  "LIKE BINARY ?",
	"endswith":    "LIKE BINARY ?",
	"istartswith": "LIKE ?",
	"iendswith":   "LIKE ?",
}

// mysql dbBaser implementation.
type dbBaseMysql struct {
	dbBase
}

var _ dbBaser = new(dbBaseMysql)

// get mysql operator.
func (d *dbBaseMysql) OperatorSQL(operator string) string {
	return mysqlOperators[operator]
}

// create new mysql dbBaser.
func newdbBaseMysql() dbBaser {
	b := new(dbBaseMysql)
	b.ins = b
	return b
}
