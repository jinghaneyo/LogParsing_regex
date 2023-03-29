package Output

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/jlaffaye/ftp"
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

	conn, err := ftp.Dial(_task.Ftp.Server+":"+strconv.Itoa(_task.Ftp.Port), ftp.DialWithTimeout(time.Duration(_task.Ftp.Connect_timeout*int(time.Second))))
	if err != nil {
		fmt.Printf("[Output_URL] Connect >> Err = %v", err)
		return
	}

	err = conn.Login(_task.Ftp.Id, _task.Ftp.Pwd)
	if err != nil {
		fmt.Printf("[Output_URL] Login >> Err = %v", err)
		return
	}

	file, err := os.OpenFile(
		_task.Ftp.LocalPath,
		os.O_RDONLY,
		os.FileMode(0644),
	)
	if err != nil {
		fmt.Printf("[Output_URL] File Open Fail(%s) >> Err = %v", _task.Ftp.LocalPath, err)
		return
	}

	r := bufio.NewReader(file)
	err = conn.Stor(_task.Ftp.RemotePath, r)
	if err != nil {
		fmt.Printf("[Output_URL] Upload Fail(%s) >> Err = %v", _task.Ftp.LocalPath, err)
		return
	}
}
