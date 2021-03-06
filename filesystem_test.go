package goblet

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

var (
	f = []*File{
		{
			Path:       "/tmp/test",
			Data:       nil,
			FileMode:   os.ModeDir,
			ModifiedAt: time.Now(),
		},
		{
			Path:       "/tmp/test/hoge.txt",
			Data:       []byte("hogehoge"),
			FileMode:   0x800001ed,
			ModifiedAt: time.Now(),
		},
		{
			Path:       "/tmp/test/fuga.txt",
			Data:       []byte("fuga"),
			FileMode:   0x800001ed,
			ModifiedAt: time.Now(),
		},
	}
	files = map[string]*File{
		"/tmp/test":          f[0],
		"/tmp/test/hoge.txt": f[1],
		"/tmp/test/fuga.txt": f[2],
	}
	dirs = map[string][]string{
		"/tmp/test": {
			"hoge.txt", "fuga.txt", "not_exists.png",
		},
	}
)

func TestNewFS(t *testing.T) {

	// expected data
	expected := &FileSystem{
		Dirs:       dirs,
		Files:      files,
		pathPrefix: "",
	}

	// get actual
	actual := NewFS(dirs, files)

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected and actual are not equal, expected %#+v and actual %#+v", expected, actual)
	}
}

func TestWithPrefix(t *testing.T) {
	path := "/static/"

	base := NewFS(dirs, files)

	// get actual
	actual := base.WithPrefix(path)

	// expected data
	expected := &FileSystem{
		Dirs:       dirs,
		Files:      files,
		pathPrefix: path,
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected and actual are not equal, expected %#+v and actual %#+v", expected, actual)
	}

	if reflect.DeepEqual(base, actual) {
		t.Errorf("expected and actual are accidentally equal, expected %#+v and actual %#+v", expected, actual)
	}

}

func TestWithPrefix_PanicWhenNil(t *testing.T) {
	var fs *FileSystem
	defer func() {
		if err := recover(); err != nil {
			return
		}
		t.Fail()
	}()
	fs.WithPrefix("")
	t.Fail()
}

func TestWithIgnorePrefix(t *testing.T) {
	path := "/static/"

	base := NewFS(dirs, files)

	// get actual
	actual := base.WithIgnoredPrefix(path)

	// expected data
	expected := &FileSystem{
		Dirs:          dirs,
		Files:         files,
		ignoredPrefix: path,
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected and actual are not equal, expected %#+v and actual %#+v", expected, actual)
	}

	if reflect.DeepEqual(base, actual) {
		t.Errorf("expected and actual are accidentally equal, expected %#+v and actual %#+v", expected, actual)
	}
}

func TestWithIgnoredPrefix_PanicWhenNil(t *testing.T) {
	var fs *FileSystem
	defer func() {
		if err := recover(); err != nil {
			return
		}
		t.Fail()
	}()
	fs.WithIgnoredPrefix("")
	t.Fail()
}
func TestResolute(t *testing.T) {
	prefix := "/static"
	path := "/static/tmp/test/fuga.txt"
	fs := NewFS(dirs, files).WithPrefix(prefix)

	actual := fs.resolute(path)

	expected := "/tmp/test/fuga.txt"

	if actual != expected {
		t.Errorf("path FileSystem.resolute returned: %s does not equal expected: %s \n", actual, expected)
	}

	fs = NewFS(dirs, files)

	actual = fs.resolute(path)

	expected = path

	if actual != expected {
		t.Errorf("path FileSystem.resolute returned: %s does not equal expected: %s \n", actual, expected)
	}

	ignored := "/static/"
	fs = NewFS(dirs, files).WithIgnoredPrefix(ignored)
	path = "/tmp/test/fuga.txt"
	actual = fs.resolute(path)

	expected = filepath.Join(ignored, path)

	if actual != expected {
		t.Errorf("path FileSystem.resolute returned: %s does not equal expected: %s \n", actual, expected)
	}
}

func TestExists(t *testing.T) {
	fs := NewFS(dirs, files)
	path := "/tmp/test/hoge.txt"
	actual := fs.Exists(path)
	expected := true

	if actual != expected {
		t.Errorf("FileSystem.Exists should return %t but returned %t", expected, actual)
	}

	path = "/tmp/test/notexists.txt"
	actual = fs.Exists(path)
	expected = false

	if actual != expected {
		t.Errorf("FileSystem.Exists should return %t but returned %t", expected, actual)
	}
}

func TestOpen(t *testing.T) {
	fs := NewFS(dirs, files)
	path := "/tmp/test/hoge.txt"
	result, err := fs.Open(path)
	if err != nil || result == nil {
		t.Error("FileSystem.Open should not return error if path exists")
	}

	path = "/tmp/test/notexists.txt"
	result, err = fs.Open(path)

	if err == nil || result != nil {
		t.Error("FileSystem.Open should return error if path doesn't exist")
	}
}

func TestFile(t *testing.T) {
	fs := NewFS(dirs, files)
	path := "/tmp/test/hoge.txt"
	result, err := fs.File(path)
	if err != nil || result == nil {
		t.Error("FileSystem.File should not return error if path exists")
	}

	path = "/tmp/test/notexists.txt"
	result, err = fs.File(path)

	if err == nil || result != nil {
		t.Error("FileSystem.File should return error if path doesn't exist")
	}
}

func TestReadFile(t *testing.T) {
	fs := NewFS(dirs, files)
	path := "/tmp/test/hoge.txt"
	result, err := fs.ReadFile(path)
	if err != nil || result == nil || string(result) != "hogehoge" {
		t.Error("FileSystem.File should not return error if path exists")
	}

	path = "/tmp/test/notexists.txt"
	result, err = fs.ReadFile(path)

	if err == nil || result != nil {
		t.Error("FileSystem.File should return error if path doesn't exist")
	}
}

func TestReadDir(t *testing.T) {
	fs := NewFS(dirs, files)
	path := "/tmp/test"
	result, err := fs.ReadDir(path)
	if err != nil || result == nil || len(result) != 2 {
		t.Error("FileSystem.File should not return error if path exists")
	}

	path = "/tmp/test/notexists"
	result, err = fs.ReadDir(path)

	if err == nil || result != nil {
		t.Error("FileSystem.File should return error if path doesn't exist")
	}
}
