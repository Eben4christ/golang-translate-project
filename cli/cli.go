package cli

import (
	"log"
	"net/http"
	"sync"

	"github.com/Jeffail/gabs"
)

type RequestBody struct {
	SourceLang string
	TargetLang string
	SourceText string
}

const translateUrl = "https://translate.googleapis.com/translate_a/single"

func RequestTranslate(body *RequestBody, str chan string, wg *sync.WaitGroup)  {
	
	client := &http.Client{}
	req, err := http.NewRequest("GET", translateUrl, nil)

	query := req.URL.Query()
	
	query.Add("client", "gtx")

	query.Add("sl", body.SourceLang)
	query.Add("tl", body.TargetLang)
	query.Add("dt", "t")
	query.Add("q", body.SourceText)

	req.URL.RawQuery = query.Encode()

	if err != nil {
		log.Fatal("1. There was an error:%s", err)
	}
	
	res, err := client.Do(req)
	if err != nil {
		log.Fatal("2. There was an error:%s", err)
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusTooManyRequests{
		str <- "You have been rate limited, try again later."
		wg.Done()
		return
	}

	parsedJSON, err := gabs.ParseJSONBuffer(res.Body)

	if err != nil {
		log.Fatalf("3. There was an error - %s", err)
	}

	nestOne, err := parsedJSON.ArrayElement(0)

	if err != nil {
		log.Fatalf("4. There was an error - %s", err)
	}

	nestTwo, err := nestOne.ArrayElement(0)
	if err != nil {
		log.Fatalf("5. There was an error - %s", err)
	}

	translatedStr, err := nestTwo.ArrayElement(0)
	if err != nil {
		log.Fatalf("6. There was an error - %s", err)
	}

	str <- translatedStr.Data().(string)
	wg.Done()
}









