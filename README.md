
`ebooktool` is command line tool to fill the gaps in building ebooks
by using `pandoc` and `calibre` tools. In other words, it is a supplyment
to the other tools.

Currently, `ebooktool` may convert the markdown files in a directory into
* an epub file.
* one single html file.
* multiple html files.
* a pdf files.
* a mobi file.
* an azw3 file.

Currently,
* the markdown files are sorted by alphabetical order.
* images are not handled.

An `.ini` file is needed to setup the conversion arguments.
The folowing are some example config files for a book (written in markdown):

The `ebooktool-common.inc` file:
```ini
; comment 1
# comment 2

book.title: Go Optimizations 101
book.author: Tapir Liu

; optional
; Recommanded ration: 1 : 1.6
; Recommanded sizes: 1000x1600, 1600x2560
book.cover-image: cover.png

; optional.
; Blank means trying to retrieve from git repo.
book.version: 

cover-text.version: -= {{.Version}}-{{.ReleaseDate}} =-
; x anchor: left | center | right
cover-text.version.center: 500
; y anchor: top | middle | bottom
cover-text.version.middle: 328
; in pixels (? to confirm)
cover-text.version.font-size: 29
; hex, default: #000
cover-text.version.color: #000
```

The `ebooktool-epub.ini` file:
```ini
; comment 1
# comment 2

; now, the included file and the including file must
; be within the same folder if there are relative
; file path values specified in the included file.
include: ebooktool-common.inc

; mobi | epub | azws | pdf | html | htmls
output.format: epub
; required. Auto suffixed with output.format.
output.path: _output/optimizations-101-{{.Version}}-{{.ReleaseDate}}

; only supports md (markdown) now
input.format: md
# the path to the directory containing the markdown files
input.path: .

# 
book.style: style.css
```

The `ebooktool-html.ini` file:
```ini
; comment 1
# comment 2

; now, the included file and the including file must
; be within the same folder if there are relative
; file path values specified in the included file.
include: ebooktool-common.inc

; mobi | epub | azws | pdf | html | htmls
output.format: html
; required. Auto suffixed with output.format.
output.path: _output/optimizations-101-{{.Version}}-{{.ReleaseDate}}

; only supports md (markdown) now
input.format: md
# the path to the directory containing the markdown files
input.path: .

# 
book.style: style.css
```

The `ebooktool-htmls.ini` file:
```
; comment 1
# comment 2

; now, the included file and the including file must
; be within the same folder if there are relative
; file path values specified in the included file.
include: ebooktool-common.inc

; mobi | epub | azws | pdf | html | htmls
# The generated htmls are the parts between <body> and </body>.
output.format: htmls
; required. Auto suffixed with output.format.
output.path: _output/optimizations-101-{{.Version}}-{{.ReleaseDate}}

; only supports md (markdown) now
input.format: md
# the path to the directory containing the markdown files
input.path: .

```

The `ebooktool-pdf.ini` file:
```ini
; comment 1
# comment 2

; now, the included file and the including file must
; be within the same folder if there are relative
; file path values specified in the included file.
include: ebooktool-common.inc

; mobi | epub | azws | pdf | html | htmls
output.format: pdf
; required. Auto suffixed with output.format.
output.path: _output/optimizations-101-{{.Version}}-{{.ReleaseDate}}

; only supports md (markdown) now
input.format: md
# the path to the directory containing the markdown files
input.path: .

# 
book.style: style.css

; An epub file will produced before generating the final pdf file,
; then the epub file will be converted to pdf by using 3rd tools.
; Supported tools: pandoc | calibre
ebook.convertor: pandoc
ebook.convertor.font.main: AR PL KaitiM GB
ebook.convertor.font-size.main: 15
ebook.convertor.toc-title: 目录
```

The `ebooktool-azw3.ini` file:
```ini
; comment 1
# comment 2

; now, the included file and the including file must
; be within the same folder if there are relative
; file path values specified in the included file.
include: ebooktool-common.inc

; mobi | epub | azw3 | pdf | html | htmls
output.format: azw3
; required. Auto suffixed with output.format.
output.path: _output/optimizations-101-{{.Version}}-{{.ReleaseDate}}

; only supports md (markdown) now
input.format: md
# the path to the directory containing the markdown files
input.path: .

# 
book.style: style.css

; An epub file will produced before generating the final azw3 file,
; then the epub file will be converted to azw3 by using 3rd tools.
; Supported tools: pandoc | calibre
ebook.convertor: pandoc


```

The `ebooktool-mobi.ini` file:
```ini
; comment 1
# comment 2

; now, the included file and the including file must
; be within the same folder if there are relative
; file path values specified in the included file.
include: ebooktool-common.inc

; mobi | epub | azws | pdf | html | htmls
output.format: mobi
; required. Auto suffixed with output.format.
output.path: _output/optimizations-101-{{.Version}}-{{.ReleaseDate}}

; only supports md (markdown) now
input.format: md
# the path to the directory containing the markdown files
input.path: .

# 
book.style: style.css

; An epub file will produced before generating the final mobi file,
; then the epub file will be converted to mobi by using 3rd tools.
; Supported tools: pandoc | calibre
ebook.convertor: pandoc
```