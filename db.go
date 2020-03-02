package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/tealeg/xlsx"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	urls     = make(map[string]*sql.DB)
	xlsFile  *xlsx.File
	cacheDB  *sql.DB
	cacheUrl string
	avgSheet float64
	progress float64
)

func exporting(export Export) {
	xlsFile = xlsx.NewFile()
	avgSheet = float64(100) / float64(len(export.Sheets))
	go func() {
		exportEntry.PromptLabels[0].SetText("Exporting data, please be patient...")
		for _, sheet := range export.Sheets {
			if success, msg := dbOperation(sheet, &export); !success {
				enableExportBtn()
				showMessage(msg)
				return
			}
			cacheUrl = sheet.URL
		}
		savedPath := export.Download + string(os.PathSeparator) + export.FileName
		err := xlsFile.Save(savedPath)
		if err != nil {
			enableExportBtn()
			showMessage(fmt.Sprintf("Save The File Error. Detail: %+v", err))
			return
		}
		_ = cacheDB.Close()
		delete(urls, cacheUrl)
		progress = 0
		export.Window.progressFinish()
		showMessage(fmt.Sprintf("Successful Export. All data has been exported, and the file is stored in the %s%s%s directory",
			export.Download, string(os.PathSeparator), export.FileName))
		enableExportBtn()
	}()
}

func showMessage(msg string) {
	length := len([]rune(msg))
	for i := 0; i < len(prompts); i++ {
		end := (i + 1) * 93
		if end > length {
			end = length
		}
		txt := msg[i*93 : end]
		exportEntry.PromptLabels[i].SetText(strings.TrimSpace(txt))
	}
}

func dbOperation(es SingleSheet, export *Export) (success bool, msg string) {
	tableName := matchSql(es.SQL)[4]
	if !strings.HasPrefix(exportEntry.PromptLabels[0].Text(), "Exporting data") {
		exportEntry.PromptLabels[0].SetText(exportEntry.PromptLabels[1].Text())
	}
	exportEntry.PromptLabels[1].SetText("Current exporting data for the " + tableName + " table...")
	if es.SheetName == "" {
		es.SheetName = tableName
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
		export.Window.setProgress(int(progress))
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
