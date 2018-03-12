package main

import (
	"github.com/nicolasmanic/tefter/cmd"
	"github.com/nicolasmanic/tefter/repository"
)

func main() {
	//FIXME:  add dbpath to external property file
	dbPath := "tefter.db"
	noteDB := repository.NewNoteRepository(dbPath)
	notebookDB := repository.NewNotebookRepository(dbPath)

	cmd.NoteDB = noteDB
	cmd.NotebookDB = notebookDB
	cmd.Execute()
}
