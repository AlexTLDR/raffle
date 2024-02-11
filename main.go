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

	// Create a slice to hold the rows
	var rows [][]string

	// Print the values from the response
	if len(resp.Values) > 0 {
		fmt.Println("Data from sheet:")
		for _, row := range resp.Values {
			strRow := make([]string, len(row))
			// Check if the row is empty
			isEmpty := true
			for i, cell := range row {
				str, ok := cell.(string)
				if ok && str != "" {
					isEmpty = false
					strRow[i] = str
				}
			}
			// If the row is not empty, add it to the slice
			if !isEmpty {
				rows = append(rows, strRow)
				for _, cell := range row {
					fmt.Printf("%v\t", cell)
				}
				fmt.Println()
			}
		}
	} else {
		fmt.Println("No data found.")
	}

	// Print the rows slice
	fmt.Println(rows)
}
