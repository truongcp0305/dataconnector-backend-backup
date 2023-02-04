package log

import (
	"os"
	"runtime"

	"github.com/sirupsen/logrus"
)

/*
	create by: Hoangnd
	create at: 2021-08-07
	des: service xử lý log cho framework
*/
type SymperLog struct {
	Message string
	Data    map[string]interface{}
}
type SymperLogInterface interface {
	sWarn()
	sInfo()
	sError()
}
type STrace struct {
	File string
	Line int
	Func string
}

/*
	create by: Hoangnd
	create at: 2021-08-07
	des: đẩy log vào file
*/
func OpenFile(t string) *os.File {
	filePath := "./log/symper_log.log"
	if t == "error" {
		filePath = "./log/symper_error_log.log"
	}
	if t == "warn" {
		filePath = "./log/symper_warn_log.log"
	}
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		logrus.Println(err)
	}
	logrus.SetOutput(file)
	logrus.SetFormatter(&logrus.JSONFormatter{})
	return file
}

func (log SymperLog) sWarn() {
	file := OpenFile("warn")
	defer file.Close()
	logrus.WithFields(log.Data).Warn(log.Message)
}
func (log SymperLog) sInfo() {
	file := OpenFile("info")
	defer file.Close()
	logrus.WithFields(log.Data).Info(log.Message)
}
func (log SymperLog) sError() {
	file := OpenFile("error")
	defer file.Close()
	logrus.WithFields(log.Data).Error(log.Message)
}
func Warn(message string, data map[string]interface{}) {
	log := new(SymperLog)
	log.Message = message
	log.Data = data
	log.sWarn()
}
func Info(message string, data map[string]interface{}) {
	log := new(SymperLog)
	log.Message = message
	log.Data = data
	log.sInfo()
}
func Error(message string, data map[string]interface{}) {
	log := new(SymperLog)
	log.Message = message
	log.Data = data
	log.sError()
}
func Trace() *STrace {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	t := new(STrace)
	t.File = frame.File
	t.Func = frame.Function
	t.Line = frame.Line
	return t
}
