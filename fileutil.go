// Copyright 2014 racklin@gmail.com .

// Package fileutil implements some File utility functions.
package fileutil

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"time"
)

const (
	NEW_FILE_PERM = 0644
)

type FileInfo os.FileInfo
type FilesInfo []FileInfo

func (f FilesInfo) Len() int      { return len(f) }
func (f FilesInfo) Swap(i, j int) { f[i], f[j] = f[j], f[i] }

// GetFileInfo returns a FileInfo describing the named file
func GetFileInfo(name string) (FileInfo, error) {
	fi, err := os.Stat(name)
	return fi, err
}

// Exists checks if the given filename exists
func Exists(filename string) (bool, error) {
	if _, err := GetFileInfo(filename); err != nil {
		if os.IsNotExist(err) {
			return false, err
		}
	}
	return true, nil
}

// Size return the size of the given filename.
// Returns -1 if the file does not exist or if the file size cannot be determined.
func Size(filename string) (int64, error) {
	if fi, err := GetFileInfo(filename); err != nil {
		return -1, err
	} else {
		return fi.Size(), nil
	}
}

// ModTime return the Last Modified Time of the given filename.
func ModTime(filename string) (time.Time, error) {
	if fi, err := GetFileInfo(filename); err != nil {
		return time.Unix(0, 0), err
	} else {
		return fi.ModTime(), nil
	}
}

// ModTimeUnix return the Last Modified Unix Timestamp of the given filename.
// Returns -1 if the file does not exist or if the file modtime cannot be determined.
func ModTimeUnix(filename string) (int64, error) {
	if ft, err := ModTime(filename); err != nil {
		return -1, err
	} else {
		return ft.Unix(), nil
	}
}

// ModTimeUnixNano return the Last Modified Unix Time of nanoseconds of the given filename.
// Returns -1 if the file does not exist or if the file modtime cannot be determined.
func ModTimeUnixNano(filename string) (int64, error) {
	if ft, err := ModTime(filename); err != nil {
		return -1, err
	} else {
		return ft.UnixNano(), nil
	}
}

// Mode return the FileMode of the given filename.
// Returns 0 if the file does not exist or if the file mode cannot be determined.
func Mode(filename string) (os.FileMode, error) {
	if fi, err := GetFileInfo(filename); err != nil {
		return 0, err
	} else {
		return fi.Mode(), nil
	}
}

// Perm return the Unix permission bits of the given filename.
// Returns 0 if the file does not exist or if the file mode cannot be determined.
func Perm(filename string) (os.FileMode, error) {
	if fi, err := GetFileInfo(filename); err != nil {
		return 0, err
	} else {
		return fi.Mode().Perm(), nil
	}
}

// Read reads the file named by filename and returns the contents.
func Read(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}

// Write writes data to a file named by filename.
func Write(filename string, content []byte) error {
	return ioutil.WriteFile(filename, content, NEW_FILE_PERM)
}

// ReadString reads the file named by filename and returns the contents as string.
func ReadString(filename string) (string, error) {
	buf, err := Read(filename)
	if err != nil {
		return "", err
	} else {
		return string(buf), nil
	}
}

// WriteString writes the contents of the string to filename.
func WriteString(filename, content string) error {
	return Write(filename, []byte(content))
}

// AppendString appends the contents of the string to filename.
func AppendString(filename, content string) error {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, NEW_FILE_PERM)
	if err != nil {
		return err
	}
	data := []byte(content)
	n, err := f.Write(data)
	if err == nil && n < len(data) {
		err = io.ErrShortWrite
	}
	if err1 := f.Close(); err == nil {
		err = err1
	}
	return err
}

// GetContents reads the file named by filename and returns the contents as string.
// GetContents is equivalent to ReadString.
func GetContents(filename string) (string, error) {
	return ReadString(filename)
}

