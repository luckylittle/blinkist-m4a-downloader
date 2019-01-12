// WORK IN PROGRESS

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/gocolly/colly"
)

// variables
var email = "XXXX"    // !!! change to your e-mail address
var password = "XXXX" // !!! enter your password here

func main() {
	// open the file
	file, err := os.Open("books_urls.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// read line in the file
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		// instantiate default collector
		c := colly.NewCollector()

		// get the authenticity token
		c.OnHTML("div.loginV2--login div[class=loginV2__form] input[type=hidden][name=authenticity_token]", func(e *colly.HTMLElement) {
			authenticityToken := e.Attr("value")
			fmt.Println("Logging in...")
			fmt.Println("authenticity_token for", scanner.Text(), "is", authenticityToken)

			// authenticate
			err := c.Post("https://www.blinkist.com/en/nc/login/", map[string]string{"utf8": "&#x2713;", "authenticity_token": authenticityToken, "login[google_id_token]": "", "login[facebook_access_token]": "", "login[email]": email, "login[password]": password})
			if err != nil {
				log.Fatal(err)
			}
		})

		// visit login page
		c.Visit("https://www.blinkist.com/en/nc/login/")

		// attach callbacks after login
		c.OnResponse(func(r *colly.Response) {
			log.Print("Login response received HTTP", r.StatusCode)
			log.Println("Visited", r.Request.URL)
		})

		// attach error after login
		c.OnError(func(_ *colly.Response, err error) {
			log.Println("Something went wrong:", err)
		})

		// create another collector to scrape book details
		bookCollector := c.Clone()

		// scrape it baby scrape it!
		// read book title
		bookCollector.OnHTML("div[class=reader__container__share] a[class=share__facebook-icon]", func(f *colly.HTMLElement) {
			dataTitle := f.Attr("data-title")
			fmt.Println("Book title is:", dataTitle)
		})

		// read book ID
		bookCollector.OnHTML("div[class=reader__container]", func(g *colly.HTMLElement) {
			dataBookID := g.Attr("data-book-id")
			fmt.Println("Book ID is:", dataBookID)
		})

		// read chapters and corresponding IDs
		bookCollector.OnHTML("div.chapter", func(h *colly.HTMLElement) {
			dataChapterNo := h.Attr("data-chapterno")
			dataChapterID := h.Attr("data-chapterid")
			fmt.Println("Chapter", dataChapterNo, "has ID:", dataChapterID)
		})

		// attach callbacks after data title
		bookCollector.OnResponse(func(f *colly.Response) {
			log.Print("Book page response received HTTP", f.StatusCode)
		})

		// visit
		fmt.Println("Visiting", scanner.Text())
		bookCollector.Visit(scanner.Text())

		bookCollector.OnScraped(func(r *colly.Response) {
			fmt.Println("Finished", r.Request.URL)
		})
	}

	// if reading the file fails
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
