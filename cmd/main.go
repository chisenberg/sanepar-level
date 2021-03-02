package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"sanepar-level/domain/entity"
	"sanepar-level/infra/dynamo"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/aws/aws-lambda-go/lambda"
)

// Handler function Using AWS Lambda Proxy Request
func Handler(ctx context.Context) (string, error) {

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

	report := entity.Report{}

	// Find all dams
	doc.Find(".views-field-field-nivel-reserv-icone-fid").Each(func(i int, s *goquery.Selection) {
		if s == nil {
			return
		}
		parent := s.Parent()
		title := parent.Find(".views-field-title").Find("span").Text()
		level := parent.Find(".views-field-body").Find("p").Text()
		// Remove % and replace , with .
		level = strings.ReplaceAll(level, "%", "")
		level = strings.ReplaceAll(level, ",", ".")
		levelFloat, _ := strconv.ParseFloat(level, 32)

		report.Dams = append(report.Dams, entity.Dam{
			Name:  title,
			Level: levelFloat,
		})
	})

	// find updated at "Atualizado em: 01/03/2021 - 08:03"
	updatedAt := doc.Find(".nivel-reserv-data").Text()
	// rmeove spaces and clean
	updatedAt = strings.ReplaceAll(updatedAt, "Atualizado em: ", "")
	updatedAt = strings.TrimSpace(updatedAt)

	// Parse as time
	location, _ := time.LoadLocation("America/Sao_Paulo")
	report.UpdatedAt, _ = time.ParseInLocation("02/01/2006 - 15:04", updatedAt, location)

	// save
	err = dynamo.SaveReport(report)
	if err != nil {
		return "", err
	}

	return "Ok", nil
}

func main() {
	log.SetOutput(os.Stdout)
	lambda.Start(Handler)
}
