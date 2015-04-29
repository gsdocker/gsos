package gsos

import (
	"errors"
	"io"
	"io/ioutil"
	"os"

	"github.com/gsdocker/gserrors"
)

// Errors
var (
	ErrFS = errors.New("file system error")
)

// IsExist check file entry if exist
func IsExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

// SameFile .
func SameFile(f1, f2 string) bool {
	s1, err := os.Stat(f1)

	if err != nil {
		return false
	}

	s2, err := os.Stat(f2)

	if err != nil {
		return false
	}

	return os.SameFile(s1, s2)
}

// IsDir check file entry if a directory entry
func IsDir(filename string) bool {

	fi, err := os.Stat(filename)

	if err != nil {
		return false
	}
	return fi.IsDir()
}

// CopyFile Copies file source to destination dest.
func CopyFile(source string, dest string) (err error) {
	sf, err := os.Open(source)
	if err != nil {
		return err
	}
	defer sf.Close()
	df, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer df.Close()
	_, err = io.Copy(df, sf)
	if err == nil {
		si, err := os.Stat(source)
		if err != nil {
			err = os.Chmod(dest, si.Mode())
		}

	}

	return
}

// CopyDir Recursively copies a directory tree, attempting to preserve permissions.
// Source directory must exist, destination directory must *not* exist.
func CopyDir(source string, dest string) (err error) {

	// get properties of source dir
	fi, err := os.Stat(source)
	if err != nil {
		return err
	}

	if !fi.IsDir() {
		return gserrors.Newf(ErrFS, "Source is not a directory")
	}

	// ensure dest dir does not already exist

	_, err = os.Open(dest)
	if !os.IsNotExist(err) {
		return gserrors.Newf(ErrFS, "Destination already exists")
	}

	// create dest dir

	err = os.MkdirAll(dest, fi.Mode())
	if err != nil {
		return err
	}

	entries, err := ioutil.ReadDir(source)

	for _, entry := range entries {

		sfp := source + "/" + entry.Name()
		dfp := dest + "/" + entry.Name()
		if entry.IsDir() {
			err = CopyDir(sfp, dfp)
			if err != nil {
				return err
			}
		} else {
			// perform copy
			err = CopyFile(sfp, dfp)
			if err != nil {
				return err
			}
		}

	}
	return
}
