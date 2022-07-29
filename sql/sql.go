package sql

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func CSVToUpsertSQL(csvPath string, outputDir string) error {
	records, err := readCSV(csvPath)
	if err != nil {
		return err
	}

	csvFileName := filepath.Base(csvPath)
	csvExtName := filepath.Ext(csvFileName)
	sqlExtName := ".sql"
	sqlFileName := strings.Replace(csvFileName, csvExtName, sqlExtName, 1)
	sqlFilePath := outputDir + sqlFileName

	f, err := os.Create(sqlFilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := os.Chmod(sqlFilePath, 0666); err != nil {
		return err
	}

	data := NewUpsertContentData()
	tableName := strings.Replace(csvFileName, csvExtName, "", 1)
	data.RegisterContentMap(tableName, records)

	tmpl, err := template.New("tmpl").Parse(data.contentsTemplate)
	if err != nil {
		return err
	}

	var doc bytes.Buffer
	if err = tmpl.Execute(&doc, data); err != nil {
		return err
	}

	str := doc.String()

	if _, err := fmt.Fprintf(f, "%v", str); err != nil {
		return err
	}
	fmt.Printf("%sの作成が完了しました\n", sqlFileName)
	return nil
}

type UpsertContentData struct {
	TableName        string
	ColumnNames      string
	Values           string
	Pk               string
	UpdateSetting    string
	contentsTemplate string
}

func NewUpsertContentData() *UpsertContentData {
	// upsert句のテンプレ
	contentsTemplate := `INSERT INTO "{{ .TableName }}" ({{ .ColumnNames }})
VALUES
{{ .Values }}
ON CONFLICT ({{ .Pk }})
DO UPDATE SET
{{ .UpdateSetting }}
;
`
	return &UpsertContentData{
		contentsTemplate: contentsTemplate,
	}
}

func (c *UpsertContentData) RegisterContentMap(tableName string, records [][]string) {
	// ヘッダーの置き換え
	// 1行目はヘッダーのはず

	columnNames := strings.Join(records[0], ",")
	var values []string
	for i, record := range records {
		// ヘッダーはスキップ
		if i < 2 {
			continue
		}
		var texts []string
		for _, v := range record {
			if v == "" {
				texts = append(texts, "null")
			} else {
				texts = append(texts, v)
			}

		}
		var splitToken string
		if i == len(records)-1 {
			splitToken = ""
		} else {
			splitToken = ","
		}
		values = append(values, fmt.Sprintf("  (%s)%s", strings.Join(texts, ","), splitToken))
	}

	var updateSetting []string
	for i, v := range records[0] {
		var splitToken string
		if i == len(records[0])-1 {
			splitToken = ""
		} else {
			splitToken = ","
		}
		updateSetting = append(updateSetting, fmt.Sprintf("  %q = EXCLUDED.%s%s", v, v, splitToken))
	}

	var pkColumns []string

	for i, columnName := range records[0] {
		if records[2][i] == "TRUE" && columnName != "" {
			pkColumns = append(pkColumns, columnName)
		}
	}
	c.TableName = tableName
	c.ColumnNames = columnNames
	c.Values = strings.Join(values, "\n")
	c.Pk = strings.Join(pkColumns, ",")
	c.UpdateSetting = strings.Join(updateSetting, "\n")
}

func readCSV(filePath string) ([][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer file.Close()
	reader := csv.NewReader(file)
	var records [][]string
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	return records, nil
}
