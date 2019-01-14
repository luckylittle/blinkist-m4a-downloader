package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"

	"github.com/gocolly/colly"
)

// variables
var (
	email    = "XXXX" // !!! change to your e-mail address
	password = "XXXX" // !!! enter your password here
)

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
			log.Println("Visited", r.Request.URL)
		})

		// attach error after login
		c.OnError(func(_ *colly.Response, err error) {
			log.Println("Something went wrong:", err)
		})

		// create another collector to scrape book details
		bookCollector := c.Clone()

		// defining variables outside the function and then assign the value in the callback later
		var dataTitle string
		var dataBookID string
		var s []string

		type URL struct {
			URL string
		}

		// scrape it baby scrape it!
		// read book title
		bookCollector.OnHTML("div[class=reader__container__share] a[class=share__facebook-icon]", func(f *colly.HTMLElement) {
			dataTitle = f.Attr("data-title")
			fmt.Println("Book title is:", dataTitle)
			return
		})

		// read book ID
		bookCollector.OnHTML("div[class=reader__container]", func(g *colly.HTMLElement) {
			dataBookID = g.Attr("data-book-id")
			//fmt.Println("Book ID is:", dataBookID)
			return
		})

		// read chapters and corresponding IDs
		bookCollector.OnHTML("div.chapter", func(h *colly.HTMLElement) {
			var link URL
			dataChapterNo := h.Attr("data-chapterno")
			dataChapterID := h.Attr("data-chapterid")
			apiLink := "https://www.blinkist.com/api/books/" + dataBookID + "/chapters/" + dataChapterID + "/audio"
			fmt.Println("API Link:" + apiLink + " for chapter " + dataTitle + "/" + dataChapterNo)
			bookCollector.OnResponse(func(r *colly.Response) {
				if r.StatusCode == 200 {
					json.Unmarshal([]byte(r.Body), &link)
				} else {
					log.Println("Doesn't contain audio!")
				}
			})
			bookCollector.Visit(apiLink)
			s = append(s, link.URL)
			return
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

		// download the files
		for index, element := range s {
			if _, err := os.Stat(dataTitle); os.IsNotExist(err) {
				os.Mkdir(dataTitle, 0700)
			} else {
				log.Println("Folder already exist, investigate!")
			}
			var url = element
			var path = dataTitle + "/" + "0" + strconv.Itoa(index) + ".m4a"
			wget(url, path)
		}

		// display collector's statistics
		log.Println(bookCollector)
	}

	// if reading the file fails
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

// call shell command `wget`` to download from URL. `wget` needs to be installed already and in $PATH
func wget(url, filepath string) error {
	// run shell `wget <URL> -O <FILEPATH>`
	cmd := exec.Command("wget", url, "-O", filepath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
