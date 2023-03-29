package Output

import (
	"LogParsing_regex/Task"
	"database/sql"
	"strconv"
	"strings"

	_ "github.com/alexbrainman/odbc"
)

type Output_Database struct {
}

func New_Output_Database() *Output_Database {
	task := new(Output_Database)
	task.Init()
	return task
}

func (This *Output_Database) Init() {
	Task.LogInst().WriteLog("OUTPUT_DATABASE", "call Output_Database Init")
}

func (This *Output_Database) DataOut(_task *Task_Output, _thread_data *[][]map[string]string) {

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
		Task.LogInst().WriteLog("OUTPUT_DATABASE", "[FAIL] Db Open is fail >> ERR = %s", err)
		return
	}

	if len(_task.Db.Sql.First) > 0 {
		This.ExecSQL(db, &_task.Db.Sql.First)
	}

	bFirst := true
	bulk_sql := ""
	var bulk_len uint64
	// 자료구조 : [고루틴별][라인수별]map[태그별]실데이터
	for thr_index := range *_thread_data {
		for row_index := range (*_thread_data)[thr_index] {

			for i := range _task.Db.Sql.Data {
				sql := ""
				if bFirst == bool(true) {
					sql = _task.Db.Sql.Data[i]
				} else {
					sql = _task.Db.Sql.Values_Sql[i]
				}

				for k, v := range (*_thread_data)[thr_index][row_index] {
					sql = strings.ReplaceAll(sql, k, v)
				}

				if _task.Db.Sql.CanBulkInsert[i] == bool(true) {
					if bFirst == bool(true) {
						bulk_sql = sql
						bFirst = false
						bulk_len += uint64(len(sql))
					} else {
						// 설정 크기를 넘어가면 쿼리 실행
						if _task.Db.Max_Bulksize < (bulk_len + uint64(len(sql))) {
							This.ExecSQL(db, &bulk_sql)

							// 다시 insert 구문부터 구성
							bulk_sql = _task.Db.Sql.Insert_Sql[i] + sql

							bulk_len = uint64(len(bulk_sql))
						} else {
							// value 구문만 추가
							bulk_sql += ","
							bulk_sql += sql

							bulk_len += uint64(len(sql)) + 1
						}
					}
				} else {
					This.ExecSQL(db, &sql)
				}
			}
		}
	}

	// 남은 쿼리 실행
	if len(bulk_sql) > 0 {
		This.ExecSQL(db, &bulk_sql)
	}

	if len(_task.Db.Sql.Finish) > 0 {
		This.ExecSQL(db, &_task.Db.Sql.Finish)
	}
}

func (This *Output_Database) ExecSQL(_db *sql.DB, _sql *string) int64 {

	affect, err := _db.Exec(*_sql)
	if err != nil {
		Task.LogInst().WriteLog("OUTPUT_DATABASE", "[FAIL] ExecSQL is fail >> ERR = %s | SQL = %s", err, _sql)
		return 0
	}

	n, err := affect.RowsAffected()
	if err == nil {
		Task.LogInst().WriteLog("OUTPUT_DATABASE", "[SUCC] rows affected %d", n)
	} else {
		Task.LogInst().WriteLog("OUTPUT_DATABASE", "[FAIL] affected is fail >> ERR = %s | SQL = %s", err, _sql)
	}

	return n
}
