package file

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/bingoohuang/gou/lang"
	"github.com/pkg/errors"
)

// SingleFileExists 检查文件是否存在，并且不是目录
func SingleFileExists(name string) error {
	mode, err := GetFileMode(name)
	if err != nil {
		return err
	}

	if mode == SingleFileMode {
		return nil
	}

	return fmt.Errorf("file %s is a directory", name)
}

// DirMode means file is single file or directory.
type DirMode int

const (
	// UnknownDirMode means unknown file or directory.
	UnknownDirMode DirMode = iota
	// DirectoryMode means directory.
	DirectoryMode
	// SingleFileMode means single file.
	SingleFileMode
)

// GetFileMode tells the name is a directory or not
func GetFileMode(name string) (DirMode, error) {
	if fi, err := os.Stat(name); err != nil {
		return UnknownDirMode, err
	} else if fi.IsDir() {
		return DirectoryMode, nil
	}

	return SingleFileMode, nil
}

// WriteBytes writes byte slice to file.
func WriteBytes(filePath string, b []byte) (int, error) {
	if err := os.MkdirAll(path.Dir(filePath), os.ModePerm); err != nil {
		return 0, err
	}

	fw, err := os.Create(filePath)
	if err != nil {
		return 0, err
	}

	defer fw.Close()

	return fw.Write(b)
}

// WriteString writes string to file.
func WriteString(filePath string, s string) (int, error) {
	return WriteBytes(filePath, []byte(s))
}

const (
	// TimeFormat defines the format of time to save to the file.
	TimeFormat = "2006-01-02 15:04:05"
)

// ReadTime reads the time.Time from the given file.
func ReadTime(filename string, defaultValue string) (time.Time, error) {
	v, err := ReadValue(filename, defaultValue)
	if err != nil {
		return time.Time{}, err
	}

	return lang.ParseTime(TimeFormat, v), nil
}

// WriteTime writes the time.Time to the given file.
func WriteTime(filename string, v time.Time) error {
	return WriteValue(filename, v.Format(TimeFormat))
}

func ReadValue(filename, defaultValue string) (string, error) {
	stat, err := StatE(filename)
	if err != nil {
		return "", errors.Wrapf(err, "file.Stat %s", filename)
	}

	if stat == NotExists || stat == Unknown {
		if err := WriteValue(filename, defaultValue); err != nil {
			return "", err
		}
	}

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", errors.Wrapf(err, "ioutil.ReadFile %s", filename)
	}

	return string(content), nil
}

// WritValue writes a string value to the file.
func WriteValue(filename string, value string) error {
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return errors.Wrapf(err, "MkdirAll %s", dir)
	}

	if err := ioutil.WriteFile(filename, []byte(value), 0644); err != nil {
		return errors.Wrapf(err, "WriteFile %s", filename)
	}

	return nil
}
