package controller

/*
	create by: Hoangnd
	create at: 2021-08-07
	des: Controller chịu trách nhiệm xử lý các api về apiquery
*/

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"data-connector/library"
	"data-connector/log"
	"data-connector/model"
	"data-connector/service"
	"data-connector/utils"

	"github.com/labstack/echo/v4"
)

/*
create by: Hoangnd
create at: 2021-08-07
des: Khởi tạo interface chứa các phương thức của controller này
*/
type ApiqueryControllerInterface interface {
	GetListApiQuery(c echo.Context) error
	DetailApiQuery(c echo.Context) error
	CreateApiQuery(c echo.Context) error
	EditApiQuery(c echo.Context) error
	ExtractData(c echo.Context) error
	LoadData(c echo.Context) error
	DeleteApiQuery(c echo.Context) error
	ExecuteJob(c echo.Context) error
	GetProgress(c echo.Context) error
	GetListPartner(c echo.Context) error
}

type ApiqueryController struct {
	ControllerInterface ControllerInterface
}

/*
	create by: Hoangnd
	create at: 2021-08-07
	des: Xử lý lấy danh sách api
*/

func (apiqueryController ApiqueryController) ExecuteJob(c echo.Context) error {
	if apiqueryController.ControllerInterface.CheckParameter(c, []string{"uuid", "status"}) {
		apiQueryModel := new(model.ApiQuery)
		c.Bind(apiQueryModel)
		status, err := strconv.Atoi(c.FormValue("status"))
		data, err1 := apiQueryModel.GetDetailApiQuery()
		if err1 != nil {
			apiqueryController.ControllerInterface.Output(c, "", err)
			return err1
		}
		apiQueryModel.InitByArray(data.(map[string]string))
		if apiQueryModel.CronjobConfig != "" {
			job := new(model.Cronjob)
			if id, err := service.GetCurrentSupporterId(); err == nil {
				job.UserCreate = id
			}
			err = job.SaveCronjob(apiQueryModel, status)
			service.InitJob()
		} else {
			err = errors.New("require cronjob config to run job")
		}
		apiqueryController.ControllerInterface.Output(c, apiQueryModel.CronjobConfig, err)
	}
	return nil
}

/*
	create by: Hoangnd
	create at: 2021-08-07
	des: Xử lý lấy danh sách api
*/

func (apiqueryController ApiqueryController) GetListApiQuery(c echo.Context) error {
	apiQueryModel := new(model.ApiQuery)
	listData, err := apiQueryModel.GetListApiQuery()
	columns := getListColumnShowList()
	resData := make(map[string]interface{})
	resData["columns"] = columns
	resData["listObject"] = listData
	resData["total"] = apiQueryModel.CountData()
	apiqueryController.ControllerInterface.Output(c, resData, err)
	return err
}

/*
create by: truongnv
Hàm lấy danh sách partner
*/
func (apiqueryController ApiqueryController) GetListPartner(c echo.Context) error {
	partnerModel := new(model.Partner)
	listData, err := partnerModel.GetListPartner()
	if err != nil {
		fmt.Print("error roi")
		return err
	}
	apiqueryController.ControllerInterface.Output(c, listData, err)
	return err
}

/*
Hàm lấy danh sách cột trả về cho showList
*/
func getListColumnShowList() []interface{} {
	return []interface{}{
		map[string]interface{}{
			"name": "uuid", "title": "Uuid", "type": "text",
		},
		map[string]interface{}{
			"name": "name", "title": "Tên Api", "type": "text",
		},
		map[string]interface{}{
			"name": "url", "title": "Url", "type": "richtext",
		},
		map[string]interface{}{
			"name": "type", "title": "Kiểu", "type": "text",
		},
		map[string]interface{}{
			"name": "status", "title": "Trạng thái", "type": "richtext",
		},
		map[string]interface{}{
			"name": "status_job", "title": "Trạng thái định kỳ", "type": "richtext",
		},
		map[string]interface{}{
			"name": "node", "title": "Ghi chú", "type": "text",
		},
		map[string]interface{}{
			"name": "jobType", "title": "Kiểu chạy", "type": "richtext",
		},
		map[string]interface{}{
			"name": "last_run_at", "title": "Thời gian chạy lần cuối", "type": "text",
		},
		map[string]interface{}{
			"name": "create_at", "title": "Thời gian tạo", "type": "text",
		},
		map[string]interface{}{
			"name": "update_at", "title": "Thời gian cập nhật", "type": "text",
		},
	}
}

/*
create by: Hoangnd
create at: 2021-08-07
des: Xử lý thêm mới api
*/
func (apiqueryController ApiqueryController) CreateApiQuery(c echo.Context) error {
	if apiqueryController.ControllerInterface.CheckParameter(c, []string{"mappings"}) {
		apiQueryModel := new(model.ApiQuery)
		c.Bind(apiQueryModel)
		apiQueryId, err := apiQueryModel.CreateApiQuery()
		if id, err := service.GetCurrentSupporterId(); err == nil {
			apiQueryModel.UserCreate = id
		}
		apiqueryController.ControllerInterface.Output(c, apiQueryId, err)
	}
	return nil
}

/*
create by: Hoangnd
create at: 2021-08-07
des: Xử lý xem chi tiết bản ghi
*/
func (apiqueryController ApiqueryController) DetailApiQuery(c echo.Context) error {
	if apiqueryController.ControllerInterface.CheckParameter(c, []string{"uuid"}) {
		apiQueryModel := new(model.ApiQuery)
		c.Bind(apiQueryModel)
		data, err := apiQueryModel.GetDetailApiQuery()
		apiqueryController.ControllerInterface.Output(c, data, err)
	}
	return nil
}

/*
create by: Hoangnd
create at: 2021-08-07
des: Xử lý sửa bản ghi
*/
func (apiqueryController ApiqueryController) EditApiQuery(c echo.Context) error {
	if apiqueryController.ControllerInterface.CheckParameter(c, []string{"uuid"}) {
		apiQueryModel := new(model.ApiQuery)
		apiQueryModel.Uuid = c.Param("uuid")
		data, err := apiQueryModel.GetDetailApiQuery()
		if err != nil {
			apiqueryController.ControllerInterface.Output(c, "Not found", err)
			return nil
		} else {
			apiQueryModel.InitByArray(data.(map[string]string))
			c.Bind(apiQueryModel)
			err1 := apiQueryModel.EditApiQuery()
			service.InitJob()
			apiqueryController.ControllerInterface.Output(c, "Successful!", err1)
		}

	}
	return nil
}

/*
create by: Hoangnd
create at: 2021-08-07
des: Xử lý xóa bản ghi
*/
func (apiqueryController ApiqueryController) DeleteApiQuery(c echo.Context) error {
	if apiqueryController.ControllerInterface.CheckParameter(c, []string{"uuid"}) {
		apiQueryModel := new(model.ApiQuery)
		c.Bind(apiQueryModel)
		err := apiQueryModel.DeleteApiQuery()
		if err == nil {
			cronjob := new(model.Cronjob)
			cronjob.ApiQueryUuid = apiQueryModel.Uuid
			cronjob.DeleteCronjob()
		}
		apiqueryController.ControllerInterface.Output(c, "Successful!", err)
	}
	return nil
}

