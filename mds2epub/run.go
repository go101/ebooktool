package mds2epub

import (
	"fmt"
	"os"

	"go101.org/ebooktool/internal"
	"go101.org/ebooktool/internal/nstd"
)

func Run(bookInfo *internal.BookInfo) error {
	epubFile := bookInfo.OutputPath
	mdsDir := bookInfo.InputPath
	cssFile := bookInfo.StyleCSS
	coverFile := bookInfo.CoverImage
	bookTitle := bookInfo.Title
	bookAuthor := bookInfo.Author

	epubBook := internal.NewEpubFile()

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
		of := nstd.String(f.Name).ToLower().TrimSuffix(".md").String()
		of = internal.Filename2ID(of) + ".xhtml"
		md.SetFilename(f.Name)
		md.SetOutputFile(of)
		mdFiles[i] = *md
		md.ConvertCodeLineLeadingTabsToSpaces(3)
	}

	imageFiles := internal.CollectMarkdownImageFiles(mdFiles, mdsDir, false)
	for path := range imageFiles {
		internalPath, err := epubBook.AddImage(path)
		if err != nil {
			return fmt.Errorf("add epub image file (%s) error: %s", path, err)
		}

		imageFiles[path] = []byte(internalPath)
	}

	internal.ReplaceMarkdownImageHrefs(mdFiles, mdsDir, imageFiles)

	// render

	renderedFiles := internal.RenderMarkdownFiles(mdFiles)
	if len(renderedFiles) != len(mdFiles) {
		for outfile, files := range renderedFiles {
			println(outfile)
			for _, f := range files {
				println("   ", f.Name)
			}
		}
		return fmt.Errorf("expected %d output files, but got %d", len(mdFiles), len(renderedFiles))
	}

	// assemble

	sections := make([]internal.EpubSection, len(mdFiles))
	for i, md := range mdFiles {
		title := md.Title()
		if title == "" {
			title = fmt.Sprintf("chapter %d", i)
		}
		sections[i] = internal.EpubSection{
			File: internal.File{
				Name:    md.OutputFile(),
				Content: renderedFiles[md.OutputFile()][0].Content,
			},
			Title: title,
		}
	}

	// write

	epubTitle := bookTitle
	if bookInfo.ReleaseDate != "" {
		epubTitle += " (" + bookInfo.ReleaseDate + ")"
	}

	switch c := bookInfo.EBookConvertor; c {
	case "":
		return epubBook.CreateEpubFile(epubFile, epubTitle, bookAuthor, cssFile, coverFile, sections)
	case "pandoc":
		break
	default:
		return fmt.Errorf("ebook converter %s is not supported to convert epub to epub now", c)
	}

	tempEpubFile := epubFile + ".temp-" + internal.RandomString(8) + ".epub"
	err = epubBook.CreateEpubFile(tempEpubFile, epubTitle, bookAuthor, cssFile, coverFile, sections)
	if err != nil {
		return err
	}
	defer func() {
		os.Remove(tempEpubFile)
	}()

	return internal.PandocX2Y(epubFile, tempEpubFile)
}
