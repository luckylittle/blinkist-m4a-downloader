// WORK IN PROGRESS

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/gocolly/colly"
)

func main() {
	file, err := os.Open("books_urls.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		c := colly.NewCollector()
		c.OnHTML("body table tbody tr td", func(e *colly.HTMLElement) {
			dataBookID := e.Attr("data-book-id")

			// Print
			fmt.Println(dataBookID)
		})
		c.Visit(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
