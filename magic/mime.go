package magic

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var magicFiles map[string]bool
var mutex sync.Mutex

func init() {
	magicFiles = make(map[string]bool)
}

func compileToMgc(file string) {
	cookie := Open(MAGIC_NONE)
	defer Close(cookie)
	Compile(cookie, file)
}

func compileMagicFiles(dir string, files []string) {
	// for some reason libmagic puts compiled files in the current working dir
	// instead of the dir of the source file, so switch to the source dir,
	// compile, then switch back
	pwd, err := os.Getwd()
	if err != nil {
		return
	}
	err = os.Chdir(dir)
	if err != nil {
		return
	}
	defer os.Chdir(pwd)

	for _, f := range files {
		compileToMgc(f)
	}
}

/* Add a directory for libmagic to search for .mgc databases. */
func AddMagicDir(dir string) error {
	var err error

	dir, err = filepath.Abs(dir)
	if err != nil {
		return err
	}
	fi, err := os.Stat(dir)
	if err != nil {
		return err
	}
	if fi.IsDir() == false {
		return fmt.Errorf("Not a directory: %s", dir)
	}

	// get list of .magic files that need to be compiled to .mgc
	var srcFiles []string
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, fi = range files {
		if filepath.Ext(fi.Name()) == ".magic" {
			mgcSrc := filepath.Join(dir, fi.Name())
			_, err := os.Stat(mgcSrc + ".mgc")
			if err != nil {
				srcFiles = append(srcFiles, mgcSrc)
			}
		}
	}
	// compile .magic files
	if len(srcFiles) > 0 {
		compileMagicFiles(dir, srcFiles)
	}

	files, err = ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	mutex.Lock()
	for _, fi = range files {
		if filepath.Ext(fi.Name()) == ".mgc" {
			mgcFile := filepath.Join(dir, fi.Name())
			magicFiles[mgcFile] = true
		}
	}
	mutex.Unlock()

	return nil
}

/* Get mimetype from a file. */
func MimeFromFile(path string) string {
	cookie := Open(MAGIC_ERROR | MAGIC_MIME_TYPE)
	defer Close(cookie)
	mutex.Lock()
	var mf []string
	for f := range magicFiles {
		mf = append(mf, f)
	}
	mutex.Unlock()
	ret := Load(cookie, strings.Join(mf, ":"))
	if ret != 0 {
		return "application/octet-stream"
	}
	r := File(cookie, path)
	return r
}

/* Get mimetype from a buffer. */
func MimeFromBytes(b []byte) string {
	cookie := Open(MAGIC_ERROR | MAGIC_MIME_TYPE)
	defer Close(cookie)
	mutex.Lock()
	var mf []string
	for f := range magicFiles {
		mf = append(mf, f)
	}
	mutex.Unlock()
	ret := Load(cookie, strings.Join(mf, ":"))
	if ret != 0 {
		return "application/octet-stream"
	}
	r := Buffer(cookie, b)
	return r
}
