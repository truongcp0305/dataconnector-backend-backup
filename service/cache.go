package service

import "sync"

var DataCache = map[string]interface{}{}
var mutex = &sync.Mutex{}

func GetCache(key string) interface{} {
	return DataCache[key]
}

func SetCache(key string, data interface{}) {
	mutex.Lock()
	DataCache[key] = data
	mutex.Unlock()
}
