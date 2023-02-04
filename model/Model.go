package model

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"data-connector/log"

	"github.com/google/uuid"
)

/*
	create by: Hoangnd
	create at: 2021-08-07
	des: hàm chuyển các row query được sang dạng map[string]string
*/

func PackageData(rows *sql.Rows) []map[string]string {
	columns, err := rows.Columns()
	if err != nil {
		log.Error(err.Error(), map[string]interface{}{
			"scope": log.Trace(),
		})
		return nil
	}
	if len(columns) == 0 {
		return nil
	}
	values := make([]sql.RawBytes, len(columns))

	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	data := make([]map[string]string, 0)
	for rows.Next() {
		newRow := make(map[string]string)
		err := rows.Scan(scanArgs...)
		if err != nil {
			log.Error(err.Error(), map[string]interface{}{
				"scope": log.Trace(),
			})
			return nil
		}
		var value string
		for i, col := range values {
			if col == nil {
				value = ""
			} else {
				value = string(col)
			}
			newRow[columns[i]] = value
		}
		data = append(data, newRow)
	}
	return data
}

/*
	create by: Hoangnd
	create at: 2021-08-07
	des: Tạo mới uuid
*/
func CreateUuid() string {
	return uuid.New().String()
}

/*
	create by: Hoangnd
	create at: 2021-08-07
	des: Hàm xử lý insert 1 bản ghi, đầu vào là 1 struct
*/
func Insert(v interface{}) (interface{}, error) {
	tableName, err := GetStructTableName(v)
	returningValue := ""
	if err != nil {
		return nil, err
	}
	db := getConnection()
	defer db.Close()
	val := getStructVal(v)
	columnArray := make([]string, 0)
	valueArray := make([]string, 0)
	for i := 0; i < val.NumField(); i++ {
		valueField := GetStructField(val, i)
		typeField := val.Type().Field(i)
		colTag := typeField.Tag.Get("json")
		colType := typeField.Tag.Get("type")
		columnArray = append(columnArray, colTag)
		value := valueField.Interface()
		colPrimary := typeField.Tag.Get("primary")
		colName := typeField.Tag.Get("db")
		v := fmt.Sprintf("%v", value)
		if colPrimary == "true" {
			v = "'" + v + "'"
			returningValue = " RETURNING " + colName
		} else {
			if (value == "" || len(v) == 0) && colType != "string" {
				v = "NULL"
			} else {
				v = strings.ReplaceAll(v, "'", "''")
				v = "'" + v + "'"
			}
		}
		valueArray = append(valueArray, v)
		// if valueField.Kind() == reflect.Struct {
		// 	GetDataInsertSql(valueField)
		// }
	}
	sqlInsert := "INSERT INTO " + tableName + " (" + strings.Join(columnArray[:], ",") + ") VALUES (" + strings.Join(valueArray[:], ",") + ") " + returningValue
	if returningValue != "" {
		id := ""
		err1 := db.QueryRow(sqlInsert).Scan(&id)
		if err1 != nil {
			return nil, err1
		}
		return id, nil
	} else {
		_, err1 := db.Exec(sqlInsert)
		return "", err1
	}
}

/*
	create by: Hoangnd
	create at: 2021-08-07
	des: Hàm xử lý update 1 bản ghi, đầu vào là 1 struct
*/
func Update(v interface{}) error {
	db := getConnection()
	defer db.Close()
	tableName, err := GetStructTableName(v)
	if err != nil {
		return err
	}
	sql := "UPDATE " + tableName + " SET "
	where := ""
	valueArray := make([]string, 0)
	val := getStructVal(v)
	for i := 0; i < val.NumField(); i++ {
		typeField := val.Type().Field(i)
		valueField := GetStructField(val, i)
		colType := typeField.Tag.Get("type")
		colPrimary := typeField.Tag.Get("primary")
		colName := typeField.Tag.Get("db")
		value := valueField.Interface()
		v := fmt.Sprintf("%v", value)
		if colPrimary == "true" {
			where = " WHERE " + colName + " ='" + v + "'"
		} else {
			if (value == "" || len(v) == 0) && colType != "string" {
				v = colName + " = NULL"
			} else {
				v = strings.ReplaceAll(v, "'", "''")
				v = colName + " ='" + v + "'"
			}
			valueArray = append(valueArray, v)
		}
	}
	if where == "" {
		return errors.New("could not find primary key in object " + tableName)
	}
	if len(valueArray) > 0 {
		sql = sql + strings.Join(valueArray[:], ",") + where
		_, err := db.Exec(sql)
		return err
	}
	return errors.New("can not update")
}

/*
	create by: Hoangnd
	create at: 2021-08-07
	des: hàm kiểm tra thông tin của 1 field trong struct
*/
func GetStructField(val reflect.Value, index int) reflect.Value {
	valueField := val.Field(index)
	if valueField.Kind() == reflect.Interface && !valueField.IsNil() {
		elm := valueField.Elem()
		if elm.Kind() == reflect.Ptr && !elm.IsNil() && elm.Elem().Kind() == reflect.Ptr {
			valueField = elm
		}
	}
	if valueField.Kind() == reflect.Ptr {
		valueField = valueField.Elem()
	}
	return valueField
}

