package magic

import (
    "fmt"
    "io/ioutil"
    "os"
    "path"
    "strings"
)

var magicFiles []string

func compileToMgc(file string) error {
    cookie := Open(MAGIC_NONE)
    defer Close(cookie)
    ret := Compile(cookie, file)
    if ret != 0 {
        return fmt.Errorf("Error compiling magic file: %s", file)
    }
    return nil
}

/* Add a directory for libmagic to search for .mgc databases. */
func AddMagicDir(dir string) error {
    pwd, err := os.Getwd()
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
    err = os.Chdir(dir)
    if err != nil {
        return err
    }
    files, err := ioutil.ReadDir(dir)
    if err != nil {
        return err
    }
    for _, fi = range files {
        if path.Ext(fi.Name()) == ".magic" {
            mgcSrc := path.Join(dir, fi.Name())
            _, err := os.Stat(mgcSrc + ".mgc")
            if err != nil {
                compileToMgc(path.Join(dir, fi.Name()))
            }
        }
    }
    files, err = ioutil.ReadDir(dir)
    if err != nil {
        return err
    }
    for _, fi = range files {
        if path.Ext(fi.Name()) == ".mgc" {
            mgcFile := path.Join(dir, fi.Name())
            magicFiles = append(magicFiles, mgcFile)
        }
    }
    err = os.Chdir(pwd)
    if err != nil {
        return err
    }
    return nil
}

/* Get mimetype from a file. */
func MimeFromFile(path string) string {
    cookie := Open(MAGIC_ERROR | MAGIC_MIME_TYPE)
    defer Close(cookie)
    ret := Load(cookie, strings.Join(magicFiles, ":"))
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
    ret := Load(cookie, strings.Join(magicFiles, ":"))
    if ret != 0 {
        return "application/octet-stream"
    }
    r := Buffer(cookie, b)
    return r
}
