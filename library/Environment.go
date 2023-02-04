package library

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var Env struct {
	Environment string                 `json:"environment"`
	Db          map[string]interface{} `json:"db"`
}

func InitEnvironment() error {
	configFile, err := ioutil.ReadFile("./env.json")
	if err != nil {
		fmt.Println("opening config file", err.Error())
		return err
	}
	err1 := json.Unmarshal([]byte(string(configFile)), &Env)
	return err1
}

func GetPrefixEnvironment() string {
	if Env.Environment != "" {
		return Env.Environment + "-"
	} else {
		return ""
	}
}
func GetEnv() string {
	return Env.Environment
}
func GetDataBaseConfig() map[string]interface{} {
	return Env.Db
}
