package repository

import (
	"github.com/nicolasmanic/tefter/model"
	"os"
	"testing"
)

func TestSaveNotebook(t *testing.T) {
	testRepo := NewNotebookRepository("test.db")
	//tear down test
	defer func() {
		testRepo.CloseDB()
		os.Remove("test.db")
	}()

	mockNotebook := model.NewNotebook("testTitle")
	id, err := testRepo.SaveNotebook(mockNotebook)

	if err != nil {
		t.Errorf("Could not save notebook to DB, error msg: %v", err)
	}

	if id != mockNotebook.ID {
		t.Error("Could not save correctly notebook to DB")
	}
}

func TestGetNotebooks(t *testing.T) {
	testRepo := NewNotebookRepository("test.db")
	testNoteRepo := NewNoteRepository("test.db")
	//tear down test
	defer func() {
		testRepo.CloseDB()
		testNoteRepo.CloseDB()
		os.Remove("test.db")
	}()

	mockNotebook2 := model.NewNotebook("notebook 2")
	mockNotebook3 := model.NewNotebook("notebook 3")
	testRepo.SaveNotebook(mockNotebook2)
	testRepo.SaveNotebook(mockNotebook3)

	mockNote1 := model.NewNote("testTitle", "test Memo", 2, []string{"testTag1", "testTag2"})
	mockNote2 := model.NewNote("testTitle", "test Memo", 3, []string{"testTag3", "testTag4"})

	testNoteRepo.SaveNote(mockNote1)
	testNoteRepo.SaveNote(mockNote2)

	notebooks, err := testRepo.GetNotebooks([]int64{2, 3})

	if err != nil {
		t.Errorf("Could not retrieve notebooks from DB, error msg: %v", err)
	}

	if len(notebooks) != 2 {
		t.Error("Could not retrieve notebooks from DB")
	}

	for _, notebook := range notebooks {
		if len(notebook.Notes) != 1 {
			t.Error("Could not retrieve notebooks from DB")
		}
	}

}

func TestUpdateNotebook(t *testing.T) {
	testRepo := NewNotebookRepository("test.db")
	//tear down test
	defer func() {
		testRepo.CloseDB()
		os.Remove("test.db")
	}()

	mockNotebook := model.NewNotebook("test Title")
	id, _ := testRepo.SaveNotebook(mockNotebook)

	mockNotebook.Title = "Updated Title"

	err := testRepo.UpdateNotebook(mockNotebook)
	if err != nil {
		t.Errorf("Could not update notebook, error msg: %v", err)
	}
	notebook, _ := testRepo.GetNotebook(id)

	if notebook.Title != "Updated Title" {
		t.Error("Could not update notebook")
	}
}

func TestDeleteNotebooks(t *testing.T) {
	testRepo := NewNotebookRepository("test.db")
	testNoteRepo := NewNoteRepository("test.db")
	//tear down test
	defer func() {
		testRepo.CloseDB()
		testNoteRepo.CloseDB()
		os.Remove("test.db")
	}()

	mockNotebook2 := model.NewNotebook("notebook 2")
	mockNotebook3 := model.NewNotebook("notebook 3")
	testRepo.SaveNotebook(mockNotebook2)
	testRepo.SaveNotebook(mockNotebook3)

	mockNote1 := model.NewNote("testTitle", "test Memo", 2, []string{"testTag1", "testTag2"})
	mockNote2 := model.NewNote("testTitle", "test Memo", 3, []string{"testTag3", "testTag4"})

	testNoteRepo.SaveNote(mockNote1)
	testNoteRepo.SaveNote(mockNote2)

	err := testRepo.DeleteNotebooks([]int64{2, 3})
	if err != nil {
		t.Errorf("Could not delete notebooks from DB, error msg: %v", err)
	}

	notes, _ := testNoteRepo.GetNotes([]int64{mockNote1.ID, mockNote2.ID})
	if len(notes) != 0 {
		t.Errorf("Could not delete notes of deleted notebook from DB")
	}

	notebooks, _ := testRepo.GetNotebooks([]int64{2, 3})
	if len(notebooks) != 0 {
		t.Errorf("Could not delete notebook from DB")
	}
}
