// Code generated by ModelQ
// key_column_usage.go contains model for the database table [information_schema.key_column_usage]

package postgres

import (
	"database/sql"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/kitekjk/modelq/gmq"
	"strings"
)

type KeyColumnUsage struct {
	ConstraintCatalog          string `json:"constraint_catalog"`
	ConstraintSchema           string `json:"constraint_schema"`
	ConstraintName             string `json:"constraint_name"`
	TableCatalog               string `json:"table_catalog"`
	TableSchema                string `json:"table_schema"`
	TableName                  string `json:"table_name"`
	ColumnName                 string `json:"column_name"`
	OrdinalPosition            int    `json:"ordinal_position"`
	PositionInUniqueConstraint int    `json:"position_in_unique_constraint"`
}

// Start of the KeyColumnUsage APIs.

func (obj KeyColumnUsage) String() string {
	if data, err := json.Marshal(obj); err != nil {
		return fmt.Sprintf("<KeyColumnUsage>")
	} else {
		return string(data)
	}
}

func (obj KeyColumnUsage) Insert(dbtx gmq.DbTx) (KeyColumnUsage, error) {
	_, err := KeyColumnUsageObjs.Insert(obj).Run(dbtx)
	return obj, err
}

func (obj KeyColumnUsage) Update(dbtx gmq.DbTx) (int64, error) {
	return 0, gmq.ErrNoPrimaryKeyDefined
}

func (obj KeyColumnUsage) Delete(dbtx gmq.DbTx) (int64, error) {
	return 0, gmq.ErrNoPrimaryKeyDefined
}

// Start of the inner Query Api

type _KeyColumnUsageQuery struct {
	gmq.Query
}

func (q _KeyColumnUsageQuery) Where(f gmq.Filter) _KeyColumnUsageQuery {
	q.Query = q.Query.Where(f)
	return q
}

func (q _KeyColumnUsageQuery) OrderBy(by ...string) _KeyColumnUsageQuery {
	tBy := make([]string, 0, len(by))
	for _, b := range by {
		sortDir := ""
		if b[0] == '-' || b[0] == '+' {
			sortDir = string(b[0])
			b = b[1:]
		}
		if col, ok := KeyColumnUsageObjs.fcMap[b]; ok {
			tBy = append(tBy, sortDir+col)
		}
	}
	q.Query = q.Query.OrderBy(tBy...)
	return q
}

func (q _KeyColumnUsageQuery) GroupBy(by ...string) _KeyColumnUsageQuery {
	tBy := make([]string, 0, len(by))
	for _, b := range by {
		if col, ok := KeyColumnUsageObjs.fcMap[b]; ok {
			tBy = append(tBy, col)
		}
	}
	q.Query = q.Query.GroupBy(tBy...)
	return q
}

func (q _KeyColumnUsageQuery) Limit(offsets ...int64) _KeyColumnUsageQuery {
	q.Query = q.Query.Limit(offsets...)
	return q
}

func (q _KeyColumnUsageQuery) Page(number, size int) _KeyColumnUsageQuery {
	q.Query = q.Query.Page(number, size)
	return q
}

func (q _KeyColumnUsageQuery) Run(dbtx gmq.DbTx) (sql.Result, error) {
	return q.Query.Exec(dbtx)
}

type KeyColumnUsageRowVisitor func(obj KeyColumnUsage) error

func (q _KeyColumnUsageQuery) Iterate(dbtx gmq.DbTx, functor KeyColumnUsageRowVisitor) error {
	return q.Query.SelectList(dbtx, func(columns []gmq.Column, rb []interface{}) error {
		obj := KeyColumnUsageObjs.toKeyColumnUsage(columns, rb)
		return functor(obj)
	})
}

