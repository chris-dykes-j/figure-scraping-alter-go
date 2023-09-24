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

type FigureData struct {
	Name      string
	TableData []string
	Material  string
	URL       string
	BlogLinks string
	Brand     string
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

func addCharacterToCsv(link string, root string, name string, brand string, fileName string) {
	data := structToArray(collectData(link, root, name, brand))
	if len(data) == 0 {
		log.Fatalf("Error: No data to add")
	}
	if len(data) != 12 {
		log.Fatalf("Error: column number does not match: %d", len(data))
	}
	addDataToCsv(fileName, data)
}

func collectData(link string, root string, name string, brand string) FigureData {
	userAgent := "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/116.0"
	c := colly.NewCollector()
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", userAgent)
	})

	var data FigureData

	data.Name = name
	data.URL = fmt.Sprintf("https://alter-web.jp%s", link)
	data.Brand = brand

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
		e.ForEach("td", func(_ int, el *colly.HTMLElement) {
			text, err := el.DOM.Html()
			if err != nil {
				fmt.Printf("%s\n", err)
			}
			text = strings.ReplaceAll(text, "<br/>", " ")
			data.TableData = append(data.TableData, strings.Join(strings.Fields(text), " "))
		})
	})
	defer c.OnHTMLDetach(".tbl-01 > tbody")

	// Get Material
	c.OnHTML(".spec > .txt", func(h *colly.HTMLElement) {
		data.Material = strings.Join(strings.Fields(h.Text), " ")
	})
	defer c.OnHTMLDetach(".spec > .txt")

	// Get Blog links
	var blogLinks []string
	c.OnHTML(".imgtxt-type-b", func(h *colly.HTMLElement) {
		h.ForEach("a", func(_ int, el *colly.HTMLElement) {
			blogLink := el.Attr("href")
			blogLink = fmt.Sprintf("https://alter-web.jp%s", blogLink)
			blogLinks = append(blogLinks, blogLink)
		})
	})
	defer c.OnHTMLDetach(".imgtxt-type-b")

	url := fmt.Sprintf("https://alter-web.jp%s", link)
	err := c.Visit(url)
	if err != nil {
		log.Fatalf("Error: %s, %s", err, link)
	}

	// Add blogLinks
	data.BlogLinks = strings.Join(blogLinks, ",")

	sleepShort()
	return data
}

func addDataToCsv(fileName string, data []string) {
	csvFile, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Can't read file: %s", err)
	}
	defer csvFile.Close()

	fmt.Printf("Adding %s to file...\n", data[0])

	for _, entry := range data {
		fmt.Println(entry)
	}
	csvWriter := *csv.NewWriter(csvFile)
	csvWriter.Write(data)
	csvWriter.Flush()
}

func structToArray(figureData FigureData) []string {
	// Remember: columnNames := []string{"name", "series", "character", "release", "price", "size", "sculptor", "painter", "material", "brand", "url", "blog_url"}
	return []string{
		figureData.Name,
		figureData.TableData[0],
		figureData.TableData[1],
		figureData.TableData[2],
		figureData.TableData[3],
		figureData.TableData[4],
		figureData.TableData[5],
		figureData.TableData[6],
        figureData.Material,
		figureData.Brand,
		figureData.URL,
		figureData.BlogLinks,
	}
}
