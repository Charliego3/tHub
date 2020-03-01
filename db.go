package main

import (
	"database/sql"
	"fmt"
	"github.com/ProtonMail/ui"
	_ "github.com/go-sql-driver/mysql"
	"github.com/tealeg/xlsx"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	window   *ui.Window
	urls     = make(map[string]*sql.DB)
	xlsFile  *xlsx.File
	cacheDB  *sql.DB
	cacheUrl string
	avgSheet float64
	progress float64
)

func exporting(export Export) {
	window = export.Window.Window
	xlsFile = xlsx.NewFile()
	avgSheet = float64(100) / float64(len(export.Sheets))
	go func() {
		for _, sheet := range export.Sheets {
			if success, msg := dbOperation(sheet, &export); !success {
				enableExportBtn()
				showErr(msg)
				return
			}
			cacheUrl = sheet.URL
		}
		savedPath := export.Download + string(os.PathSeparator) + export.FileName
		err := xlsFile.Save(savedPath)
		if err != nil {
			enableExportBtn()
			showErr(fmt.Sprintf("Save The File Error. Detail: %+v", err))
			return
		}
		export.Window.progressFinish()
		_ = cacheDB.Close()
		delete(urls, cacheUrl)
		progress = 0
		enableExportBtn()
	}()
}

// TODO show err text
func showErr(msg string) {
	for i := 0; i < int(math.Ceil(float64(len([]rune(msg))/85))); i++ {
		exportEntry.PromptLabels[i].SetText(msg[:(i+1)*85])
		exportEntry.PromptLabels[i].Show()
	}
}

func dbOperation(es SingleSheet, export *Export) (success bool, msg string) {
	if es.SheetName == "" {
		es.SheetName = matchSql(es.SQL)[4]
	}
	db, msg := getConnection(es.URL)
	if db == nil {
		return false, msg
	}
	if cacheDB != nil && db != cacheDB {
		_ = cacheDB.Close()
		delete(urls, cacheUrl)
	}
	cacheDB = db

	sheet, err := xlsFile.AddSheet(es.SheetName)
	if err != nil {
		msg = fmt.Sprintf("Create Sheet failed, SheetName is %s, Error Detail: %+v", es.SheetName, err)
		return
	}

	//Sheet Args
	var args []interface{}
	if es.Args != "" {
		for _, arg := range strings.Split(es.Args, ",") {
			args = append(args, strings.TrimSpace(arg))
		}
	}

	queryCount := getCount(es.SQL, args, db)
	preLine := avgSheet / float64(queryCount)

	// DB Query
	rows, msg := execute(es.SQL, db, args...)
	if rows == nil {
		return false, msg
	}
	columns, err := rows.Columns()
	if err != nil {
		_ = rows.Close()
		_ = db.Close()
		return false, fmt.Sprintf("Get columns error. Detail: %+v", err)
	}

	// Completed Titles
	var ts []string
	if es.Titles != "" {
		for _, t := range strings.Split(es.Titles, ",") {
			ts = append(ts, strings.TrimSpace(t))
		}
	}
	ts = append(ts, columns[len(ts):]...)

	addTitle(sheet, ts)

	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	var count, page = 0, 0
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			msg = fmt.Sprintf("Scan the result error. Detail: %+v", err)
			return
		}
		row := sheet.AddRow()
		for _, col := range values {
			var value string
			switch col.(type) {
			case []byte, sql.RawBytes:
				value = string(col.([]byte))
			case nil:
				value = "NULL"
			default:
				value = fmt.Sprintf("%v", col)
			}
			row.AddCell().Value = value
		}
		progress += preLine
		if progress >= 100 {
			progress = 99
		}
		fmt.Printf("Progress: %v\n", progress)
		export.Window.setProgress(int(progress))
		//time.Sleep(time.Millisecond * 100)
		count++
		if count%65534 == 0 {
			page++
			sheet, err = xlsFile.AddSheet(es.SheetName + "_" + strconv.Itoa(page))
			if err != nil {
				msg = fmt.Sprintf("Add Sheet Error. Detail: %+v", err)
				return
			}
			addTitle(sheet, ts)
		}
	}

	success = true
	return
}

func addTitle(sheet *xlsx.Sheet, ts []string) {
	// Add Sheet Title
	titleRow := sheet.AddRow()
	for _, t := range ts {
		titleRow.AddCell().Value = t
	}
}

func getCount(query string, args []interface{}, db *sql.DB) int {
	var count int
	_ = db.QueryRow("select count(1) from ("+query+") t", args...).Scan(&count)
	return count
}

func execute(querySql string, db *sql.DB, args ...interface{}) (*sql.Rows, string) {
	var rows *sql.Rows
	var err error
	if args != nil && len(args) > 0 {
		rows, err = db.Query(querySql, args...)
	} else {
		rows, err = db.Query(querySql)
	}
	if err != nil {
		return nil, fmt.Sprintf("Execute SQL Error. Detail: %+v", err)
	}
	return rows, ""
}

func matchSql(sql string) []string {
	regex := regexp.MustCompile(SQLRegex)
	return regex.FindStringSubmatch(sql)
}

func getConnection(url string) (*sql.DB, string) {
	if db, ok := urls[url]; ok {
		return db, ""
	}
	errMsg := "Database connect failed, please check URL is correct, Error Detail: %+v"
	connection, err := sql.Open("mysql", url)
	if err != nil {
		return nil, fmt.Sprintf(errMsg, err)
	}
	err = connection.Ping()
	if err != nil {
		return nil, fmt.Sprintf(errMsg, err)
	}
	urls[url] = connection
	return connection, ""
}

func showExportError(msg string) {
	ui.MsgBoxError(window, "Error exporting Excel", msg)
}
