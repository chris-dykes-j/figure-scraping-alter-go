package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
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

    columnNames := []string { "name","series","character","release","price","size","sculptor","painter","material", "blog_url", "brand", "url"}
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

func visitFirstPage(page string, c *colly.Collector) []string {
	var years []string
	c.OnHTML("#changeY", func(e *colly.HTMLElement) {
		years = e.ChildAttrs("option", "value")
	})
	url := fmt.Sprintf("https://alter-web.jp/%s", page)
	c.Visit(url)
    sleepShort()
    c.OnHTMLDetach("#changeY")
	return years
}

func visitPageByYear(year string, page string, c *colly.Collector) {
    root := createYearDir(year)

    // Seems to behave linearly, so it should be fine.
    c.OnHTML(".type-a", func(e *colly.HTMLElement) {
        var figureNames []string
        e.ForEach("figcaption", func(i int, h *colly.HTMLElement) {
            figureNames = append(figureNames, h.Text)
        })
        createFigureDirs(root, figureNames)

        images := e.ChildAttrs("img", "src")
        for i, image := range images {
            figureDir := filepath.Join(root, figureNames[i])
            downloadImg(image, figureDir, "profile.jpg")
        }

        links := e.ChildAttrs("a", "href")
        for i, link := range links {
            sleepShort()
            addCharacterToCsv(link, root, figureNames[i])
        }
    })

    sleepShort()
	url := fmt.Sprintf("https://alter-web.jp/%s/?yy=%s&mm=", page, year)
	c.Visit(url)
    c.OnHTMLDetach(".type-a")
}

func downloadImg(imageUrl string, directory string, imageName string) {
    // Ensure the directory exists
    if err := os.MkdirAll(directory, 0755); err != nil {
        log.Fatal(err)
    }

    // Get image bytes
    url := fmt.Sprintf("https://alter-web.jp%s", imageUrl)
    res, err := http.Get(url)
    if err != nil {
        log.Fatal(err)
    }
    defer res.Body.Close()

    // Make image file
    imagePath := filepath.Join(directory, imageName)
	imageFile, err := os.Create(imagePath)
	if err != nil {
        log.Fatal(err)
	}
	defer imageFile.Close()

    // Save bytes to image file
    _, err = io.Copy(imageFile, res.Body)
	if err != nil {
        log.Fatal(err)
	}

    fmt.Printf("Saved image %s to %s", imageName, directory)
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
    if (len(blogLinks) == 0) {
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

func sleepShort() {
	randomNumber := rand.Float64()*(4-2) + 2 // MATH
	time.Sleep(time.Duration(randomNumber) * time.Second)
}

func sleepLong() {
	randomNumber := rand.Float64()*(10-5) + 5
	time.Sleep(time.Duration(randomNumber) * time.Second)
}

func createYearDir(year string) string {
    root := os.Getenv("FIGURE_IMAGE_ROOT")
    if root == "" {
        log.Fatal("FIGURE_IMAGE_ROOT environment variable not set")
    }

    parentDir := root + "/" + year
    if _, err := os.Stat(parentDir); errors.Is(err, os.ErrNotExist) {
        err := os.Mkdir(parentDir, os.ModePerm)
        if err != nil {
            log.Fatal(err)
        }
        fmt.Printf("Made %s directory\n", parentDir)
    }

    return parentDir
}

func createFigureDirs(root string, figures []string) {
    for _, figure := range figures {
        dir := root + "/" + figure
        if _, err := os.Stat(dir); errors.Is(err, os.ErrNotExist) {
            err := os.Mkdir(dir, os.ModePerm)
            if err != nil {
                log.Fatal(err)
            }
            fmt.Printf("Made %s directory\n", dir)
        }
    }
}
