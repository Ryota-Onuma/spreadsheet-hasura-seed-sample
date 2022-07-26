package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Ryota-Onuma/sample/google"
	"github.com/Ryota-Onuma/sample/sql"
)

func main() {
	ctx := context.Background()
	client, err := google.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	svc, err := google.NewSpreadSheetService(ctx, client)
	if err != nil {
		log.Fatal(err)
	}
	sheets := svc.GetSheets("1OgASmnTGImaaTbddo7l3ZfaOKoKvJS8VOxlFtz9SIDk")
	sheet := svc.FindSheetByName(sheets, "master_users")

	// CSV生成
	const csvOutputDir = "/go/src/csv/"
	if err := google.SheetToCSV(sheet, csvOutputDir); err != nil {
		log.Fatal(err)
	}

	// SQL生成
	sqlOutputPath := "/go/src/hasura/seeds/"
	csvFilePath := fmt.Sprintf("%s%s.csv", csvOutputDir, sheet.Properties.Title)
	if err := sql.CSVToUpsertSQL(csvFilePath, sqlOutputPath); err != nil {
		log.Fatal(err)
	}
}
