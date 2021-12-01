package mds2html

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html"
	"os"
	"path/filepath"

	"go101.org/ebooktool/internal"
	"go101.org/ebooktool/internal/nstd"
)

func Run(bookInfo *internal.BookInfo, forPDF bool) error {
	htmlFile := bookInfo.OutputPath
	mdsDir := bookInfo.InputPath
	cssFile := bookInfo.StyleCSS
	coverFile := bookInfo.CoverImage
	bookTitle := bookInfo.Title
	bookAuthor := bookInfo.Author

	// ...

	outputFile := filepath.Base(htmlFile)
	outputDir := filepath.Dir(htmlFile)

	// load

	files, err := internal.ReadFiles(mdsDir, func(filename string) bool {
		return nstd.String(filename).ToLower().HasSuffix(".md")
	})
	if err != nil {
		return err
	}

	// parse

	mdFiles := make([]internal.Markdown, len(files))
	for i, f := range files {
		md, err := internal.ParseMarkdown(f.Content)
		if err != nil {
			return fmt.Errorf("parse markdown file (%s) error: %w", f.Name, err)
		}
		md.SetFilename(f.Name)
		md.SetOutputFile(outputFile)
		mdFiles[i] = *md
	}

	// render

	renderedFiles := internal.RenderMarkdownFiles(mdFiles)
	if len(renderedFiles) != 1 {
		for outfile, files := range renderedFiles {
			println(outfile)
			for _, f := range files {
				println("   ", f.Name)
			}
		}
		return fmt.Errorf("expected one output file, but got %d", len(renderedFiles))
	}

	// write

	var authorMeta string
	if bookAuthor != "" {
		authorMeta = fmt.Sprintf(`<meta name="author" content="%s">`, html.EscapeString(bookAuthor))
	}

	var coverImageHTML bytes.Buffer
	if coverFile != "" {
		var imgType string
		if nstd.String(coverFile).ToLower().HasSuffix(".png") {
			imgType = "png"
		} else if nstd.String(coverFile).ToLower().HasSuffix(".gif") {
			imgType = "gif"
		} else if nstd.String(coverFile).ToLower().HasSuffix(".jpg") ||
			nstd.String(coverFile).ToLower().HasSuffix(".jpeg") {
			imgType = "jpeg"
		} else {
			return fmt.Errorf("unsupported image mime type: %s", coverFile)
		}

		imgData, err := os.ReadFile(coverFile)
		if err != nil {
			return err
		}
		base64Data := make([]byte, (len(imgData)+2)/3*4)
		base64.StdEncoding.Encode(base64Data, imgData)

		if forPDF {
			coverImageHTML.WriteString(`<img style="width: auto; max-height: 100%;" src="data:image/`)
		} else {
			coverImageHTML.WriteString(`<img style="max-width: 100%; height: auto;" src="data:image/`)
		}
		coverImageHTML.WriteString(imgType)
		coverImageHTML.WriteString(`;base64,`)
		coverImageHTML.Write(base64Data)
		coverImageHTML.WriteString(`"></img>`)
	}

	var tocContent bytes.Buffer

	var cssContent []byte
	if cssFile != "" {
		cssContent, err = os.ReadFile(cssFile)
		if err != nil {
			return fmt.Errorf("failed to load css file %s: %w", cssFile, err)
		}
	}

	var start bytes.Buffer
	fmt.Fprintf(&start, htmlStart, html.EscapeString(bookTitle), authorMeta, cssContent)

	bookContent := internal.MergeFileContents(renderedFiles[outputFile]...)
	content := nstd.MergeByteSlices(start.Bytes(), coverImageHTML.Bytes(), tocContent.Bytes(), bookContent, []byte(htmlEnd))
	outputFiles := map[string][]byte{
		outputFile: content,
	}

	return internal.WriteFiles(outputDir, outputFiles)
}

var htmlStart = `
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">
<head>
<title>%s</title>
%s
<meta charset="utf-8" />
<style>
%s
</style>
</head>
<body>
`

const htmlEnd = `
</body>
</html>
`
