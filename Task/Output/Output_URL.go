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

	bulk_data := ""
	var bulk_len uint64
	// 자료구조 : [고루틴별][라인수별]map[태그별]실데이터
	for thr_index := range *_thread_data {
		for row_index := range (*_thread_data)[thr_index] {

			format := _task.Url.Body.Repeat_data_value

			for k, v := range (*_thread_data)[thr_index][row_index] {
				format = strings.ReplaceAll(format, k, v)
			}

			if bulk_len == 0 {
				bulk_data = _task.Url.Body.Repeat_data_first + format

				bulk_len += uint64(len(bulk_data))
			} else if _task.Url.Body.Bulk_max_size < (bulk_len + uint64(len(format)+len(_task.Url.Body.Repeat_data_last))) {
				// 설정 크기를 넘어가면 쿼리 실행
				bulk_data += _task.Url.Body.Repeat_data_last
				This.Send_Data(_task, &bulk_data)

				// 다시 data 구문부터 구성
				bulk_data = _task.Url.Body.Repeat_data_first + format

				bulk_len = uint64(len(bulk_data))
			} else {
				// value 구문만 추가
				bulk_data += ","
				bulk_data += format

				bulk_len += uint64(len(bulk_data)) + 1
			}
		}
	}

	// 남은 쿼리 실행
	if len(bulk_data) > 0 {
		bulk_data += _task.Url.Body.Repeat_data_last
		This.Send_Data(_task, &bulk_data)
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
		Task.LogInst().WriteLog("OUTPUT_URL", "[FAIL] Status Code = %d >> Err = %v\n>> Data = %s",
			resp.StatusCode, err, *_data)
		return false
	}

	fmt.Println(string(return_body))
	Task.LogInst().WriteLog("OUTPUT_URL", "[SUCC] RETURN >> %s", string(return_body))

	return true
}
