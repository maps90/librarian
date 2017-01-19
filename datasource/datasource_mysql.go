package datasource

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"reflect"
	"regexp"
	"strings"
	"time"
)

type MysqlAccess struct {
	db *sqlx.DB
}

func NewMysqlRepository(server string) (*MysqlAccess, error) {
	db, err := sqlx.Open("mysql", server)
	//set max idle conn to preserve connection pool
	db.SetMaxIdleConns(10)
	//limit connection pool
	db.SetMaxOpenConns(10)
	return &MysqlAccess{db}, err
}

func fieldName(field reflect.StructField) string {
	if t := field.Tag.Get("db"); t != "" {
		return t
	}
	return underscore(field.Name)
}

var camel = regexp.MustCompile("(^[^A-Z]*|[A-Z]*)([A-Z][^A-Z]+|$)")

func underscore(s string) string {
	var a []string
	for _, sub := range camel.FindAllStringSubmatch(s, -1) {
		if sub[1] != "" {
			a = append(a, sub[1])
		}
		if sub[2] != "" {
			a = append(a, sub[2])
		}
	}
	return strings.ToLower(strings.Join(a, "_"))
}

func isZero(v reflect.Value, dv interface{}) bool {

	switch v.Kind() {
	case reflect.Func, reflect.Map, reflect.Slice:
		return v.IsNil()
	case reflect.Array:
		z := true
		for i := 0; i < v.Len(); i++ {
			z = z && isZero(v.Index(i), v)
		}
		return z
	case reflect.Struct:
		z := true
		if v.Type().String() == "time.Time" {
			zeroTime := time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)
			comparedTime, found := v.Interface().(time.Time)
			if found {
				return zeroTime.Equal(comparedTime)
			}
		}
		for i := 0; i < v.NumField(); i++ {
			z = z && isZero(v.Field(i), v)
		}
		return z
	}

	// Other types is compared directly
	z := reflect.Zero(v.Type())
	return dv == z.Interface()
}

func extractStructs(data Data) (map[string]interface{}, error) {
	val := reflect.Indirect(reflect.ValueOf(data))
	t := val.Type()

	if val.Kind() == reflect.Struct {
		result := make(map[string]interface{})
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			fieldVal := reflect.Indirect(val.Field(i))
			value := fieldVal.Interface()
			if !isZero(fieldVal, value) {
				fieldStr := fieldName(field)
				result[fieldStr] = value
			}
		}
		return result, nil
	} else {
		return nil, errors.New("Submitted type must in struct type")
	}
}

func (r *MysqlAccess) Insert(data Data) (interface{}, error) {
	db := r.db

	columnsMap, err := extractStructs(data)
	if err != nil {
		return nil, err
	}

	buffColumnNames := bytes.NewBufferString("")
	buffValueNames := bytes.NewBufferString("")
	index := 0
	values := []interface{}{}
	for k, v := range columnsMap {
		if index < (len(columnsMap) - 1) {
			buffColumnNames.WriteString(k + ",")
			buffValueNames.WriteString("?,")
		} else {
			buffColumnNames.WriteString(k)
			buffValueNames.WriteString("?")
		}
		values = append(values, v)
		index++
	}

	queryString := fmt.Sprintf("INSERT INTO %s (%s) VALUES(%s)", data.PersistenceName(), buffColumnNames.String(), buffValueNames.String())
	fmt.Println("Query string", queryString)
	result := db.MustExec(queryString, values...)

	count, err := result.RowsAffected()

	if err != nil {
		return nil, err
	}

	msg := fmt.Sprintf("%d row is inserted successfully", count)

	return struct{ msg string }{msg: msg}, nil
}

