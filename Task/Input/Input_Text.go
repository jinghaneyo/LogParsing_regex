package Input

import (
	"LogParsing_regex/Task"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

type Input_Text struct {
}

func New_Input_Text() *Input_Text {
	task := new(Input_Text)
	task.Init()
	return task
}

func (This *Input_Text) Init() {
	Task.LogInst().WriteLog("INPUT", "call Input_Text Init")
}

func (This *Input_Text) Load(_task *Task_Input) *[][]map[string]string {

	//*
	// 파일 사이즈를 구해서 동시 진행할 스레드 수에서 작업할 바이트를 구하여 시작점과 종료점을 넘긴다
	stat, err := os.Stat(_task.File.Path)
	if err != nil {
		Task.LogInst().WriteLog("INPUT", "[Load] ERR = %s", err)
		return nil
	}

	// todo 추후 10을 설정으로 빼자
	ThreadCount := 10
	WorkSize := stat.Size() / int64(ThreadCount)
	var wg sync.WaitGroup

	// 고루틴별 추출 데이터
	// 자료구조 : [고루틴별][라인수별]map[태그별]실데이터
	thread_data := make([][]map[string]string, ThreadCount)

	var nStart int64
	nStart = 0
	for i := 0; i < ThreadCount; i++ {
		wg.Add(1)
		go This.Extracter(&wg, i, _task, nStart, WorkSize, &thread_data[i])
		nStart += WorkSize + 1
	}

	wg.Wait()
	return &thread_data
}

func (This *Input_Text) Extracter(
	_wg *sync.WaitGroup,
	_thread_num int,
	_task *Task_Input,
	_nStart int64,
	_nSize int64,
	_dataOut *[]map[string]string) {

	defer _wg.Done()

	f, err := os.OpenFile(_task.File.Path, os.O_RDONLY, 644)
	if err != nil {
		Task.LogInst().WriteLog("INPUT", "[Extracter] ERR = %s", err)
		return
	}
	defer f.Close()

	data := make([]map[string]string, 0)
	var field map[string]string

	var nTotal_Bytes int64

	// 파일 포인터 위치 읽어들여야할 위치로 이동
	f.Seek(_nStart, 0)

	// 일단 2000 읽어 들이자
	nRead_Len := 2000
	if nRead_Len > int(_nSize) {
		nRead_Len = int(_nSize)
	}

	var nRemaning int64
	last_buf := ""
	for {

		// 고루틴 한개가 읽어들여야할 바이트 수만큼 읽어들이기 위해서
		nRemaning = _nSize - nTotal_Bytes
		if nRemaning < int64(nRead_Len) {
			nRead_Len = int(nRemaning)
		}
		if nRead_Len < 1 {
			nRead_Len = 4000
		}

		line_bytes := make([]byte, nRead_Len)
		_, err := f.Read(line_bytes)
		if err != nil {
			*_dataOut = data
			return
		}

		line_string := string(line_bytes)
		lines := strings.Split(line_string, "\n")

		// 한줄씩 파싱(정규식)을 한다
		for i := range lines {

			// _nStart 가 0보다 크면 라인의 첫번째 글자가 아니라 라인의 중간이란 애기이므로
			// 이전 턴에서의 마지막 라인하고 합치자
			if i == 0 {
				if len(last_buf) > 0 {
					lines[i] = last_buf + lines[i]
				}
			}

			// 마지막 라인은 완전한 한줄이 아니므로 임시변수에 저장한 후에
			// 다음턴 첫번째랑 합쳐서 파싱하도록 한다
			if i == len(lines)-1 {
				last_buf = lines[i]
				break
			} else {
				nTotal_Bytes += int64(len(lines[i]))
				nTotal_Bytes += 1 // (\n)개행도 더해준다
			}

			ret, tag, value := This.Parsing(_task, lines[i])

			if nTotal_Bytes <= _nSize {

				if ret == bool(true) {
					// 필드 셋이 새로 시작
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
			} else {

				if ret == bool(true) {
					if _task.File.Start_tag == tag {
						if field != nil {
							data = append(data, field)
						}

						*_dataOut = data
						return
					}

					if field != nil {
						field[tag] = value
					}
				}
			}

			time.Sleep(1 * time.Microsecond)
		}
	}
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
