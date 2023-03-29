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
		Task.LogInst().WriteLog("OUTPUT_FTP", "[FAIL] Connect >> Err = %v", err)
		return
	}

	err = conn.Login(_task.Ftp.Id, _task.Ftp.Pwd)
	if err != nil {
		Task.LogInst().WriteLog("OUTPUT_FTP", "[FAIL] Login >> Err = %v", err)
		return
	}

	file, err := os.OpenFile(
		_task.Ftp.LocalPath,
		os.O_RDONLY,
		os.FileMode(0644),
	)
	if err != nil {
		Task.LogInst().WriteLog("OUTPUT_FTP", "[FAIL] File Open Fail(%s) >> Err = %v", _task.Ftp.LocalPath, err)
		return
	}

	r := bufio.NewReader(file)
	err = conn.Stor(_task.Ftp.RemotePath, r)
	if err != nil {
		Task.LogInst().WriteLog("OUTPUT_FTP", "[FAIL] Upload Fail(%s) >> Err = %v", _task.Ftp.LocalPath, err)
		return
	} else {
		Task.LogInst().WriteLog("OUTPUT_FTP", "[SUCC] Upload succ(%s) >> Err = %v", _task.Ftp.LocalPath, err)
	}
}

func Upload_sFtp(_task *Task_Output, _thread_data *[][]map[string]string) {

	config := &ssh.ClientConfig{
		User: _task.Ftp.Id,
		Auth: []ssh.AuthMethod{
			ssh.Password(_task.Ftp.Pwd),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	conn, err := ssh.Dial("tcp", _task.Ftp.Server+":"+strconv.Itoa(_task.Ftp.Port), config)
	if err != nil {
		Task.LogInst().WriteLog("OUTPUT_SFTP", "[FAIL] Connect Fail >> ERR = %s", err)
	}
	defer conn.Close()

	client, err := sftp.NewClient(conn)
	if err != nil {
		Task.LogInst().WriteLog("OUTPUT_SFTP", "[FAIL] NewClient is Fail to make instance >> ERR = %s", err)
	}
	defer client.Close()

	dstFile, err := client.Create(_task.Ftp.RemotePath)
	if err != nil {
		Task.LogInst().WriteLog("OUTPUT_SFTP", "[FAIL] Client is Fail to create >> ERR = %s", err)
	}
	defer dstFile.Close()

	srcFile, err := os.Open(_task.Ftp.LocalPath)
	if err != nil {
		Task.LogInst().WriteLog("OUTPUT_SFTP", "[FAIL] Loacal file is fail to open>> ERR = %s", err)
	}

	bytes, err := io.Copy(dstFile, srcFile)
	if err != nil {
		Task.LogInst().WriteLog("OUTPUT_SFTP", "[FAIL] Upload is fail >> ERR = %s | %s => %s", err, srcFile, dstFile)
	} else {
		Task.LogInst().WriteLog("OUTPUT_SFTP", "[SUCC] Upload is SUCC (bytes = %d) | %s => %s", bytes, srcFile, dstFile)
	}
}
