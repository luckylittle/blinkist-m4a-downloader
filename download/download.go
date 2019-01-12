// WORK IN PROGRESS

package main

import (
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

	// authenticate
	fmt.Println("Logging in...")
	// get the authenticity token
	c := colly.NewCollector()
	c.OnHTML("div.loginV2--login div[class=loginV2__form] input[type=hidden][name=authenticity_token]", func(e *colly.HTMLElement) {
		authenticityToken := e.Attr("value")
		fmt.Println("authenticity_token is", authenticityToken)
		err := c.Post("https://www.blinkist.com/en/nc/login", map[string]string{"login[email]": email, "login[password]": password, "authenticity_token": authenticityToken})
		if err != nil {
			log.Fatal(err)
		}
	})
	// attach callbacks after login
	c.OnResponse(func(r *colly.Response) {
		log.Print("response received HTTP", r.StatusCode)
	})
	c.Visit("https://www.blinkist.com/en/nc/login")

	// read line in the file
	//scanner := bufio.NewScanner(file)
	//for scanner.Scan() {
	// create a new collector
	//c := colly.NewCollector()

	// scrape it baby scrape it
	// c.OnHTML("body", func(e *colly.HTMLElement) {
	//e.ForEach("span", func(_ int, elem *colly.HTMLElement) {
	//if strings.Contains(elem.Text, "") {
	// 	dataBookID := e.Text
	// 	fmt.Println(dataBookID)
	// })

	// visit
	// c.Visit(scanner.Text())
	//}

	//if err := scanner.Err(); err != nil {
	//	log.Fatal(err)
	//}
}
