package controller

import (
	"data-connector/utils"

	"github.com/labstack/echo/v4"
)

type ControllerInterface interface {
	Output(c echo.Context, data interface{}, err error)
	CheckParameter(c echo.Context, params []string) bool
}
type Controller struct {
	IsRequireLogin bool
}

/*
	create by: Hoangnd
	create at: 2021-08-07
	des: Xử lý output cho api
*/
func (controller Controller) Output(c echo.Context, data interface{}, err error) {
	dataResponse := make(map[string]interface{})
	dataResponse["data"] = data
	if err != nil {
		dataResponse["data"] = err.Error()
		dataResponse["status"] = utils.STATUS_BAD_REQUEST
		dataResponse["message"] = utils.STORE_STATUS[utils.STATUS_BAD_REQUEST]
	} else {
		dataResponse["status"] = utils.STATUS_OK
		dataResponse["message"] = utils.STORE_STATUS[utils.STATUS_OK]
	}
	c.JSON(dataResponse["status"].(int), dataResponse)
}

/*
	create by: Hoangnd
	create at: 2021-08-07
	des: Kiểm tra các tham số đầu vào
*/
func (controller Controller) CheckParameter(c echo.Context, params []string) bool {
	// for i := 0; i < len(params); i++ {
	// 	if c.FormValue(params[i]) == "" && c.Param(params[i]) == "" {
	// 		controller.Output(c, "", errors.New("tham số không chính xác"))
	// 		return false
	// 	}
	// 	fmt.Println("hoangnnn")
	// 	fmt.Println(c.FormValue("hoang"))
	// }
	return true
}
