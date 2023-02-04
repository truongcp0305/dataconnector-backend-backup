package model

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"data-connector/library"
	"data-connector/log"
	"data-connector/utils"
)

type ApiQueryInterFace interface {
	InitByArray(data map[string]string)
	GetTableName() string
	GetListApiQuery() ([]map[string]string, error)
	CountData() string
	CreateApiQuery() (string, error)
	EditApiQuery() error
	UpdateStatus() error
	DeleteApiQuery() error
	StopApiQuery() error
	GetDetailApiQuery() (interface{}, error)
	ExtractData() interface{}
	GetMappingData() Mappings
	GetDataFromPath(data string) []map[string]interface{}
}

/*
create by: Hoangnd
create at: 2021-08-07
des: Khởi tạo struct chứa thông tin của apiquery và các cấu hình về thông tin của cột trong data base ( type, primary, db)
*/
type ApiQuery struct {
	Uuid             string `json:"uuid" db:"uuid" param:"uuid" form:"uuid" type:"uuid" primary:"true"`
	Url              string `json:"url" db:"url" form:"url" type:"string"`
	Name             string `json:"name" db:"name" form:"name" type:"string"`
	Note             string `json:"note" db:"note" form:"note" type:"string"`
	Type             int    `json:"type" db:"type" form:"type" type:"number"`
	Headers          string `json:"headers" db:"headers" form:"headers" type:"string"`
	Body             string `json:"body" db:"body" form:"body" type:"string"`
	Status           int    `json:"status" db:"status" form:"status" type:"string"`
	Params           string `json:"params" db:"params" form:"params" type:"string"`
	PathToData       string `json:"path_to_data" db:"path_to_data" form:"pathToData" type:"string"`
	PathToTotal      string `json:"path_to_total" db:"path_to_total" form:"pathToTotal" type:"string"`
	DeleteCondition  string `json:"delete_conditions" db:"delete_conditions" form:"deleteConditions" type:"string"`
	Method           string `json:"method" db:"method" form:"method" type:"string"`
	Mappings         string `json:"mappings" db:"mappings" form:"mappings" type:"string"`
	UpdateColumnInfo string `json:"update_column_info" db:"update_column_info" form:"updateColumnInfo" type:"string"`
	CronjobConfig    string `json:"cronjob_config" db:"cronjob_config" form:"cronjobConfig" type:"string"`
	UpdateBykey      string `json:"update_by_key" db:"update_by_key" form:"updateByKey" type:"string"`
	UserCreate       int    `json:"user_create" db:"user_create" form:"userCreate" type:"number"`
	CreateAt         string `json:"create_at" db:"create_at" form:"createAt" type:"datetime"`
	UpdateAt         string `json:"update_at" db:"update_at" form:"updateAt" type:"datetime"`
	LastRunAt        string `json:"last_run_at" db:"last_run_at" form:"lastRunAt" type:"datetime"`
	Partner          string `json:"partner" db:"partner" form:"partner" type:"string"`
}

/*
create by: Hoangnd
create at: 2021-08-07
des: struct lưu thông tin của phần mappings api vào doc
*/
type Mappings struct {
	DocumentName   string      `json:"documentName"`
	MappingColumns interface{} `json:"mappingColumns"`
	Data           interface{} `json:"data"`
	Total          int         `json:"total"`
}

func (apiQuery ApiQuery) GetTableName() string {
	return "api_query"
}

/*
create by: Hoangnd
create at: 2021-08-07
des: Hàm khởi tạo api với đầu vào là 1 mảng
*/
func (apiQuery *ApiQuery) InitByArray(data map[string]string) {
	apiQuery.Uuid = data["uuid"]
	apiQuery.Url = data["url"]
	apiQuery.Name = data["name"]
	apiQuery.Note = data["note"]
	typeA, _ := strconv.Atoi(data["type"])
	apiQuery.Type = typeA
	apiQuery.Headers = data["headers"]
	apiQuery.Body = data["body"]
	apiQuery.Params = data["params"]
	apiQuery.Status, _ = strconv.Atoi(data["status"])
	apiQuery.PathToData = data["path_to_data"]
	apiQuery.PathToTotal = data["path_to_total"]
	apiQuery.Method = data["method"]
	apiQuery.Mappings = data["mappings"]
	apiQuery.UpdateColumnInfo = data["update_column_info"]
	apiQuery.CronjobConfig = data["cronjob_config"]
	apiQuery.DeleteCondition = data["delete_conditions"]
	apiQuery.UpdateBykey = data["update_by_key"]
	userCreate, _ := strconv.Atoi(data["user_create"])
	apiQuery.UserCreate = userCreate
	apiQuery.CreateAt = data["create_at"]
	apiQuery.UpdateAt = data["update_at"]
	apiQuery.LastRunAt = data["last_run_at"]
	apiQuery.Partner = data["partner"]
}

