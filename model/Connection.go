package model

import (
	"database/sql"
	"fmt"

	"data-connector/library"
	"data-connector/log"

	_ "github.com/lib/pq"
)

/*
	create by: Hoangnd
	create at: 2021-08-07
	des: Cấu hình kết nối database
*/
func connectPostgreSQL() (*sql.DB, error) {
	config := library.GetDataBaseConfig()
	pgConfig := config["postgresql"].(map[string]interface{})
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", pgConfig["host"].(string), 5432, pgConfig["username"].(string), pgConfig["password"].(string), pgConfig["dbname"].(string))
	db, err := sql.Open("postgres", psqlInfo)
	return db, err
}
func ConnectSql() *sql.DB {
	db, err := connectPostgreSQL()
	if err != nil {
		log.Error(err.Error(), map[string]interface{}{
			"scope": log.Trace(),
		})
	}
	if err = db.Ping(); err != nil {
		log.Error(err.Error(), map[string]interface{}{
			"scope": log.Trace(),
		})
	} else {
		fmt.Println("DB Connected...")
	}
	return db
}
func getConnection() *sql.DB {
	return ConnectSql()
}
