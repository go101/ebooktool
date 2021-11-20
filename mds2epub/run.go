package mds2epub

import (
	"fmt"

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
		md.SetOutputFile(nstd.String(f.Name).ToLower().ReplaceSuffix(".md", ".xhtml").String())
		mdFiles[i] = *md
	}

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
	return internal.CreateEpubFile(epubFile, epubTitle, bookAuthor, cssFile, coverFile, sections)
}
