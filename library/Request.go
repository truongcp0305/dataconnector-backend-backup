package library

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"data-connector/log"
)

var HEADER = map[string]string{"Authorization": "Bearer eyJhbGciOiJSUzI1NiIsInR5cGUiOiJKV1QifQ==.eyJpZCI6IjEyMDciLCJmaXJzdE5hbWUiOiJBZG1pbiIsImxhc3ROYW1lIjoiU3lzdGVtIiwidXNlck5hbWUiOiJzeXN0ZW1fYWRtaW4iLCJkaXNwbGF5TmFtZSI6IlN5c3RlbSBBZG1pbiIsImVtYWlsIjoiYWRtaW5fdXNlckBzeW1wZXIudm4iLCJwaG9uZSI6IiIsInN0YXR1cyI6IjEiLCJhdmF0YXIiOiIiLCJ0ZW5hbnRJZCI6IjEiLCJ0eXBlIjoidXNlciIsImlwIjoiMjEwLjI0NS4xMDQuMTU3IiwidXNlckFnZW50IjoiTW96aWxsYVwvNS4wIChNYWNpbnRvc2g7IEludGVsIE1hYyBPUyBYIDEwXzE1XzcpIEFwcGxlV2ViS2l0XC81MzcuMzYgKEtIVE1MLCBsaWtlIEdlY2tvKSBDaHJvbWVcLzEwMy4wLjAuMCBTYWZhcmlcLzUzNy4zNiIsImxvY2F0aW9uIjpbXSwiZGVsZWdhdGVkQnkiOnsiaWQiOiIxOSIsImVtYWlsIjoiaG9hbmduZEBzeW1wZXIudm4iLCJuYW1lIjoiTmd1eVx1MWVjNW4gXHUwMTEwXHUwMGVjbmggSG9cdTAwZTBuZyJ9LCJ0ZW5hbnQiOnsiaWQiOiIxIn0sImlzX2Nsb3VkIjpmYWxzZSwidGVuYW50X2RvbWFpbiI6InZ0aG1ncm91cC52biIsInJvbGUiOiJhdXRvIn0=.MTc4ZTJiMmIxNDQwOTlkYjYzNDY5ZjQ0ZGNiNDc4ZDFlYmI0Mjc0NmU0ZDRhZjk0YTFkOGZiY2UxN2JiMGE1YTUzMDRiODIyMDE5MWFiYjNhMjYwZGE0YmU4ZDE3YWJiNzI2ZGVjNzlkZDk5NzA1MWY0NWYxNjRjYWFlNTVlYTZmODJkOTJkZTRlZTZmOGM2MmQ3YjU1NmQ3YTU4YTY0ODIzZGFmOWEzNjZjZDg4MTUyM2ZiMGUxMWZmNTgyYjM3NmYyMDExYTVlZWUxMjEzNDBlMjcyNzQ4YWY0YmViZGZlZTlmNjZhMDJlMDM0MWYwZDVhZjI2MTY5NjNjMzY0OTRiZDU2ZjEyNTgzY2QzNDhlZWRlODY1YjBiYTA5ZTFjZGVkMGRjNmNjMzI1YTQ2NWY2ZjM1NTg1OWRmMjgzZmJkYmEzMGQ4Njc5NDAzNWJhY2ViODQ4MDE0NGE1MDA1ZTFmYzlhODgwY2Q5MmZlM2JlNmI0NmYxYzk2NTJjODM0YjdjNTQwYjNlY2FkOWU3YzZhYzUyZWVjNTZmNDNhNDJiZDQ0ZWQ5NWNkZDNhMmEwZGU3OTQyM2RlZGYyNTFmZTg2OTVjOTdlYzQ1NjE3ZDMwMzllMzNmYzExYjcwYjE4YmMxYzMyY2UyZDhiOWFiYzU5NzM4YWI4NmE4MzMzZDg="}

type Request struct {
	Header            map[string]string
	Body              map[string]string
	Url               string
	Method            string
	SuppressParseData bool
}
type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
type RequestInterface interface {
	Send() (string, error)
}

func (r Request) Send() (Response, error) {
	if r.Method == "" {
		r.Method = "GET"
	}
	dataBody := url.Values{}
	for k, v := range r.Body {
		dataBody.Set(k, v)
	}
	req, err := http.NewRequest(r.Method, r.Url, strings.NewReader(dataBody.Encode()))
	if err != nil {
		fmt.Println("errerrerrerr1")
		fmt.Println(err)
		log.Error(err.Error(), map[string]interface{}{
			"scope": log.Trace(),
			"res":   err,
		})
	}
	r.addHeader(req)
	req.Header.Add("Content-Length", strconv.Itoa(len(dataBody.Encode())))
	req.Header.Set("Connection", "Keep-Alive")
	req.Header.Set("Accept-Language", "en-US")
	req.Header.Set("Accept", "application/json")
	req.Close = true
	client := &http.Client{}
	response, err := client.Do(req)

	if err != nil {
		fmt.Println("errerrerrerr")
		fmt.Println(err)
		log.Error(err.Error(), map[string]interface{}{
			"scope": log.Trace(),
			"res":   err,
		})
		var x Response
		return x, err
	}
	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)
	if !r.SuppressParseData {
		var dataResponse Response
		err1 := json.Unmarshal(data, &dataResponse)
		return dataResponse, err1
	} else {
		dataResponse := Response{
			Data:    string(data),
			Status:  response.StatusCode,
			Message: "",
		}
		return dataResponse, nil
	}

}

func (r Request) addHeader(req *http.Request) {
	for key, value := range r.Header {
		req.Header.Set(key, value)
	}
}
