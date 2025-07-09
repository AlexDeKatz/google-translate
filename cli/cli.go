package cli

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/Jeffail/gabs/v2"
)

type RequestBody struct {
	SourceLang string
	SourceText string
	TargetLang string
}

const translateURL = "https://translate.googleapis.com/translate_a/single"

func RequestTranslate(reqBody *RequestBody, str chan string, wg *sync.WaitGroup) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", translateURL, nil)
	if err != nil {
		log.Fatalf("1. There was a problem: %s", err)
	}

	query := req.URL.Query()
	query.Add("client", "gtx")
	query.Add("sl", reqBody.SourceLang)
	query.Add("tl", reqBody.TargetLang)
	query.Add("dt", "t")
	query.Add("q", reqBody.SourceText)

	req.URL.RawQuery = query.Encode()

	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("2. There was a problem: %s", err)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusTooManyRequests {
		str <- "Error: Too many requests. Please try again later."
		wg.Done()
		return
	}

	parsedJSON, err := gabs.ParseJSONBuffer(res.Body)

	if err != nil {
		log.Fatalf("3. There was a problem: %s", err)
	}

	fmt.Println("Data parsed successfully: ", parsedJSON)

	nestOne, err := parsedJSON.ArrayElement(0)

	if err != nil {
		log.Fatalf("4. There was a problem: %s", err)
	}

	nestTwo, err := nestOne.ArrayElement(0)
	if err != nil {
		log.Fatalf("5. There was a problem: %s", err)
	}

	translatedString, err := nestTwo.ArrayElement(0)

	if err != nil {
		log.Fatalf("6. There was a problem: %s", err)
	}

	str <- translatedString.Data().(string)
	wg.Done()

}
