package model

import (
	"encoding/json"
	"errors"
	"strconv"

	"data-connector/utils"
)

type Cronjob struct {
	Uuid         string `json:"uuid" db:"uuid" form:"uuid" type:"uuid" primary:"true"`
	ApiQueryUuid string `json:"api_query_uuid" db:"api_query_uuid" form:"apiQueryUuid" type:"uuid"`
	Status       int    `json:"status" db:"status" form:"status" type:"number"`
	Config       string `json:"config" db:"config" form:"config" type:"string"`
	UserCreate   int    `json:"user_create" db:"user_create"  form:"userCreate" type:"number"`
	CreateAt     string `json:"create_at" db:"create_at" form:"createAt" type:"datetime"`
	UpdateAt     string `json:"update_at" db:"update_at" form:"updateAt" type:"datetime"`
}
type CronjobModelInterface interface {
	CreateCronjob() error
	ToJson() error
	SaveCronjob() string
	GetCronJob() error
	InitByArray()
	GetListCronjob() ([]map[string]string, error)
	DeleteCronjob() error
}

func (cronjob Cronjob) GetTableName() string {
	return "cronjob"
}

/*
	create by: Hoangnd
	create at: 2021-08-07
	des: Hàm khởi tạo cronjob với đầu vào là 1 mảng
*/
func (cronjob *Cronjob) InitByArray(data map[string]string) {
	cronjob.Uuid = data["uuid"]
	statusInt, _ := strconv.Atoi(data["status"])
	cronjob.ApiQueryUuid = data["api_query_uuid"]
	cronjob.Status = statusInt
	cronjob.Config = data["config"]
	userCreate, _ := strconv.Atoi(data["user_create"])
	cronjob.UserCreate = userCreate
	cronjob.CreateAt = data["create_at"]
	cronjob.UpdateAt = data["update_at"]

}
func (cronjob Cronjob) ToJson() string {
	data, _ := json.Marshal(cronjob)
	return string(data)
}

/*
	create by: Hoangnd
	create at: 2021-08-07
	des: Hàm insert cronjob vào db
*/
func (cronjob *Cronjob) CreateCronjob() error {
	uuid := CreateUuid()
	cronjob.Uuid = uuid

	cronjob.CreateAt = utils.GetCurrentTimeStamp()
	cronjob.UpdateAt = utils.GetCurrentTimeStamp()
	_, err := Insert(cronjob)
	return err
}

/*
	create by: Hoangnd
	create at: 2021-08-07
	des: Kiểm tra có cấu hình cronjob thì run cronjob đó ( với status = 1)
*/
func (cronjob Cronjob) SaveCronjob(apiQuery *ApiQuery, status int) error {
	data, err := GetByCondition(cronjob, " api_query_uuid = '"+apiQuery.Uuid+"'", "", "update_at DESC")
	config := apiQuery.CronjobConfig
	if config == "" {
		return errors.New("require cronjob config to run job")
	} else {
		configMap := make(map[string]string)
		json.Unmarshal([]byte(config), &configMap)
		config = configMap["cronjob"]
		if data != nil && err == nil {
			cronjob.InitByArray(data.(map[string]string))
			cronjob.Status = status
			cronjob.Config = config
			err := Update(cronjob)
			return err
		} else {
			cronjob.Config = config
			cronjob.Status = status
			cronjob.ApiQueryUuid = apiQuery.Uuid
			err := cronjob.CreateCronjob()
			return err
		}
	}
}

/*
	create by: Hoangnd
	create at: 2021-08-07
	des: Hàm lấy ra history theo condition
*/
func (cronjob *Cronjob) GetCronJob(condition string) error {
	data, err := GetByCondition(cronjob, condition, "", "update_at DESC")
	if data != nil {
		cronjob.InitByArray(data.(map[string]string))
	}
	return err
}

/*
	create by: Hoangnd
	create at: 2021-08-07
	des: Query lấy danh sách bản ghi
*/
func (cronjob Cronjob) GetListCronjob(status string) ([]map[string]string, error) {
	db := getConnection()
	defer db.Close()
	rows, err := db.Query("SELECT * FROM cronjob WHERE status = " + status)
	if err != nil {
		return nil, err
	}
	values := PackageData(rows)
	return values, nil
}

/*
	create by: Hoangnd
	create at: 2021-08-07
	des: Query lấy danh sách bản ghi
*/
func (cronjob Cronjob) DeleteCronjob() error {
	err := DeleteByCondition(cronjob, "api_query_uuid = '"+cronjob.ApiQueryUuid+"'")
	return err
}
