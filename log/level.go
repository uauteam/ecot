package log


type Level uint8


const (
	OFF Level = iota
	TRACE
	DEBUG
	INFO
	WARN
	ERROR
)