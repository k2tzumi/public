package btstrpr

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestRender(t *testing.T) {
	os.Chdir("golden")
	files, err := ioutil.ReadDir(".")
	if err != nil {
		t.Fatal(err)
	}

	for _, file := range files {
		var got bytes.Buffer
		cmd := exec.Command("go", "run", filepath.Join(file.Name(), "got.go"))
		cmd.Stdout = &got
		cmd.Run()

		expected, err := ioutil.ReadFile(filepath.Join(file.Name(), "expect.html"))
		if err != nil {
			t.Fatal(err)
		}

		if result := bytes.Compare(got.Bytes(), expected); result != 0 {
			t.Error(file.Name(), "error")
			t.Log("got:", got.String())
			t.Log("len:", len(got.String()))
			t.Log("expected:", string(expected))
			t.Log("len:", len(string(expected)))
		}
	}
}
