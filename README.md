# go-magic #

## Go library for getting MIME type using libmagic ##

### Installing ###

```
go get github.com/vimeo/go-magic/magic
```

### Dependencies ###

**libmagic**<br />
*URL*: [http://www.darwinsys.com/file/](http://www.darwinsys.com/file/)<br />
*Ubuntu*: `apt-get install libmagic-dev`<br />
*CentOS*: `yum install file-devel`<br />

### Usage ###

- Create some custom magic files (e.g. ~/magicfiles)
- Add the default system magic file dir
    - magic.AddMagicDir(magic.GetDefaultDir())
- Add the custom magic file dir
    - magic.AddMagicDir("~/magicfiles")
- Get MIME type with either one of:
    - magic.MimeFromFile(filename)
    - magic.MimeFromBytes(data)

### API Documentation ###

[http://godoc.org/github.com/vimeo/go-magic/magic](http://godoc.org/github.com/vimeo/go-magic/magic)
