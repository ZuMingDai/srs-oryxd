package core

import "fmt"

const (
	Major     = 0
	Minor     = 0
	Reversion = 2
)

var Version = fmt.Sprintf("%v.%v.%v", Major, Minor, Reversion)
