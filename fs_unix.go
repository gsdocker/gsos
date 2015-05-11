// +build !windows

package gsos

import "os"

// windows special const variable defines
const (
	ExeSuffix = ""
)

// RemoveAll .
var RemoveAll = os.RemoveAll
