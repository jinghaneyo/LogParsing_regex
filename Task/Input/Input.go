package Input

type Input_File struct {
	Path      string            `yaml:"path"`
	Regex     string            `yaml:"file_regex"`
	Start_tag string            `yaml:"start_tag"`
	Field_tag map[string]string `yaml:"file_field_tag"`
}

type Task_Input struct {
	Type string     `yaml:"type"`
	File Input_File `yaml:"file"`
}

type Input interface {
	// []row > []col 형태로 데이터를 넘기자
	Load(_task *Task_Input) []map[string]string
}
