package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gocolly/colly"
)

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

func addCharacterToCsv(link string, root string, name string) {
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

	// Add figure images
	figureDir := filepath.Join(root, name)
	c.OnHTML(".item-mainimg figure img", func(h *colly.HTMLElement) {
		downloadImg(h.Attr("src"), figureDir, "1.jpg")
	})
	i := 1
	c.OnHTML(".imgset li img", func(h *colly.HTMLElement) {
		i++
		downloadImg(h.Attr("src"), figureDir, fmt.Sprintf("%d.jpg", i))
	})
	defer c.OnHTMLDetach(".item-mainimg figure img")
	defer c.OnHTMLDetach(".imgset li img")

	// Get Figure Table
	c.OnHTML(".tbl-01 > tbody", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(_ int, el *colly.HTMLElement) {
			text := el.ChildText("td")
			if text == "" {
				data = append(data, "null")
			} else {
				data = append(data, strings.Join(strings.Fields(text), " "))
			}
		})
	})

	// Get Material
	c.OnHTML(".spec > .txt", func(h *colly.HTMLElement) {
		data = append(data, strings.Join(strings.Fields(h.Text), " "))
	})

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

	url := fmt.Sprintf("https://alter-web.jp%s", link)
	err = c.Visit(url)
	if err != nil {
		log.Fatalf("Error: %s, %s", err, link)
	}

	// Handle empty blogLinks
	if len(blogLinks) == 0 {
		data = append(data, "null")
	}

	// Add Brand
	data = append(data, brand)

	// Add Url
	data = append(data, url)

	fmt.Printf("Adding %s to file...\n", name)
	for _, entry := range data {
		fmt.Println(entry)
	}
	csvWriter := *csv.NewWriter(csvFile)
	csvWriter.Write(data)
	csvWriter.Flush()

	c.OnHTMLDetach(".hl06")
	c.OnHTMLDetach(".tbl-01 > tbody")
	c.OnHTMLDetach(".spec > .txt")
	c.OnHTMLDetach(".imgtxt-type-b")
	c.OnHTMLDetach(".item-mainimg > img")
	c.OnHTMLDetach(".imgset > li")
	sleepShort()
}
