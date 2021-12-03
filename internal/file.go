package internal

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"go101.org/ebooktool/internal/nstd"
)

type File struct {
	// absolute or relative path, by context
	Name    string
	Content []byte
}

func ReadDirFiles(dir string, filter func(filename string) bool) ([]File, error) {
	dirfs := os.DirFS(dir)
	entries, err := fs.ReadDir(dirfs, ".")
	if err != nil {
		return nil, fmt.Errorf("fs.ReadDir(%s) error: %w", dir, err)
	}

	files := make([]File, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		if filter(e.Name()) {
			content, err := fs.ReadFile(dirfs, e.Name())
			if err != nil {
				return nil, fmt.Errorf("fs.ReadFile(%s) error: %w", e.Name(), err)
			}

			files = append(files, File{
				Name:    e.Name(),
				Content: content,
			})
		}
	}

	return files, nil
}

func ReadFiles(filepaths []string) ([]File, error) {
	files := make([]File, 0, len(filepaths))
	for _, path := range filepaths {
		content, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("os.ReadFile(%s) error: %w", path, err)
		}

		files = append(files, File{
			Name:    path,
			Content: content,
		})
	}

	return files, nil
}

func WriteFiles(outputDir string, outputFiles map[string][]byte) error {
	for outFile, content := range outputFiles {
		fileFullPath := filepath.Join(outputDir, outFile)
		err := os.WriteFile(fileFullPath, content, 0600)
		if err != nil {
			return err
		}
	}
	return nil
}

func MergeFileContents(files ...File) []byte {
	n := 0
	for _, f := range files {
		n += len(f.Content)
	}

	all := make([]byte, 0, n)
	for _, f := range files {
		all = append(all, f.Content...)
	}

	return all
}

func patternizeFilename(filename string) string {
	k := nstd.String(filename).LastIndex(".")
	if k >= 0 {
		return filename[:k] + "-*" + filename[k:]
	}
	return filename + "-*"
}
