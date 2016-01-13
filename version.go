package main

import "fmt"

const (
	major     = 0
	minor     = 0
	reversion = 2
)

func Version() string {
	return fmt.Sprintf("%v.%v.%v", major, minor, reversion)
}