/*
create by: Hoangnd
create at: 2021-08-07
des: Xử lý lấy tiến trình chạy api
*/
func (apiqueryController ApiqueryController) GetProgress(c echo.Context) error {
	if apiqueryController.ControllerInterface.CheckParameter(c, []string{"uuid"}) {
		uuid := c.Param("uuid")
		apiqueryController.ControllerInterface.Output(c, progress[uuid], nil)
	}
	return nil
}
func (apiqueryController ApiqueryController) HelthCheck(c echo.Context) error {
	apiqueryController.ControllerInterface.Output(c, nil, nil)
	return nil
}

/*
create by: Hoangnd
create at: 2021-08-07
des: Xử lý lấy data từ 1 api nào đó
*/
func (apiqueryController ApiqueryController) ExtractData(c echo.Context) error {
	if apiqueryController.ControllerInterface.CheckParameter(c, []string{"url", "headers"}) {
		apiQueryModel := new(model.ApiQuery)
		c.Bind(apiQueryModel)
		newHeader, newBody, newParams := GetApiQueryParams(apiQueryModel.Headers, apiQueryModel.Body, apiQueryModel.Params)
		token := GetTokenMisa(apiQueryModel.Partner)
		if token != "" {
			newHeader["Authorization"] = fmt.Sprintf("Bearer %s", token)
		}
		if len(newParams) > 0 {
			queryParams := paramsToUrl(newParams)
			url := apiQueryModel.Url
			url = url + "?" + queryParams
			url = strings.ReplaceAll(url, "{page}", "1")
			apiQueryModel.Url = url
		}
		data, err := apiQueryModel.ExtractData(newHeader, newBody, apiQueryModel.Url)
		apiqueryController.ControllerInterface.Output(c, data, err)
	}
	return nil
}

