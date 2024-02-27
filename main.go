package main

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"go-aws-pass/internal/db"
	sm "go-aws-pass/internal/secretsmanager"
	ssm "go-aws-pass/internal/systemsmanager"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Person struct {
	Email string
	Name  string
	Age   int
}

func ReadFromConsole(consoleType string) string {
	fmt.Println("Type /exit to exit")
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		text := scanner.Text()

		if text == "/exit" {
			os.Exit(0)
		}

		if consoleType == "table" || consoleType == "ssm" {
			return text
		}

		fmt.Println(scanner.Text()) // Println will add back the final '\n'
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

	return ""
}

func main() {
	table := "persons"

	featureCreateDatabase := false

	if featureCreateDatabase {
		d := db.CreateDatabase(table)
		email := "kevin@homelabwithkevin.com"

		p := Person{email, "kevin", 69}
		db.WriteToDatabase(d, table, db.Person(p))
		db.ReadFromDatabase(d, table, email)
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		panic(err)
	}

	// Get user Input to Create SSM Parameter if flag is true
	featureSSM := false
	if featureSSM {
		fmt.Printf("\nWhat SSM parameter name?\n")
		parameterName := ReadFromConsole("ssm")

		fmt.Printf("\nWhat SSM value?\n")
		parameterValue := ReadFromConsole("ssm")

		ssm.CreateParameter(cfg, parameterName, parameterValue)
	}

	// Retrieve Secret if Feature Flag is True
	featureSecretsManager := false
	if featureSecretsManager {
		result := sm.ListSecrets(cfg)

		for _, v := range result.SecretList {
			name := string(*v.Name)

			secrets := sm.GetSecret(cfg, name)

			fmt.Printf("\n----------------------------\n")
			fmt.Printf("Secret Name: %s", name)
			fmt.Printf("\n----------------------------\n")

			for i, v := range secrets {
				fmt.Printf("Secret Key: %s %v", i, "\n")
				fmt.Printf("Secret Value: %s %v", v, "\n\n")
			}
		}
	}

	featureTView := true

	if featureTView {
		app := tview.NewApplication()

		grid := tview.NewGrid().
			SetRows(3, 0, 3).
			SetColumns(30, 0, 30).
			SetBorders(true)

		// Header
		searchPrefx := tview.NewInputField().SetLabel("Enter a search prefix: ")

		searchPrefx.SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEsc {
				app.Stop()
			}
		})
		grid.AddItem(searchPrefx, 0, 0, 1, 3, 0, 0, true)

		// Main
		table := tview.NewTable().
			SetFixed(1, 6).
			SetSelectable(true, false)
		table.
			SetBorder(true).
			SetTitle("  Browser ").
			SetBorderPadding(0, 0, 1, 1)

		headers := [][]string{{"Name", "Type", "Version", "Last modified"}}
		for row, line := range headers {
			for col, text := range line {
				cell := tview.NewTableCell(text)
				table.SetCell(row+1, col, cell)
			}
		}

		// Footer
		leftFooterItem := tview.NewTextView()
		grid.AddItem(leftFooterItem, 2, 0, 1, 3, 0, 0, false)
		leftFooterItem.SetText("ESC/CTRL+C=Exit | TAB=Switch focus | ENTER=See details C=Copy value to clipboard | X=Copy name to clipboard")

		// Layout for screens narrower than 100 cells (menu and side bar are hidden).
		// Layout for screens wider than 100 cells.
		grid.AddItem(table, 1, 0, 1, 3, 0, 0, false)

		if err := app.SetRoot(grid, true).SetFocus(grid).Run(); err != nil {
			panic(err)
		}
	}
}
