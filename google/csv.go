package google

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	gsheets "google.golang.org/api/sheets/v4"
)

func SheetToCSV(s *gsheets.Sheet, csvOutputDir string) error {
	var records [][]string
	if len(s.Data[0].RowData) == 0 {
		return nil
	}
	fileName := s.Properties.Title

	// headerの名前を取得
	var columnNames []string
	const headerRowIndex = 0
	for i, value := range s.Data[0].RowData[headerRowIndex].Values {
		if i == 0 {
			continue
		}
		columnNames = append(columnNames, value.FormattedValue)
	}
	records = append(records, columnNames)

	// column型を取得
	var columnTypes []string
	const columnTypeRowIndex = 2
	for i, columnType := range s.Data[0].RowData[columnTypeRowIndex].Values {
		if i == 0 {
			continue
		}
		columnTypes = append(columnTypes, columnType.FormattedValue)
	}
	if len(columnNames) != len(columnTypes) {
		return errors.New("ヘッダーかカラムタイプの長さが一致しません")
	}

	// pkの設定
	var pks []string
	const pkRowIndex = 3
	for i, columnType := range s.Data[0].RowData[pkRowIndex].Values {
		if i == 0 {
			continue
		}
		pks = append(pks, columnType.FormattedValue)
	}
	records = append(records, pks)

	// レコードの差し替え

	// レコードが入っているのは4行目以降
	const dataBeginRowIndex = 4
	for _, row := range s.Data[0].RowData[dataBeginRowIndex:] {
		if len(row.Values) == 0 {
			continue
		}

		const dataBeginColumnIndex = 1
		// 先頭がfalseの行はCSV出力しない
		if len(row.Values) > 0 && row.Values[dataBeginColumnIndex-1].FormattedValue != "TRUE" {
			continue
		}
		var record []string

		values := row.Values[dataBeginColumnIndex:]

		for i, value := range values {
			// 改行あるとCSVが壊れるので余計な改行コードを削除する(任意に文字列データに改行入れたいときはエスケープ文字構文に対応させる)
			rowValue := strings.Replace(value.FormattedValue, "\n", "", -1)
			v := generateCSVValue(columnTypes[i], rowValue)
			// {{ hoge }} のマスタッシュ部分だけ削除
			v = strings.TrimPrefix(v, "{{ ")
			v = strings.TrimSuffix(v, " }}")
			record = append(record, v)
		}
		records = append(records, record)
	}

	fileFullPath := fmt.Sprintf("%s%s", csvOutputDir, fileName)
	if err := writeCSV(fileFullPath, records); err != nil {
		return err
	}
	fmt.Printf("%v.csvの書き込みを完了しました\n", fileName)
	return nil
}

func generateCSVValue(dataType string, v string) string {
	var result string
	switch dataType {
	case "integer":
		num, err := strconv.Atoi(v)
		if err != nil || v == "" {
			// integerにできない値がvalueに入っていたら0で返す
			result = "{{ 0 }}"
		} else {
			result = fmt.Sprintf("{{ %v }}", num)
		}

	case "string":
		if v == "" {
			result = "{{ '' }}"
		} else {
			result = fmt.Sprintf("{{ '%s' }}", v)
		}

	case "boolean":
		val, err := strconv.ParseBool(v)
		if err != nil || v == "" {
			// 変換できなかった時はfalseにする
			result = "{{ false }}"
		} else {
			result = fmt.Sprintf("{{ %v }}", val)
		}
	}
	return result
}

func writeCSV(filePath string, records [][]string) error {
	file, err := os.Create(fmt.Sprintf("%s.csv", filePath))
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	for _, record := range records {
		writer.Write(record)
	}
	writer.Flush()
	return nil
}
