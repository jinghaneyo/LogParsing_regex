package main

import (
	"LogParsing_regex/Task"
	"LogParsing_regex/Task/Input"
	"LogParsing_regex/Task/Output"
	"io/ioutil"
	"log"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

type Task_Confg struct {
	Task_Parser struct {
		Input  Input.Task_Input     `yaml:"input"`
		Output []Output.Task_Output `yaml:"output"`
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

	// 쿼리일 경우 bulk insert 가능 여부 확인
	for i := range *task {
		for j := range (*task)[i].Task_Parser.Output {
			out := &(*task)[i].Task_Parser.Output[j]

			// bulk insert 가능 여부 판단
			CanBulkInsert := make([]bool, 0)

			// bulk insert 의 value 추출 구문
			Value_sql := make([]string, 0)
			// bulk insert 의 insert 추출 구문
			Insert_sql := make([]string, 0)

			for index := range out.Db.Sql.Data {

				out.Db.Sql.Data[index] = strings.Trim(out.Db.Sql.Data[index], " ")

				// 마지막에 ; 있으면 제거
				if out.Db.Sql.Data[index][len(out.Db.Sql.Data[index])-1] == ';' {
					out.Db.Sql.Data[index] = out.Db.Sql.Data[index][0 : len(out.Db.Sql.Data[index])-1]
				}

				r, _ := regexp.Compile(`(INSERT[\s]+INTO[\s]+.+)VALUES[\s]*(\(.+\))`)
				ret := r.FindStringSubmatch(out.Db.Sql.Data[index])
				if len(ret) > 2 {
					if strings.EqualFold(out.Db.Sql.Data[index], "DUPLICATE KEY") == bool(false) {
						CanBulkInsert = append(CanBulkInsert, true)

						Insert_sql = append(Value_sql, ret[1])
						Value_sql = append(Value_sql, ret[2])
					} else {
						CanBulkInsert = append(CanBulkInsert, false)
					}
				}
			}

			out.Db.Sql.CanBulkInsert = CanBulkInsert
			out.Db.Sql.Insert_Sql = Insert_sql
			out.Db.Sql.Values_Sql = Value_sql
		}
	}

	return task
}

func main() {

	task := Load_Task()
	if task == nil {
		return
	}

	Task.LogInst().Enable_Console(false)
	Task.LogInst().Enable_LogDate(false)

	for _, t := range *task {

		// 추출 로드 및 파싱
		// 자료구조 : [고루틴별][라인수별]map[태그별]실데이터
		var thread_data *[][]map[string]string

		if t.Task_Parser.Input.Type == "file" {
			input := Input.New_Input_Text()
			thread_data = input.Load(&t.Task_Parser.Input)
		}

		// 결과 출력
		if thread_data != nil {
			for i := range t.Task_Parser.Output {
				out := t.Task_Parser.Output[i]
				var output Output.Output
				if out.Type == "file" {
					output = Output.New_Output_Text()
				} else if out.Type == "ftp" {
					output = Output.New_Output_Ftp()
				} else if out.Type == "db" {
					output = Output.New_Output_Database()
				} else if out.Type == "url" {
					output = Output.New_Output_URL()
				}

				if output != nil {
					output.DataOut(&out, thread_data)
				}
			}

			*thread_data = nil
		}
	}
}
