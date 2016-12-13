package tail

import "github.com/hpcloud/tail"

// Tail implements an Unix-like tail call which polls a given file for its
// contents.
type Tail struct {
	tail *tail.Tail
	line string
}

// Scan reads one line of the tail-read file. Meant to be used in tandem with
// Text().
func (t *Tail) Scan() bool {
	line, ok := <-t.tail.Lines
	if line == nil {
		return false
	}
	t.line = line.Text
	return ok
}

// Text returns last Scan()'d file.
func (t *Tail) Text() string {
	return t.line
}

// Open creates a Tail for a given file. It bubbles up the errors from the
// underlying file monitor package (github.com/hpcloud/tail).
func Open(fn string) (*Tail, error) {
	t, err := tail.TailFile(fn, tail.Config{Follow: true, ReOpen: true, Logger: tail.DiscardingLogger})
	return &Tail{
		tail: t,
	}, err
}
