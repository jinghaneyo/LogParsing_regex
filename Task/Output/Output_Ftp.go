package Output

import (
	"LogParsing_regex/Task"
	"bufio"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/jlaffaye/ftp"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type Output_Ftp struct {
}

func New_Output_Ftp() *Output_Ftp {
	task := new(Output_Ftp)
	task.Init()
	return task
}

func (This *Output_Ftp) Init() {
	Task.LogInst().WriteLog("OUTPUT_FTP", "call Output_Ftp Init")
}

func (This *Output_Ftp) DataOut(_task *Task_Output, _thread_data *[][]map[string]string) {

	if _task.Ftp.Sftp == bool(true) {
		Upload_sFtp(_task, _thread_data)
	} else {
		Upload_Ftp(_task, _thread_data)
	}
}

func Upload_Ftp(_task *Task_Output, _thread_data *[][]map[string]string) {

	conn, err := ftp.Dial(_task.Ftp.Server+":"+strconv.Itoa(_task.Ftp.Port), ftp.DialWithTimeout(time.Duration(_task.Ftp.Connect_timeout*int(time.Second))))
	if err != nil {
		Task.LogInst().WriteLog("OUTPUT_FTP", "[FAIL][FTP] Connect >> Err = %v", err)
		return
	}

	err = conn.Login(_task.Ftp.Id, _task.Ftp.Pwd)
	if err != nil {
		Task.LogInst().WriteLog("OUTPUT_FTP", "[FAIL][FTP] Login >> Err = %v", err)
		return
	}

	file, err := os.OpenFile(
		_task.Ftp.LocalPath,
		os.O_RDONLY,
		os.FileMode(0644),
	)
	if err != nil {
		Task.LogInst().WriteLog("OUTPUT_FTP", "[FAIL][FTP] File Open Fail(%s) >> Err = %v", _task.Ftp.LocalPath, err)
		return
	}

	r := bufio.NewReader(file)
	err = conn.Stor(_task.Ftp.RemotePath, r)
	if err != nil {
		Task.LogInst().WriteLog("OUTPUT_FTP", "[FAIL][FTP] Upload Fail(%s) >> Err = %v", _task.Ftp.LocalPath, err)
	} else {
		Task.LogInst().WriteLog("OUTPUT_FTP", "[SUCC][FTP] Upload succ(%s) >> Err = %v", _task.Ftp.LocalPath, err)
	}
}

func Upload_sFtp(_task *Task_Output, _thread_data *[][]map[string]string) {

	// 먼저 파일로 저장
	file := New_Output_Text()
	ret := file.Write(_task.Ftp.LocalPath, _task.Ftp.Format, "OUTPUT_FTP", _thread_data)
	if ret == bool(false) {
		return
	}

	config := &ssh.ClientConfig{
		User: _task.Ftp.Id,
		Auth: []ssh.AuthMethod{
			ssh.Password(_task.Ftp.Pwd),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	conn, err := ssh.Dial("tcp", _task.Ftp.Server+":"+strconv.Itoa(_task.Ftp.Port), config)
	if err != nil {
		Task.LogInst().WriteLog("OUTPUT_FTP", "[FAIL][SFTP] Connect Fail >> ERR = %s", err)
		return
	}
	defer conn.Close()

	client, err := sftp.NewClient(conn)
	if err != nil {
		Task.LogInst().WriteLog("OUTPUT_FTP", "[FAIL][SFTP] NewClient is Fail to make instance >> ERR = %s", err)
		return
	}
	defer client.Close()

	dstFile, err := client.Create(_task.Ftp.RemotePath)
	if err != nil {
		Task.LogInst().WriteLog("OUTPUT_FTP", "[FAIL][SFTP] Client is Fail to create >> ERR = %s", err)
		return
	}
	defer dstFile.Close()

	srcFile, err := os.Open(_task.Ftp.LocalPath)
	if err != nil {
		Task.LogInst().WriteLog("OUTPUT_FTP", "[FAIL][SFTP] Loacal file is fail to open>> ERR = %s", err)
		return
	}

	bytes, err := io.Copy(dstFile, srcFile)
	if err != nil {
		Task.LogInst().WriteLog("OUTPUT_FTP", "[FAIL][SFTP] Upload is fail >> ERR = %s | %s => %s", err, srcFile, dstFile)
	} else {
		Task.LogInst().WriteLog("OUTPUT_FTP", "[SUCC][SFTP] Upload is SUCC (bytes = %d) | %s => %s", bytes, srcFile, dstFile)
	}
}
