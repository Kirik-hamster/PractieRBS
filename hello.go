package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
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
		log.Fatal("Source file path is not specified. Use --src flag to specify the file path.")
	}

	srcFile, err := os.Open(*src)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer srcFile.Close()

	content, err := io.ReadAll(srcFile)
	if err != nil {
		log.Fatalf("Error reading file %v", err)
	}

	url := strings.TrimSpace(string(content))

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error fetching URL: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatalf("Unexpected status code: %d", resp.StatusCode)
	}

	fileName := "test.html"
	filePath := filepath.Join(*dst, fileName)

	dstFile, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("Error creating destination file: %v", err)
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, resp.Body)

	fmt.Printf("Content of reading file:\n\n%s\n", content)
	fmt.Printf("file copied successfully to \n%s\n", filePath)

	elapsed := time.Since(start)
	fmt.Printf("\nProgram execution time: %s\n", elapsed)
}