func (q _KeyColumnUsageQuery) One(dbtx gmq.DbTx) (KeyColumnUsage, error) {
	var obj KeyColumnUsage
	err := q.Query.SelectOne(dbtx, func(columns []gmq.Column, rb []interface{}) error {
		obj = KeyColumnUsageObjs.toKeyColumnUsage(columns, rb)
		return nil
	})
	return obj, err
}

func (q _KeyColumnUsageQuery) List(dbtx gmq.DbTx) ([]KeyColumnUsage, error) {
	result := make([]KeyColumnUsage, 0, 10)
	err := q.Query.SelectList(dbtx, func(columns []gmq.Column, rb []interface{}) error {
		obj := KeyColumnUsageObjs.toKeyColumnUsage(columns, rb)
		result = append(result, obj)
		return nil
	})
	return result, err
}

// Start of the model facade Apis.

type _KeyColumnUsageObjs struct {
	fcMap map[string]string
}

func (o _KeyColumnUsageObjs) Names() (schema, tbl, alias string) {
	return "information_schema", "key_column_usage", "KeyColumnUsage"
}

func (o _KeyColumnUsageObjs) Select(fields ...string) _KeyColumnUsageQuery {
	q := _KeyColumnUsageQuery{}
	if len(fields) == 0 {
		fields = []string{"ConstraintCatalog", "ConstraintSchema", "ConstraintName", "TableCatalog", "TableSchema", "TableName", "ColumnName", "OrdinalPosition", "PositionInUniqueConstraint"}
	}
	q.Query = gmq.Select(o, o.columns(fields...))
	return q
}

func (o _KeyColumnUsageObjs) Insert(obj KeyColumnUsage) _KeyColumnUsageQuery {
	q := _KeyColumnUsageQuery{}
	q.Query = gmq.Insert(o, o.columnsWithData(obj, "ConstraintCatalog", "ConstraintSchema", "ConstraintName", "TableCatalog", "TableSchema", "TableName", "ColumnName", "OrdinalPosition", "PositionInUniqueConstraint"))
	return q
}

func (o _KeyColumnUsageObjs) Update(obj KeyColumnUsage, fields ...string) _KeyColumnUsageQuery {
	q := _KeyColumnUsageQuery{}
	q.Query = gmq.Update(o, o.columnsWithData(obj, fields...))
	return q
}

func (o _KeyColumnUsageObjs) Delete() _KeyColumnUsageQuery {
	q := _KeyColumnUsageQuery{}
	q.Query = gmq.Delete(o)
	return q
}

///// Managed Objects Filters definition

func (o _KeyColumnUsageObjs) FilterConstraintCatalog(op string, p string, ps ...string) gmq.Filter {
	params := make([]interface{}, 1+len(ps))
	params[0] = p
	for i := range ps {
		params[i+1] = ps[i]
	}
	return o.newFilter("constraint_catalog", op, params...)
}

func (o _KeyColumnUsageObjs) FilterConstraintSchema(op string, p string, ps ...string) gmq.Filter {
	params := make([]interface{}, 1+len(ps))
	params[0] = p
	for i := range ps {
		params[i+1] = ps[i]
	}
	return o.newFilter("constraint_schema", op, params...)
}

func (o _KeyColumnUsageObjs) FilterConstraintName(op string, p string, ps ...string) gmq.Filter {
	params := make([]interface{}, 1+len(ps))
	params[0] = p
	for i := range ps {
		params[i+1] = ps[i]
	}
	return o.newFilter("constraint_name", op, params...)
}

func (o _KeyColumnUsageObjs) FilterTableCatalog(op string, p string, ps ...string) gmq.Filter {
	params := make([]interface{}, 1+len(ps))
	params[0] = p
	for i := range ps {
		params[i+1] = ps[i]
	}
	return o.newFilter("table_catalog", op, params...)
}

func (o _KeyColumnUsageObjs) FilterTableSchema(op string, p string, ps ...string) gmq.Filter {
	params := make([]interface{}, 1+len(ps))
	params[0] = p
	for i := range ps {
		params[i+1] = ps[i]
	}
	return o.newFilter("table_schema", op, params...)
}

