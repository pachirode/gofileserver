package staticserver

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func checkFilename(filename string) error {
	if strings.ContainsAny(filename, "\\/:*<>|") {
		return errors.New("Name should not contain \\/:*<>|")
	}
	return nil
}

func sanitizedName(filename string) string {
	if len(filename) > 1 && filename[1] == ':' && runtime.GOOS == "windows" {
		filename = filename[2:]
	}

	filename = strings.TrimLeft(strings.Replace(filename, `\`, "/", -1), `/`)
	filename = filepath.ToSlash(filename)
	filename = filepath.Clean(filename)
	return filename
}

func isFile(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.Mode().IsRegular()
}

func isDir(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.Mode().IsDir()
}
