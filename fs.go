package gsos

import "os"

// IsExist check file entry if exist
func IsExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

// IsDir check file entry if a directory entry
func IsDir(filename string) bool {

	fi, err := os.Stat(filename)

	if err != nil {
		return false
	}
	return fi.IsDir()
}