/*
create by: Hoangnd
create at: 2021-08-07
des: Hàm lấy thông tin mapping của api đẩy vào struct mapping
*/
func (apiQuery ApiQuery) GetMappingData() Mappings {
	var sMapping Mappings
	err := json.Unmarshal([]byte(apiQuery.Mappings), &sMapping)
	if err != nil {
		log.Error(err.Error(), map[string]interface{}{
			"scope": log.Trace(),
			"res":   err,
		})
	}
	return sMapping
}

func (apiQuery ApiQuery) CountData() string {
	db := getConnection()
	defer db.Close()
	rows, err := db.Query("SELECT count(*) as count FROM api_query")
	if err != nil {
		return "0"
	}
	values := PackageData(rows)
	if len(values) > 0 {
		return values[0]["count"]
	}
	return "0"
}

/*
create by: Hoangnd
create at: 2021-08-07
des: Query lấy danh sách bản ghi
*/
func (apiQuery ApiQuery) GetListApiQuery() ([]map[string]string, error) {
	db := getConnection()
	defer db.Close()
	rows, err := db.Query("SELECT api_query.*,cronjob.status as status_job FROM api_query left join cronjob on api_query.uuid = cronjob.api_query_uuid order by create_at DESC")
	if err != nil {
		return nil, err
	}
	values := PackageData(rows)
	for _, item := range values {
		if item["type"] == "0" {
			item["type"] = "Thêm mới"
		} else if item["type"] == "1" {
			item["type"] = "Ghi đè"
		} else if item["type"] == "2" {
			item["type"] = "Cập nhật"
		}
		if item["status"] == "1" {
			item["status"] = "<font color=\"green\">Đang chạy</font>"
		} else {
			item["status"] = "<font color=\"red\">Đã dừng</font>"
		}
		if item["cronjob_config"] != "" && item["cronjob_config"] != "{}" {
			item["jobType"] = "<font color=\"orange\">Định kì</font>"
			if item["status_job"] == "1" {
				item["status_job"] = "<div style=\"width:12px;height:12px;border-radius:50%;background:green;\"></div>"
			} else if item["status_job"] == "0" {
				item["status_job"] = "<div style=\"width:12px;height:12px;border-radius:50%;background:yellow;\"></div>"
			}
		} else {
			item["jobType"] = "Thủ công"
			item["status_job"] = "<div style=\"width:12px;height:12px;border-radius:50%;background:gray;\"></div>"
		}

		if item["create_at"] != "" {
			item["create_at"] = utils.FormatDatabaseDateTime(item["create_at"])
		}
		if item["update_at"] != "" {
			item["update_at"] = utils.FormatDatabaseDateTime(item["update_at"])
		}
		if item["last_run_at"] != "" {
			item["last_run_at"] = utils.FormatDatabaseDateTime(item["last_run_at"])
		}
	}
	return values, nil
}

/*
create by: Hoangnd
create at: 2021-08-07
des: Hàm trả về chi tiết apiquery
*/
func (apiQuery ApiQuery) GetDetailApiQuery() (interface{}, error) {
	data, err := FindById(apiQuery)
	return data, err
}

/*
create by: Hoangnd
create at: 2021-08-07
des: Hàm thực thi tạo mới apiquery
*/
func (apiQuery ApiQuery) CreateApiQuery() (string, error) {
	uuid := CreateUuid()
	apiQuery.Uuid = uuid

	apiQuery.CreateAt = utils.GetCurrentTimeStamp()
	apiQuery.UpdateAt = utils.GetCurrentTimeStamp()
	returnValue, err := Insert(apiQuery)
	if err != nil {
		return "", err
	}
	return returnValue.(string), nil

}

