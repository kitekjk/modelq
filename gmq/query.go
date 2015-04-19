package gmq

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"
)

type _Columns []Column

func (c _Columns) fieldsAndParams(alias string, driverName string) ([]string, []interface{}) {
	fields := make([]string, len(c))
	params := make([]interface{}, 0, len(c))
	for i, col := range c {
		fields[i] = nameWithAlias(col.Name, alias, driverName)
		if col.Value != nil {
			params = append(params, col.Value)
		}
	}
	return fields, params
}

type _Query struct {
	model   TableModel
	columns _Columns
	where   Filter
	orderBy []string
	groupBy []string
	limit   []int64
}

func (q _Query) Exec(dbtx DbTx) (sql.Result, error)                  { return nil, ErrNotSupportedCall }
func (q _Query) SelectOne(dbtx DbTx, functor QueryRowVisitor) error  { return ErrNotSupportedCall }
func (q _Query) SelectList(dbtx DbTx, functor QueryRowVisitor) error { return ErrNotSupportedCall }

func (q _Query) sqlRemains(alias string, driverName string) (string, []interface{}) {
	statements := make([]string, 0)
	params := make([]interface{}, 0)
	if q.where != nil {
		statements = append(statements, fmt.Sprintf("WHERE %s", q.where.SqlString(alias, driverName)))
		params = append(params, q.where.Params()...)
	}
	if q.groupBy != nil && len(q.groupBy) > 0 {
		fields := make([]string, len(q.groupBy))
		for i, gb := range q.groupBy {
			fields[i] = nameWithAlias(gb, alias, driverName)
		}
		groupBy := fmt.Sprintf("GROUP BY %s", strings.Join(fields, ", "))
		statements = append(statements, groupBy)
	}
	if q.orderBy != nil && len(q.orderBy) > 0 {
		fields := make([]string, len(q.orderBy))
		for i, ob := range q.orderBy {
			sortDir := "ASC"
			switch ob[0] {
			case '-':
				sortDir = "DESC"
				ob = ob[1:]
			case '+':
				ob = ob[1:]
			}
			fields[i] = nameWithAlias(ob, alias, driverName) + " " + sortDir
		}
		orderBy := fmt.Sprintf("ORDER BY %s", strings.Join(fields, ", "))
		statements = append(statements, orderBy)
	}
	// FIXME: different limit for different driver
	if q.limit != nil && len(q.limit) == 2 {
		if driverName == "postgres" {
			statements = append(statements, "LIMIT ? OFFSET ?")
			offset := q.limit[0]
			if offset >= 1 {
				offset--
			}
			params = append(params, q.limit[1], offset)
		} else {
			statements = append(statements, "LIMIT ?, ?")
			params = append(params, q.limit[0], q.limit[1])
		}
	}
	return strings.Join(statements, " "), params
}

func (q _Query) queryOne(dbtx DbTx, query string, params []interface{}, functor QueryRowVisitor) error {
	rowCount := 0
	var err error
	err = q.query(dbtx, query, params, func(cols []Column, r []interface{}) error {
		rowCount++
		if rowCount >= 2 {
			err = ErrMultipleRowReturned
			return err
		} else {
			return functor(cols, r)
		}
	})
	if rowCount == 0 && err == nil {
		err = sql.ErrNoRows
	}
	return err
}

func (q _Query) query(dbtx DbTx, query string, params []interface{}, functor QueryRowVisitor) error {
	start := time.Now().UnixNano()
	defer func() {
		if Debug {
			log.Printf("Query SQL - [%s], params=%v, duration=%dms", query, params, (time.Now().UnixNano()-start)/1e6)
		}
	}()

	if stmt, err := dbtx.Prepare(query); err != nil {
		return err
	} else {
		defer stmt.Close()
		if rows, err := stmt.Query(params...); err != nil {
			return err
		} else {
			defer rows.Close()
			cols, err := rows.Columns()
			if err != nil {
				return err
			}
			vals := make([]interface{}, len(cols))
			ints := make([]interface{}, len(cols))
			for i := range ints {
				ints[i] = &vals[i]
			}
			for rows.Next() {
				if err := rows.Scan(ints...); err != nil {
					log.Println(err)
					return err
				}
				if err := functor(q.columns, vals); err != nil {
					break
				}
			}
		}
	}
	return nil
}

