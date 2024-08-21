package migration

import (
	"database/sql"
	"reflect"
	"strings"
	"time"
)

type (
	DbField struct {
		Name string
		Type string
		Tags []string
	}
	DbFields []DbField
	// DbTable  struct {
	// 	Name   string
	// 	Fields DbFields
	// }
	DbTables    map[string]DbFields
	PgMigration struct {
		tables DbTables
	}
)

func NewPgMigration() *PgMigration {
	pm := &PgMigration{}
	pm.tables = DbTables{}
	return pm
}

func (pm *PgMigration) AddTable(name string, model interface{}) {
	// pm.tables[name] = append(pm.tables, DbTable{Name: name, Fields: getModelFields(model)})
	pm.tables[name] = getModelFields(model)
}

func (pm *PgMigration) GetTables() DbTables {
	return pm.tables
}

func unpackFieldName(fieldName string) string {
	fn := ""
	delimiter := ""
	b := false
	for _, r := range fieldName {
		if r >= 'A' && r <= 'Z' {
			if !b {
				fn += delimiter
			}
			fn += strings.ToLower(string(r))
			b = true
		} else {
			fn += string(r)
			b = false
		}
		delimiter = "_"
	}
	return fn
}

func getModelFields(model interface{}) DbFields {
	t := reflect.TypeOf(model)
	dbFields := DbFields{}
	strucName := t.Name()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldName := strings.ToLower(strucName) + "_" + unpackFieldName(field.Name)
		tags := strings.Split(field.Tag.Get("db"), ";")

		if fieldName == "" {
			fieldName = strings.ToLower(field.Name)
		}
		fieldType := field.Type.Kind()
		sqlType := ""
		if fieldType != reflect.Struct {
			switch fieldType {
			case reflect.String:
				sqlType = "TEXT"
			case reflect.Ptr:
				sqlType = "PTR"
			case reflect.Int, reflect.Int64, reflect.Uint, reflect.Uint64:
				sqlType = "BIGINT"
			case reflect.Int32, reflect.Int16, reflect.Int8:
				sqlType = "INTEGER"
			case reflect.Float64:
				sqlType = "REAL"
			case reflect.Bool:
				sqlType = "BOOLEAN"
			case reflect.Slice:
				sqlType = "SLICE"
			default:
				sqlType = "TEXT"
			}
			dbFields = append(dbFields, DbField{Name: fieldName, Type: sqlType, Tags: tags})
		} else {
			switch field.Type {
			case reflect.TypeOf(time.Time{}), reflect.TypeOf(sql.NullTime{}):
				dbFields = append(dbFields, DbField{Name: fieldName, Type: "TIMESTAMP"})
			default:
				dbFields = append(dbFields, getModelFields(reflect.New(field.Type).Elem().Interface())...)
			}
		}
	}
	return dbFields
}
