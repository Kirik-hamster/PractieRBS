package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

func main() {
	start := time.Now()

	src := flag.String("src", "", "Source file path")
	flag.Parse()

	if *src == "" {
		log.Fatal("Source file path is not specified. Use --src flag to specify the file path.")
	}

	file, err := os.Open(*src)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}

	content, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("Error reading file %v", err)
	}

	fmt.Printf("Error reaing file:\n\n%s\n", content)

	elapsed := time.Since(start)
	fmt.Printf("\nProgram execution time: %s\n", elapsed)
}
