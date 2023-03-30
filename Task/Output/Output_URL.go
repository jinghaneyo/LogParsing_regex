package Output

import (
	"LogParsing_regex/Task"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Output_URL struct {
}

func New_Output_URL() *Output_URL {
	task := new(Output_URL)
	task.Init()
	return task
}

func (This *Output_URL) Init() {
	fmt.Println("call Output_URL Init")
}

func (This *Output_URL) DataOut(_task *Task_Output, _thread_data *[][]map[string]string) {

	// 자료구조 : [고루틴별][라인수별]map[태그별]실데이터
	for thr_index := range *_thread_data {
		for row_index := range (*_thread_data)[thr_index] {

			format := _task.Url.Body.Data

			for k, v := range (*_thread_data)[thr_index][row_index] {
				format = strings.ReplaceAll(format, k, v)
			}

			This.Send_Data(_task, &format)
		}
	}
}

func (This *Output_URL) Send_Data(_task *Task_Output, _data *string) bool {

	params := url.Values{}
	params.Add("", *_data)

	var body *strings.Reader
	if _task.Url.Body.Url_encode == bool(true) {
		body = strings.NewReader(params.Encode())
	} else {
		body = strings.NewReader(params.Get(""))
	}

	// tls skip
	if _task.Url.Tls == bool(true) {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	req, err := http.NewRequest(_task.Url.Method, _task.Url.Url, body)
	if err != nil {
		Task.LogInst().WriteLog("OUTPUT_URL", "[FAIL] Connect >> Err = %v", err)
		return false
	}
	//req.SetBasicAuth("banana", "coconuts")
	for k, v := range _task.Url.Header {
		req.Header.Set(k, v)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		Task.LogInst().WriteLog("OUTPUT_URL", "[FAIL] Connect >> Err = %v", err)
		return false
	}
	defer resp.Body.Close()

	return_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Task.LogInst().WriteLog("OUTPUT_URL", "[FAIL] Request >> Err = %v", err)
		return false
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		Task.LogInst().WriteLog("OUTPUT_URL", "[FAIL] Status Code = %d >> Err = %v", resp.StatusCode, err)
		return false
	}

	fmt.Println(string(return_body))
	Task.LogInst().WriteLog("OUTPUT_URL", "[SUCC] RETURN >> %s", string(return_body))

	return true
}