func (q _Query) exec(dbtx DbTx, query string, params []interface{}) (sql.Result, error) {
	start := time.Now().UnixNano()
	defer func() {
		if Debug {
			log.Printf("Running SQL - [%s], params=%v, duration=%dms", query, params, (time.Now().UnixNano()-start)/1e6)
		}
	}()

	var result sql.Result
	runInTx := func(tx *Tx) error {
		if stmt, txErr := tx.Prepare(query); txErr != nil {
			return txErr
		} else {
			result, txErr = stmt.Exec(params...)
			return txErr
		}
	}
	if db, ok := dbtx.(*Db); ok {
		err := WithinTx(db, func(tx *Tx) error {
			return runInTx(tx)
		})
		return result, err
	} else if tx, ok := dbtx.(*Tx); ok {
		err := runInTx(tx)
		return result, err
	}
	return result, ErrNotDbTxObject
}

////// Select Query

type _SelectQuery struct {
	_Query
}

func (q _SelectQuery) Where(f Filter) Query {
	q.where = f
	return q
}

func (q _SelectQuery) OrderBy(by ...string) Query {
	q.orderBy = by
	return q
}

func (q _SelectQuery) GroupBy(by ...string) Query {
	q.groupBy = by
	return q
}

func (q _SelectQuery) Limit(offsets ...int64) Query {
	var start, size int64
	if len(offsets) > 0 {
		if len(offsets) == 1 {
			start, size = 0, offsets[0]
		} else {
			start, size = offsets[0], offsets[1]
		}
		q.limit = []int64{start, size}
	}
	return q
}

func (q _SelectQuery) Page(number, size int) Query {
	start := int64((number-1)*size)
	end := int64(size)
	return q.Limit(start, end)
}

func (q _SelectQuery) SelectOne(dbtx DbTx, functor QueryRowVisitor) error {
	if len(q.columns) == 0 {
		return ErrNotEnoughColumns
	}
	query, params := q.sqlStringAndParam(dbtx.DriverName())
	return q.queryOne(dbtx, query, params, functor)
}

func (q _SelectQuery) SelectList(dbtx DbTx, functor QueryRowVisitor) error {
	if len(q.columns) == 0 {
		return ErrNotEnoughColumns
	}
	query, params := q.sqlStringAndParam(dbtx.DriverName())
	return q.query(dbtx, query, params, functor)
}

func (q _SelectQuery) sqlStringAndParam(driverName string) (string, []interface{}) {
	schema, table, alias := q.model.Names()
	fields, params := q.columns.fieldsAndParams(alias, driverName)
	query := fmt.Sprintf("SELECT %s FROM %s",
		strings.Join(fields, ", "),
		tableNamewithAlias(schema, table, alias, driverName))
	if remains, extras := q.sqlRemains(alias, driverName); remains != "" && len(extras) > 0 {
		query = fmt.Sprintf("%s %s", query, remains)
		params = append(params, extras...)
	}
	return rebindSqlParams(query, driverName), params
}

func (q _SelectQuery) Explain(driverName string) string {
	query, params := q.sqlStringAndParam(driverName)
	return fmt.Sprintf("[%s], params=%v", query, params)
}

func (q _SelectQuery) String() string {
	query, params := q.sqlStringAndParam("mysql")
	return fmt.Sprintf("[%s], params=%v", query, params)
}

////// Insert Query

type _InsertQuery struct {
	_Query
}

func (q _InsertQuery) Exec(dbtx DbTx) (sql.Result, error) {
	if len(q.columns) == 0 {
		return nil, ErrNotEnoughColumns
	}
	query, params := q.sqlStringAndParam(dbtx.DriverName())
	return q.exec(dbtx, query, params)
}

