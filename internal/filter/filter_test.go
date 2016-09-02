package filter

import (
	"bufio"
	"os"
	"testing"
)

func TestBloomfilter(t *testing.T) {
	b := New(2048, 16)
	lines, err := readLines("wordlist")
	if err != nil {
		t.Fatal("Dictionary could not be loaded")
	}
	for _, line := range lines {
		b.Add(line)
	}

	if b.Has("THIS WORD SHOULD NOT BE FOUND") {
		t.Errorf("Should not find this word: THIS WORD SHOULD NOT BE FOUND")
	}
	if !b.Has("three") {
		t.Errorf("Should find this word: three")
	}
	if !b.Has("used") {
		t.Errorf("Should find this word: used")
	}

	if b.Saturation() != 0.69140625 {
		t.Errorf("change in algorithm made a change in saturation")
	}
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func TestBloomfilterDel(t *testing.T) {
	const word = "word"
	b := New(10, 8)
	b.Add(word)

	if !b.Has("word") {
		t.Errorf("Word should have been found")
	}

	b.Del(word)
	if b.Has("word") {
		t.Errorf("Word should have not been found after deletion")
	}

	b.Del(word)
	if b.Has("word") {
		t.Errorf("Word should have not been found after deletion")
	}

	for _, v := range b.bitspaceT {
		if v < 0 {
			t.Fatalf("bitspaceT corrupted")
		}
	}
	for _, v := range b.bitspaceF {
		if v < 0 {
			t.Fatalf("bitspaceF corrupted")
		}
	}
}
