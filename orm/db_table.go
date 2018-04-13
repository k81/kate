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

import (
	"fmt"
	"strings"
	"time"
)

// table info struct.
type dbTable struct {
	mi   *modelInfo
	base dbBaser
}

// parse orm model struct field tag expression.
func (t *dbTable) parseExprs(mi *modelInfo, exprs []string) (fi *fieldInfo, operator string, success bool) {
	if len(exprs) == 0 {
		return
	}

	if fi, success = mi.fields.GetByAny(exprs[0]); success {
		if len(exprs) > 1 {
			operator = exprs[1]
		} else {
			operator = "exact"
		}
	}

	return
}

// generate condition sql.
func (t *dbTable) getCondSQL(cond *Condition, sub bool, tz *time.Location) (where string, params []interface{}) {
	if cond == nil || cond.IsEmpty() {
		return
	}

	mi := t.mi

	for i, p := range cond.params {
		if i > 0 {
			if p.isOr {
				where += "OR "
			} else {
				where += "AND "
			}
		}
		if p.isNot {
			where += "NOT "
		}
		if p.isCond {
			w, ps := t.getCondSQL(p.cond, true, tz)
			if w != "" {
				w = fmt.Sprintf("( %s) ", w)
			}
			where += w
			params = append(params, ps...)
		} else {
			fi, operator, suc := t.parseExprs(mi, p.exprs)
			if suc == false {
				panic(fmt.Errorf("unknown field/column name `%s`", strings.Join(p.exprs, ExprSep)))
			}

			operSQL, args := t.base.GenerateOperatorSQL(mi, fi, operator, p.args, tz)

			leftCol := t.base.quote(fi.column)
			t.base.GenerateOperatorLeftCol(fi, operator, &leftCol)

			where += fmt.Sprintf("%s %s ", leftCol, operSQL)
			params = append(params, args...)

		}
	}

	if sub == false && where != "" {
		where = "WHERE " + where
	}

	return
}

// generate order sql.
func (t *dbTable) getOrderSQL(orders []string) (orderSQL string) {
	if len(orders) == 0 {
		return
	}

	orderSqls := make([]string, 0, len(orders))
	for _, order := range orders {
		asc := "ASC"
		if order[0] == '-' {
			asc = "DESC"
			order = order[1:]
		}
		exprs := strings.Split(order, ExprSep)

		fi, _, suc := t.parseExprs(t.mi, exprs)
		if suc == false {
			panic(fmt.Errorf("unknown field/column name `%s`", strings.Join(exprs, ExprSep)))
		}

		orderSqls = append(orderSqls, fmt.Sprintf("%s %s", t.base.quote(fi.column), asc))
	}

	orderSQL = fmt.Sprintf("ORDER BY %s ", strings.Join(orderSqls, ", "))
	return
}

// generate group sql.
func (t *dbTable) getGroupSQL(groups []string) (groupSQL string) {
	if len(groups) == 0 {
		return
	}

	groupSqls := make([]string, 0, len(groups))
	for _, group := range groups {
		exprs := strings.Split(group, ExprSep)

		fi, _, suc := t.parseExprs(t.mi, exprs)
		if suc == false {
			panic(fmt.Errorf("unknown field/column name `%s`", strings.Join(exprs, ExprSep)))
		}

		groupSqls = append(groupSqls, fmt.Sprintf("%s", t.base.quote(fi.column)))
	}

	groupSQL = fmt.Sprintf("GROUP BY %s ", strings.Join(groupSqls, ", "))
	return
}

// generate limit sql.
func (t *dbTable) getLimitSQL(mi *modelInfo, offset int64, limit int64) (limits string) {
	if limit == 0 {
		limit = int64(DefaultRowsLimit)
	}
	if limit < 0 {
		// no limit
		if offset > 0 {
			maxLimit := t.base.MaxLimit()
			if maxLimit == 0 {
				limits = fmt.Sprintf("OFFSET %d", offset)
			} else {
				limits = fmt.Sprintf("LIMIT %d, %d", offset, maxLimit)
			}
		}
	} else if offset <= 0 {
		limits = fmt.Sprintf("LIMIT %d", limit)
	} else {
		limits = fmt.Sprintf("LIMIT %d, %d", offset, limit)
	}
	return
}

// crete new tables collection.
func newDbTable(mi *modelInfo, base dbBaser) *dbTable {
	table := &dbTable{}
	table.mi = mi
	table.base = base
	return table
}
