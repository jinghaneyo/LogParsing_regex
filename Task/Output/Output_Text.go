package Output

import (
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
		fmt.Printf("[Output_Text] DataOut >> Err = %v", err)
		return
	}

	// 자료구조 : [고루틴별][라인수별]map[태그별]실데이터
	for thr_index := range *_thread_data {
		for row_index := range (*_thread_data)[thr_index] {

			format := _task.File.Format

			for k, v := range (*_thread_data)[thr_index][row_index] {
				format = strings.ReplaceAll(format, k, v)
			}

			n3, err := f.WriteString(format)
			if err != nil {
				fmt.Printf("[Output_Text] WriteString >> Err = %v, %v", err, n3)
				return
			}
			f.WriteString("\n")
		}
	}
}
