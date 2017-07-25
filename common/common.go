// Package common declares common type which are used in both
// server and client
package common

const (
	TEXT   = 1
	BINARY = 2
)

const (
	STATUS_OK    = 1
	STATUS_ERROR = 2
)

type DataType int
type Status int
