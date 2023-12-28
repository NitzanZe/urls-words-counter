package models

type errorLevel string

const (
	WARN  errorLevel = "warning"
	FATAL errorLevel = "fatal"
)

type Error struct {
	Err   error
	Level errorLevel
}
