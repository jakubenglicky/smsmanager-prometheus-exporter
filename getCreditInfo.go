package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type Account struct {
	Name   string `json:"name"`
	ApiKey string `json:"apiKey"`
	Credit string
}

func LoadCreditInfo(file string) []Account {

	jsonFile, err := os.ReadFile(file)
	if err != nil {
		fmt.Errorf("Error reading file: %v", err)
	}

	accounts := []Account{}
	err = json.Unmarshal(jsonFile, &accounts)
	if err != nil {
		fmt.Errorf("Error unmarshalling: %v", err)
	}

	creditAccounts := []Account{}
	for _, account := range accounts {
		account.Credit = getCreditInfo(account.ApiKey)
		creditAccounts = append(creditAccounts, account)
	}

	return creditAccounts
}

func getCreditInfo(apiKey string) string {

	response, err := http.Get("https://http-api.smsmanager.cz/GetUserInfo?apikey=" + apiKey)

	if err != nil {
		fmt.Errorf("Error sending request: %v", err)
	}

	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Errorf("Error reading response: %v", err)
	}

	credit := strings.Split(string(data), "|")[0]

	return credit
}
