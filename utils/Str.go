package utils

import (
	"fmt"
	"net/url"
	"runtime"
	"strconv"
	"strings"
	"time"

	"data-connector/log"
)

func GetCurrentTimeStamp() string {
	currentTime := time.Now()
	return currentTime.Format("2006-01-02 15:04:05")

}
func FormatDatabaseDateTime(dateTime string) string {
	var dt string
	t, err := time.Parse("2006-01-02T15:04:05Z", dateTime)
	if err != nil {
		fmt.Println(err)
		log.Error(err.Error(), map[string]interface{}{
			"scope": log.Trace(),
		})
		return dateTime
	}
	dt = t.Format("2006-01-02 15:04:05")
	return dt
}

func FormatDatabaseDateTimeWithUtc(dateTime string) string {
	var dt string
	t, err := time.Parse(time.RFC3339, dateTime)
	if err != nil {
		fmt.Println(err)
		log.Error(err.Error(), map[string]interface{}{
			"scope": log.Trace(),
		})
		return dateTime
	}
	dt = t.Format("2006-01-02 15:04:05")
	return dt
}
func FormatDateTimeByPattern(dateTime string, oldPattern string, pattern string) string {
	var dt string
	t, err := time.Parse(oldPattern, dateTime)
	if err != nil {
		log.Error(err.Error(), map[string]interface{}{
			"scope": log.Trace(),
		})
		return dateTime
	}
	dt = t.Format(pattern)
	return dt
}

// Hàm kiểm tra 1 phần tử có trong slice hay không
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func GetTimeMiliseconds(dateTime string) int64 {
	if dateTime == "" {
		return UnixMilli(time.Now())
	}
	dateTimeS := strings.Split(dateTime, " ")
	dates := dateTimeS[0]
	times := dateTimeS[1]
	dates1 := strings.Split(dates, "-")
	times1 := strings.Split(times, ":")
	y, _ := strconv.Atoi(dates1[0])
	m, _ := strconv.Atoi(dates1[1])
	d, _ := strconv.Atoi(dates1[2])
	hh, _ := strconv.Atoi(times1[0])
	mm, _ := strconv.Atoi(times1[1])
	ss, _ := strconv.Atoi(times1[2])
	loc, _ := time.LoadLocation("Asia/Ho_Chi_Minh")
	t := time.Date(y, time.Month(m), d, hh, mm, ss, 0, loc)
	mili := UnixMilli(t)
	return mili
}

func UnixMilli(t time.Time) int64 {
	return t.Round(time.Millisecond).UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}

// Hàm lấy ra các phần tử khác nhau của 2 slice
func Difference(a, b []string) []string {
	mb := make(map[string]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff []string
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}

func IsUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
