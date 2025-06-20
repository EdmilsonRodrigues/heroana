package files_test

import (
	"fmt"
	"io/fs"
	"reflect"
	"testing"
	"testing/fstest"

	"github.com/EdmilsonRodrigues/melo-project/src/melo/files"
)

func TestScanModule(t *testing.T) {
	fs := fstest.MapFS{
		"root":         {Mode: fs.ModeDir},
		"root/go.mod":  {Data: []byte("module example.com\n\ngo 1.24.0\n")},
		"root/main.go": {Data: []byte("package main\n\nfunc main() {}")},

		"root/exported_package":               {Mode: fs.ModeDir},
		"root/exported_package/other.go":      {Data: []byte(fmt.Sprintf("%smypackage.exported_package\n\npackage exported_package\n\nfunc main() {}", files.GoExportedDirective))},
		"root/exported_package/other_test.go": {Data: []byte(fmt.Sprintf("%smypackage.exported_package\n\npackage exported_package\n\nfunc TestMain(t *testing.T) {}", files.GoExportedDirective))},

		"root/unexported_package":          {Mode: fs.ModeDir},
		"root/unexported_package/other.go": {Data: []byte("package unexported_package\n\nfunc main() {}")},

		"root/mixed_package":                                              {Mode: fs.ModeDir},
		"root/mixed_package/not_exported_file.go":                         {Data: []byte("package mixed_package\n\nfunc main() {}")},
		"root/mixed_package/exported_file.go":                             {Data: []byte(fmt.Sprintf("%smypackage.mixed_package\n\npackage mixed_package\n\nfunc main() {}", files.GoExportedDirective))},
		"root/mixed_package/not_exported_file_directive_after_package.go": {Data: []byte("package mixed_package\n\n%smypackage.mixed_package\n\nfunc main() {}")},

		"root/package_with_different_name":          {Mode: fs.ModeDir},
		"root/package_with_different_name/bacon.go": {Data: []byte(fmt.Sprintf("%smypackage.baconpackage\n\npackage baconpackage\n\nfunc main() {}", files.GoExportedDirective))},
	}

	t.Run("should return exported packages and not return unexported packages", func(t *testing.T) {
		exportedPackages, err := files.ScanModule(fs, "root", "example.com")
		if err != nil {
			t.Errorf("ScanModule should not return error, got %v", err)
		}

		expected := []files.ExportedPackage{
			{
				GoPath:      "example.com/exported_package",
				PythonPath:  "mypackage.exported_package",
				PackageName: "",
			},
			{
				GoPath:      "example.com/mixed_package",
				PythonPath:  "mypackage.mixed_package",
				PackageName: "",
			},
			{
				GoPath:      "example.com/package_with_different_name",
				PythonPath:  "mypackage.baconpackage",
				PackageName: "baconpackage",
			},
		}

		if len(exportedPackages) != len(expected) {
			t.Errorf("ScanModule should return %d exported packages, got %d", len(expected), len(exportedPackages))
		}

		if !reflect.DeepEqual(exportedPackages, expected) {
			t.Errorf("ScanModule should return %+v, got %+v", expected, exportedPackages)
		}
	})
}