func (o _KeyColumnUsageObjs) FilterTableName(op string, p string, ps ...string) gmq.Filter {
	params := make([]interface{}, 1+len(ps))
	params[0] = p
	for i := range ps {
		params[i+1] = ps[i]
	}
	return o.newFilter("table_name", op, params...)
}

func (o _KeyColumnUsageObjs) FilterColumnName(op string, p string, ps ...string) gmq.Filter {
	params := make([]interface{}, 1+len(ps))
	params[0] = p
	for i := range ps {
		params[i+1] = ps[i]
	}
	return o.newFilter("column_name", op, params...)
}

func (o _KeyColumnUsageObjs) FilterOrdinalPosition(op string, p int, ps ...int) gmq.Filter {
	params := make([]interface{}, 1+len(ps))
	params[0] = p
	for i := range ps {
		params[i+1] = ps[i]
	}
	return o.newFilter("ordinal_position", op, params...)
}

func (o _KeyColumnUsageObjs) FilterPositionInUniqueConstraint(op string, p int, ps ...int) gmq.Filter {
	params := make([]interface{}, 1+len(ps))
	params[0] = p
	for i := range ps {
		params[i+1] = ps[i]
	}
	return o.newFilter("position_in_unique_constraint", op, params...)
}

///// Managed Objects Columns definition

func (o _KeyColumnUsageObjs) ColumnConstraintCatalog(p ...string) gmq.Column {
	var value interface{}
	if len(p) > 0 {
		value = p[0]
	}
	return gmq.Column{"constraint_catalog", value}
}

func (o _KeyColumnUsageObjs) ColumnConstraintSchema(p ...string) gmq.Column {
	var value interface{}
	if len(p) > 0 {
		value = p[0]
	}
	return gmq.Column{"constraint_schema", value}
}

func (o _KeyColumnUsageObjs) ColumnConstraintName(p ...string) gmq.Column {
	var value interface{}
	if len(p) > 0 {
		value = p[0]
	}
	return gmq.Column{"constraint_name", value}
}

func (o _KeyColumnUsageObjs) ColumnTableCatalog(p ...string) gmq.Column {
	var value interface{}
	if len(p) > 0 {
		value = p[0]
	}
	return gmq.Column{"table_catalog", value}
}

func (o _KeyColumnUsageObjs) ColumnTableSchema(p ...string) gmq.Column {
	var value interface{}
	if len(p) > 0 {
		value = p[0]
	}
	return gmq.Column{"table_schema", value}
}

func (o _KeyColumnUsageObjs) ColumnTableName(p ...string) gmq.Column {
	var value interface{}
	if len(p) > 0 {
		value = p[0]
	}
	return gmq.Column{"table_name", value}
}

func (o _KeyColumnUsageObjs) ColumnColumnName(p ...string) gmq.Column {
	var value interface{}
	if len(p) > 0 {
		value = p[0]
	}
	return gmq.Column{"column_name", value}
}

func (o _KeyColumnUsageObjs) ColumnOrdinalPosition(p ...int) gmq.Column {
	var value interface{}
	if len(p) > 0 {
		value = p[0]
	}
	return gmq.Column{"ordinal_position", value}
}

func (o _KeyColumnUsageObjs) ColumnPositionInUniqueConstraint(p ...int) gmq.Column {
	var value interface{}
	if len(p) > 0 {
		value = p[0]
	}
	return gmq.Column{"position_in_unique_constraint", value}
}

////// Internal helper funcs

func (o _KeyColumnUsageObjs) newFilter(name, op string, params ...interface{}) gmq.Filter {
	if strings.ToUpper(op) == "IN" {
		return gmq.InFilter(name, params)
	}
	return gmq.UnitFilter(name, op, params[0])
}

