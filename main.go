package main

import (
	"LogParsing_regex/Task"
	"LogParsing_regex/Task/Input"
	"LogParsing_regex/Task/Output"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Task_Confg struct {
	Task_Parser struct {
		Input  Input.Task_Input   `yaml:"input"`
		Output Output.Task_Output `yaml:"output"`
	} `yaml:"task"`
}

func Load_Task() *[]Task_Confg {

	buf, err := ioutil.ReadFile("task.yaml")
	if err != nil {
		log.Printf("[FAIL] ReadFile (%s) >> ERR = %v", "task.yaml", err)
		return nil
	}

	task := &[]Task_Confg{}
	err = yaml.Unmarshal(buf, task)
	if err != nil {
		log.Fatalf("[FAIL] Unmarshal: %v", err)
		return nil
	}

	return task
}

func main() {

	task := Load_Task()
	if task == nil {
		return
	}

	Task.GetInst().Enable_Console(false)
	Task.GetInst().Enable_LogDate(false)

	for _, t := range *task {

		// 자료구조 : [고루틴별][라인수별]map[태그별]실데이터
		var thread_data *[][]map[string]string

		if t.Task_Parser.Input.Type == "file" {
			input := Input.New_Input_Text()
			thread_data = input.Load(&t.Task_Parser.Input)
		}

		if thread_data != nil {
			if t.Task_Parser.Output.Type == "file" {
				output := Output.New_Output_Text()
				output.DataOut(&t.Task_Parser.Output, thread_data)
			}

			*thread_data = nil
		}
	}
}
