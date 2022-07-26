package google

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func NewClient(ctx context.Context) (*http.Client, error) {
	const gcsAuthJSONFilePath = "/go/src/google/gcpClientSecret.json"
	b, err := os.ReadFile(gcsAuthJSONFilePath)
	if err != nil {
		return nil, fmt.Errorf("googleとの接続に失敗しました: %v", err.Error())
	}
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets.readonly")
	if err != nil {
		return nil, fmt.Errorf("googleとの接続に失敗しました: %v", err.Error())
	}

	token, err := getTokenFromWeb(config)
	if err != nil {
		return nil, fmt.Errorf("GoogleDriveとの接続に失敗しました: %v", err.Error())
	}

	client := config.Client(context.Background(), token)
	if err != nil {
		return nil, fmt.Errorf("GoogleDriveとの接続に失敗しました: %v", err.Error())
	}
	return client, nil
}

func getTokenFromWeb(config *oauth2.Config) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("\nブラウザで以下のURLを開いてauthentication codeを取得してください: \n\n%v\n\n", authURL)
	var authCode string
	fmt.Print("authorization code: ")
	if _, err := fmt.Scan(&authCode); err != nil {
		return nil, fmt.Errorf("authorization codeを読み取れませんでした\n: %v", err)
	}

	token, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		return nil, fmt.Errorf("トークンを取得できませんでした: %v\n", err)
	}
	fmt.Printf("\n取得したトークン> %s\n\n", token)
	return token, nil
}
