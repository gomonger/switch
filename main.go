package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var list *tview.List
var app *tview.Application

func Picked() {
	//list
	itemIdx := list.GetCurrentItem()
	itemTxt, note := list.GetItemText(itemIdx)
	app.Stop()
	fmt.Println("picked item:", itemIdx, itemTxt, note)
}

func RunList() {
	app = tview.NewApplication()

	list = tview.NewList().
		AddItem("Server 1", "Item 1", 'a', Picked).
		AddItem("Server 2", "Item 2", 'b', Picked).
		AddItem("Server 3", "Item 3", 'c', Picked).
		AddItem("Server 4", "Item 4", 'd', Picked).
		AddItem("Quit", "press to exit", 'q', func() {
			app.Stop()
		})

	list.ShowSecondaryText(false)

	list.SetBorder(true).SetTitle("Switch").SetTitleAlign(tview.AlignLeft)

	if err := app.SetRoot(list, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

func TableSetupFilter(table *tview.Table, term string, data [][]string) {
	if term == "clear" {
		table.Clear()
		return
	}
	table.Clear()
	time.Sleep(10 * time.Millisecond)
	tableRow := 0
	// Add some rows and columns to the table.
	for row := 0; row < len(data); row++ {
		color := tcell.ColorWhite

		if term != "" && row > 0 {
			rowText := fmt.Sprintf("%s %s %s", data[row][0], data[row][1], data[row][2])
			if !strings.Contains(rowText, term) {
				continue
			}
		}

		for col := 0; col < 3; col++ {
			if row == 0 {
				color = tcell.ColorYellow
			}
			align := tview.AlignLeft
			if col > 0 {
				align = tview.AlignCenter
			}
			cell := tview.NewTableCell(fmt.Sprintf("%s", data[row][col])).
				SetAlign(align).
				SetTextColor(color).
				SetSelectable(row != 0)

			table.SetCell(tableRow, col, cell)
		}
		tableRow++
	}

}

func RunTable(inputDataFile string) {

	var data [][]string

	dataType := "servers"
	selectFormat := "server: %-15s %-10s %-10s"

	if inputDataFile != "" {

		if strings.Contains(inputDataFile, "users") {
			dataType = "users"
			selectFormat = "user: %-6s,%-10s,%-10s"
		}
		fileData, err := os.ReadFile(inputDataFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading file: %v\n", err)
			os.Exit(1)
		}
		lines := strings.Split(string(fileData), "\n")
		for _, line := range lines {
			if line == "" {
				continue
			}
			parts := strings.Split(line, ",")
			data = append(data, parts)
		}

	} else {

		data = [][]string{
			{"Server Name", "Environment", "Appication"},
			{"servera", "prod", "acme"},
			{"serverb", "test", "acme"},
			{"serverc", "dev", "acme"},
			{"server001", "prod", "fudge"},
			{"server002", "test", "fudge"},
			{"server003", "dev", "fudge"},
			{"server004", "prod", "zbox"},
			{"server005", "test", "zbox"},
			{"server006", "dev", "zbox"},
		}
	}

	app = tview.NewApplication()

	flex := tview.NewFlex().SetDirection(tview.FlexRow)

	table := tview.NewTable().
		SetBorders(true).
		SetFixed(1, 1)

	searchBox := tview.NewInputField().
		SetLabel("Search: ").
		SetFieldWidth(30)

	searchBox.SetDoneFunc(func(key tcell.Key) {
		// Implement search functionality here if needed
		if key == tcell.KeyEnter {
			TableSetupFilter(table, searchBox.GetText(), data)
			app.SetFocus(table)
		}

	})

	searchBox.SetChangedFunc(func(text string) {
		TableSetupFilter(table, text, data)
	})

	searchBox.SetBorder(true).SetTitle("Switch").SetTitleAlign(tview.AlignLeft)

	TableSetupFilter(table, "", data)

	table.SetBorder(true).SetTitle(dataType).SetTitleAlign(tview.AlignLeft)
	table.SetSelectable(true, false).SetSelectedFunc(func(row int, column int) {
		app.Stop()

		cell0 := table.GetCell(row, 0).Text
		cell1 := table.GetCell(row, 1).Text
		cell2 := table.GetCell(row, 2).Text
		fmt.Printf(selectFormat, cell0, cell1, cell2)
	})

	flex.AddItem(searchBox, 5, 1, false)
	flex.AddItem(table, 0, 1, true)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			if app.GetFocus() == table {
				app.SetFocus(searchBox)
				return nil
			} else if app.GetFocus() == searchBox {
				app.SetFocus(table)
				return nil
			}
		}
		return event
	})

	if err := app.SetRoot(flex, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

func main() {
	var tableMode bool
	var tableDataFile string
	flag.BoolVar(&tableMode, "t", false, "table mode")
	flag.StringVar(&tableDataFile, "f", "", "table data file")
	flag.String("h", "", "help")
	flag.Parse()

	if tableMode {
		RunTable(tableDataFile)
		return
	}
	RunList()
}
