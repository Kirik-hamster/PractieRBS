package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"bufio"
	"strings"
	"time"
)

func main() {
	start := time.Now()

	src := flag.String("src", "", "Source file path")
	dst := flag.String("dst", "", "Destination file path")
	flag.Parse()

	if *src == "" || *dst == "" {
		log.Fatal("Source file path is not specified. Use --src flag to specify the file path.")
	}

	srcFile, err := os.Open(*src)
	if err != nil {
		log.Fatalf("Error opening file: %v\n", err)
	}
	defer srcFile.Close()

	scanner := bufio.NewScanner(srcFile)

	for scanner.Scan() {
		url := strings.TrimSpace(scanner.Text())
		if url == "" {
			continue
		}

		resp, err := http.Get(url)
		if err != nil {
			log.Printf("Error fetching URL: %v\n", err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			log.Printf("Unexpected status code for URL %s: %d\n", url, resp.StatusCode)
			continue
		}

		fileName, err := getFileNameFromURL(url) 
		if fileName == "" {
			continue
		}
		if err != nil {
			log.Printf("err: %v\n", err)
		}
		fileName += ".html"
		filePath := filepath.Join(*dst, fileName)

		dstFile, err := os.Create(filePath)
		if err != nil {
			log.Fatalf("Error creating destination file: %v\n", err)
		}
		_, err = io.Copy(dstFile, resp.Body)
		defer dstFile.Close()
		
		if err != nil {
			log.Fatalf("Error copying content: %v\n", err)
		}
		fmt.Printf("File copied successfully to \n%s\n", filePath)
	}


	elapsed := time.Since(start)
	fmt.Printf("\nProgram execution time: %s\n", elapsed)
}

func getFileNameFromURL(siteURL string) (string, error) {
	parsedURL, err := url.Parse(siteURL)
	if err != nil {
		return "", err
	}

	domain := strings.TrimPrefix(parsedURL.Host, "www.")

	if strings.Contains(domain, ":") {
		domain = strings.Split(domain, ":")[0]
	}

	parts := strings.SplitN(domain, ".", 2)
	if len(parts) < 2 {
		return domain, nil
	}


	return parts[0], nil
}