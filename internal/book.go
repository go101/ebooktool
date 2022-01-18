package internal

import (
	"fmt"
	"image/color"
	"path/filepath"
	"strings"
	"text/template"

	"go101.org/ebooktool/internal/nstd"
)

type BookInfo struct {
	OutputFormat string
	OutputPath   string // file path
	InputFormat  string
	InputPath    string // file path

	Title       string
	Author      string
	Version     string
	ReleaseDate string
	CoverImage  string // file path
	StyleCSS    string // file path

	VersionOnCover TextPlacement

	// For pdf format only
	EBookConvertor string
	MainFont       string
	MainFontSize   int32
	TocTitle       string
}

type TextPlacement struct {
	Text             string
	X, Y             int32
	AnchorX, AnchorY int8
	FontSize         int32
	Color            color.RGBA
}

func BuildBookInfoFromConfig(cfg *Config) (*BookInfo, error) {
	info := &BookInfo{}

	info.Title, _ = cfg.String("book.title")
	info.Author, _ = cfg.String("book.author")
	info.Version, _ = cfg.String("book.version")
	if info.Version == "" {
		info.Version, info.ReleaseDate = GetVersionAndDateFromGit(filepath.Dir(cfg.path))
	}

	info.CoverImage, _ = cfg.Path("book.cover-image")
	info.StyleCSS, _ = cfg.Path("book.style")

	info.OutputFormat, _ = cfg.String("output.format")
	info.OutputPath, _ = cfg.Path("output.path")
	info.InputFormat, _ = cfg.String("input.format")
	info.InputPath, _ = cfg.Path("input.path")

	info.EBookConvertor, _ = cfg.String("ebook.convertor")
	info.MainFont, _ = cfg.String("ebook.convertor.font.main")
	info.MainFontSize, _ = cfg.Int32("ebook.convertor.font-size.main")
	info.TocTitle, _ = cfg.String("ebook.convertor.toc-title")

	var okVersionOnCover bool
	info.VersionOnCover.Text, okVersionOnCover = cfg.String("cover-text.version")
	if okVersionOnCover && info.VersionOnCover.Text != "" {
		var err error
		info.VersionOnCover.X, info.VersionOnCover.AnchorX, err = cfg.CoordinateX("cover-text.version")
		if err != nil {
			return nil, err
		}
		info.VersionOnCover.Y, info.VersionOnCover.AnchorY, err = cfg.CoordinateY("cover-text.version")
		if err != nil {
			return nil, err
		}
		info.VersionOnCover.FontSize, _ = cfg.Int32("cover-text.version.font-size")
		info.VersionOnCover.Color, _ = cfg.Color("cover-text.version.color")
	}

	if info.OutputPath != "" && (info.Version != "" || info.ReleaseDate != "") {
		date := nstd.String(info.ReleaseDate).ReplaceAll("/", "-").String()
		txt, err := ExecuteTextTemplate(info.OutputPath, map[string]string{"Version": info.Version, "ReleaseDate": date})
		if err != nil {
			return nil, err
		}
		info.OutputPath = txt
	}

	// ToDo: remove the speciality
	switch info.OutputFormat {
	case "htmls":
	default:
		info.OutputPath += "." + info.OutputFormat
	}

	if info.VersionOnCover.Text != "" {
		if info.Version != "" || info.ReleaseDate != "" {
			txt, err := ExecuteTextTemplate(info.VersionOnCover.Text, map[string]string{"Version": info.Version, "ReleaseDate": info.ReleaseDate})
			if err != nil {
				return nil, err
			}
			info.VersionOnCover.Text = txt
		}
	}

	if info.CoverImage != "" && info.VersionOnCover.Text != "" {
		newPath, err := CreateImageWithOverlayTexts(info.CoverImage, info.VersionOnCover)
		if err != nil {
			return nil, fmt.Errorf("draw version text on cover error: %w", err)
		} else {
			info.CoverImage = newPath
		}
	}

	return info, nil
}

func ExecuteTextTemplate(text string, kvs map[string]string) (string, error) {
	tmpl, err := template.New("").Parse(text)
	if err != nil {
		return "", err
	}

	var b strings.Builder
	err = tmpl.Execute(&b, kvs)
	if err != nil {
		return "", err
	}

	return b.String(), nil
}
