package Input

const (
	Tag_1_in_line = iota
	Tag_all_in_line
)

type Input_File struct {
	Path         string            `yaml:"path"`
	Extract_type string            `yaml:"extract_type"`
	Tag_in_Line  string            `yaml:"tag_in_line"`
	Split_word   string            `yaml:"split_word"`
	Start_tag    string            `yaml:"start_tag"`
	Field_tag    map[string]string `yaml:"file_field_tag"`
}

type Input_Db_Config struct {
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

type Task_Input struct {
	Type string          `yaml:"type"`
	File Input_File      `yaml:"file"`
	Db   Input_Db_Config `yaml:"db"`
}

type Input interface {
	// []row > []col 형태로 데이터를 넘기자
	Load(_task *Task_Input) []map[string]string
}
