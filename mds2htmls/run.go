package mds2htmls

import (
	"fmt"
	"os"

	"go101.org/ebooktool/internal"
	"go101.org/ebooktool/internal/nstd"
)

func Run(bookInfo *internal.BookInfo) error {
	htmlsDir := bookInfo.OutputPath
	mdsDir := bookInfo.InputPath

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
		of = internal.ValidateIdentifier(of) + ".html"
		md.SetFilename(f.Name)
		md.SetOutputFile(of)
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

	// write

	if err := os.MkdirAll(htmlsDir, 0700); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", htmlsDir, err)
	}

	outputFiles := make(map[string][]byte, len(mdFiles))
	for _, md := range mdFiles {
		outputFiles[md.OutputFile()] = renderedFiles[md.OutputFile()][0].Content
	}

	return internal.WriteFiles(htmlsDir, outputFiles)
}
