package internal

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/bmaupin/go-epub"
)

// Better to fork the epub lib?
// 1. totally memory operations
// 2. section content is []byte instead of string
// 3. avoid duplicate cover when converting to pdf
// 4. multi-level table of content

type EpubFile struct {
	f *epub.Epub
}

func NewEpubFile() *EpubFile {
	return &EpubFile{
		f: epub.NewEpub(""),
	}
}

func (ef *EpubFile) AddImage(path string) (string, error) {
	return ef.f.AddImage(path, "")
}

type EpubSection struct {
	File
	Title string
}

func (ef *EpubFile) CreateEpubFile(outputPath, title, author, cssFile, coverFile string, sections []EpubSection) error {
	//epubBook := epub.NewEpub(title)
	epubBook := ef.f
	epubBook.SetTitle(title)
	epubBook.SetAuthor(author)
	var cssPath string
	if cssFile != "" {
		var err error
		cssPath, err = epubBook.AddCSS(cssFile, filepath.Base(cssFile))
		if err != nil {
			return fmt.Errorf("add epub css file %s err: %w", cssFile, err)
		}
	}
	if coverFile != "" {
		coverPath, err := epubBook.AddImage(coverFile, filepath.Base(coverFile))
		if err != nil {
			return fmt.Errorf("add epub image file %s err: %w", coverFile, err)
		}
		epubBook.SetCover(coverPath, "")
	}

	for _, s := range sections {
		_, err := epubBook.AddSection(string(s.Content), s.Title, s.Name, cssPath)
		if err != nil {
			return fmt.Errorf("add epub section file %s err: %w", s.Name, err)
		}
	}

	if err := epubBook.Write(outputPath); err != nil {
		return fmt.Errorf("write epub file %s err: %w", outputPath, err)
	}

	return nil
}

func RemovePagesFromEpub(epubFile string, pagesToRemove ...string) error {
	r, err := zip.OpenReader(epubFile)
	if err != nil {
		return err
	}
	defer r.Close()

	os.Remove(epubFile)

	outputFile, err := os.Create(epubFile)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	w := zip.NewWriter(outputFile)

	shouldRemove := map[string]bool{}
	for _, page := range pagesToRemove {
		shouldRemove[page] = true
	}

	// Iterate through the files in the archive,
	// printing some of their contents.
	for _, f := range r.File {
		//log.Printf("Contents of %s:\n", f.Name)
		if shouldRemove[f.Name] {
			continue
		}

		if err := func() error {
			rc, err := f.Open()
			if err != nil {
				return err
			}
			defer rc.Close()

			of, err := w.Create(f.Name)
			if err != nil {
				return err
			}

			_, err = io.Copy(of, rc)
			if err != nil {
				return err
			}

			return nil
		}(); err != nil {
			return err
		}
	}

	err = w.Close()
	if err != nil {
		return err
	}

	return outputFile.Sync()
}
