package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func main() {

	credentialsFile := "gsheets-16Xh3_KIqn9A6MGE2cLijKym7bvf9dR4W8tK3qcFBCng.json"

	spreadsheetID := "16Xh3_KIqn9A6MGE2cLijKym7bvf9dR4W8tK3qcFBCng"

	// Replace 'Sheet1!A1:B10' with the range of cells you want to read.
	readRange := "Form Responses 1!B2:C1000"

	// Initialize Google Sheets API
	ctx := context.Background()
	srv, err := sheets.NewService(ctx, option.WithCredentialsFile(credentialsFile), option.WithScopes(sheets.SpreadsheetsScope))
	if err != nil {
		log.Fatalf("Unable to initialize Sheets API: %v", err)
	}

	// Read data from the specified range
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, readRange).Context(ctx).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	// Print the values from the response
	if len(resp.Values) > 0 {
		fmt.Println("Data from sheet:")
		for _, row := range resp.Values {
			for _, cell := range row {
				fmt.Printf("%v\t", cell)
			}
			fmt.Println()
		}
	} else {
		fmt.Println("No data found.")
	}
}
