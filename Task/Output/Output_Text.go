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

	f, err := os.Create(_task.File.Path)
	if err != nil {
		Task.LogInst().WriteLog("OUTPUT_TEXT", "[FAIL] Open file is fail >> ERR = %s | File = %s", err, _task.File.Path)
		return
	}

	// 자료구조 : [고루틴별][라인수별]map[태그별]실데이터
	for thr_index := range *_thread_data {
		for row_index := range (*_thread_data)[thr_index] {

			format := _task.File.Format

			for k, v := range (*_thread_data)[thr_index][row_index] {
				format = strings.ReplaceAll(format, k, v)
			}

			_, err := f.WriteString(format)
			if err != nil {
				Task.LogInst().WriteLog("OUTPUT_TEXT", "[FAIL] WriteString is fail >> ERR = %s | File = %s", err, _task.File.Path)
				return
			}
			f.WriteString("\n")
		}
	}
}
