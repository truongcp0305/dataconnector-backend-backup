package router

import (
	"errors"
	"fmt"
	"net/http"

	"data-connector/controller"
	"data-connector/library"
	"data-connector/service"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func InitRoutes() {
	e := echo.New()
	basecontroller := &controller.Controller{IsRequireLogin: true}
	var baseControllerInterface controller.ControllerInterface = basecontroller
	apiqueryController := &controller.ApiqueryController{ControllerInterface: baseControllerInterface}
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	config := getMiddleWareConfig()
	e.Use(middleware.JWTWithConfig(config))
	e.Static("/", "static/")
	e.GET("/apiQueries", apiqueryController.GetListApiQuery)
	e.GET("/apiQueries/:uuid", apiqueryController.DetailApiQuery)
	e.POST("/apiQueries", apiqueryController.CreateApiQuery)
	e.POST("/apiQueries/extractData", apiqueryController.ExtractData)
	e.POST("/apiQueries/loadData", apiqueryController.LoadData)
	e.PUT("/apiQueries/:uuid", apiqueryController.EditApiQuery)
	e.DELETE("/apiQueries/:uuid", apiqueryController.DeleteApiQuery)
	e.POST("/apiQueries/executeJob", apiqueryController.ExecuteJob)
	e.GET("/apiQueries/progress/:uuid", apiqueryController.GetProgress)
	e.GET("/helthcheck", apiqueryController.HelthCheck)
	e.GET("/partner", apiqueryController.GetListPartner)
	e.Logger.Fatal(e.Start(":1323"))
}

/*
Hàm lấy thông tin token và kiểm tra
- hợp lệ
- Xác thực
*/
func getMiddleWareConfig() middleware.JWTConfig {
	return middleware.JWTConfig{
		ParseTokenFunc: func(token string, c echo.Context) (interface{}, error) {
			library.HEADER["Authorization"] = "Bearer " + token
			auth := new(service.Auth)
			auth.Token = token
			result := auth.VerifyJwt()
			if !result {
				c.JSON(http.StatusBadRequest, map[string]string{
					"status":  fmt.Sprint(http.StatusBadRequest),
					"message": "Token không chính xác",
				})
				return nil, errors.New("invalid token")
			}
			return auth, nil
		},
	}
}
