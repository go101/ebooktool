package mds2more

import (
	"fmt"
	"log"
	"os"
	"time"

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

		return epub2more_pandoc(outputFile, tempHtmlFile)

	case "calibre":
		// create a temp epub file
		tempEpubFile := outputFile + ".temp-" + internal.RandomString(8) + ".epub"
		bookInfo.OutputPath = tempEpubFile
		defer func() {
			bookInfo.OutputPath = outputFile
			os.Remove(tempEpubFile)
		}()

		err := mds2epub.Run(bookInfo)
		if err != nil {
			return err
		}

		return epub2more_calibre(outputFile, tempEpubFile)
	}
}

func epub2more_pandoc(outputFilename, tempEpubFile string) error {
	conversionParameters := make([]string, 0, 32)
	pushParams := func(params ...string) {
		conversionParameters = append(conversionParameters, params...)
	}

	pushParams("pandoc", "-o", outputFilename, tempEpubFile)

	output, err := internal.ExecCommand(time.Minute*5, ".", nil, conversionParameters...)
	if err != nil {
		log.Printf("%s\n%s", output, err)
		return err
	}

	return nil
}

func epub2more_calibre(outputFilename, inputFilename string) error {
	conversionParameters := make([]string, 0, 32)
	pushParams := func(params ...string) {
		conversionParameters = append(conversionParameters, params...)
	}
	pushParams("ebook-convert", inputFilename, outputFilename)

	output, err := internal.ExecCommand(time.Minute*5, ".", nil, conversionParameters...)
	if err != nil {
		log.Printf("%s\n%s", output, err)
		return err
	}

	return nil
}
