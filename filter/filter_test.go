package filter

import (
	"bufio"
	"os"
	"testing"
)

func TestBloomfilter(t *testing.T) {
	b := newBloomfilter(2048, 8)
	lines, err := readLines("wordlist")
	if err != nil {
		t.Fatal("Dictionary could not be loaded")
	}
	for _, line := range lines {
		b.add(line)
	}

	if b.has("THIS WORD SHOULD NOT BE FOUND") {
		t.Errorf("Should not find this word: THIS WORD SHOULD NOT BE FOUND")
	}
	if !b.has("three") {
		t.Errorf("Should find this word: three")
	}
	if !b.has("used") {
		t.Errorf("Should find this word: used")
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
