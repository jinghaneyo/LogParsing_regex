package Output

import (
	"LogParsing_regex/Task"
	"io"
	"os"
	"strconv"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type Output_sFtp struct {
}

func New_Output_sFtp() *Output_sFtp {
	task := new(Output_sFtp)
	task.Init()
	return task
}

func (This *Output_sFtp) Init() {
	Task.LogInst().WriteLog("OUTPUT_SFTP", "call Output_sFtp Init")
}

func (This *Output_sFtp) DataOut(_task *Task_Output, _thread_data *[][]map[string]string) {

	config := &ssh.ClientConfig{
		User: _task.Sftp.Id,
		Auth: []ssh.AuthMethod{
			ssh.Password(_task.Sftp.Pwd),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	conn, err := ssh.Dial("tcp", _task.Sftp.Server+":"+strconv.Itoa(_task.Sftp.Port), config)
	if err != nil {
		Task.LogInst().WriteLog("OUTPUT_SFTP", "[FAIL] Connect Fail >> ERR = %s", err)
	}
	defer conn.Close()

	client, err := sftp.NewClient(conn)
	if err != nil {
		Task.LogInst().WriteLog("OUTPUT_SFTP", "[FAIL] NewClient is Fail to make instance >> ERR = %s", err)
	}
	defer client.Close()

	dstFile, err := client.Create(_task.Sftp.RemotePath)
	if err != nil {
		Task.LogInst().WriteLog("OUTPUT_SFTP", "[FAIL] Client is Fail to create >> ERR = %s", err)
	}
	defer dstFile.Close()

	srcFile, err := os.Open(_task.Sftp.LocalPath)
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
