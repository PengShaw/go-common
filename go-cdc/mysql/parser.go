package mysql

import (
	"errors"
	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/go-mysql-org/go-mysql/schema"
	"reflect"
	"strings"
	"time"
)

type binlogParser struct{}

func (p *binlogParser) getRowItem(e *canal.RowsEvent, rowIndex int, table Table) (interface{}, error) {
	// 通过反射创建新对象
	if table.Model == nil {
		return nil, errors.New("table model is nil" + table.Name)
	}
	t := reflect.TypeOf(table.Model)
	if t.Kind() == reflect.Ptr { // 指针类型的情况
		t = t.Elem()
	}
	rowItem := reflect.New(t).Elem()
	t = rowItem.Type()
	for i := 0; i < t.NumField(); i++ {
		tags := parseTags(t.Field(i).Tag)
		columnName, ok := tags["COLUMN"]
		if !ok || columnName == "COLUMN" {
			continue
		}

		fieldType := rowItem.Field(i).Type()
		fieldTypeName := fieldType.Name()

		switch fieldTypeName {
		case "int8", "int16", "int32", "int64", "int":
			rowItem.Field(i).SetInt(p.parseInt64(e, rowIndex, columnName))
		case "uint8", "uint16", "uint32", "uint64", "uint":
			rowItem.Field(i).SetUint(p.parseUint64(e, rowIndex, columnName))
		case "bool":
			rowItem.Field(i).SetBool(p.parseBool(e, rowIndex, columnName))
		case "float32", "float64":
			rowItem.Field(i).SetFloat(p.parseFloat(e, rowIndex, columnName))
		case "string":
			rowItem.Field(i).SetString(p.parseString(e, rowIndex, columnName))
		case "Time":
			rowItem.Field(i).Set(reflect.ValueOf(p.parseTime(e, rowIndex, columnName)))
		}

		// *time.Time
		if fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
			columnId, ok := p.getColumnIdByName(e, columnName)
			if e.Rows[rowIndex][columnId] == nil || !ok {
				continue
			}
			fieldTypeName = fieldType.Name()
			switch fieldTypeName {
			case "Time":
				v := p.parseTime(e, rowIndex, columnName)
				rowItem.Field(i).Set(reflect.ValueOf(&v))
			}
		}

	}
	return rowItem.Interface(), nil
}

func parseTags(tags reflect.StructTag) map[string]string {
	settings := map[string]string{}
	for _, str := range []string{tags.Get("binlog")} {
		tags := strings.Split(str, ";")
		for _, value := range tags {
			v := strings.Split(value, ":")
			k := strings.TrimSpace(strings.ToUpper(v[0]))
			if len(v) >= 2 {
				settings[k] = strings.Join(v[1:], ":")
			} else {
				settings[k] = k
			}
		}
	}
	return settings
}

func (p *binlogParser) parseString(e *canal.RowsEvent, n int, columnName string) string {
	columnId, ok := p.getColumnIdByName(e, columnName)
	if !ok {
		return ""
	}
	if e.Table.Columns[columnId].Type == schema.TYPE_ENUM {
		values := e.Table.Columns[columnId].EnumValues
		if len(values) == 0 {
			return ""
		}
		if e.Rows[n][columnId] == nil {
			return ""
		}

		return values[e.Rows[n][columnId].(int64)-1]
	}

	value := e.Rows[n][columnId]

	switch value := value.(type) {
	case []byte:
		return string(value)
	case string:
		return value
	}
	return ""
}

func (p *binlogParser) parseFloat(e *canal.RowsEvent, n int, columnName string) float64 {
	columnId, ok := p.getColumnIdByName(e, columnName)
	if e.Table.Columns[columnId].Type != schema.TYPE_FLOAT || !ok {
		return float64(0)
	}

	switch e.Rows[n][columnId].(type) {
	case float32:
		return float64(e.Rows[n][columnId].(float32))
	case float64:
		return e.Rows[n][columnId].(float64)
	}
	return float64(0)
}

func (p *binlogParser) parseTime(e *canal.RowsEvent, n int, columnName string) time.Time {
	columnId, ok := p.getColumnIdByName(e, columnName)

	cond := e.Table.Columns[columnId].Type == schema.TYPE_DATETIME ||
		e.Table.Columns[columnId].Type == schema.TYPE_TIMESTAMP ||
		e.Table.Columns[columnId].Type == schema.TYPE_DATE ||
		e.Table.Columns[columnId].Type == schema.TYPE_TIME
	if !cond || !ok {
		return time.Time{}
	}

	t, _ := time.Parse("2006-01-02 15:04:05", e.Rows[n][columnId].(string))
	return t
}

func (p *binlogParser) parseBool(e *canal.RowsEvent, n int, columnName string) bool {
	val := p.parseInt64(e, n, columnName)
	if val == 1 {
		return true
	}
	return false
}

func (p *binlogParser) parseInt64(e *canal.RowsEvent, n int, columnName string) int64 {
	columnId, ok := p.getColumnIdByName(e, columnName)
	if e.Table.Columns[columnId].Type != schema.TYPE_NUMBER || !ok {
		return 0
	}

	switch e.Rows[n][columnId].(type) {
	case int8:
		return int64(e.Rows[n][columnId].(int8))
	case int16:
		return int64(e.Rows[n][columnId].(int16))
	case int32:
		return int64(e.Rows[n][columnId].(int32))
	case int64:
		return e.Rows[n][columnId].(int64)
	case int:
		return int64(e.Rows[n][columnId].(int))
	case uint8:
		return int64(e.Rows[n][columnId].(uint8))
	case uint16:
		return int64(e.Rows[n][columnId].(uint16))
	case uint32:
		return int64(e.Rows[n][columnId].(uint32))
	case uint64:
		return int64(e.Rows[n][columnId].(uint64))
	case uint:
		return int64(e.Rows[n][columnId].(uint))
	}
	return 0
}

func (p *binlogParser) parseUint64(e *canal.RowsEvent, n int, columnName string) uint64 {
	columnId, ok := p.getColumnIdByName(e, columnName)
	if e.Table.Columns[columnId].Type != schema.TYPE_NUMBER || !ok {
		return 0
	}

	switch e.Rows[n][columnId].(type) {
	case int8:
		return uint64(e.Rows[n][columnId].(int8))
	case int16:
		return uint64(e.Rows[n][columnId].(int16))
	case int32:
		return uint64(e.Rows[n][columnId].(int32))
	case int64:
		return uint64(e.Rows[n][columnId].(int64))
	case int:
		return uint64(e.Rows[n][columnId].(int))
	case uint8:
		return uint64(e.Rows[n][columnId].(uint8))
	case uint16:
		return uint64(e.Rows[n][columnId].(uint16))
	case uint32:
		return uint64(e.Rows[n][columnId].(uint32))
	case uint64:
		return e.Rows[n][columnId].(uint64)
	case uint:
		return uint64(e.Rows[n][columnId].(uint))
	}
	return 0
}

func (p *binlogParser) getColumnIdByName(e *canal.RowsEvent, columnName string) (int, bool) {
	for id, value := range e.Table.Columns {
		if value.Name == columnName {
			return id, true
		}
	}
	return 0, false
}
