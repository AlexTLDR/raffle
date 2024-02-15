package main

import (
	"bufio"
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

// Number of winners to select
const numWinners = 6
const specialPrize = "Ultimate Go Bundle"

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
					// Convert the name to lowercase
					if i == 0 {
						str = strings.ToLower(str)
					}
					strRow[i] = str
				}
			}
			// If the row is not empty and the name is unique, add it to the slice
			if !isEmpty && !names[strRow[0]] {
				names[strRow[0]] = true
				rows = append(rows, strRow)
				// Print only the name
				fmt.Println(strRow[0])
			}
		}
	} else {
		fmt.Println("No data found.")
	}

	// Print the rows slice
	//fmt.Println(rows)

	// Seed the random number generator
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Clearing the print statements from the terminal
	clear()

	// Check if there are enough rows
	if len(rows) >= numWinners {
		// Create a CSV file
		file, err := os.Create("winners.csv")
		if err != nil {
			log.Fatalf("Failed to create file: %v", err)
		}
		defer file.Close()

		writer := csv.NewWriter(file)
		defer writer.Flush()

		// Select the winners
		for i := 0; i < numWinners; i++ {
			fmt.Println("Drum roll please...")
			time.Sleep(5 * time.Second)

			index := r.Intn(len(rows))
			winner := rows[index]

			// Check if this is the last winner
			if i == numWinners-1 {
				fmt.Printf("The %s goes to: %s\n\n", specialPrize, winner[0])
				// Add a mark to indicate that this is a special winner
				winner = append(winner, specialPrize)
			} else {
				fmt.Printf("Winner %d is: %s\n\n", i+1, winner[0])
			}

			// Write the winner to the CSV file
			err := writer.Write(winner)
			if err != nil {
				log.Fatalf("Failed to write to file: %v", err)
			}

			// Remove the winner from the rows slice
			rows = append(rows[:index], rows[index+1:]...)
		}
	} else {
		fmt.Println("Not enough rows to select", numWinners, "winner(s) from.")
	}
}

// clear clears the terminal
func clear() {
	//time.Sleep(5 * time.Second)

	fmt.Print("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

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
