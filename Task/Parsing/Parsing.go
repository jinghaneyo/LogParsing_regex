package Parsing

import "fmt"

type Parser_regex struct {
}

func New_Parser_regex() *Parser_regex {
	task := new(Parser_regex)
	task.Init()
	return task
}

func (This *Parser_regex) Init() {
	fmt.Println("call Parser_regex Init")
}

func (This *Parser_regex) Parsing() {
	fmt.Println("call Parser_regex Init")
}
