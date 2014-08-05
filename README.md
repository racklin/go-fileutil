# Package fileutil implements some File utility functions.
Simple and without get into nested if statement hell.

## Using
```go
import fileutil "github.com/racklin/go-fileutil"
```

## Check File Exists
```go
func Exists(filename string) (bool, error)
```
Example:
```go
    exists, err := fileutil.Exists("/etc/passwd")
```

## FileInfo

### Filesize
```go
func Size(filename string) (int64, error)
```
Example:
```go
    size, err := fileutil.Size("/etc/passwd")
```

### File Modified Time
```go
func ModTime(filename string) (time.Time, error)
```
or
```go
func ModTimeUnix(filename string) (int64, error)
```
or
```go
func ModTimeUnixNano(filename string) (int64, error)
```
Example:
```go
    mtime, err := fileutil.ModTimeUnix("/etc/passwd")
```

### File Mode and Permission
```go
func Mode(filename string) (os.FileMode, error)
func Perm(filename string) (os.FileMode, error)
```

## File Read And Write
### Bytes Operations
```go
func Read(filename string) ([]byte, error)
func Write(filename string, content []byte) error
```
### String Operations
```go
func ReadString(filename string) (string, error)
func WriteString(filename, content string) error
func AppendString(filename, content string) error
```
Example:
```go
    func Log(message string) {
        fileutil.AppendString("/var/log/test.log", message)
    }
```
### Alias
If you are PHPer , feels like coming home (file_put_contents/file_get_contents).
```go
func GetContents(filename string) (string, error)
func PutContents(filename, content string) error
func AppendContents(filename, content string) error
```

## File Copy
Copies a file from source to destination.
If source and destination files exist and the same, return success.
Otherise, using hard link between the two files. If that fail, copy the file contents.
```go
func Copy(source, dest string) (err error)
```
Example:
```go
err := fileutil.Copy("/etc/passwd", "/tmp/passwd")
```

## File Glob and Sort
Find returns the FilesInfo([]FileInfo) of all files matching pattern
```go
func Find(pattern string) (FilesInfo, error)
```
Sorting FilesInfo:
```go
func (fis FilesInfo) SortByName()
func (fis FilesInfo) SortBySize()
func (fis FilesInfo) SortByModTime()
func (fis FilesInfo) SortByNameReverse()
func (fis FilesInfo) SortBySizeReverse()
func (fis FilesInfo) SortByModTimeReverse()
```
Example1: List file of /tmp
```go
    files, err := fileutil.Find("/tmp/*")
```
Example2: Find go file in /usr/local directory
```go
    files, err := fileutil.Find("/usr/local/*/go")
```
Example3: Sorting result files
```go
    files, err := fileutil.Find("/tmp/*")
    files.SortByName()
```
