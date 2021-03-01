package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func main() { // Request the HTML page.
	res, err := http.Get("http://site.sanepar.com.br")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find(".views-field-field-nivel-reserv-icone-fid").Each(func(i int, s *goquery.Selection) {
		parent := s.Parent()
		title := parent.Find(".views-field-title").Find("span").Text()
		level := parent.Find(".views-field-body").Find("p").Text()
		fmt.Println(title, level)
	})
}
