package Input

import (
	"LogParsing_regex/Task"
	"os"
	"regexp"
	"strconv"
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

	// 파일 사이즈를 구해서 동시 진행할 스레드 수에서 작업할 바이트를 구하여 시작점과 종료점을 넘긴다
	stat, err := os.Stat(_task.File.Path)
	if err != nil {
		Task.LogInst().WriteLog("INPUT", "[Load] ERR = %s", err)
		return nil
	}

	// todo 추후 10을 설정으로 빼자
	ThreadCount := _task.Worker_count
	WorkSize := stat.Size() / int64(ThreadCount)
	var wg sync.WaitGroup

	// 고루틴별 추출 데이터
	// 자료구조 : [고루틴별][라인수별]map[태그별]실데이터
	thread_data := make([][]map[string]string, ThreadCount)

	var nStart int64
	nStart = 0
	for i := 0; i < ThreadCount; i++ {
		wg.Add(1)

		if _task.File.Extract_type == "regex" {
			go This.Extracter(&wg, i, _task, nStart, WorkSize, &thread_data[i])
		} else if _task.File.Extract_type == "split" {
			go This.Split(&wg, i, _task, nStart, WorkSize, &thread_data[i])
		}
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
	nRead_Len := 4000
	if nRead_Len > int(_nSize) {
		nRead_Len = int(_nSize)
	}

	static_field := make(map[string]string, 0)
	bFirst_Block := true
	bFirst_Field := true

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

			ret, tag, value, type_valiable := This.Parsing_Regex(_task, &lines[i], _thread_num)

			if nTotal_Bytes <= _nSize {
				if ret == bool(true) {

					// 파일일 읽어들이고 블럭 시작 태그를 아직 찾지 못했다
					// 즉, 현재 필드들은 다른 고루틴의 블럭의 필드들이다
					if bFirst_Block == bool(true) {
						if _task.File.Start_block_tag == tag {
							// 블럭 시작
							bFirst_Block = false
						} else {
							continue
						}
					}

					// 필드 셋이 새로 시작
					if _task.File.Start_field_tag == tag {

						// 첫 블럭이 시작된 후 "시작 필드 태그" 는 아직 필드셋이 검색되지 않았으므로
						// 건너 뛰고 다음 "시작 필드 태그" 일 경우 append를 하도록 한다
						if bFirst_Field == bool(true) {
							bFirst_Field = false
							continue
						}

						// 필드 셋이 새로 시작을 하니 이전 셋은 데이터 배열에 넣도록 한다
						if field != nil {

							for k, v := range static_field {
								_, ok := field[k]
								if ok == bool(false) {
									field[k] = v
								}
							}

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
					if _task.File.Start_field_tag == tag {
						if field != nil {

							for k, v := range static_field {
								_, ok := field[k]
								if ok == bool(false) {
									field[k] = v
								}
							}

							data = append(data, field)
						}

						// 필드 셋을 새로 만들어주자
						field = make(map[string]string, 0)
					}

					// 블럭 시작 태그가 있다면 이전 필드들은 종료되었으므로 더이상 읽어들이지 않는다
					if _task.File.Start_block_tag == tag {
						*_dataOut = data
						return
					}

					if field != nil {
						field[tag] = value
					}
				}
			}

			// 스태틱 필드값 유지
			if type_valiable == "static" {
				if ret == bool(true) {
					static_field[tag] = value
				}
			}

			// cpu 과부하 방지
			time.Sleep(1 * time.Millisecond)
		}
	}
}

func (This *Input_Text) Parsing_Regex(
	_task *Task_Input,
	_line *string,
	_Thread_Num int) (bool, string, string, string) {

	for tag, reg := range _task.File.Field_tag {

		// 0번째는 static(유지), local(새로 갱신)
		// 1번째는 필터링할 정규식
		r := regexp.MustCompile(reg[1])
		result := r.FindAllStringSubmatch(*_line, -1)

		if result != nil {

			// return 1 : 존재 여부
			// return 2 : 태그명
			// return 3 : 정규식을 통한 필터링 값
			// return 4 : static, local 여부
			if len(result[0]) > 1 {
				// 정규식에 그룹이 없으면 찾은 라인의 전체가 value, 있으면 첫번째 그룹이 value가 된다
				//Task.LogInst().WriteLog("INPUT", "[Parsing_Regex][THR %d][TAG = %s][VAL = %s][LINE = %s]", _Thread_Num, tag, result[0][1], *_line)
				return true, tag, result[0][1], reg[0]
			} else {
				// 정규식에 그룹이 없으면 찾은 라인의 전체가 value, 있으면 첫번째 그룹이 value가 된다
				//Task.LogInst().WriteLog("INPUT", "[Parsing_Regex][THR %d][TAG = %s][VAL = %s][LINE = %s]", _Thread_Num, tag, result[0][0], *_line)
				return true, tag, result[0][0], reg[0]
			}
		}
	}

	return false, "", "", ""
}

func (This *Input_Text) Split(
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
	nRead_Len := 4000
	if nRead_Len > int(_nSize) {
		nRead_Len = int(_nSize)
	}

	bFirstRead := true
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
				// 다중스레드이기 때문에 스레드 번호가 0 보다 크다면 최초 한줄은 완전한 한줄이 아니라 중간이 잘린 데이터이다.
				// 해서 첫줄은 일단 건너뛰고 다음줄 부터 정상적으로 추출하자
				if bFirstRead == bool(true) {
					bFirstRead = false
					if _thread_num > 0 {
						nTotal_Bytes += int64(len(lines[i]))
						continue
					}
				}

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

			temp := strings.Split(lines[i], _task.File.Split_word)

			field = make(map[string]string, 0)
			for i := range temp {

				for tag, val := range _task.File.Field_tag {

					// 0번째는 static(유지), local(새로 갱신)
					// 1번째는 정규식
					if val[1] == strconv.Itoa(i) {
						field[tag] = temp[i]
					}
				}
			}
			data = append(data, field)

			if nTotal_Bytes > _nSize {
				*_dataOut = data
				return
			}

			time.Sleep(1 * time.Microsecond)
		}
	}
}
