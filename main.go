package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/gocolly/colly"
)

func main() {
	c := colly.NewCollector()

	userAgent := "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/116.0"
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", userAgent)
	})

    fileName := "alter-jp.csv"
    columnNames := []string { "name","series","character","release","price","size","sculptor","painter","material","brand","url","blog_url" }
    createCsvFile(fileName, columnNames)

	years := visitFirstPage(c) // Get years and scrape the data
	for _, year := range years {
		fmt.Println("Scraping from year: ", year)
		visitPageByYear(year, c)
		sleepLong()
	}
}

func visitFirstPage(c *colly.Collector) []string {
	var years []string
	c.OnHTML("#changeY", func(e *colly.HTMLElement) {
		years = e.ChildAttrs("option", "value")
	})

	c.OnHTML(".type-a", getFigureLinks)

	sleepShort()
	url := fmt.Sprintf("https://alter-web.jp/figure")
	c.Visit(url)

	return years
}

func visitPageByYear(year string, c *colly.Collector) {
	c.OnHTML(".type-a", getFigureLinks)
	sleepShort()
	url := fmt.Sprintf("https://alter-web.jp/figure/?yy=%s&mm=", year)
	c.Visit(url)
}

func getFigureLinks(e *colly.HTMLElement) {
	links := e.ChildAttrs("a", "href")
	for _, link := range links {
        addCharacterToCsv(link)
	}
}

func createCsvFile(fileName string, fileHeader []string) {
    _, err := os.Stat(fileName)
    if os.IsNotExist(err) {
        csvFile, err2 := os.Create(fileName)
        if err2 != nil {
            log.Fatalf("csv file creation failed: %s", err)
        }
        csvWriter := csv.NewWriter(csvFile)
        csvWriter.Write(fileHeader)

        csvWriter.Flush()
        csvFile.Close()
    }
}

func addCharacterToCsv(link string) {
    
}

func sleepShort() {
	randomNumber := rand.Float64()*(4-2) + 2 // MATH
	time.Sleep(time.Duration(randomNumber) * time.Second)
}

func sleepLong() {
	randomNumber := rand.Float64()*(10-5) + 5
	time.Sleep(time.Duration(randomNumber) * time.Second)
}

