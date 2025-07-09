package main

import (
	"flag"
	"fmt"
	"google-translate/cli"
	"os"
	"strings"
	"sync"
)

var wg sync.WaitGroup

var sourceLang string
var targetLang string
var sourceText string

func init() {
	flag.StringVar(&sourceLang, "source", "en", "Source language[en]")
	flag.StringVar(&targetLang, "target", "fr", "Target language[fr]")
	flag.StringVar(&sourceText, "text", "", "Text to translate")
}

func main() {
	flag.Parse()

	if flag.NFlag() == 0 {
		fmt.Println("Options:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	strChan := make(chan string)

	reqBody := &cli.RequestBody{
		SourceLang: sourceLang,
		SourceText: sourceText,
		TargetLang: targetLang,
	}

	wg.Add(1)
	go cli.RequestTranslate(reqBody, strChan, &wg)

	processedStr := strings.ReplaceAll(<-strChan, "+", " ")
	fmt.Printf("%s\n", processedStr)
	close(strChan)
	wg.Wait()
}
