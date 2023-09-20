package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func createYearDir(year string) string {
	root := os.Getenv("FIGURE_IMAGE_ROOT")
	if root == "" {
		log.Fatal("FIGURE_IMAGE_ROOT environment variable not set")
	}

	parentDir := filepath.Join(root, year)
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
		dir := filepath.Join(root, strings.ReplaceAll(figure, "/", "-"))
		if _, err := os.Stat(dir); errors.Is(err, os.ErrNotExist) {
			err := os.Mkdir(dir, os.ModePerm)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Made %s directory\n", dir)
		}
	}
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

	fmt.Printf("Saved image %s to %s\n", imageName, directory)
}
