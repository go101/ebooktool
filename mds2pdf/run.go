package mds2pdf

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"go101.org/ebooktool/internal"
	"go101.org/ebooktool/mds2epub"
	"go101.org/ebooktool/mds2html"
)

func Run(bookInfo *internal.BookInfo) error {
	pdfFile := bookInfo.OutputPath

	// convert epub file to pdf

	switch bookInfo.EBookConvertor {
	default:
		return fmt.Errorf("ebook convertor (epub->pdf) is unspecified or unrecognized")

	case "pandoc":
		// create a temp html file
		tempHtmlFile := pdfFile + ".temp-" + internal.RandomString(8) + ".html"
		bookInfo.OutputPath = tempHtmlFile
		defer func() {
			bookInfo.OutputPath = pdfFile
			os.Remove(tempHtmlFile)
		}()

		err := mds2html.Run(bookInfo, true, true)
		if err != nil {
			return err
		}

		return epub2pdf_pandoc(false, pdfFile, tempHtmlFile, bookInfo.MainFont, bookInfo.MainFontSize, bookInfo.TocTitle)

	case "calibre":
		// create a temp epub file
		tempEpubFile := pdfFile + ".temp-" + internal.RandomString(8) + ".epub"
		bookInfo.OutputPath = tempEpubFile
		defer func() {
			bookInfo.OutputPath = pdfFile
			os.Remove(tempEpubFile)
		}()

		err := mds2epub.Run(bookInfo)
		if err != nil {
			return err
		}

		// Remove the cover page to avoid dupliated covers in pdf.
		err = internal.RemovePagesFromEpub(tempEpubFile, "EPUB/xhtml/cover.xhtml")
		if err != nil {
			return err
		}
		return epub2pdf_calibre(false, pdfFile, tempEpubFile, bookInfo.MainFont, bookInfo.MainFontSize, "", 0)
	}
}

func epub2pdf_pandoc(forPrinting bool, outputFilename, tempEpubFile string, mainFont string, mainFontSize int32, tocTitle string) error {
	conversionParameters := make([]string, 0, 32)
	pushParams := func(params ...string) {
		conversionParameters = append(conversionParameters, params...)
	}
	pushParams("pandoc", "-s", "--pdf-engine", "xelatex", "--toc", "--number-sections")
	//pushParams("--number-offset", "-1,0") // useless
	pushParams("-V", tocTitle)
	//pushParams("-V", "papersize:a4")
	pushParams("-V", "documentclass=report")
	pushParams("-V", "colorlinks")
	pushParams("-V", "linkcolor=blue")
	pushParams("-V", "urlcolor=blue")
	pushParams("-V", "toccolor=blue")
	pushParams("-V", "geometry: top=3cm, bottom=3cm, left=3.9cm, right=3.9cm")
	pushParams("--tab-stop", "5")
	if mainFont != "" {
		pushParams("-V", fmt.Sprintf(`CJKmainfont=%s`, mainFont))
		pushParams("-V", fmt.Sprintf(`mainfont=%s`, mainFont))
	}
	pushParams("-V", fmt.Sprintf(`fontsize=%spt`, mainFontSize))

	pushParams("-o", outputFilename, tempEpubFile)

	output, err := internal.ExecCommand(time.Minute*5, ".", nil, conversionParameters...)
	if err != nil {
		log.Printf("%s\n%s", output, err)
		return err
	}

	return nil
}

func epub2pdf_calibre(forPrinting bool, outputFilename, inputFilename string, serifFont string, serifFontSize int32, monoFoint string, monoFontSize int32) error {
	conversionParameters := make([]string, 0, 32)
	pushParams := func(params ...string) {
		conversionParameters = append(conversionParameters, params...)
	}
	pushParams("ebook-convert", inputFilename, outputFilename)
	//if contentIndexTitle != "" {
	//	pushParams("--toc-title", contentIndexTitle)
	//}
	//if false {
	//	// It looks the latest calibre version doesn't correctly center the texts.
	//	// And it becomes extremelu slow to generate Chinese book versions.
	//	pushParams("--pdf-header-template", `<div style='text-align: center; font-size: small;'>_SECTION_</div>`)
	pushParams("--pdf-footer-template", `<div style='text-align: center; font-size: small;'>_PAGENUM_</div>`)
	//}

	//pushParams("--pdf-page-numbers")
	pushParams("--paper-size", "a4")
	if serifFont != "" {
		pushParams("--pdf-serif-family", serifFont)
	}
	//pushParams("--pdf-sans-family", serifFont)
	if monoFoint != "" {
		pushParams("--pdf-mono-family", monoFoint) // "Liberation Mono")
	}
	if serifFontSize == 0 {
		serifFontSize = 15
	}
	pushParams("--pdf-default-font-size", strconv.Itoa(int(serifFontSize)))
	if monoFontSize == 0 {
		monoFontSize = serifFontSize - 1
	}
	pushParams("--pdf-mono-font-size", strconv.Itoa(int(monoFontSize)))

	//pushParams("--pdf-page-margin-top", "36")
	//pushParams("--pdf-page-margin-bottom", "36")
	//if forPrinting {
	//	pushParams("--pdf-add-toc")
	//	pushParams("--pdf-page-margin-left", "72")
	//	pushParams("--pdf-page-margin-right", "72")
	//} else {
	//	pushParams("--pdf-page-margin-left", "36")
	//	pushParams("--pdf-page-margin-right", "36")
	//}
	pushParams("--pdf-page-margin-left", "72")
	pushParams("--pdf-page-margin-right", "72")
	pushParams("--pdf-page-margin-top", "72")
	pushParams("--pdf-page-margin-bottom", "72")
	pushParams("--preserve-cover-aspect-ratio")
	//pushParams("--use-auto-toc")

	output, err := internal.ExecCommand(time.Minute*5, ".", nil, conversionParameters...)
	if err != nil {
		log.Printf("%s\n%s", output, err)
		return err
	}

	return nil
}
