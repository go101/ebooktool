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
		tempHtmlFile := outputFile + ".temp-" + internal.RandomString(8) + ".html"
		bookInfo.OutputPath = tempHtmlFile
		defer func() {
			bookInfo.OutputPath = outputFile
			os.Remove(tempHtmlFile)
		}()

		err := mds2html.Run(bookInfo, false, false)
		if err != nil {
			return err
		}

		return internal.PandocX2Y(outputFile, tempHtmlFile)

	case "calibre":
		// create a temp epub file
		tempEpubFile := outputFile + ".temp-" + internal.RandomString(8) + ".epub"
		bookInfo.OutputPath = tempEpubFile
		defer func() {
			bookInfo.OutputPath = outputFile
			os.Remove(tempEpubFile)
		}()

		bookInfo.EBookConvertor = "" // to genetate non-validated epub file
		err := mds2epub.Run(bookInfo)
		if err != nil {
			return err
		}

		return internal.CalibreX2Y(outputFile, tempEpubFile)
	}
}
