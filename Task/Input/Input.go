package Input

const (
	Tag_1_in_line = iota
	Tag_all_in_line
)

type Input_File struct {
	Path         string            `yaml:"path"`
	Extract_type string            `yaml:"extract_type"`
	Tag_in_Line  string            `yaml:"tag_in_line"`
	Start_tag    string            `yaml:"start_tag"`
	Field_tag    map[string]string `yaml:"file_field_tag"`
}

type Task_Input struct {
	Type string     `yaml:"type"`
	File Input_File `yaml:"file"`
}

type Input interface {
	// []row > []col 형태로 데이터를 넘기자
	Load(_task *Task_Input) []map[string]string
}
