package internal

import (
	"log"
	"time"
)

func PandocX2Y(outputFilename, tempEpubFile string) error {
	conversionParameters := make([]string, 0, 32)
	pushParams := func(params ...string) {
		conversionParameters = append(conversionParameters, params...)
	}

	pushParams("pandoc", "-o", outputFilename, tempEpubFile)

	output, err := ExecCommand(time.Minute*5, ".", nil, conversionParameters...)
	if err != nil {
		log.Printf("%s\n%s", output, err)
		return err
	}

	return nil
}

func CalibreX2Y(outputFilename, inputFilename string) error {
	conversionParameters := make([]string, 0, 32)
	pushParams := func(params ...string) {
		conversionParameters = append(conversionParameters, params...)
	}
	pushParams("ebook-convert", inputFilename, outputFilename)

	output, err := ExecCommand(time.Minute*5, ".", nil, conversionParameters...)
	if err != nil {
		log.Printf("%s\n%s", output, err)
		return err
	}

	return nil
}
