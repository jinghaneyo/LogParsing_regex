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

func (This *Output_Text) DataOut(_task *Task_Output, _data *[]map[string]string) {

	f, err := os.Create(_task.File.Path)
	if err != nil {
		fmt.Printf("[Output_Text] DataOut >> Err = %v", err)
		return
	}

	for i := range *_data {

		format := _task.File.Format

		for k, v := range (*_data)[i] {
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
