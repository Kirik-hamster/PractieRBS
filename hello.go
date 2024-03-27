package main

import (
	"flag"
	"fmt"
	"time"
)

func main() {
	start := time.Now()
	nameTo := flag.String("name", "world", "name to sell hi")
	flag.Parse()

	fmt.Printf("Hello, %s!\n", *nameTo)

	elapsed := time.Since(start)
	fmt.Printf("Program execution time: %s\n", elapsed)
}