/*
create by: Hoangnd
create at: 2021-08-07
des: Hàm thực thi sửa apiquery
*/
func (apiQuery ApiQuery) EditApiQuery() error {
	apiQuery.UpdateAt = utils.GetCurrentTimeStamp()
	return Update(apiQuery)
}

func (apiQuery ApiQuery) UpdateStatus(status string) error {
	if apiQuery.Uuid != "" {
		db := getConnection()
		defer db.Close()
		tableName := apiQuery.GetTableName()
		lastRunAt := ""
		if status == "1" {
			lastRunAt = ", last_run_at = '" + utils.GetCurrentTimeStamp() + "'"
		}
		_, err := db.Exec("Update " + tableName + " set status = " + status + lastRunAt + " where uuid = '" + apiQuery.Uuid + "'")
		return err
	}
	return errors.New("can not find api")

}

func ResetStatus() {
	db := getConnection()
	defer db.Close()
	db.Exec("Update api_query set status = 0")

}

/*
create by: Hoangnd
create at: 2021-08-07
des: Hàm thực thi xóa apiquery
*/
func (apiQuery ApiQuery) DeleteApiQuery() error {
	return Delete(apiQuery)
}

/*
create by: Hoangnd
create at: 2021-08-07
des: Hàm thực thi dừng apiquery đang chạy
*/
func (apiQuery ApiQuery) StopApiQuery() error {
	return nil
}

/*
create by: Hoangnd
create at: 2021-08-07
des: Hàm thực thi dừng apiquery đang chạy
*/
func (apiQuery ApiQuery) ExtractData(newHeader map[string]string, newBody map[string]string, url string) (interface{}, error) {
	req := new(library.Request)
	req.Url = url
	req.Header = newHeader
	req.Body = newBody
	req.SuppressParseData = true
	res, err := req.Send()
	if err == nil {
		data := apiQuery.GetDataFromPath(res.Data.(string))
		total := apiQuery.GetTotalFromPath(res.Data.(string))
		return map[string]interface{}{
			"data":  data,
			"total": total,
		}, nil
	}
	log.Error(err.Error(), map[string]interface{}{
		"req":   req,
		"err":   err,
		"scope": log.Trace(),
	})
	return map[string]interface{}{
		"data":  nil,
		"total": nil,
	}, err
}

/*
create by: Hoangnd
create at: 2021-08-07
des: Hàm lấy dữ liệu từ path
*/
func (apiQuery ApiQuery) GetDataFromPath(data string) interface{} {
	if apiQuery.PathToData != "" {
		var dataJson map[string]interface{}
		json.Unmarshal([]byte(data), &dataJson)
		if strings.Contains(apiQuery.PathToData, "/") {
			listPath := strings.Split(apiQuery.PathToData, "/")
			var newData interface{}
			for i := 0; i < len(listPath); i++ {
				if i == 0 {
					newData = dataJson[listPath[i]]
				} else {
					newData = newData.(map[string]interface{})[listPath[i]]
				}
			}
			return newData
		} else {
			return dataJson[apiQuery.PathToData]
		}

	} else {
		var dataJson []map[string]interface{}
		json.Unmarshal([]byte(data), &dataJson)
		return dataJson
	}
}

/*
create by: Hoangnd
create at: 2021-08-07
des: Hàm lấy tổng bản ghi từ path
*/
func (apiQuery ApiQuery) GetTotalFromPath(data string) interface{} {
	if apiQuery.PathToTotal != "" {
		var dataJson map[string]interface{}
		json.Unmarshal([]byte(data), &dataJson)
		if strings.Contains(apiQuery.PathToTotal, "/") {
			listPath := strings.Split(apiQuery.PathToTotal, "/")
			var newData interface{}
			for i := 0; i < len(listPath); i++ {
				if i == len(listPath)-1 {
					newData = newData.((map[string]interface{}))[listPath[i]]
				} else {
					newData = dataJson[listPath[i]].(map[string]interface{})
				}
			}
			return newData
		} else {
			return dataJson[apiQuery.PathToTotal]
		}
	}
	return 0
}
