package Input

import (
	"LogParsing_regex/Task"
	"database/sql"
	"strconv"
	"strings"
	"time"

	_ "github.com/alexbrainman/odbc"
)

type Input_Database struct {
}

func New_Input_Database() *Input_Database {
	task := new(Input_Database)
	task.Init()
	return task
}

func (This *Input_Database) Init() {
	Task.LogInst().WriteLog("Input_Database", "call Input_Database Init")
}

func (This *Input_Database) Load(_task *Task_Input) *[][]map[string]string {

	var db *sql.DB
	var err error

	connStr := "DSN=${DSN};UID=${ID};PWD=${PWD};DB=${DB};Port=${PORT};CHARSET=utf8"
	connStr = strings.ReplaceAll(connStr, "${DSN}", _task.Db.Odbc)
	connStr = strings.ReplaceAll(connStr, "${ID}", _task.Db.Id)
	connStr = strings.ReplaceAll(connStr, "${PWD}", _task.Db.Pwd)
	connStr = strings.ReplaceAll(connStr, "${DB}", _task.Db.Database)
	connStr = strings.ReplaceAll(connStr, "${PORT}", strconv.Itoa(_task.Db.Port))

	db, err = sql.Open("odbc", connStr)

	if err != nil {
		Task.LogInst().WriteLog("Input_Database", "[FAIL] Db Open is fail >> ERR = %s", err)
		return nil
	}

	rows, err := db.Query(_task.Db.Sql_Select)
	if err != nil {
		Task.LogInst().WriteLog("Input_Database", "[FAIL] SQL is fail >> ERR = %s | SQL = %s", err, _task.Db.Sql_Select)
		return nil
	}
	defer rows.Close()

	thread_data_set := make([][]map[string]string, 1)
	data_set := make([]map[string]string, 0)

	//ilog.Writelog(_target.LogName, ilog.I, "[Get_RowList] STEP >> %s", _step.Desc)

	colType, _ := rows.ColumnTypes()

	colTag := ""

	Row_Count := 0
	for rows.Next() == bool(true) {

		col_count := len(colType)
		ptr := make([]interface{}, col_count)
		cols := make([]interface{}, col_count)
		for i := range ptr {
			ptr[i] = &cols[i]
		}

		err := rows.Scan(ptr...)
		if err != nil {
			//ilog.Writelog(_target.LogName, ilog.W, "[Get_RowList] Read table Fail Rows (err:%v)", err)
			thread_data_set[0] = data_set
			return &thread_data_set
		} else {

			Row_Count++
			col_val := ""
			temp := ""
			var err2 error

			mapData := make(map[string]string)

			for i, d := range cols {

				// 포맷이 json 형태이므로 내용중 쌍따옴표가 존재하면 역슬레쉬(\)을 붙여주자
				temp = *This.ConvertStr(d)
				temp = strings.ReplaceAll(temp, "\"", "\\\"")

				col_val = temp

				if err2 == nil {
					colTag = "${" + colType[i].Name() + "}"
					mapData[colTag] = col_val
				}
			}

			data_set = append(data_set, mapData)
		}
	}

	// if len(step_data_set) == 0 {
	// 	ilog.Writelog(_target.LogName, ilog.I, "[Get_RowList] No Data !!")
	// }

	// err := rows.Scan(&id, &name)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(id, name)

	thread_data_set[0] = data_set
	return &thread_data_set
}

func (This *Input_Database) ConvertStr(_v interface{}) *string {

	tempVal := ""

	switch v := _v.(type) {
	case int:
		tempVal = strconv.Itoa(v)
	case int8:
		tempVal = strconv.FormatInt(int64(v), 10)
	case int16:
		tempVal = strconv.FormatInt(int64(v), 10)
	case int32:
		tempVal = strconv.FormatInt(int64(v), 10)
	case int64:
		tempVal = strconv.FormatInt(int64(v), 10)
	case uint:
		tempVal = strconv.FormatInt(int64(v), 10)
	case uint8:
		tempVal = strconv.FormatInt(int64(v), 10)
	case uint16:
		tempVal = strconv.FormatInt(int64(v), 10)
	case uint32:
		tempVal = strconv.FormatInt(int64(v), 10)
	case uint64:
		tempVal = strconv.FormatInt(int64(v), 10)
	case float32:
		tempVal = strconv.FormatFloat(float64(v), 'f', -1, 64)
	case float64:
		tempVal = strconv.FormatFloat(float64(v), 'f', -1, 64)
	case []uint8:
		tempVal = string(v)
	case time.Time:
		tempVal = v.Format("2006-01-02 15:04:05")
	default:
		tempVal = v.(string)
	}

	return &tempVal
}
