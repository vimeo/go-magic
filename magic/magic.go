package magic

/*
 #cgo LDFLAGS: -lmagic
 #include <magic.h>
 #include <stdlib.h>
*/
import "C"
import (
    "os"
    "path"
    "unsafe"
)

const (
    MAGIC_NONE              = C.MAGIC_NONE
    MAGIC_DEBUG             = C.MAGIC_DEBUG
    MAGIC_SYMLINK           = C.MAGIC_SYMLINK
    MAGIC_COMPRESS          = C.MAGIC_COMPRESS
    MAGIC_DEVICES           = C.MAGIC_DEVICES
    MAGIC_MIME_TYPE         = C.MAGIC_MIME_TYPE
    MAGIC_CONTINUE          = C.MAGIC_CONTINUE
    MAGIC_CHECK             = C.MAGIC_CHECK
    MAGIC_PRESERVE_ATIME    = C.MAGIC_PRESERVE_ATIME
    MAGIC_RAW               = C.MAGIC_RAW
    MAGIC_ERROR             = C.MAGIC_ERROR
    MAGIC_MIME_ENCODING     = C.MAGIC_MIME_ENCODING
    MAGIC_MIME              = C.MAGIC_MIME
    MAGIC_APPLE             = C.MAGIC_APPLE
    MAGIC_NO_CHECK_COMPRESS = C.MAGIC_NO_CHECK_COMPRESS
    MAGIC_NO_CHECK_TAR      = C.MAGIC_NO_CHECK_TAR
    MAGIC_NO_CHECK_SOFT     = C.MAGIC_NO_CHECK_SOFT
    MAGIC_NO_CHECK_APPTYPE  = C.MAGIC_NO_CHECK_APPTYPE
    MAGIC_NO_CHECK_ELF      = C.MAGIC_NO_CHECK_ELF
    MAGIC_NO_CHECK_TEXT     = C.MAGIC_NO_CHECK_TEXT
    MAGIC_NO_CHECK_CDF      = C.MAGIC_NO_CHECK_CDF
    MAGIC_NO_CHECK_TOKENS   = C.MAGIC_NO_CHECK_TOKENS
    MAGIC_NO_CHECK_ENCODING = C.MAGIC_NO_CHECK_ENCODING
    MAGIC_NO_CHECK_ASCII    = C.MAGIC_NO_CHECK_ASCII
    MAGIC_NO_CHECK_FORTRAN  = C.MAGIC_NO_CHECK_FORTRAN
    MAGIC_NO_CHECK_TROFF    = C.MAGIC_NO_CHECK_TROFF
)
const (
    MAGIC_NO_CHECK_BUILTIN  = MAGIC_NO_CHECK_COMPRESS |
                              MAGIC_NO_CHECK_TAR      |
                              MAGIC_NO_CHECK_APPTYPE  |
                              MAGIC_NO_CHECK_ELF      |
                              MAGIC_NO_CHECK_TEXT     |
                              MAGIC_NO_CHECK_CDF      |
                              MAGIC_NO_CHECK_TOKENS   |
                              MAGIC_NO_CHECK_ENCODING
)

type Magic C.magic_t

var system_mgc_locations = []string {
    "/usr/share/misc/magic.mgc",
    "/usr/share/file/magic.mgc",
    "/usr/share/magic/magic.mgc",
}

/* Find the real magic file location */
func GetDefaultDir() string {
    var f string

    found_mgc := false
    for _, f = range system_mgc_locations {
        fi, err := os.Lstat(f)
        if err == nil && fi.Mode() & os.ModeSymlink != os.ModeSymlink {
            found_mgc = true
            break
        }
    }
    if found_mgc {
        return path.Dir(f)
    } else {
        return ""
    }
}

func Open(flags int) Magic {
    cookie := (Magic)(C.magic_open(C.int(flags)))
    return cookie;
}

func Close(cookie Magic) {
    C.magic_close((C.magic_t)(cookie))
}

func Error(cookie Magic) string {
    s := (C.magic_error((C.magic_t)(cookie)))
    return C.GoString(s)
}

func Errno(cookie Magic) int {
    return (int)(C.magic_errno((C.magic_t)(cookie)))
}

func File(cookie Magic, filename string) string {
    cfilename := C.CString(filename)
    defer C.free(unsafe.Pointer(cfilename))
    return C.GoString(C.magic_file(cookie, cfilename))
}

func Buffer(cookie Magic, b []byte) string {
    length := C.size_t(len(b))
    return C.GoString(C.magic_buffer(cookie, unsafe.Pointer(&b[0]), length))
}

func SetFlags(cookie Magic, flags int) int {
    return (int)(C.magic_setflags(cookie, C.int(flags)))
}

func Check(cookie Magic, filename string) int {
    cfilename := C.CString(filename)
    defer C.free(unsafe.Pointer(cfilename))
    return (int)(C.magic_check(cookie, cfilename))
}

func Compile(cookie Magic, filename string) int {
    cfilename := C.CString(filename)
    defer C.free(unsafe.Pointer(cfilename))
    return (int)(C.magic_compile(cookie, cfilename))
}

func Load(cookie Magic, filename string) int {
    if filename == "" {
        return (int)(C.magic_load(cookie, (*C.char)(unsafe.Pointer(uintptr(0)))))
    }
    cfilename := C.CString(filename)
    defer C.free(unsafe.Pointer(cfilename))
    return (int)(C.magic_load(cookie, cfilename))
}
