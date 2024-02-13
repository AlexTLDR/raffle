package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"time"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

// Number of winners to select
const numWinners = 5

func main() {

	// For getting the credentials, use the below article
	// https://medium.com/@a.marenkov/how-to-get-credentials-for-google-sheets-456b7e88c430
	credentialsFile := "gsheets-16Xh3_KIqn9A6MGE2cLijKym7bvf9dR4W8tK3qcFBCng.json"
	spreadsheetID := "16Xh3_KIqn9A6MGE2cLijKym7bvf9dR4W8tK3qcFBCng"

	// The range of data to read
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

	// Create a slice to hold the rows. In my case, I will work with a slice of slices of strings instead of a slice of slices of interface{}
	// I believe that for this example, slices of interface{} would complicate things unnecessarily
	var rows [][]string

	// Create a map to hold the unique names
	names := make(map[string]bool)

	// Print the values from the response
	if len(resp.Values) > 0 {
		fmt.Println("Data from sheet:")
		for _, row := range resp.Values {
			// Convert the row to a slice of strings
			strRow := make([]string, len(row))
			isEmpty := true
			for i, cell := range row {
				str, ok := cell.(string)
				if ok && str != "" {
					isEmpty = false
					strRow[i] = str
				}
			}
			// If the row is not empty and the name is unique, add it to the slice
			if !isEmpty && !names[strRow[0]] {
				names[strRow[0]] = true
				rows = append(rows, strRow)
				for _, cell := range strRow {
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

	// Seed the random number generator
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Clearing the print statements from the terminal
	clear()

	// Check if there are enough rows
	if len(rows) >= numWinners {
		// Select the winners
		for i := 0; i < numWinners; i++ {
			fmt.Println("Drum roll please...")
			time.Sleep(5 * time.Second)

			index := r.Intn(len(rows))
			winner := rows[index]
			fmt.Printf("Winner %d is: username: %s, email: %s\n", i+1, winner[0], winner[1])

			// Remove the winner from the rows slice
			rows = append(rows[:index], rows[index+1:]...)
		}
	} else {
		fmt.Println("Not enough rows to select", numWinners, "winner(s) from.")
	}
}

// clear clears the terminal
func clear() {
	time.Sleep(5 * time.Second)
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}