/*
Create by truongnv

	Hàm gọi API để lấy token cho MISA
*/
func GetTokenMisa(partner string) string {
	type partnerModel struct {
		NamePartner  string `json:"namePartner"`
		ClientId     string `json:"clientId"`
		ClientSecret string `json:"clientSecret"`
	}
	type Response struct {
		Data string `json:"data"`
	}
	partnerJson := new(partnerModel)
	json.Unmarshal([]byte(partner), &partnerJson)
	fmt.Println(partnerJson.ClientId)
	if partnerJson.NamePartner == "Misa" {
		url := "https://crmconnect.misa.vn/api/v1/Account"
		bodyJson := fmt.Sprintf(`{"client_id":"%s", "client_secret":"%s"}`, partnerJson.ClientId, partnerJson.ClientSecret)
		jsonBody := []byte(bodyJson)
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
		if err != nil {
			fmt.Println("error in request token")
		}
		req.Header.Set("accept", "*/*")
		req.Header.Set("Content-Type", "application/json-patch+json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		fmt.Println("response Status:", resp.Status)
		body, _ := ioutil.ReadAll(resp.Body)
		var dataResp Response
		json.Unmarshal(body, &dataResp)
		fmt.Println("response Body:", dataResp.Data)
		return dataResp.Data
	}
	return ""
}

/*
create by: Hoangnd
create at: 2021-08-07
des: Hàm parse các thông tin header, params, body sang định dạng sử dụng cho request
*/
func GetApiQueryParams(header string, body string, params string) (map[string]string, map[string]string, []string) {
	headerJson := map[string]map[string]string{}
	json.Unmarshal([]byte(header), &headerJson)
	bodyJson := map[string]map[string]string{}
	json.Unmarshal([]byte(body), &bodyJson)
	paramsJson := map[string]map[string]string{}
	json.Unmarshal([]byte(params), &paramsJson)
	req := new(library.Request)
	req.Url = "https://" + library.GetPrefixEnvironment() + "syql.symper.vn/formulas/get-data"
	newHeader := map[string]string{}
	newBody := map[string]string{}
	newParams := []string{}
	req.Header = library.HEADER
	req.Method = "post"
	for key, value := range headerJson {
		if value["value"] != "" {
			newHeader[key] = value["value"]
		}
		if value["formula"] != "" {
			formula := value["formula"]
			req.Body = map[string]string{"formula": formula}
			res, err := req.Send()
			if err == nil {
				formulaValue := getDataSyqlFromResponse(res)
				newHeader[key] = formulaValue
			}
		}
	}
	for key, value := range bodyJson {
		if value["value"] != "" {
			newBody[key] = value["value"]
		}
		if value["formula"] != "" {
			formula := value["formula"]
			req.Body = map[string]string{"formula": formula}
			res, err := req.Send()
			if err == nil {
				formulaValue := getDataSyqlFromResponse(res)
				newBody[key] = formulaValue
			}
		}
	}
	for key, value := range paramsJson {
		if value["value"] != "" {
			// newParams[key] = value["value"]
			newParams = append(newParams, key+"="+value["value"])
		}
		if value["formula"] != "" {
			formula := value["formula"]
			req.Body = map[string]string{"formula": formula}
			res, err := req.Send()
			if err == nil {
				formulaValue := getDataSyqlFromResponse(res)
				newParams = append(newParams, key+"="+formulaValue)
			}
		}
	}
	return newHeader, newBody, newParams
}

/*
create by: Hoangnd
create at: 2021-08-07
des: Hàm lấy ra giá trị của công thức sau khi chạy
*/
func getDataSyqlFromResponse(res library.Response) string {
	if res.Status != 200 {
		return ""
	} else {
		dataSql := res.Data.(map[string]interface{})
		for _, value := range dataSql["data"].([]interface{})[0].(map[string]interface{}) {
			return value.(string)
		}
	}
	return ""
}

func paramsToUrl(params []string) string {
	return strings.Join(params[:], "&")
}

/*
create by: Hoangnd
create at: 2021-08-07
des: Xử lý lấy data từ 1 api nào đó
*/
func (apiqueryController ApiqueryController) LoadData(c echo.Context) error {
	if apiqueryController.ControllerInterface.CheckParameter(c, []string{"uuid"}) {
		apiQueryModel := new(model.ApiQuery)
		c.Bind(apiQueryModel)
		apiQuery, err := apiQueryModel.GetDetailApiQuery()
		if err != nil {
			progress[apiQueryModel.Uuid] = nil
			apiqueryController.ControllerInterface.Output(c, apiQuery, err)
		} else {
			progress[apiQueryModel.Uuid] = map[string]interface{}{
				"status":   "Chuẩn bị cho quá trình lấy dữ liệu từ API",
				"progress": 10,
			}
			apiQueryModel.UpdateStatus("1")
			err1 := PrepareToRun(apiQuery, apiQueryModel)
			if err1 != nil {
				progress[apiQueryModel.Uuid] = map[string]interface{}{
					"status":   "error",
					"progress": 100,
					"error":    err1,
				}
			}
			progress[apiQueryModel.Uuid] = map[string]interface{}{}
			apiQueryModel.UpdateStatus("0")
			apiqueryController.ControllerInterface.Output(c, "Successful", err1)
		}

	}
	return nil
}

var currentPrimaryObject = make(map[string]map[string][]string)

func PrepareToRun(apiQuery interface{}, apiQueryModel *model.ApiQuery) error {
	apiQueryjson := apiQuery.(map[string]string)
	if apiQueryjson["url"] != "" {

	} else {
		return errors.New("url not found")
	}
	apiQueryModel.InitByArray(apiQueryjson)
	mappings := apiQueryModel.GetMappingData()
	if mappings.DocumentName == "" {
		return errors.New("mappings data not valid")
	}
	if !utils.IsUrl(apiQueryModel.Url) {
		return errors.New("url not valid")
	}
	newHeader, newBody, newParams := GetApiQueryParams(apiQueryModel.Headers, apiQueryModel.Body, apiQueryModel.Params)

	fmt.Println(apiQueryModel.Partner)
	token := GetTokenMisa(apiQueryModel.Partner)
	if token != "" {
		newHeader["Authorization"] = fmt.Sprintf("Bearer %s", token)
	}
	if len(newParams) > 0 {
		sort.Strings(newParams)
		queryParams := paramsToUrl(newParams)
		url := apiQueryjson["url"]
		url = url + "?" + queryParams
		apiQueryModel.Url = url
	}
	errorCount[apiQueryModel.Uuid] = map[int]int{}
	jobUuid := ""
	e := checkValidApi(apiQueryModel)
	if e != nil {
		return e
	}
	history := new(model.History)
	documentName := mappings.DocumentName
	deleteCurrentData(documentName)
	history.GetHistory("url = '" + apiQueryModel.Url + "' AND document_name = '" + documentName + "'")
	if apiQueryModel.Type == 2 {
		progress[apiQueryModel.Uuid] = map[string]interface{}{
			"status":   "Chuẩn bị cho quá trình lấy dữ liệu từ API",
			"progress": 50,
		}
		// d, err := getCurrentDocumentObjectPrimaryId(apiQueryModel.DeleteCondition, &mappings, apiQueryModel.Env)
		// currentPrimaryObject[apiQueryModel.Uuid] = d
		isFirstRun := true
		// for k, v := range currentPrimaryObject[apiQueryModel.Uuid] {
		// 	fmt.Println(k, len(v))
		// 	if len(v) > 0 {
		// 		isFirstRun = false
		// 	}
		// 	log.Info("Check id từ sdocument", map[string]interface{}{
		// 		"docName":              k,
		// 		"currentPrimaryObject": len(v),
		// 		"err":                  err,
		// 		"apiName":              apiQueryModel.Name,
		// 		"trace":                log.Trace(),
		// 	})
		// }
		// if err != nil {
		// 	log.Error("Lỗi", map[string]interface{}{
		// 		"err":     err,
		// 		"apiName": apiQueryModel.Name,
		// 		"trace":   log.Trace(),
		// 	})
		// 	return err
		// }
		err1 := handleLoadDataForUpdate(apiQueryModel, history, newHeader, newBody, &mappings, isFirstRun)
		if err1 != nil {
			log.Error("Lỗi", map[string]interface{}{
				"err":     err1,
				"apiName": apiQueryModel.Name,
				"trace":   log.Trace(),
			})
			return err1
		}
		// currentPrimaryObject[apiQueryModel.Uuid] = nil
		// if !isFirstRun {
		// 	deleteChildRecordUnnecessary(&mappings, apiQueryModel.Env)
		// }
		history := new(model.History)
		history.CreateHistory("", jobUuid, *apiQueryModel)
	} else if apiQueryModel.Type == 0 {
		err1 := handleInsertApiType(apiQueryModel, newHeader, newBody, jobUuid)
		if err1 != nil {
			return err1
		}
		history := new(model.History)
		history.CreateHistory("", jobUuid, *apiQueryModel)
	} else {
		return errors.New("kiểu ghi đè hiện không còn được hỗ trợ")
	}
	return nil
}

func deleteCurrentData(docName string) {
	req := new(library.Request)
	req.Url = "https://" + library.GetPrefixEnvironment() + "sdocument-management.symper.vn/documents/objects"
	body := map[string]string{
		"type":         "all",
		"documentName": docName,
		"isTruncate":   "1",
	}
	req.Header = library.HEADER
	req.Body = body
	req.Method = "DELETE"
	req.Send()
}

/*
create by: Hoangnd
create at: 2021-08-07
Xử lý cho api kiểu thêm mới
*/
func handleInsertApiType(apiQueryModel *model.ApiQuery, newHeader map[string]string, newBody map[string]string, jobUuid string) error {
	chanel := make(chan interface{})
	totalPage, err := extractGetTotalPage(apiQueryModel, newHeader, newBody)
	if err != nil {
		return err
	}
	for i := 1; i <= totalPage+1; i++ {
		go handleLoadDataForInsert(apiQueryModel, newHeader, newBody, i, chanel)
	}
	allData := make([][]interface{}, 0)
	listData := make([]interface{}, 0)
	for i := 0; i <= totalPage; i++ {
		dataLoad := <-chanel
		progress[apiQueryModel.Uuid] = map[string]interface{}{
			"status":   "Đang trích xuất dữ liệu từ Api",
			"progress": (float64(i) / float64(totalPage)) * 100,
		}
		if dataLoad != nil {
			listData = append(listData, dataLoad.([]interface{})...)
		}
		if len(listData) > 3000 {
			allData = append(allData, listData)
			listData = listData[:0]
		}
	}
	if len(listData) > 0 {
		allData = append(allData, listData)
		listData = nil
	}
	for i := 0; i < len(allData); i++ {
		progress[apiQueryModel.Uuid] = map[string]interface{}{
			"status":   "Đang đẩy dữ liệu vào Document",
			"progress": (float64(i) / float64(len(allData))) * 100,
		}
		loadBatchData(allData[i], apiQueryModel, jobUuid, 0)
	}
	return nil
}

/*
create by: Hoangnd
create at: 2021-08-07
Phân chia luồng đẩy dữ liệu sang sang document
*/
func loadBatchData(dataLoad interface{}, apiQuery *model.ApiQuery, jobUuid string, index float64) {
	c := make(chan bool)
	data := dataLoad.([]interface{})
	l := len(data)
	if l == 300 {
		l = l - 1
	}
	for i := 0; i <= l/300; i++ {
		if i == l/300 {
			go transform(data[i*300:l], apiQuery, jobUuid, c)
		} else {
			go transform(data[i*300:(i+1)*300], apiQuery, jobUuid, c)
		}
	}
	for i := 0; i <= l/300; i++ {
		<-c
		progress[apiQuery.Uuid] = map[string]interface{}{
			"status":   "Đang đẩy dữ liệu vào Document",
			"progress": float64(i) + (index * 100),
		}
	}
	c = nil
}

func transform(dataLoad interface{}, apiQuery *model.ApiQuery, jobUuid string, c chan bool) {
	transformData(dataLoad, apiQuery, jobUuid)
	c <- true
}

/*
create by: Hoangnd
create at: 2021-08-07
Kiểm tra tính hợp lệ của api trước khi chạy
*/
func checkValidApi(apiQueryModel *model.ApiQuery) error {
	if apiQueryModel.PathToTotal == "" || apiQueryModel.PathToData == "" {
		return errors.New("cấu hình api không chính xác, vui lòng Kiểm tra lại đường dẫn data")
	}
	if apiQueryModel.Type == 2 {
		if apiQueryModel.UpdateBykey == "" {
			return errors.New("cấu hình api không chính xác, vui lòng Kiểm tra lại thông tin UpdateByKey")
		}
		if apiQueryModel.UpdateColumnInfo == "" {
			return errors.New("cấu hình api không chính xác, vui lòng Kiểm tra lại thông tin UpdateColumnInfo")
		}
	}

	return nil
}

var progress = make(map[string]map[string]interface{})

/*
create by: Hoangnd
create at: 2021-08-07
des: Xử lý api có type là cập nhật
*/
func handleLoadDataForUpdate(apiQuery *model.ApiQuery, history *model.History, newHeader map[string]string, newBody map[string]string, mappings *model.Mappings, isFirstRun bool) error {
	updateColumnInfo := apiQuery.UpdateColumnInfo
	var updateInfoJson map[string]interface{}
	err := json.Unmarshal([]byte(updateColumnInfo), &updateInfoJson)
	if err != nil {
		log.Error("Lỗi parse json", map[string]interface{}{
			"apiName": apiQuery.Name,
			"err":     err,
			"trace":   log.Trace(),
		})
		return err
	}
	condition := updateInfoJson["conditions"].(map[string]interface{})
	chanel := make(chan []interface{})
	totalPage, err1 := extractGetTotalPage(apiQuery, newHeader, newBody)
	if err1 != nil {
		log.Error("Lỗi lấy total page", map[string]interface{}{
			"apiName":   apiQuery.Name,
			"totalPage": totalPage,
			"err":       err,
			"trace":     log.Trace(),
		})
		return err1
	}
	allDataFromApi := make([]interface{}, 0)
	for i := 1; i <= totalPage+1; i++ {
		go queryApi(apiQuery, newHeader, newBody, i, chanel)
	}
	for i := 0; i <= totalPage; i++ {
		chanCallBack := <-chanel
		allDataFromApi = append(allDataFromApi, chanCallBack...)
		progress[apiQuery.Uuid] = map[string]interface{}{
			"status":   "Đang trích xuất dữ liệu từ API",
			"progress": (float64(i) / float64(totalPage)) * 100,
		}
	}
	chanel = nil
	if isFirstRun {
		l1 := len(allDataFromApi)
		progress[apiQuery.Uuid] = map[string]interface{}{
			"status":   "Đang đẩy dữ liệu vào document",
			"progress": 1,
		}
		log.Warn("log-data", map[string]interface{}{
			"len data": l1,
			"scope":    log.Trace(),
		})
		for i := 0; i <= l1/3000; i++ {
			if i == l1/3000 {
				loadBatchData(allDataFromApi[i*3000:l1], apiQuery, "", float64(i)/float64(l1/3000+1))
			} else {
				loadBatchData(allDataFromApi[i*3000:(i+1)*3000], apiQuery, "", float64(i)/float64(l1/3000+1))
			}
		}
		progress[apiQuery.Uuid] = map[string]interface{}{
			"status":   "Done",
			"progress": 100,
		}
		return nil
	} else {
		progress[apiQuery.Uuid] = map[string]interface{}{
			"status":   "Đang kiểm tra cập nhật dữ liệu",
			"progress": 10,
		}
		res := getSqlUpdate(allDataFromApi, apiQuery, history, condition)
		sql := res["listSql"].([]string)
		listDataInsert := res["listDataInsert"].([]interface{})
		listChildDataInsert := res["listChildDataInsert"].(map[string][]interface{})
		listRowIds := res["listRowId"].(map[string][]string)
		res = nil
		findAndDeleteRecordFromData(listRowIds, apiQuery.Uuid)
		l := len(listDataInsert)
		// trường hợp có các bản ghi thêm mới thì call sang insert bản ghi
		if l > 0 {
			for i := 0; i <= l/3000; i++ {
				if i == l/3000 {
					loadBatchData(listDataInsert[i*3000:l], apiQuery, "", float64(i)/float64(l/3000+1))
				} else {
					loadBatchData(listDataInsert[i*3000:(i+1)*3000], apiQuery, "", float64(i)/float64(l/3000+1))
				}
			}
		}
		progress[apiQuery.Uuid] = map[string]interface{}{
			"status":   "Đang cập nhật dữ liệu",
			"progress": 30,
		}
		l1 := len(sql)
		l2 := len(listChildDataInsert)
		log.Info("log", map[string]interface{}{
			"data insert":      l,
			"sql update":       l1,
			"len child insert": l2,
			"trace":            log.Trace(),
		})
		fmt.Println("data insert", l)
		fmt.Println("sql update", l1)
		fmt.Println("len child insert", l2)
		if l1 > 0 {
			executeSql(sql, false, "0")
		}
		progress[apiQuery.Uuid] = map[string]interface{}{
			"status":   "Đang cập nhật dữ liệu",
			"progress": 70,
		}
		if l2 > 0 {
			v, _ := json.Marshal(listChildDataInsert)
			queryInsertChildData(apiQuery, string(v))
		}
	}

	return nil
}

/*
create by: Hoangnd
create at: 2021-08-07
So khớp data hiện tại trong doc và data của api để tìm ra các bản ghi cần xóa đi
*/
func findAndDeleteRecordFromData(listRowIds map[string][]string, apiQueryUuid string) {
	mapDocumentWithKey := make(map[string]string)
	sTmp := make(map[string]map[string]struct{})
	for k, v := range listRowIds {
		item := strings.Split(k, ":")
		mb := make(map[string]struct{}, len(v))
		for _, x := range v {
			mb[x] = struct{}{}
		}
		sTmp[item[0]] = mb
		mapDocumentWithKey[item[0]] = item[1]
	}
	rowDelete := make(map[string][]string)
	for k, v := range currentPrimaryObject[apiQueryUuid] {
		mb := sTmp[k]
		for i := 0; i < len(v); i++ {
			if _, found := mb[v[i]]; !found {
				rowDelete[k] = append(rowDelete[k], v[i])
			}
		}
	}
	for k, v := range rowDelete {
		fmt.Println("xoa", len(v))
		if len(v) > 0 {
			ids := getConditionDeleteRecord(v)
			sql := []string{"DELETE FROM " + k + " WHERE " + mapDocumentWithKey[k] + " IN (" + ids + ") "}
			executeSql(sql, true, "0")
		}
	}

}
func getConditionDeleteRecord(rowId []string) string {
	r := make([]string, len(rowId))
	for i := 0; i < len(rowId); i++ {
		id := rowId[i]
		if strings.Contains(id, ":::::") {
			item := strings.Split(id, ":::::")
			if len(item) > 1 {
				id = item[1]
			}
		}
		r[i] = "'" + id + "'"
	}
	return strings.Join(r[:], ",")
}

/*
create by: Hoangnd
create at: 2021-08-07
des: Có dòng mới trong bảng con thì call api này để insert vào doc
*/
func queryInsertChildData(apiQuery *model.ApiQuery, childData string) error {
	req := new(library.Request)
	req.Url = "https://" + library.GetPrefixEnvironment() + "sdocument-management.symper.vn/documents/submit-childs"
	body := map[string]string{
		"values": childData,
	}
	req.Header = library.HEADER
	req.Body = body
	req.Method = "POST"
	_, err := req.Send()
	return err
}

/*
create by: Hoangnd
create at: 2021-08-07
des: Hàm lấy toàn bộ id của các bản ghi hiện tại đang có trong csdl, dựa vào rowId api cung cấp
*/
func getCurrentDocumentObjectPrimaryId(condition string, mappings *model.Mappings) (map[string][]string, error) {
	mappingColumns := mappings.MappingColumns
	dataMappings := getMappingColumnByKey(mappingColumns.([]interface{}))
	primaryKeys := dataMappings["primaryKey"]
	if len(primaryKeys) == 0 {
		return map[string][]string{}, errors.New("require primary key for update api")
	}
	dataReturn := make(map[string][]string)
	w := " WHERE " + condition
	if condition == "" {
		w = ""
	}
	for k, v := range primaryKeys {
		value := v.(map[string]interface{})
		if value["to"] == nil {
			return dataReturn, errors.New("require primary key for update api")
		}
		if k == "parent" {
			tablename := "document_" + mappings.DocumentName
			dataReturn[tablename] = getDocumentObjectData("document_"+mappings.DocumentName, value["to"].(string), "", w)
		} else {
			tablename := "document_child_" + mappings.DocumentName + "_" + k
			dataReturn[tablename] = getDocumentObjectData(mappings.DocumentName, value["to"].(string), primaryKeys["parent"].(map[string]interface{})["to"].(string), w)
		}
	}
	return dataReturn, nil
}

/*
create by: Hoangnd
create at: 2021-08-07
des: Xóa dữ liệu của các dòng trong bảng bị dư thừa
*/
func deleteChildRecordUnnecessary(mappings *model.Mappings) {
	mappingColumns := mappings.MappingColumns
	docName := mappings.DocumentName
	dataMappings := getMappingColumnByKey(mappingColumns.([]interface{}))
	primaryKeys := dataMappings["primaryKey"]
	sql := make([]string, 0)
	for k, _ := range primaryKeys {
		if k != "parent" {
			tablename := "document_child_" + mappings.DocumentName + "_" + k
			s := "DELETE FROM " + tablename + " WHERE document_object_parent_id NOT IN (SELECT document_object_id FROM document_" + docName + ")"
			sql = append(sql, s)
		}
	}
	executeSql(sql, false, "0")
}

/*
create by: Hoangnd
create at: 2021-08-07
des: Hàm lấy ra dữ liệu là các key unique của các bản ghi hiện có ( key unique là key định danh được cấu hình)
*/
func getDocumentObjectData(documentName, columnPrimary, parentPrimary, condition string) []string {
	countObject := getCountObject(documentName, columnPrimary, condition)
	allData := make([]string, 0)
	for j := 0; j <= countObject/20000; j++ {
		columnQuery := columnPrimary
		if parentPrimary != "" {
			columnQuery += "," + parentPrimary
		}
		offset := strconv.Itoa(j * 20000)
		sql := []string{"SELECT " + columnQuery + " from " + documentName + " " + condition + " ORDER BY " + columnPrimary + " LIMIT 20000 OFFSET " + offset}
		res, err := executeSql(sql, false, "1")
		if err != nil {
			return []string{}
		}
		dataSql := res.Data.(map[string]interface{})
		l := len(dataSql["data"].([]interface{}))
		for i := 0; i < l; i++ {
			item := dataSql["data"].([]interface{})[i]
			v := ""
			if item.(map[string]interface{})[columnPrimary] != nil {
				v = item.(map[string]interface{})[columnPrimary].(string)
			}
			if v != "" {
				if parentPrimary != "" && item.(map[string]interface{})[parentPrimary] != nil {
					v1 := v
					v2 := item.(map[string]interface{})[parentPrimary].(string)
					v = v2 + ":::::" + v1
					allData = append(allData, v)
				}
				if parentPrimary == "" {
					allData = append(allData, v)
				}
			}
		}
	}
	return allData
}

/*
create by: Hoangnd
create at: 2021-08-07
Hàm lấy tổng số bản ghi của doc hiện tại, để lấy ra thông tin cho việc call lấy dữ liệu theo page
*/
func getCountObject(tableName, columnPrimary, condition string) int {
	sql := []string{"SELECT count(" + columnPrimary + ") from " + tableName + condition}
	res, err := executeSql(sql, false, "0")
	if err != nil {
		return 0
	}
	dataSql := res.Data.(map[string]interface{})
	l := dataSql["data"].([]interface{})[0].(map[string]interface{})["count"]
	count, _ := strconv.Atoi(l.(string))
	return count
}

/*
create by: Hoangnd
create at: 2021-08-07
des: các luồng thực thi lấy về sql update và chạy sql ( trường hợp api update )
*/
func queryApi(apiQuery *model.ApiQuery, newHeader map[string]string,
	newBody map[string]string, currentPage int, c chan []interface{}) {
	url := apiQuery.Url
	if strings.Contains(url, "{page}") {
		url = strings.ReplaceAll(url, "{page}", strconv.Itoa(currentPage))
	}
	dataExtract, errData := getDataFromApi(apiQuery, newHeader, newBody, url)
	dataLoad := dataExtract["data"]
	dataExtract = nil
	if dataLoad != nil {
		c <- dataLoad.([]interface{})
		dataLoad = nil
	} else {
		errorCount[apiQuery.Uuid][currentPage]++
		log.Error("Lỗi không response từ api", map[string]interface{}{
			"currentPage": currentPage,
			"apiName":     apiQuery.Name,
			"errData":     errData,
			"trace":       log.Trace(),
		})
		if errorCount[apiQuery.Uuid][currentPage] < 10 {
			go queryApi(apiQuery, newHeader, newBody, currentPage, c)
		} else {
			c <- nil
		}
	}
}

/*
	create by: Hoangnd
	create at: 2021-08-07
	des: trả về mảng sql để update
*/

func getSqlUpdate(dataLoad interface{}, apiQuery *model.ApiQuery,
	history *model.History, condition map[string]interface{}) map[string]interface{} {
	mappings := apiQuery.GetMappingData()
	documentName := "document_" + mappings.DocumentName
	mappingColumns := mappings.MappingColumns
	dataMappings := getMappingColumnByKey(mappingColumns.([]interface{}))
	mapKeyToName := dataMappings["mapKeyToName"]
	mapNameToType := dataMappings["mapNameToType"]
	primaryKey := dataMappings["primaryKey"]
	listRowId := make(map[string][]string)
	dataToSql := dataToSql(dataLoad, mapKeyToName, mapNameToType, primaryKey, apiQuery, history, condition, documentName, "", listRowId, "")
	listSql := dataToSql["listSql"]
	listDataInsert := dataToSql["listDataInsert"]
	listChildDataInsert := dataToSql["listChildInsert"]
	dataToSql = nil
	return map[string]interface{}{
		"listRowId":           listRowId,
		"listSql":             listSql,
		"listChildDataInsert": listChildDataInsert,
		"listDataInsert":      listDataInsert,
	}
}

func getWithoutKey(condition map[string]interface{}) []string {
	data := make([]string, 0)
	for _, val := range condition {
		data = append(data, val.(string))
	}
	return data
}

/*
create by: Hoangnd
create at: 2021-08-07
des: Hàm lấy ra dữ liệu là các key unique của các row trả về từ api để check có thêm mới hay không
*/
func checkRowDeleted(listRow interface{}, documentName string, primaryKeys map[string]interface{}, inTable string, parentPrimaryValue string, apiUuid string) map[string]interface{} {
	rowInsert := make(map[string]struct{})
	listRowIdByRowKey := make([]string, 0)
	key := ""
	if inTable == "" {
		docPrimaryKey := primaryKeys["parent"].(map[string]interface{})
		key = docPrimaryKey["from"].(string)
	} else {
		tablePrimaryKey := primaryKeys[inTable].(map[string]interface{})
		key = tablePrimaryKey["from"].(string)
	}
	rawData := listRow.([]interface{})
	mb := make(map[string]struct{})
	for _, x := range currentPrimaryObject[apiUuid][documentName] {
		mb[x] = struct{}{}
	}
	for i := 0; i < len(rawData); i++ {
		rowData := rawData[i].(map[string]interface{})
		if rowData[key] != "" && rowData[key] != nil {
			str := fmt.Sprintf("%v", rowData[key])
			v := str
			if parentPrimaryValue != "" {
				v = parentPrimaryValue + ":::::" + rowData[key].(string)
			}
			listRowIdByRowKey = append(listRowIdByRowKey, v)
			if _, found := mb[v]; !found {
				rowInsert[v] = struct{}{}
			}
		}
	}
	return map[string]interface{}{
		"rowInsert":         rowInsert,
		"listRowIdByRowKey": listRowIdByRowKey,
	}

}

/*
create by: Hoangnd
create at: 2021-08-07
hàm lấy ra tên cột định danh của bản ghi từ api( dựa trên thông tin cấu hình )
*/
func getRowDataKeyDefinition(primaryKey map[string]interface{}, inTable string) map[string]interface{} {
	if inTable == "" {
		return primaryKey["parent"].(map[string]interface{})
	} else {
		return primaryKey[inTable].(map[string]interface{})
	}
}

func getChildRowDataByMapping(rowData map[string]interface{}, mapKeyToName map[string]interface{}) map[string]interface{} {
	newData := make(map[string]interface{})
	for k, v := range rowData {
		if mapKeyToName[k] != nil {
			newData[mapKeyToName[k].(string)] = v
		}
	}
	return newData

}

/*
	create by: Hoangnd
	create at: 2021-08-07
	des: Xử lý tìm các dòng đủ điều kiện và update trong data của api, với các trường hợp được update thì trả về 1 mảng các câu lệnh sql
	các trường hợp insert thì mảng data
*/

func dataToSql(dataLoad interface{}, mapKeyToName map[string]interface{},
	mapNameToType map[string]interface{}, primaryKey map[string]interface{}, apiQuery *model.ApiQuery, history *model.History,
	condition map[string]interface{}, documentName string, inTable string, listRowId map[string][]string,
	parentRowId string) map[string]interface{} {
	withoutKey := getWithoutKey(condition)
	listDataInsert := make([]interface{}, 0)
	listChildInsert := make(map[string][]interface{})
	listSql := make([]string, 0)
	updateByKey := apiQuery.UpdateBykey
	rawData := dataLoad.([]interface{})
	checkRow := checkRowDeleted(dataLoad, documentName, primaryKey, inTable, parentRowId, apiQuery.Uuid)
	rowInsert := checkRow["rowInsert"].(map[string]struct{})
	listRowIdByRowKey := checkRow["listRowIdByRowKey"].([]string)
	objectRowKeyDefi := getRowDataKeyDefinition(primaryKey, inTable)
	rowApiKey := objectRowKeyDefi["from"].(string)
	rowDocKey := objectRowKeyDefi["to"].(string)
	listRowId[documentName+":"+rowDocKey] = append(listRowId[documentName+":"+rowDocKey], listRowIdByRowKey...)
	for i := 0; i < len(rawData); i++ {
		rowData := rawData[i].(map[string]interface{})
		isRowUpdate := false
		vSet := make([]string, 0)
		vWhere := make([]string, 0)
		s := fmt.Sprintf("%v", rowData[rowApiKey])
		rowId := s
		if inTable != "" {
			rowId = parentRowId + ":::::" + rowId
		}
		if _, found := rowInsert[rowId]; found { // nếu là bản ghi mới thì đưa ra list để insert
			newRow := rowData
			if inTable != "" {
				newRow = getChildRowDataByMapping(newRow, mapKeyToName)
			}
			listDataInsert = append(listDataInsert, newRow)
			newRow = nil
			continue
		}
		isRowNotUpdate := false
		if updateByKey != "" && rowData[updateByKey] == nil {
			log.Error("Bản ghi ko có dữ liệu", map[string]interface{}{
				"updateByKey": updateByKey,
				"rowData":     rowData,
				"trace":       log.Trace(),
			})
			continue
		}
		if updateByKey != "" && rowData[updateByKey].(string) != "" && !checkUpdateItem(history, apiQuery, rowData[updateByKey].(string)) {
			isRowNotUpdate = true
		}
		for key, value := range rowData {
			if mapKeyToName[key] != "" && mapKeyToName[key] != nil && !utils.StringInSlice(key, withoutKey) {
				switch mapKeyToName[key].(type) {
				case map[string]interface{}:
					tableMap := mapKeyToName[key].(map[string]interface{})
					if mapNameToType[tableMap["tablename"].(string)] == "table" {
						documentName1 := documentName
						childTableName := strings.ReplaceAll(documentName1, "document_", "")
						childDocName := "document_child_" + childTableName + "_" + tableMap["tablename"].(string)
						listSqlChildData := dataToSql(value, tableMap["columns"].(map[string]interface{}),
							mapNameToType, primaryKey, apiQuery, history, condition, childDocName,
							tableMap["tablename"].(string), listRowId, rowId)
						listSqlChild := listSqlChildData["listSql"].([]string)
						listDataInsertChild := listSqlChildData["listDataInsert"].([]interface{})
						if len(listDataInsertChild) > 0 {
							childRowKey := listSqlChildData["childRowKey"]
							k := childTableName + ":" + tableMap["tablename"].(string) + ":" + rowDocKey + ":" + childRowKey.(string) + ":" + rowId
							listChildInsert[k] = append(listChildInsert[k], listDataInsertChild...)
							childRowKey = nil
						}
						if len(listSqlChild) > 0 {
							listSql = append(listSql, listSqlChild...)
						}
						listSqlChildData = nil
						listSqlChild = nil
						listDataInsertChild = nil
					}
				default:
					if !isRowNotUpdate {
						if mapNameToType[mapKeyToName[key].(string)] == "fileUpload" {
							switch value.(type) {
							case []interface{}:
								listFile := make([]map[string]interface{}, 0)
								for j := 0; j < len(value.([]interface{})); j++ {
									link := value.([]interface{})[j].(string)
									fileName := getFileNameFromLink(link)
									v := map[string]interface{}{
										"id":         0,
										"uid":        "",
										"name":       fileName,
										"type":       "link",
										"serverPath": link,
										"size":       0,
									}
									listFile = append(listFile, v)
								}
								listFileJsonStr, _ := json.Marshal(listFile)
								v1 := mapKeyToName[key].(string) + " = '" + string(listFileJsonStr) + "'"
								vSet = append(vSet, v1)
							}
						} else {
							v1 := mapKeyToName[key].(string) + " = '" + fmt.Sprint(value) + "'"
							vSet = append(vSet, v1)
						}
					}
				}
			}
		}
		if isRowNotUpdate {
			continue
		}
		if !isRowUpdate {
			for key, value := range condition {
				st := fmt.Sprintf("%v", rowData[value.(string)])
				if st != "" {
					w := key + " = " + "'" + st + "'"
					vWhere = append(vWhere, w)
				}
			}
			whereStr := ""
			if len(vWhere) > 0 {
				whereStr = " WHERE " + strings.Join(vWhere[:], " AND ")
			}
			if len(vSet) > 0 {
				sql := "UPDATE " + documentName + " SET " + strings.Join(vSet[:], ",") + whereStr
				if !utils.StringInSlice(sql, listSql) {
					listSql = append(listSql, sql)
				}
			}
		}
	}
	return map[string]interface{}{
		"listSql":         listSql,
		"listDataInsert":  listDataInsert,
		"listChildInsert": listChildInsert,
		"childRowKey":     rowDocKey,
	}
}

/*
create by: Hoangnd
create at: 2021-08-07
des: Kiểm tra xem row mới có phải là row được update từ khách hàng sau lần gọi cuối cùng
*/
func checkUpdateItem(history *model.History, apiQuery *model.ApiQuery, lastUpdateTime string) bool {
	lastUpdateMili := utils.GetTimeMiliseconds(lastUpdateTime)
	if history.UpdateAt == "" {
		return false
	}
	lastHistoryMili := utils.GetTimeMiliseconds(utils.FormatDatabaseDateTime(history.UpdateAt))
	return lastUpdateMili > lastHistoryMili
}

/*
create by: Hoangnd
create at: 2021-08-07
des: query câu lệnh sql update data
*/
func executeSql(sql []string, SuppressParseData bool, distinct string) (library.Response, error) {
	if len(sql) > 0 {
		sql := strings.Join(sql[:], ";")
		dataPostToService := map[string]string{
			"formula":  sql,
			"distinct": distinct,
		}
		req := new(library.Request)
		req.Url = "https://" + library.GetPrefixEnvironment() + "syql.symper.vn/formulas/get-data"
		req.Header = library.HEADER
		req.Body = dataPostToService
		req.Method = "POST"
		req.SuppressParseData = SuppressParseData
		res, err := req.Send()
		return res, err
	}
	return library.Response{}, nil
}

/*
Query tạm để lấy thông tin về pageCount
*/
func extractGetTotalPage(apiQuery *model.ApiQuery, newHeader map[string]string, newBody map[string]string) (int, error) {
	url := apiQuery.Url
	if strings.Contains(url, "{page}") {
		url = strings.ReplaceAll(url, "{page}", "1")
	}
	dataExtract, err := getDataFromApi(apiQuery, newHeader, newBody, url)
	fmt.Println(apiQuery, newHeader, newBody, url)
	dataLoad := dataExtract["data"]
	total := dataExtract["total"]
	fmt.Println("LOG CHECK LOAD API")
	fmt.Println("dataLoad: ", dataLoad)
	fmt.Println("total: ", total)
	fmt.Println("err: ", err)
	if dataLoad == nil || total == nil || err != nil {
		return 0, errors.New("can not get data from " + url)
	} else {
		totalPage := int(total.(float64) / float64(len(dataLoad.([]interface{}))))
		return totalPage, nil
	}
}

var errorCount = map[string]map[int]int{}

/*
create by: Hoangnd
create at: 2021-08-07
des: Xử lý phân chia luồng call api
*/
func handleLoadDataForInsert(apiQuery *model.ApiQuery, newHeader map[string]string, newBody map[string]string, currentPage int, c chan interface{}) {
	url := apiQuery.Url
	if strings.Contains(url, "{page}") {
		url = strings.ReplaceAll(url, "{page}", strconv.Itoa(currentPage))
	}
	dataExtract, errData := getDataFromApi(apiQuery, newHeader, newBody, url)
	dataLoad := dataExtract["data"]
	dataExtract = nil
	if dataLoad != nil {
		c <- dataLoad
		dataLoad = nil
	} else {
		errorCount[apiQuery.Uuid][currentPage]++
		log.Error("Lỗi không response từ api", map[string]interface{}{
			"currentPage": currentPage,
			"apiName":     apiQuery.Name,
			"errData":     errData,
			"trace":       log.Trace(),
		})
		if errorCount[apiQuery.Uuid][currentPage] < 10 {
			go handleLoadDataForInsert(apiQuery, newHeader, newBody, currentPage, c)
		} else {
			c <- nil
		}
	}
}

/*
create by: Hoangnd
create at: 2021-08-07
des: Lấy dữ liệu dạng raw từ api
*/
func getDataFromApi(apiQuery *model.ApiQuery, newHeader map[string]string, newBody map[string]string, url string) (map[string]interface{}, error) {
	data, err := apiQuery.ExtractData(newHeader, newBody, url)
	return data.(map[string]interface{}), err
}

/*
create by: Hoangnd
create at: 2021-08-07
des: Hàm thực hiện đẩy dữ liệu vào doc, trường hợp ghi đè thì có xóa các bản ghi cũ
*/
func transformData(dataLoad interface{}, apiQuery *model.ApiQuery, jobUuid string) (library.Response, error) {
	mappings := apiQuery.GetMappingData()
	documentName := mappings.DocumentName
	mappingColumns := mappings.MappingColumns
	dataMappings := getMappingColumnByKey(mappingColumns.([]interface{}))
	mapKeyToName := dataMappings["mapKeyToName"]
	mapNameToType := dataMappings["mapNameToType"]
	listChild := make([]string, 0)
	dataToLoad := getDataByMapping(dataLoad, mapKeyToName, mapNameToType, &listChild)
	headers := library.HEADER
	if apiQuery.Type == 1 {
		if jobUuid != "" {
			history := new(model.History)
			history.GetHistory("api_query_uuid = '" + apiQuery.Uuid + "' AND job_uuid = '" + jobUuid + "'")
			if history.Uuid != "" {
				listChildStr, _ := json.Marshal(listChild)
				dataPostToService := map[string]string{
					"ids":           history.ObjectId,
					"childrenNames": string(listChildStr),
				}
				req := new(library.Request)
				req.Url = "https://" + library.GetPrefixEnvironment() + "sdocument-management.symper.vn/documents/" + documentName + "/objects"
				req.Header = headers
				req.Body = dataPostToService
				req.Method = "DELETE"
				response, err := req.Send()
				log.Error(err.Error(), map[string]interface{}{
					"data":  dataPostToService,
					"res":   response,
					"scope": log.Trace(),
				})
			}
		}
	}
	req := new(library.Request)
	req.Url = "https://" + library.GetPrefixEnvironment() + "sdocument-management.symper.vn/documents/loadDataToDocument"
	dataToLoadStr, _ := json.Marshal(dataToLoad)
	body := map[string]string{
		"documentName": documentName, "data": string(dataToLoadStr),
	}
	req.Header = headers
	req.Body = body
	req.Method = "POST"
	dataMappings = nil
	dataToLoad = nil
	res, err := req.Send()
	if err != nil {
		log.Warn(err.Error(), map[string]interface{}{
			"err":   err,
			"scope": log.Trace(),
		})
	}
	return res, err

}

/*
create by: Hoangnd
create at: 2021-08-07
des: Hàm lấy thông tin mapping api từ cấu hình
*/
func getMappingColumnByKey(mappingColumn []interface{}) map[string]map[string]interface{} {
	mapKeyToName := make(map[string]interface{})
	mapNameToType := make(map[string]interface{})
	primaryKey := make(map[string]interface{})
	for i := 0; i < len(mappingColumn); i++ {
		mappingItem := mappingColumn[i].(map[string]interface{})
		mapKeyToName[mappingItem["from"].(string)] = mappingItem["to"]
		mapNameToType[mappingItem["to"].(string)] = mappingItem["type"]
		if mappingItem["primary"] != "" && mappingItem["primary"] != nil && mappingItem["primary"] == true {
			primaryKey["parent"] = mappingItem
		}
		if mappingItem["type"].(string) == "table" {
			columnsStr := mappingItem["columns"]
			dataMapTable := getMappingColumnByKey(columnsStr.([]interface{}))
			mapKeyToName[mappingItem["from"].(string)] = map[string]interface{}{
				"tablename": mappingItem["to"],
				"columns":   dataMapTable["mapKeyToName"],
			}
			mapNameToTypeTable := dataMapTable["mapNameToType"]
			primaryKeyTable := dataMapTable["primaryKey"]
			if primaryKeyTable["parent"] != nil {
				primaryKey[mappingItem["to"].(string)] = primaryKeyTable["parent"]
			}
			for k, v := range mapNameToTypeTable {
				mapNameToType[k] = v
			}
		}
	}
	return map[string]map[string]interface{}{
		"mapKeyToName":  mapKeyToName,
		"mapNameToType": mapNameToType,
		"primaryKey":    primaryKey,
	}
}

/*
	create by: Hoangnd
	create at: 2021-08-07
	des: Convert data từ api sang dạng theo mapping để insert vào doc
*/

func getDataByMapping(dataLoad interface{}, mapKeyToName map[string]interface{}, mapNameToType map[string]interface{}, listChild *[]string) []map[string]interface{} {
	dataToLoad := make([]map[string]interface{}, 0)
	rawData := dataLoad.([]interface{})
	for i := 0; i < len(rawData); i++ {
		rowData := rawData[i].(map[string]interface{})
		rowToLoad := make(map[string]interface{})
		for keyMap, valueMap := range mapKeyToName {
			var value interface{} = ""
			if rowData[keyMap] != nil && rowData[keyMap] != "" {
				value = rowData[keyMap]
			}
			switch valueMap.(type) {
			case map[string]interface{}:
				tableMap := valueMap.(map[string]interface{})
				if mapNameToType[tableMap["tablename"].(string)] == "table" {
					if !utils.StringInSlice(tableMap["tablename"].(string), *listChild) {
						*listChild = append(*listChild, tableMap["tablename"].(string))
					}
					rowToLoad[tableMap["tablename"].(string)] = getDataByMapping(value, tableMap["columns"].(map[string]interface{}), mapNameToType, listChild)
				}
			default:
				if mapNameToType[valueMap.(string)] == "fileUpload" {
					switch value.(type) {
					case []interface{}:
						listFile := make([]map[string]interface{}, 0)
						for j := 0; j < len(value.([]interface{})); j++ {
							link := value.([]interface{})[j].(string)
							fileName := getFileNameFromLink(link)
							v := map[string]interface{}{
								"id":         0,
								"uid":        "",
								"name":       fileName,
								"type":       "link",
								"serverPath": link,
								"size":       0,
							}
							listFile = append(listFile, v)
						}
						listFileJsonStr, _ := json.Marshal(listFile)
						rowToLoad[valueMap.(string)] = string(listFileJsonStr)
					}
				} else if mapNameToType[valueMap.(string)] == "date" || mapNameToType[valueMap.(string)] == "dateTime" {
					dt := utils.FormatDatabaseDateTimeWithUtc(value.(string))
					rowToLoad[valueMap.(string)] = dt
				} else {
					rowToLoad[valueMap.(string)] = value
				}

			}

		}
		dataToLoad = append(dataToLoad, rowToLoad)
	}
	return dataToLoad

}

/*
create by: Hoangnd
create at: 2021-08-07
des: hàm lấy tên file từ link
*/
func getFileNameFromLink(link string) string {
	re := regexp.MustCompile(`(\w+\.)+\w+$`)
	match := re.FindStringSubmatch(link)
	if len(match) > 0 {
		return match[0]
	} else {
		return link
	}
}

/*
create by: Hoangnd
create at: 2021-08-07
des: Lấy thông tin object id insert được vào doc
*/
func getObjectIdsFromResponse(data interface{}) []string {
	if data != nil {
		arrayData := data.([]interface{})
		listIds := make([]string, len(arrayData))
		for i := 0; i < len(arrayData); i++ {
			rowData := arrayData[i].(map[string]interface{})
			if rowData["result"].(bool) {
				listIds = append(listIds, rowData["document_object_id"].(string))
			}
		}
		return listIds
	}
	return []string{}

}
