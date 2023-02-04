package model

import (
	"fmt"
	"strconv"

	"data-connector/utils"
)

type History struct {
	Uuid         string `json:"uuid" db:"uuid" form:"uuid" type:"uuid" primary:"true"`
	ApiQueryUuid string `json:"api_query_uuid" db:"api_query_uuid" form:"apiQueryUuid" type:"uuid"`
	JobUuid      string `json:"job_uuid" db:"job_uuid" form:"jobUuid" type:"uuid"`
	ObjectId     string `json:"object_id" db:"object_id" form:"objectId" type:"string"`
	Url          string `json:"url" db:"url" form:"url" type:"string"`
	DocumentName string `json:"document_name" db:"document_name" form:"documentName" type:"string"`
	UserCreate   int    `json:"user_create" db:"user_create"  form:"userCreate" type:"number"`
	CreateAt     string `json:"create_at" db:"create_at" form:"createAt" type:"datetime"`
	UpdateAt     string `json:"update_at" db:"update_at" form:"updateAt" type:"datetime"`
}

func (history History) GetTableName() string {
	return "history"
}

type HistoryModelInterface interface {
	CreateHistory() error
	GetHistory(condition string) error
	InitByArray(data map[string]string)
}

/*
	create by: Hoangnd
	create at: 2021-08-07
	des: Hàm khởi tạo cronjob với đầu vào là 1 mảng
*/
func (history *History) InitByArray(data map[string]string) {
	history.Uuid = data["uuid"]
	history.ApiQueryUuid = data["api_query_uuid"]
	history.JobUuid = data["job_uuid"]
	history.ObjectId = data["object_id"]
	history.Url = data["url"]
	history.DocumentName = data["document_name"]
	userCreate, _ := strconv.Atoi(data["user_create"])
	history.UserCreate = userCreate
	history.CreateAt = data["create_at"]
	history.UpdateAt = data["update_at"]

}

/*
	create by: Hoangnd
	create at: 2021-08-07
	des: Hàm thực thi tạo mới history
*/
func (history History) CreateHistory(objectIds string, jobUuid string, apiQuery ApiQuery) error {
	uuid := CreateUuid()
	history.Uuid = uuid
	history.CreateAt = utils.GetCurrentTimeStamp()
	history.UpdateAt = utils.GetCurrentTimeStamp()
	mappings := apiQuery.GetMappingData()
	history.DocumentName = mappings.DocumentName
	history.ApiQueryUuid = apiQuery.Uuid
	history.Url = apiQuery.Url
	if jobUuid != "" {
		history.JobUuid = jobUuid
	}
	history.ObjectId = objectIds

	_, err := Insert(history)
	fmt.Println(err)
	return err
}

/*
	create by: Hoangnd
	create at: 2021-08-07
	des: Hàm lấy ra history theo condition
*/
func (history *History) GetHistory(condition string) error {
	data, err := GetByCondition(history, condition, "1", "update_at DESC")
	if data != nil {
		history.InitByArray(data.(map[string]string))
	}
	return err
}
