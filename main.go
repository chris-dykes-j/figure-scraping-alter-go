package main

import (
	"fmt"
	"github.com/gocolly/colly"
)

var fileName = "alter-jp.csv"
var userAgent = "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/116.0"
var brand string

func main() {
	c := colly.NewCollector()
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", userAgent)
	})

	columnNames := []string{"name", "series", "character", "release", "price", "size", "sculptor", "painter", "material", "blog_url", "brand", "url"}
	createCsvFile(columnNames)

	var years []string

	brand = "Alter"
	years = visitFirstPage("figure", c)
	for _, year := range years {
		fmt.Println("Scraping from year: ", year)
		visitPageByYear(year, "figure", c)
		sleepLong()
	}

	brand = "Altair"
	years = visitFirstPage("altair", c)
	for _, year := range years {
		fmt.Println("Scraping from year: ", year)
		visitPageByYear(year, "altair", c)
		sleepLong()
	}

	brand = "Alter"
	years = visitFirstPage("collabo", c)
	for _, year := range years {
		fmt.Println("Scraping from year: ", year)
		visitPageByYear(year, "collabo", c)
		sleepLong()
	}
}
