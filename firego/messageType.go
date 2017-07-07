package firego

import "fmt"

type MessageType int

const (
	Log MessageType = iota
	Info
	Warn
	Error
	Table
	GroupStart
	GroupEnd
)

func (i MessageType) String() string {
	translation := []string{
		Log:        `LOG`,
		Info:       `INFO`,
		Warn:       `WARN`,
		Error:      `ERROR`,
		Table:      `TABLE`,
		GroupStart: `GROUP_START`,
		GroupEnd:   `GROUP_END`,
	}
	if i < 0 || int(i) > len(translation) {
		return fmt.Sprintf("MessageType(%d)", i)
	}
	return translation[i]
}
