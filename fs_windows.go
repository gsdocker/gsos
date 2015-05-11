// +build windows

package gsos

import (
	"os"
	"os/exec"

	"github.com/gsdocker/gserrors"
)

// windows special const variable defines
const (
	ExeSuffix = ".exe"
)

// RemoveAll .
func RemoveAll(dir string) error {

	current := CurrentDir()

	err := os.Chdir(dir)

	if err != nil {
		return err
	}

	cmd := exec.Command("attrib", "-R", dir, "/S", "/D", "/L")

	output, err := cmd.Output()

	os.Chdir(current)

	if err != nil {
		return gserrors.Newf(err, string(output))
	}

	return os.RemoveAll(dir)
}
