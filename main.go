package main

import (
	"fmt"
	"github.com/gocolly/colly"
)

func main() {
	userAgent := "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/116.0"

	c := colly.NewCollector()
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", userAgent)
	})

	columnNames := []string{"name", "series", "character", "release", "price", "size", "sculptor", "painter", "material", "brand", "url", "blog_url"}
	fileName := "alter-jp.csv"
	createCsvFile(fileName, columnNames)

	var years []string

	years = visitFirstPage("figure", c)
	for _, year := range years {
		fmt.Println("Scraping from year: ", year)
		visitPageByYear(year, "figure", "Alter", c)
		sleepLong()
	}

	years = visitFirstPage("altair", c)
	for _, year := range years {
		fmt.Println("Scraping from year: ", year)
		visitPageByYear(year, "altair", "Altair", c)
		sleepLong()
	}

	years = visitFirstPage("collabo", c)
	for _, year := range years {
		fmt.Println("Scraping from year: ", year)
		visitPageByYear(year, "collabo", "Alter", c)
		sleepLong()
	}
}
