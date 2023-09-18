package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

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

    columnNames := []string { "name","series","character","release","price","size","sculptor","painter","material","url","blog_url", "brand"}
    createCsvFile(columnNames)

    var years []string

    brand = "Alter"
    years = visitFirstPage("collabo", c)
	for _, year := range years {
		fmt.Println("Scraping from year: ", year)
		visitPageByYear(year, "collabo", c)
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
	years = visitFirstPage("figure", c) 
	for _, year := range years {
		fmt.Println("Scraping from year: ", year)
		visitPageByYear(year, "figure", c)
		sleepLong()
	}
}

func visitFirstPage(page string, c *colly.Collector) []string {
	var years []string
	c.OnHTML("#changeY", func(e *colly.HTMLElement) {
		years = e.ChildAttrs("option", "value")
	})

    c.OnHTML(".type-a", func(e *colly.HTMLElement) {
        links := e.ChildAttrs("a", "href")
        for _, link := range links {
            sleepShort()
            addCharacterToCsv(link)
        }
    })

    sleepShort()
	url := fmt.Sprintf("https://alter-web.jp/%s", page)
	c.Visit(url)

	return years
}

func visitPageByYear(year string, page string, c *colly.Collector) {
    c.OnHTML(".type-a", func(e *colly.HTMLElement) {
        links := e.ChildAttrs("a", "href")
        for _, link := range links {
            sleepShort()
            addCharacterToCsv(link)
        }
    })

    sleepShort()
	url := fmt.Sprintf("https://alter-web.jp/%s/?yy=%s&mm=", page, year)
	c.Visit(url)
}

func createCsvFile(fileHeader []string) {
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
	c := colly.NewCollector()
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", userAgent)
	})

    csvFile, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Fatalf("Can't read file: %s", err)
    }
    defer csvFile.Close()

    data := []string{}

    // Get Figure name
    var name string
    c.OnHTML(".hl06", func(e *colly.HTMLElement) {
        name = e.Text
        data = append(data, name) 
    })

    // Get Figure Table
    c.OnHTML(".tbl-01 > tbody", func(e *colly.HTMLElement) {
        e.ForEach("tr", func(_ int, el *colly.HTMLElement) {
            text := el.ChildText("td")
            data = append(data, strings.Join(strings.Fields(text), " "))
        })
    })

    // Get Material
    c.OnHTML(".spec > .txt", func(h *colly.HTMLElement) {
        data = append(data, strings.Join(strings.Fields(h.Text), " "))
    })

    // Add Url
    url := fmt.Sprintf("https://alter-web.jp%s", link)
    data = append(data, url)

    // Add Blog links
    var blogLinks []string
    c.OnHTML(".imgtxt-type-b", func(h *colly.HTMLElement) {
        blogLinks = h.ChildAttrs("a", "href")
        for i, blogLink := range blogLinks {
            println(blogLink)
            blogLinks[i] = fmt.Sprintf("https://alter-web.jp%s", blogLink)
        }
        data = append(data, strings.Join(blogLinks, ","))
    })

	err = c.Visit(url)
    if err != nil {
        log.Fatalf("Error: %s, %s", err, link)
    }

    // Handle empty blogLinks
    if (len(blogLinks) == 0) {
        data = append(data, "null")
    }

    // Add Brand
    data = append(data, brand) 

    fmt.Printf("Adding %s to file...\n", name)
    for _, entry := range data {
        fmt.Println(entry)
    }
    csvWriter := *csv.NewWriter(csvFile)
    csvWriter.Write(data)
    csvWriter.Flush()
    sleepShort()
}

func sleepShort() {
	randomNumber := rand.Float64()*(4-2) + 2 // MATH
	time.Sleep(time.Duration(randomNumber) * time.Second)
}

func sleepLong() {
	randomNumber := rand.Float64()*(10-5) + 5
	time.Sleep(time.Duration(randomNumber) * time.Second)
}

