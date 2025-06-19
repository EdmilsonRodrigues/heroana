package files_test

import (
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/EdmilsonRodrigues/melo-project/src/melo/files"
)

func TestCheckInputFolder(t *testing.T) {
	fs := fstest.MapFS{
		"root": {Mode: fs.ModeDir},

		"root/right_input_folder": {Mode: fs.ModeDir},
		"root/right_input_folder/go.mod": {Data: []byte("module example.com\n\ngo 1.24.0\n")},

		"root/wrong_input_folder": {Mode: fs.ModeDir},
		"root/wrong_input_folder/not_go.mod": {Data: []byte("module example.com\n\ngo 1.24.0\n")},  // wrong file name
	
		"root/input_folder_wrong_go_mod": {Mode: fs.ModeDir},
		"root/input_folder_wrong_go_mod/go.mod": {Data: []byte("odule example.com\n\ngo 1.24.0\n")}, // wrong module name
	}

	t.Run("should return true for when input folder has a valid go.mod", func(t *testing.T) {
		if !files.CheckInputFolder(fs, "root/right_input_folder") {
			t.Errorf("CheckInputFolder should return true for right input folder")
		}
	})

	t.Run("should return false for input folder without go.mod", func(t *testing.T) {
		if files.CheckInputFolder(fs, "root/wrong_input_folder") {
			t.Errorf("CheckInputFolder should return false for wrong input folder")
		}
	})

	t.Run("should return false for input folder with wrong go.mod", func(t *testing.T) {
		if files.CheckInputFolder(fs, "root/input_folder_wrong_go_mod") {
			t.Errorf("CheckInputFolder should return false for wrong input folder")
		}
	})
}
