package Output

type Output_File_Config struct {
	Path   string `yaml:"path"`
	Format string `yaml:"format"`
}

type Output_Ftp_Config struct {
	Sftp            bool   `yaml:"sftp"`
	Server          string `yaml:"server"`
	Id              string `yaml:"id"`
	Pwd             string `yaml:"pwd"`
	Port            int    `yaml:"port"`
	Connect_timeout int    `yaml:"connect_timeout"`
	Format          string `yaml:"format"`
	LocalPath       string `yaml:"local"`
	RemotePath      string `yaml:"remote"`
}

type Output_Db_Config struct {
	Odbc            string `yaml:"odbc"`
	Database        string `yaml:"database"`
	Id              string `yaml:"id"`
	Pwd             string `yaml:"pwd"`
	Port            int    `yaml:"port"`
	Max_Bulksize    uint64 `yaml:"bulk_insert_max_size"`
	Connect_timeout int    `yaml:"connect_timeout"`
	Sql             struct {
		First  string   `yaml:"first"`
		Data   []string `yaml:"data"`
		Finish string   `yaml:"finish"`
		// Data 의 쿼리를 보고 INSERT 구문으로만 구성되어 있어 bulk insert 쿼리로 만들수 있을지 여부
		CanBulkInsert []bool
		// bulk insert 가능일 경우 쿼리에서 insert만 추출한 구문
		Insert_Sql []string
		// bulk insert 가능일 경우 쿼리에서 values만 추출한 구문
		Values_Sql []string
	} `yaml:"sql"`
}

type Output_URL_Config struct {
	Tls    bool              `yaml:"tls"`
	Method string            `yaml:"method"`
	Url    string            `yaml:"url"`
	Header map[string]string `yaml:"header"`
	Body   struct {
		Url_encode bool   `yaml:"url_encode"`
		Data       string `yaml:"data"`
	} `yaml:"body"`
}

type Task_Output struct {
	Type string             `yaml:"type"`
	File Output_File_Config `yaml:"file"`
	Ftp  Output_Ftp_Config  `yaml:"ftp"`
	Db   Output_Db_Config   `yaml:"db"`
	Url  Output_URL_Config  `yaml:"url"`
}

type Output interface {
	DataOut(_task *Task_Output, _data *[][]map[string]string)
}
