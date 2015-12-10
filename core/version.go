package core
import "fmt"

const (
	Major = 0
	Minor = 0
	Reversion = 1
)

var Version = fmt.Sprintf("%v.%v.%v", Major,Minor, Reversion)