/*
	create by: Hoangnd
	create at: 2021-08-07
	des: Hàm lấy ra table name ( bảng của database ) nếu được cấu hình qua interface GetTableName của struct đó
*/
func GetStructTableName(v interface{}) (string, error) {
	val := reflect.ValueOf(v)
	methodVal := val.MethodByName("GetTableName").Call([]reflect.Value{})
	if len(methodVal) == 0 {
		return "", errors.New("can not find table name")
	}
	return methodVal[0].String(), nil
}

/*
	create by: Hoangnd
	create at: 2021-08-07
	des: Hàm lấy ra cấu trúc của struct với đầu vào là interface
*/
func getStructVal(v interface{}) reflect.Value {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Interface && !val.IsNil() {
		elm := val.Elem()
		if elm.Kind() == reflect.Ptr && !elm.IsNil() && elm.Elem().Kind() == reflect.Ptr {
			val = elm
		}
	}
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	return val
}

/*
	create by: Hoangnd
	create at: 2021-08-07
	des: Hàm xử lý xóa 1 bản ghi, đầu vào là 1 struct
*/
func Delete(v interface{}) error {
	db := getConnection()
	defer db.Close()
	where := ""
	tableName, err := GetStructTableName(v)
	if err != nil {
		return err
	}
	sql := "DELETE FROM " + tableName + " WHERE "
	val := getStructVal(v)
	isHasPrimaryKey := false
	for i := 0; i < val.NumField(); i++ {
		valueField := GetStructField(val, i)
		typeField := val.Type().Field(i)
		colPrimary := typeField.Tag.Get("primary")
		colName := typeField.Tag.Get("db")
		value := valueField.Interface()
		v := fmt.Sprintf("%v", value)
		if colPrimary == "true" {
			isHasPrimaryKey = true
			where = colName + " ='" + v + "'"
		}
	}
	if !isHasPrimaryKey {
		return errors.New("could not find primary key in object " + tableName)
	}
	sql += where
	fmt.Println(sql)
	_, err1 := db.Exec(sql)
	return err1
}

/*
	create by: Hoangnd
	create at: 2021-08-07
	des: Hàm xử lý xóa các bản ghi theo điều kiện
*/
func DeleteByCondition(v interface{}, where string) error {
	db := getConnection()
	defer db.Close()
	tableName, err := GetStructTableName(v)
	if err != nil {
		return err
	}
	sql := "DELETE FROM " + tableName + " WHERE "
	sql += where
	fmt.Println(sql)
	_, err1 := db.Exec(sql)
	return err1
}

/*
	create by: Hoangnd
	create at: 2021-08-07
	des: Hàm lấy 1 bản ghi dựa trên id primary, sau đó set data vào struct tương ứng
*/

func FindById(v interface{}) (interface{}, error) {
	db := getConnection()
	defer db.Close()
	where := ""
	tableName, err := GetStructTableName(v)
	fmt.Println(tableName)

	if err != nil {
		return nil, err
	}
	sql := "SELECT * FROM " + tableName + " WHERE "
	val := getStructVal(v)
	for i := 0; i < val.NumField(); i++ {
		valueField := GetStructField(val, i)
		typeField := val.Type().Field(i)
		colPrimary := typeField.Tag.Get("primary")
		colName := typeField.Tag.Get("db")
		value := valueField.Interface()
		v := fmt.Sprintf("%v", value)
		if colPrimary == "true" {
			where = colName + " ='" + v + "'"
		}
	}
	if where == "" {
		return nil, errors.New("could not find primary key in object " + tableName)
	}
	sql += where
	rows, err := db.Query(sql)
	if err != nil {
		return nil, err
	}
	values := PackageData(rows)
	if len(values) > 0 {
		return values[0], nil
	}
	return nil, nil
}

func GetByCondition(v interface{}, condition string, limit string, orderBy string) (interface{}, error) {
	db := getConnection()
	defer db.Close()
	tableName, err := GetStructTableName(v)
	if err != nil {
		return nil, err
	}
	sql := "SELECT * FROM " + tableName + " WHERE "
	if condition == "" {
		return nil, errors.New("could not find condition")
	}
	sql += condition
	if orderBy != "" {
		sql += " ORDER BY " + orderBy
	}
	if limit != "" {
		sql += " LIMIT " + limit + " OFFSET 0"
	}
	rows, err := db.Query(sql)
	if err != nil {
		return nil, err
	}
	values := PackageData(rows)
	if len(values) > 0 {
		return values[0], nil
	}
	return nil, nil
}

// chuyển dữ liệu từ map sang struct
func FillStruct(m map[string]interface{}, s interface{}) error {
	structValue := reflect.ValueOf(s).Elem()

	for name, value := range m {
		structFieldValue := structValue.FieldByName(name)

		if !structFieldValue.IsValid() {
			return fmt.Errorf("No such field: %s in obj", name)
		}

		if !structFieldValue.CanSet() {
			return fmt.Errorf("Cannot set %s field value", name)
		}

		val := reflect.ValueOf(value)
		if structFieldValue.Type() != val.Type() {
			return errors.New("Provided value type didn't match obj field type")
		}

		structFieldValue.Set(val)
	}
	return nil
}
