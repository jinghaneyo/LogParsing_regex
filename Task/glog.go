package Task

import (
	"fmt"
	"os"
	"sync/atomic"
	"time"
)

var g_EnableConsole int32

func SetLog_EnableConsole(_Enable bool) {
	if _Enable == bool(true) {
		atomic.StoreInt32(&g_EnableConsole, 1)
	} else {
		atomic.StoreInt32(&g_EnableConsole, 0)
	}
}

func WriteLog(_filename string, _format string, _msg ...interface{}) {

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

	// call := CallFunc()
	// if len(call) > 20 {
	// 	call = fmt.Sprintf("[%s][%20s] ", rst, call[len(call)-20:])
	// } else {
	// 	call = fmt.Sprintf("[%s][%20s] ", rst, call)
	// }
	call := fmt.Sprintf("[%s]", rst)

	msg := fmt.Sprintf(_format, _msg...)
	if 0 < len(msg) {
		if g_EnableConsole == 1 {
			fmt.Printf("%s%s\n", call, msg)
		}

		_, err = fmt.Fprint(logfile, call, msg, "\n")
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
