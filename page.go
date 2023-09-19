package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"path/filepath"
)

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
