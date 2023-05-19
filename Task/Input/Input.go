package Input

const (
	Tag_1_in_line = iota
	Tag_all_in_line
)

type Input_File struct {
	Path            string              `yaml:"path"`
	Extract_type    string              `yaml:"extract_type"`
	Tag_in_Line     string              `yaml:"tag_in_line"`
	Split_word      string              `yaml:"split_word"`
	Start_field_tag string              `yaml:"start_field_tag"`
	Start_block_tag string              `yaml:"start_block_tag"`
	Field_tag       map[string][]string `yaml:"file_field_tag"`
}

type Input_Db_Config struct {
	Odbc            string `yaml:"odbc"`
	Database        string `yaml:"database"`
	Id              string `yaml:"id"`
	Pwd             string `yaml:"pwd"`
	Port            int    `yaml:"port"`
	Max_Bulksize    uint64 `yaml:"bulk_insert_max_size"`
	Connect_timeout int    `yaml:"connect_timeout"`
	Sql_Select      string `yaml:"select"`
}

type Task_Input struct {
	Worker_count int             `yaml:"worker_count"`
	Type         string          `yaml:"type"`
	File         Input_File      `yaml:"file"`
	Db           Input_Db_Config `yaml:"db"`
}

type Input interface {
	// []row > []col 형태로 데이터를 넘기자
	Load(_task *Task_Input) []map[string]string
}
