package Input

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
)

type Input_Text struct {
}

func New_Input_Text() *Input_Text {
	task := new(Input_Text)
	task.Init()
	return task
}

func (This *Input_Text) Init() {
	fmt.Println("call Input_Text Init")
}

func (This *Input_Text) Load(_task *Task_Input) *[]map[string]string {

	f, err := os.OpenFile(_task.File.Path, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Fatalf("open file error: %v", err)
		return nil
	}
	defer f.Close()

	data := make([]map[string]string, 0)
	var field map[string]string

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		ret, tag, value := This.Parsing(_task, sc.Text())
		if ret == bool(true) {
			if _task.File.Start_tag == tag {
				// 필드 셋이 새로 시작을 하니 이전 셋은 데이터 배열에 넣도록 한다
				if field != nil {
					data = append(data, field)
				}
				// 필드 셋을 새로 만들어주자
				field = make(map[string]string, 0)
			}
			if field != nil {
				field[tag] = value
			}
		}
	}
	if err := sc.Err(); err != nil {
		log.Fatalf("scan file error: %v", err)
		return nil
	}

	return &data
}

func (This *Input_Text) Parsing(_task *Task_Input, _line string) (bool, string, string) {

	for tag, reg := range _task.File.Field_tag {

		r := regexp.MustCompile(reg)
		result := r.FindAllStringSubmatch(_line, -1)

		if result != nil {

			if len(result[0]) > 1 {
				// 정규식에 그룹이 없으면 찾은 라인의 전체가 value, 있으면 첫번째 그룹이 value가 된다
				return true, tag, result[0][1]
			} else {
				// 정규식에 그룹이 없으면 찾은 라인의 전체가 value, 있으면 첫번째 그룹이 value가 된다
				return true, tag, result[0][0]
			}
		}
	}

	return false, "", ""
}
