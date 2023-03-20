package Output

import (
	"fmt"
	"io"
	"log"
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
	fmt.Println("call Output_sFtp Init")
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
		log.Fatal(err)
	}
	defer conn.Close()

	client, err := sftp.NewClient(conn)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	dstFile, err := client.Create(_task.Sftp.RemotePath)
	if err != nil {
		log.Fatal(err)
	}
	defer dstFile.Close()

	srcFile, err := os.Open(_task.Sftp.LocalPath)
	if err != nil {
		log.Fatal(err)
	}

	bytes, err := io.Copy(dstFile, srcFile)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d bytes copied\n", bytes)
}
