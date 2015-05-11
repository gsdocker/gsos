// +build windows

package gsos

import (
	"os/exec"

	"github.com/gsdocker/gserrors"
)

// windows special const variable defines
const (
	ExeSuffix = ".exe"
)

// RemoveAll .
func RemoveAll(dir string) error {

	cmd := exec.Command("cmd", "/C", "rd", "/S", "/Q", dir)

	output, err := cmd.Output()

	if err != nil {
		return gserrors.Newf(err, string(output))
	}

	return nil
}
