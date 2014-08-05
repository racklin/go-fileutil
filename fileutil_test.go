// Copyright 2014 racklin@gmail.com .

// Package fileutil implements some File utility functions.
package fileutil

import (
	"os"
	"testing"
	"time"
)

const (
	TEST_FILE           = "fileutil_test.go"
	TEST_FILE_NONEXISTS = "dummy.dummy"
)

var ZERO_TIME = time.Unix(0, 0)

func TestExists(t *testing.T) {
	b, _ := Exists(TEST_FILE)
	if !b {
		t.Errorf("TestExists: %s should exists.", TEST_FILE)
	}
}

func TestNonExists(t *testing.T) {
	b, _ := Exists(TEST_FILE_NONEXISTS)
	if b {
		t.Errorf("TestNonExists: %s should not exists.", TEST_FILE_NONEXISTS)
	}
}

func TestBasename(t *testing.T) {
	b := Basename("/a/b/c")
	if b != "c" {
		t.Errorf("TestBasename: %s should been 'c'", "/a/b/c")
	}
}

func TestDirname(t *testing.T) {
	b := Dirname("/a/b/c")
	if b != "/a/b" {
		t.Errorf("TestDirname: %s should been '/a/b'", "/a/b/c")
	}
}

func TestExtname(t *testing.T) {
	b := Extname(TEST_FILE)

	if b != ".go" {
		t.Errorf("TestExtname: %s should been '.go'", TEST_FILE)
	}
}

func TestSize(t *testing.T) {
	s, _ := Size(TEST_FILE)
	if s <= 0 {
		t.Errorf("TestSize: %s should larger then 0.", TEST_FILE)
	}
}

func TestSizeNonExists(t *testing.T) {
	s, _ := Size(TEST_FILE_NONEXISTS)
	if s >= 0 {
		t.Errorf("TestSizeNonExists: %s should less then 0.", TEST_FILE)
	}
}

func TestModTime(t *testing.T) {
	ft, _ := ModTime(TEST_FILE)
	if ft.Before(ZERO_TIME) {
		t.Errorf("TestModTime: %s should after then 1970/1/1.", TEST_FILE)
	}
}

func TestModTimeNonExists(t *testing.T) {
	ft, _ := ModTime(TEST_FILE_NONEXISTS)
	if ft.After(ZERO_TIME) {
		t.Errorf("TestModTimeNonExists: %s should been 1970/1/1.", TEST_FILE_NONEXISTS)
	}
}

func TestModTimeUnix(t *testing.T) {
	ft, _ := ModTimeUnix(TEST_FILE)
	if ft <= 0 {
		t.Errorf("TestModTimeUnix: %s should larger then 0", TEST_FILE)
	}
}

func TestModTimeUnixNonExists(t *testing.T) {
	ft, _ := ModTimeUnix(TEST_FILE_NONEXISTS)
	if ft >= 0 {
		t.Errorf("TestModTimeUnixNonExists: %s should less then 0", TEST_FILE_NONEXISTS)
	}
}

func TestGetContents(t *testing.T) {
	buf, _ := GetContents(TEST_FILE)
	if len(buf) <= 0 {
		t.Errorf("TestGetContents: %s should larger then 0", TEST_FILE)
	}
}

func TestAppendContents(t *testing.T) {
	//err := AppendContents(TEST_FILE, "XXX\n\n\n")
	//fmt.Println(err)
	//if err != nil {
	//	t.Errorf("TestGetContents: %s should larger then 0", TEST_FILE)
	//}
}

func TestTempName(t *testing.T) {
	_, err := TempName()
	if err != nil {
		t.Errorf("TestTempName: ", err)
	}
}

func TestCopy(t *testing.T) {
	tmp, _ := TempName()
	err := Copy(TEST_FILE, tmp)
	defer os.Remove(tmp)
	if err != nil {
		t.Errorf("TestCopy: ", err)
	}
}

func TestFind(t *testing.T) {
	_, err := Find("/tmp/*")
	if err != nil {
		t.Errorf("TestFind: ", err)
	}
}

func TestFindAndSort(t *testing.T) {
	_, err := Find("/tmp/*")
	if err != nil {
		t.Errorf("TestFindAndSort: ", err)
	}
}

func TestExec(t *testing.T) {
	_, err := Exec("/bin/date")
	if err != nil {
		t.Errorf("TestExec: ", err)
	}
}
