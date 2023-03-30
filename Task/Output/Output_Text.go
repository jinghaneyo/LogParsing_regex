package Output

import (
	"LogParsing_regex/Task"
	"fmt"
	"os"
	"strings"
)

type Output_Text struct {
}

func New_Output_Text() *Output_Text {
	task := new(Output_Text)
	task.Init()
	return task
}

func (This *Output_Text) Init() {
	fmt.Println("call Output_Text Init")
}

func (This *Output_Text) DataOut(_task *Task_Output, _thread_data *[][]map[string]string) {

	This.Write(_task.File.Path, _task.File.Format, "OUTPUT_TEXT", _thread_data)
}

func (This *Output_Text) Write(
	_path string,
	_format string,
	_LogFile string,
	_thread_data *[][]map[string]string) bool {

	f, err := os.Create(_path)
	if err != nil {
		Task.LogInst().WriteLog(_LogFile, "[FAIL] Open file is fail >> ERR = %s | File = %s", err, _path)
		return false
	}

	// 자료구조 : [고루틴별][라인수별]map[태그별]실데이터
	for thr_index := range *_thread_data {
		for row_index := range (*_thread_data)[thr_index] {

			format := _format

			for k, v := range (*_thread_data)[thr_index][row_index] {
				format = strings.ReplaceAll(format, k, v)
			}

			_, err := f.WriteString(format)
			if err != nil {
				Task.LogInst().WriteLog(_LogFile, "[FAIL] WriteString is fail >> ERR = %s | File = %s", err, _path)
				return false
			}
			f.WriteString("\n")
		}
	}

	return true
}
