package Task

import (
	"fmt"
	"os"
	"time"
)

type Log struct {
	bEnableConsole bool
	bEnableLogDate bool
}

var g_Instance *Log

func GetInst() *Log {
	if g_Instance == nil {
		g_Instance = new(Log)
		g_Instance.bEnableConsole = true
		g_Instance.bEnableLogDate = true
	}
	return g_Instance
}

func (This *Log) Enable_Console(_Enable bool) {
	This.bEnableConsole = _Enable
}

func (This *Log) Enable_LogDate(_Enable bool) {
	This.bEnableLogDate = _Enable
}

func (This *Log) WriteLog(_filename string, _format string, _msg ...interface{}) {

	t := time.Now()
	rst := t.Format("2006-01-02 15:04:05")
	date := rst[:10]
	hour := rst[11:13]

	folderPath := fmt.Sprintf("Log/%s", date)
	os.MkdirAll(folderPath, os.ModePerm)
	fname := fmt.Sprintf("%s/%s[%s]%s.log", folderPath, date, hour, _filename)
	logfile, err := os.OpenFile(fname, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return
	}

	msg := fmt.Sprintf(_format, _msg...)
	if 0 < len(msg) {

		if This.bEnableLogDate == bool(true) {
			call := fmt.Sprintf("[%s]", rst)
			_, err = fmt.Fprint(logfile, call, msg, "\n")

			if This.bEnableConsole == bool(true) {
				fmt.Printf("%s%s\n", call, msg)
			}
		} else {
			_, err = fmt.Fprint(logfile, msg, "\n")

			if This.bEnableConsole == bool(true) {
				fmt.Printf("%s\n", msg)
			}
		}

		if err != nil {
			logfile.Close()
			return
		}
	} else {
		_, err = fmt.Fprint(logfile, "\n")
		if err != nil {
			logfile.Close()
			return
		}
	}
	logfile.Close()
}
