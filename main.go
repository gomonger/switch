package main

import (
	"flag"
	"fmt"

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

func RunTable() {

	data := [][]string{
		{"Server Name", "Environment", "Appication"},
		{"servera", "prod", "acme"},
		{"serverb", "test", "acme"},
		{"serverc", "dev", "acme"},
		{"server001", "prod", "fudge"},
		{"server002", "test", "fudge"},
		{"server003", "dev", "fudge"},
	}

	app = tview.NewApplication()

	table := tview.NewTable().
		SetBorders(true).
		SetFixed(1, 1)

	// Add some rows and columns to the table.
	for row := 0; row < len(data); row++ {
		color := tcell.ColorWhite
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

			table.SetCell(row, col, cell)
		}
	}

	table.SetBorder(true).SetTitle("Table Example").SetTitleAlign(tview.AlignLeft)
	table.SetSelectable(true, false).SetSelectedFunc(func(row int, column int) {
		app.Stop()
		server := data[row][0]
		env := data[row][1]
		fmt.Printf("selected server: %s Env: %s\n", server, env)
	})

	if err := app.SetRoot(table, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

func main() {
	var tableMode bool
	flag.BoolVar(&tableMode, "t", false, "table mode")
	flag.String("h", "", "help")
	flag.Parse()

	if tableMode {
		RunTable()
		return
	}
	RunList()
}
