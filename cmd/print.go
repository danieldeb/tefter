package cmd

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
	"log"
	"sort"
	"strconv"
	"strings"
)

var printCmd = &cobra.Command{
	Use:   "print",
	Short: "Print notes based on given ids",
	Long: "There are 4 ways to print a set of notes" +
		" 1) Give a comma separated list of note ids" +
		" 2) Give a comma separated list of notebook titles" +
		" 3) Give a comma separated list of tags," +
		" 4) If -a or --all flag is set all notes will be printed" +
		"Press Esc to exit print mode",
	Example: "print -i 1,2,... -n notebook1,notebook2,... -t tag1,tag2,... ",
	Run: func(cmd *cobra.Command, args []string) {

		ids, _ := cmd.Flags().GetIntSlice("ids")
		notebookTitles, _ := cmd.Flags().GetStringSlice("notebook")
		tags, _ := cmd.Flags().GetStringSlice("tags")
		printAll, _ := cmd.Flags().GetBool("all")

		notes := collectNotesFromDB(ids, notebookTitles, tags, printAll)
		jNotes, err := transformNotes2JSONNotes(noteMap2Slice(notes))
		if err != nil {
			log.Panicln(err)
		}
		printNotes2Terminal(jNotes)
	},
}

func init() {
	rootCmd.AddCommand(printCmd)
	printCmd.Flags().IntSliceP("ids", "i", []int{}, "Comma separated list of note ids.")
	printCmd.Flags().StringSlice("tags", []string{}, "Comma-separated tags of note.")
	printCmd.Flags().StringSliceP("notebook", "n", []string{}, "Comma separated list of notebook titles")
	printCmd.Flags().BoolP("all", "a", false, "Print all notes")
}

func printNotes2Terminal(jNotes []*jsonNote) {
	uiApp := createUI(jNotes)

	if err := uiApp.Run(); err != nil {
		log.Fatal(err)
	}
}

func createUI(jNotes []*jsonNote) *tview.Application {
	if len(jNotes) == 0 {
		return nil
	}
	//Sort By date (descending)
	sort.Slice(jNotes, func(i, j int) bool {
		return jNotes[i].LastUpdated.After(jNotes[j].LastUpdated)
	})
	app := tview.NewApplication()
	notesTable := tview.NewTable()
	notesTable.SetFixed(1, 0)
	notesTable.SetSelectable(true, false)
	notesFlex := tview.NewFlex()
	notesFlex.SetDirection(tview.FlexRow)
	numberOfVisibleRows := 5
	notesFlex.AddItem(notesTable, numberOfVisibleRows, 1, true)

	memo := tview.NewTextView()
	memo.SetBorder(true)
	memo.SetWordWrap(true)
	memo.SetBorderPadding(0, 0, 0, 0)
	memo.SetTextAlign(tview.AlignLeft)
	memo.SetText(jNotes[0].Memo)

	dates := tview.NewTable()
	dates.SetSelectable(false, false)
	dates.SetCell(0, 0, &tview.TableCell{Text: "Created: ", Align: tview.AlignLeft, Color: tcell.ColorRed, Expansion: 1})
	dates.SetCell(0, 1, &tview.TableCell{Text: jNotes[0].Created.Format("Jan 2 2006 15:04"), Align: tview.AlignLeft, Color: tcell.ColorRed, Expansion: 1})
	dates.SetCell(0, 2, &tview.TableCell{Text: "Updated: ", Align: tview.AlignLeft, Color: tcell.ColorRed, Expansion: 1})
	dates.SetCell(0, 3, &tview.TableCell{Text: jNotes[0].LastUpdated.Format("Jan 2 2006 15:04"), Align: tview.AlignLeft, Color: tcell.ColorRed, Expansion: 1})

	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexRow)
	flex.AddItem(notesFlex, numberOfVisibleRows, 1, true)
	flex.AddItem(memo, 0, 3, false)
	flex.AddItem(dates, 1, 1, false)

	//set header
	idCell := &tview.TableCell{Text: "ID", Align: tview.AlignLeft, Color: tcell.ColorBlue, Expansion: 1, NotSelectable: true}
	notesTable.SetCell(0, 0, idCell)
	notebookTitleCell := &tview.TableCell{Text: "Notebook Title", Align: tview.AlignLeft, Color: tcell.ColorBlue, Expansion: 2, NotSelectable: true}
	notesTable.SetCell(0, 1, notebookTitleCell)
	noteTitleCell := &tview.TableCell{Text: "Note title", Align: tview.AlignLeft, Color: tcell.ColorBlue, Expansion: 2, NotSelectable: true}
	notesTable.SetCell(0, 2, noteTitleCell)
	tagsCell := &tview.TableCell{Text: "Tags", Align: tview.AlignLeft, Color: tcell.ColorBlue, Expansion: 2, NotSelectable: true}
	notesTable.SetCell(0, 3, tagsCell)

	for row := 0; row < len(jNotes); row++ {
		notesTable.SetCell(row+1, 0, &tview.TableCell{Text: strconv.FormatInt(jNotes[row].ID, 10), Align: tview.AlignLeft, Color: tcell.ColorLimeGreen, Expansion: 1})
		notesTable.SetCell(row+1, 1, &tview.TableCell{Text: jNotes[row].NotebookTitle, Align: tview.AlignLeft, Color: tcell.ColorLimeGreen, Expansion: 2})
		notesTable.SetCell(row+1, 2, &tview.TableCell{Text: jNotes[row].Title, Align: tview.AlignLeft, Color: tcell.ColorLimeGreen, Expansion: 2})
		notesTable.SetCell(row+1, 3, &tview.TableCell{Text: strings.Join(jNotes[row].Tags, ","), Align: tview.AlignLeft, Color: tcell.ColorLimeGreen, Expansion: 2})

	}
	notesTable.SetSelectionChangedFunc(func(row, column int) {
		if row < 0 {
			return
		}
		selectedNote := jNotes[row-1]
		memo.SetText(selectedNote.Memo)
		dates.SetCell(0, 1, &tview.TableCell{Text: selectedNote.Created.Format("Jan 2 2006 15:04"), Align: tview.AlignLeft, Color: tcell.ColorRed, Expansion: 2})
		dates.SetCell(0, 3, &tview.TableCell{Text: selectedNote.LastUpdated.Format("Jan 2 2006 15:04"), Align: tview.AlignLeft, Color: tcell.ColorRed, Expansion: 2})
	})

	app.SetRoot(flex, true)

	return app
}
