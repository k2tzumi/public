package firego

type MessageType int

const (
	Log MessageType = iota
	Info
	Warn
	Error
	Table
)