func (q _InsertQuery) Where(f Filter) Query         { return q }
func (q _InsertQuery) OrderBy(by ...string) Query   { return q }
func (q _InsertQuery) GroupBy(by ...string) Query   { return q }
func (q _InsertQuery) Limit(offsets ...int64) Query { return q }
func (q _InsertQuery) Page(number, size int) Query  { return q }

func (q _InsertQuery) sqlStringAndParam(driverName string) (string, []interface{}) {
	schema, table, _ := q.model.Names()
	fields, params := q.columns.fieldsAndParams("", driverName)
	marks := paramMarkers(len(q.columns), driverName)
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		tableNamewithAlias(schema, table, "", driverName),
		strings.Join(fields, ", "), marks)
	return rebindSqlParams(query, driverName), params
}

func (q _InsertQuery) String() string {
	query, params := q.sqlStringAndParam("mysql")
	return fmt.Sprintf("[%s], params=%v", query, params)
}

///// Update Query

type _UpdateQuery struct {
	_Query
}

func (q _UpdateQuery) Exec(dbtx DbTx) (sql.Result, error) {
	if len(q.columns) == 0 {
		return nil, ErrNotEnoughColumns
	}
	query, params := q.sqlStringAndParam(dbtx.DriverName())
	return q.exec(dbtx, query, params)
}

func (q _UpdateQuery) Where(f Filter) Query {
	q.where = f
	return q
}

func (q _UpdateQuery) OrderBy(by ...string) Query   { return q }
func (q _UpdateQuery) GroupBy(by ...string) Query   { return q }
func (q _UpdateQuery) Limit(offsets ...int64) Query { return q }
func (q _UpdateQuery) Page(number, size int) Query  { return q }

func (q _UpdateQuery) sqlStringAndParam(driverName string) (string, []interface{}) {
	schema, table, _ := q.model.Names()
	fields, params := q.columns.fieldsAndParams("", driverName)
	for i, f := range fields {
		fields[i] = fmt.Sprintf("%s = ?", f)
	}
	query := fmt.Sprintf("UPDATE %s SET %s",
		tableNamewithAlias(schema, table, "", driverName),
		strings.Join(fields, ", "))
	if remains, extras := q.sqlRemains("", driverName); remains != "" && len(extras) > 0 {
		query = fmt.Sprintf("%s %s", query, remains)
		params = append(params, extras...)
	}
	return rebindSqlParams(query, driverName), params
}

func (q _UpdateQuery) String() string {
	query, params := q.sqlStringAndParam("mysql")
	return fmt.Sprintf("[%s], params=%v", query, params)
}

///// Delete Query

type _DeleteQuery struct {
	_Query
}

func (q _DeleteQuery) Exec(dbtx DbTx) (sql.Result, error) {
	query, params := q.sqlStringAndParam(dbtx.DriverName())
	return q.exec(dbtx, query, params)
}

func (q _DeleteQuery) Where(f Filter) Query {
	q.where = f
	return q
}

func (q _DeleteQuery) OrderBy(by ...string) Query   { return q }
func (q _DeleteQuery) GroupBy(by ...string) Query   { return q }
func (q _DeleteQuery) Limit(offsets ...int64) Query { return q }
func (q _DeleteQuery) Page(number, size int) Query  { return q }

func (q _DeleteQuery) sqlStringAndParam(driverName string) (string, []interface{}) {
	schema, table, _ := q.model.Names()
	query := fmt.Sprintf("DELETE FROM %s", tableNamewithAlias(schema, table, "", driverName))
	var params []interface{}
	if remains, extras := q.sqlRemains("", driverName); remains != "" && len(extras) > 0 {
		query = fmt.Sprintf("%s %s", query, remains)
		params = extras
	}
	return rebindSqlParams(query, driverName), params
}

func (q _DeleteQuery) String() string {
	query, params := q.sqlStringAndParam("mysql")
	return fmt.Sprintf("[%s], params=%v", query, params)
}
