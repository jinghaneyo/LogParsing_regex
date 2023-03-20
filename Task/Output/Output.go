package Output

type Output_File_Config struct {
	Path   string `yaml:"path"`
	Format string `yaml:"format"`
}

type Output_Ftp_Config struct {
	Server          string `yaml:"server"`
	Id              string `yaml:"id"`
	Pwd             string `yaml:"pwd"`
	Port            int    `yaml:"port"`
	Connect_timeout int    `yaml:"connect_timeout"`
	LocalPath       string `yaml:"local"`
	RemotePath      string `yaml:"remote"`
}

type Task_Output struct {
	Type string             `yaml:"type"`
	File Output_File_Config `yaml:"file"`
	Ftp  Output_Ftp_Config  `yaml:"ftp"`
	Sftp Output_Ftp_Config  `yaml:"sftp"`
}

type Output interface {
	DataOut(_task *Task_Output, _data *[][]map[string]string)
}