// PutContents writes the contents of the string to filename.
// PutContents is equivalent to WriteString.
func PutContents(filename, content string) error {
	return WriteString(filename, content)
}

// AppendContents appends the contents of the string to filename.
// AppendContents is equivalent to AppendString.
func AppendContents(filename, content string) error {
	return AppendString(filename, content)
}

// TempFile creates a new temporary file in the default directory for temporary files (see os.TempDir), opens the file for reading and writing, and returns the resulting *os.File.
func TempFile() (*os.File, error) {
	return ioutil.TempFile("", "")
}

// TempName creates a new temporary file in the default directory for temporary files (see os.TempDir), opens the file for reading and writing, and returns the filename.
func TempName() (string, error) {

	f, err := TempFile()
	if err != nil {
		return "", err
	}
	return f.Name(), nil

}

// Copy makes a copy of the file source to dest.
func Copy(source, dest string) (err error) {

	// checks source file is regular file
	sfi, err := GetFileInfo(source)
	if err != nil {
		return
	}
	if !sfi.Mode().IsRegular() {
		errors.New("cannot copy non-regular files.")
		return
	}

	// checks dest file is regular file or the same
	dfi, err := GetFileInfo(dest)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
	} else {
		if !dfi.Mode().IsRegular() {
			errors.New("cannot copy to non-regular files.")
			return
		}
		if os.SameFile(sfi, dfi) {
			return
		}
	}

	// hardlink source to dest
	err = os.Link(source, dest)
	if err == nil {
		return
	}

	// cannot hardlink , copy contents
	in, err := os.Open(source)
	if err != nil {
		return
	}
	defer in.Close()

	out, err := os.Create(dest)
	if err != nil {
		return
	}
	defer out.Close()

	if _, err = io.Copy(out, in); err != nil {
		return
	}

	// syncing file
	err = out.Sync()
	return
}

// Find returns the FilesInfo([]FileInfo) of all files matching pattern or nil if there is no matching file. The syntax of patterns is the same as in Match. The pattern may describe hierarchical names such as /usr/*/bin/ed (assuming the Separator is '/').
func Find(pattern string) (FilesInfo, error) {

	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}

	files := make(FilesInfo, 0)
	for _, f := range matches {
		fi, err := GetFileInfo(f)
		if err == nil {
			files = append(files, fi)
		}
	}
	return files, err
}

type ByName struct{ FilesInfo }

type BySize struct{ FilesInfo }

type ByModTime struct{ FilesInfo }

func (s ByName) Less(i, j int) bool { return s.FilesInfo[i].Name() < s.FilesInfo[j].Name() }
func (s BySize) Less(i, j int) bool { return s.FilesInfo[i].Size() < s.FilesInfo[j].Size() }
func (s ByModTime) Less(i, j int) bool {
	return s.FilesInfo[i].ModTime().Before(s.FilesInfo[j].ModTime())
}

// SortByName sorts a slice of files by filename in increasing order.
func (fis FilesInfo) SortByName() {
	sort.Sort(ByName{fis})
}

// SortBySize sorts a slice of files by filesize in increasing order.
func (fis FilesInfo) SortBySize() {
	sort.Sort(BySize{fis})
}

// SortByModTime sorts a slice of files by file modified time in increasing order.
func (fis FilesInfo) SortByModTime() {
	sort.Sort(ByModTime{fis})
}

// SortByNameReverse sorts a slice of files by filename in decreasing order.
func (fis FilesInfo) SortByNameReverse() {
	sort.Sort(sort.Reverse(ByName{fis}))
}

// SortBySizeReverse sorts a slice of files by filesize in decreasing order.
func (fis FilesInfo) SortBySizeReverse() {
	sort.Sort(sort.Reverse(BySize{fis}))
}

// SortByModTimeReverse sorts a slice of files by file modified time in decreasing order.
func (fis FilesInfo) SortByModTimeReverse() {
	sort.Sort(sort.Reverse(ByModTime{fis}))
}
