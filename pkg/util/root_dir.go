package util

import (
	"path/filepath"
	"runtime"
)

// RootDir is a function uses to get an absolute path of current package.
// It useful for getting the real path of static file which is live in go-pkg.
func RootDir() string {
	_, b, _, ok := runtime.Caller(0)
	if ok {
		return filepath.Join(filepath.Dir(b), "../..")
	}

	return ""
}
