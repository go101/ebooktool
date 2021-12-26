package mds2more

import (
	"fmt"
	"os"

	"go101.org/ebooktool/internal"
	"go101.org/ebooktool/mds2epub"
	"go101.org/ebooktool/mds2html"
)

func Run(bookInfo *internal.BookInfo) error {
	outputFile := bookInfo.OutputPath

	// convert epub file to pdf

	switch bookInfo.EBookConvertor {
	default:
		return fmt.Errorf("ebook convertor (epub->pdf) is unspecified or unrecognized")

	case "pandoc":
		// create a temp html file
		htmlBookInfo := *bookInfo
		htmlBookInfo.OutputPath = outputFile + ".temp-" + internal.RandomString(8) + ".html"
		defer os.Remove(htmlBookInfo.OutputPath)
		err := mds2html.Run(&htmlBookInfo, false, false)
		if err != nil {
			return err
		}

		return internal.PandocX2Y(outputFile, htmlBookInfo.OutputPath)

	case "calibre":
		// create a temp epub file
		epubBookInfo := *bookInfo
		epubBookInfo.EBookConvertor = "" // to genetate non-validated epub file
		epubBookInfo.OutputPath = outputFile + ".temp-" + internal.RandomString(8) + ".epub"
		defer os.Remove(epubBookInfo.OutputPath)
		err := mds2epub.Run(&epubBookInfo)
		if err != nil {
			return err
		}

		return internal.CalibreX2Y(outputFile, epubBookInfo.OutputPath)
	}
}