func (r *MysqlAccess) Update(id interface{}, data Data) (interface{}, error) {
	db := r.db

	columnsMap, err := extractStructs(data)
	if err != nil {
		return nil, err
	}

	buffColumnNames := bytes.NewBufferString("")
	index := 0
	values := []interface{}{}
	for k, v := range columnsMap {
		if index < (len(columnsMap) - 1) {
			buffColumnNames.WriteString(k + "=?,")
		} else {
			buffColumnNames.WriteString(k + "=?")
		}
		values = append(values, v)
		index++
	}

	var columnId string

	idMap, found := id.(map[string]interface{})

	if !found {
		return nil, errors.New("Id should have type map[string]interface{}")
	}

	for k, v := range idMap {
		columnId = k
		values = append(values, v)
	}
	queryString := fmt.Sprintf("UPDATE %s SET %s WHERE %s = ?", data.PersistenceName(), buffColumnNames.String(), columnId)
	result := db.MustExec(queryString, values...)

	count, err := result.RowsAffected()

	if err != nil {
		return nil, err
	}

	msg := fmt.Sprintf("%d row is updated successfully", count)

	return struct{ msg string }{msg: msg}, nil
}

func (r *MysqlAccess) Delete(id interface{}, data Data) (interface{}, error) {
	db := r.db

	var columnId string
	var value interface{}

	idMap, found := id.(map[string]interface{})

	if !found {
		return nil, errors.New("Id should have type map[string]interface{}")
	}
	for k, v := range idMap {
		columnId = k
		value = v
	}

	queryString := fmt.Sprintf("DELETE FROM %s WHERE %s=?", data.PersistenceName(), columnId)
	result := db.MustExec(queryString, value)

	count, err := result.RowsAffected()

	if err != nil {
		return nil, err
	}

	msg := fmt.Sprintf("%d row is deleted successfully", count)

	return struct{ msg string }{msg: msg}, nil
}

func (r *MysqlAccess) Find(data Data, query map[string]interface{}, order []string, results interface{}) error {
	buffWhere := bytes.NewBufferString("")
	index := 0
	values := []interface{}{}
	for k, v := range query {
		if index < (len(query) - 1) {
			buffWhere.WriteString(k + "=? AND ")
		} else {
			buffWhere.WriteString(k + "=?")
		}
		values = append(values, v)
		index++
	}
	queryString := fmt.Sprintf("SELECT * FROM %s WHERE %s ORDER BY %s", data.PersistenceName(), buffWhere.String(), strings.Join(order, ","))

	db := r.db

	q, err := db.Queryx(queryString, values...)

	if err != nil {
		return err
	}

	err = sqlx.StructScan(q, results)

	if err != nil {
		return err
	}

	return nil
}

func (r *MysqlAccess) FindById(data Data, id interface{}, result interface{}) error {
	db := r.db

	var columnId string
	var idValue interface{}

	idMap, found := id.(map[string]interface{})

	if !found {
		return errors.New("Id should have type map[string]interface{}")
	}

	for k, v := range idMap {
		columnId, idValue = k, v
	}

	queryString := fmt.Sprintf("SELECT * FROM %s WHERE %s=? LIMIT 1", data.PersistenceName(), columnId)

	return db.Get(result, queryString, idValue)
}

func (r *MysqlAccess) FindPaging(data Data, query map[string]interface{}, order []string, page, limit uint, results interface{}) error {
	db := r.db

	offset := (page - 1) * limit

	buffWhere := bytes.NewBufferString("")
	index := 0
	values := []interface{}{}
	for k, v := range query {
		if index < (len(query) - 1) {
			buffWhere.WriteString(k + "=? AND ")
		} else {
			buffWhere.WriteString(k + "=?")
		}
		values = append(values, v)
		index++
	}
	queryString := fmt.Sprintf("SELECT * FROM %s WHERE %s ORDER BY %s LIMIT %d, %d", data.PersistenceName(), buffWhere.String(), strings.Join(order, ","), offset, limit)

	fmt.Println("Query string", queryString)
	q, err := db.Queryx(queryString, values...)
	if err != nil {
		return err
	}

	err = sqlx.StructScan(q, results)

	if err != nil {
		return err
	}

	return nil
}
