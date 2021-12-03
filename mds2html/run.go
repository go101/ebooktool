package mds2html

import (
	"bytes"
	"fmt"
	"html"
	"os"
	"path/filepath"

	"go101.org/ebooktool/internal"
	"go101.org/ebooktool/internal/nstd"
)

func Run(bookInfo *internal.BookInfo, forPDF, convertCodeLeadingTabsToSpaces bool) error {
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

	files, err := internal.ReadDirFiles(mdsDir, func(filename string) bool {
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
		if convertCodeLeadingTabsToSpaces {
			md.ConvertCodeLeadingTabsToSpaces(5)
		}
	}

	imageFiles := internal.CollectMarkdownImageFiles(mdFiles, mdsDir, true)
	for path, content := range imageFiles {
		base64Data, err := internal.Base64Image(path, content)
		if err != nil {
			return err
		}

		imageFiles[path] = base64Data
	}

	internal.ReplaceMarkdownImageHrefs(mdFiles, mdsDir, imageFiles)

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
		imgData, err := os.ReadFile(coverFile)
		if err != nil {
			return err
		}

		base64Data, err := internal.Base64Image(coverFile, imgData)
		if err != nil {
			return err
		}

		if forPDF {
			coverImageHTML.WriteString(`<img style="width: auto; max-height: 100%;" src="`)
		} else {
			coverImageHTML.WriteString(`<img style="max-width: 100%; height: auto;" src="`)
		}
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
