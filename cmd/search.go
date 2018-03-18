package cmd

import (
	"github.com/nicolasmanic/tefter/model"
	"github.com/spf13/cobra"
	"log"
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search notes given a keyword",
	Long:  "Keyword is searched against title and content of the note, if no keyword is given all notes will be printed",
	Args:  cobra.ExactArgs(1),
	Run:   searchWrapper,
}

func init() {
	rootCmd.AddCommand(searchCmd)
}

func searchWrapper(cmd *cobra.Command, args []string) {
	search(args, printNotes2Terminal)
}

func search(args []string, output func(notes []*model.Note)) {
	keyword := ""
	if len(args) > 0 {
		keyword = args[0]
	}
	var notes []*model.Note
	var err error
	if len(keyword) == 0 {
		notes, err = NoteDB.GetNotes([]int64{})
	} else {
		notes, err = NoteDB.SearchNotesByKeyword(keyword)
	}
	if err != nil {
		log.Panicf("Error retrieving Notes from DB, error msg: %v", err)
	}
	output(notes)
}
