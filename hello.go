package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	start := time.Now()

	src := flag.String("src", "", "Source file path")
	dst := flag.String("dst", "", "Destination file path")
	flag.Parse()

	if *src == "" || *dst == "" {
		log.Fatal("Source file path and destination directory path must be specified. Use --src and --dst flags to specify them.")
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
		fetchAndSave(url, *dst)

	}

	elapsed := time.Since(start)
	fmt.Printf("\nProgram execution time: %s\n", elapsed)
}

// fetchAndSave() функця считывает urlStr и пытается отправть get запрос по этому url,
// затем, если успешный запрос, сохраняет полученый body с запроса и сохраняет по
// пути dst
func fetchAndSave(url, dst string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error fetching URL: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("Unexpected status code for URL %s: %d\n", url, resp.StatusCode)
		return
	}

	fileName, err := getFileNameFromURL(url)
	if fileName == "" {
		return
	}
	if err != nil {
		log.Printf("err: %v\n", err)
	}
	fileName = fmt.Sprintf("%s.html", fileName)

	if dst == "./" {
		dst = "./list"
		err := os.MkdirAll(dst, 0755)
		if err != nil {
			log.Fatalf("Error creating folder %s: %v\n", dst, err)
		}
	}
	filePath := filepath.Join(dst, fileName)

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

// getFileNameFromURL() функция получает url сайт
// и возвращает имя файла на основе доменного имени url
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
