package google

import (
	"context"
	"fmt"
	"net/http"

	"google.golang.org/api/option"
	gsheets "google.golang.org/api/sheets/v4"
)

type SpreadSheetService interface {
	GetSheets(spreadsheetID string) []*gsheets.Sheet
	FindSheetByName(sheets []*gsheets.Sheet, sheetName string) *gsheets.Sheet
}

type ssService gsheets.Service

func NewSpreadSheetService(ctx context.Context, client *http.Client) (SpreadSheetService, error) {
	sheetsService, err := gsheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("GoogleDriveとの接続に失敗しました: %v", err.Error())
	}
	svc := ssService(*sheetsService)
	return &svc, nil
}

func (ss *ssService) GetSheets(spreadsheetID string) []*gsheets.Sheet {
	res, err := ss.Spreadsheets.Get(spreadsheetID).IncludeGridData(true).Do()
	if err != nil {
		return nil
	}
	return res.Sheets
}

func (ss *ssService) FindSheetByName(sheets []*gsheets.Sheet, sheetName string) *gsheets.Sheet {
	var targetSheet *gsheets.Sheet
	for _, sheet := range sheets {
		if sheet.Properties.Title == sheetName {
			targetSheet = sheet
		}
	}
	return targetSheet
}
