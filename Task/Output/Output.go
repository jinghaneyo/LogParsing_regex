package Output

type Output_File struct {
	Path   string `yaml:"path"`
	Format string `yaml:"format"`
}

type Task_Output struct {
	Type string      `yaml:"type"`
	File Output_File `yaml:"file"`
}

type Output interface {
	DataOut(_task *Task_Output, _data *[]map[string]string) bool
}