func (o _KeyColumnUsageObjs) toKeyColumnUsage(columns []gmq.Column, rb []interface{}) KeyColumnUsage {
	obj := KeyColumnUsage{}
	if len(columns) == len(rb) {
		for i := range columns {
			switch columns[i].Name {
			case "constraint_catalog":
				obj.ConstraintCatalog = gmq.CastString(rb[i])
			case "constraint_schema":
				obj.ConstraintSchema = gmq.CastString(rb[i])
			case "constraint_name":
				obj.ConstraintName = gmq.CastString(rb[i])
			case "table_catalog":
				obj.TableCatalog = gmq.CastString(rb[i])
			case "table_schema":
				obj.TableSchema = gmq.CastString(rb[i])
			case "table_name":
				obj.TableName = gmq.CastString(rb[i])
			case "column_name":
				obj.ColumnName = gmq.CastString(rb[i])
			case "ordinal_position":
				obj.OrdinalPosition = gmq.CastInt(rb[i])
			case "position_in_unique_constraint":
				obj.PositionInUniqueConstraint = gmq.CastInt(rb[i])
			}
		}
	}
	return obj
}

func (o _KeyColumnUsageObjs) columns(fields ...string) []gmq.Column {
	data := make([]gmq.Column, 0, len(fields))
	for _, f := range fields {
		switch f {
		case "ConstraintCatalog":
			data = append(data, o.ColumnConstraintCatalog())
		case "ConstraintSchema":
			data = append(data, o.ColumnConstraintSchema())
		case "ConstraintName":
			data = append(data, o.ColumnConstraintName())
		case "TableCatalog":
			data = append(data, o.ColumnTableCatalog())
		case "TableSchema":
			data = append(data, o.ColumnTableSchema())
		case "TableName":
			data = append(data, o.ColumnTableName())
		case "ColumnName":
			data = append(data, o.ColumnColumnName())
		case "OrdinalPosition":
			data = append(data, o.ColumnOrdinalPosition())
		case "PositionInUniqueConstraint":
			data = append(data, o.ColumnPositionInUniqueConstraint())
		}
	}
	return data
}

func (o _KeyColumnUsageObjs) columnsWithData(obj KeyColumnUsage, fields ...string) []gmq.Column {
	data := make([]gmq.Column, 0, len(fields))
	for _, f := range fields {
		switch f {
		case "ConstraintCatalog":
			data = append(data, o.ColumnConstraintCatalog(obj.ConstraintCatalog))
		case "ConstraintSchema":
			data = append(data, o.ColumnConstraintSchema(obj.ConstraintSchema))
		case "ConstraintName":
			data = append(data, o.ColumnConstraintName(obj.ConstraintName))
		case "TableCatalog":
			data = append(data, o.ColumnTableCatalog(obj.TableCatalog))
		case "TableSchema":
			data = append(data, o.ColumnTableSchema(obj.TableSchema))
		case "TableName":
			data = append(data, o.ColumnTableName(obj.TableName))
		case "ColumnName":
			data = append(data, o.ColumnColumnName(obj.ColumnName))
		case "OrdinalPosition":
			data = append(data, o.ColumnOrdinalPosition(obj.OrdinalPosition))
		case "PositionInUniqueConstraint":
			data = append(data, o.ColumnPositionInUniqueConstraint(obj.PositionInUniqueConstraint))
		}
	}
	return data
}

var KeyColumnUsageObjs _KeyColumnUsageObjs

func init() {
	KeyColumnUsageObjs.fcMap = map[string]string{
		"ConstraintCatalog":          "constraint_catalog",
		"ConstraintSchema":           "constraint_schema",
		"ConstraintName":             "constraint_name",
		"TableCatalog":               "table_catalog",
		"TableSchema":                "table_schema",
		"TableName":                  "table_name",
		"ColumnName":                 "column_name",
		"OrdinalPosition":            "ordinal_position",
		"PositionInUniqueConstraint": "position_in_unique_constraint",
	}
	gob.Register(KeyColumnUsage{})
}
