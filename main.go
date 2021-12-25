package main

import (
	"flag"
	"log"

	"go101.org/ebooktool/mds2epub"
	"go101.org/ebooktool/mds2html"
	"go101.org/ebooktool/mds2htmls"
	"go101.org/ebooktool/mds2more"
	"go101.org/ebooktool/mds2pdf"

	"go101.org/ebooktool/internal"
	"go101.org/ebooktool/internal/nstd"
)

var ToolVersion = "v0.0.1"

func printVersion() {
	log.Println("ebooktool", ToolVersion)
}

var (
	hFlag    = flag.Bool("h", false, "show help")
	helpFlag = flag.Bool("help", false, "show help")

	//vFlag = flag.Bool("v", false, "verbose mode")
	//verboseFlag = flag.Bool("verbose", false, "verbose mode")
)

func printUsage() {
	log.Println(`ebooktool ini-files...

The ini-files argument is opitional. Defaulted to ebooktool.ini.
All files listed must have a .ini suffix. 
`,
	)
}

func main() {
	log.SetFlags(0)
	flag.Parse()

	if *hFlag || *helpFlag {
		printUsage()
		return
	}

	if flag.NArg() == 1 && flag.Arg(0) == "version" {
		printVersion()
		return
	}

	var iniFiles []string
	if flag.NArg() == 0 {
		iniFiles = []string{"ebooktool.ini"}
	} else {
		iniFiles = make([]string, 0, flag.NArg())
		for i := 0; i < flag.NArg(); i++ {
			if arg := flag.Arg(i); false && !nstd.String(arg).ToLower().HasSuffix(".ini") {
				log.Printf("%s is not suffixed with .ini, so it is ignored", arg)
			} else {
				iniFiles = append(iniFiles, arg)
			}
		}

		if len(iniFiles) == 0 {
			log.Fatal("no book config files are specified")
		}
	}

	for _, iniFile := range iniFiles {
		log.SetPrefix("[" + iniFile + "]: ")
		config, err := internal.LoadIniFile(iniFile)
		if err != nil {
			log.Printf("ignores ini file for parsing error: %s", err)
			continue
		}

		bookInfo, err := internal.BuildBookInfoFromConfig(config)
		if err != nil {
			log.Printf("ignores ini file for info building error: %s", err)
			continue
		}

		if bookInfo.InputFormat != "md" {
			log.Printf("ignore ini file %s for unsupported input.format", bookInfo.InputFormat)
			continue
		}

		if bookInfo.OutputPath == "" {
			log.Printf("ignore ini file for not specifying output.path")
			continue
		}

		switch bookInfo.OutputFormat {
		default:
			log.Printf("ignore ini file for unknown output.format: %s", bookInfo.OutputFormat)
			continue
		case "":
			log.Printf("ignore ini file for not specifying output.format")
			continue
		case "htmls":
			err = mds2htmls.Run(bookInfo)
		case "html":
			err = mds2html.Run(bookInfo, false, false)
		case "epub":
			err = mds2epub.Run(bookInfo)
		case "pdf":
			err = mds2pdf.Run(bookInfo)
		case "azw3", "mobi", "docx":
			err = mds2more.Run(bookInfo)
		}

		if err != nil {
			log.Println("failed to generate", bookInfo.OutputPath, ", err:", err)
		} else {
			log.Println("successfully generated", bookInfo.OutputPath)
		}
	}

	log.SetPrefix("")
	log.Println("Done.")
}
