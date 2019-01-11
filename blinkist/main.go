package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/gocolly/colly"
)

// variables
var baseURL = "https://www.blinkist.com"
var filename = "books_urls.txt"
var username = "xxxxx@xxxxxx.xx" // !!! change to your e-mail address
var password = "xxxxxxx"         // !!! enter your password here

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	// links to 27 categories containing all of the books
	categories := [27]string{
		"https://www.blinkist.com/en/nc/categories/entrepreneurship-and-small-business-en/books",
		"https://www.blinkist.com/en/nc/categories/science-en/books",
		"https://www.blinkist.com/en/nc/categories/economics-en/books",
		"https://www.blinkist.com/en/nc/categories/corporate-culture-en/books",
		"https://www.blinkist.com/en/nc/categories/money-and-investments-en/books",
		"https://www.blinkist.com/en/nc/categories/relationships-and-parenting-en/books",
		"https://www.blinkist.com/en/nc/categories/parenting-en/books",
		"https://www.blinkist.com/en/nc/categories/education-en/books",
		"https://www.blinkist.com/en/nc/categories/society-and-culture-en/books",
		"https://www.blinkist.com/en/nc/categories/politics-and-society-en/books",
		"https://www.blinkist.com/en/nc/categories/health-and-fitness-en/books",
		"https://www.blinkist.com/en/nc/categories/biography-and-history-en/books",
		"https://www.blinkist.com/en/nc/categories/management-and-leadership-en/books",
		"https://www.blinkist.com/en/nc/categories/psychology-en/books",
		"https://www.blinkist.com/en/nc/categories/technology-and-the-future-en/books",
		"https://www.blinkist.com/en/nc/categories/nature-and-environment-en/books",
		"https://www.blinkist.com/en/nc/categories/philosophy-en/books",
		"https://www.blinkist.com/en/nc/categories/career-and-success-en/books",
		"https://www.blinkist.com/en/nc/categories/marketing-and-sales-en/books",
		"https://www.blinkist.com/en/nc/categories/personal-growth-and-self-improvement-en/books",
		"https://www.blinkist.com/en/nc/categories/communication-and-social-skills-en/books",
		"https://www.blinkist.com/en/nc/categories/motivation-and-inspiration-en/books",
		"https://www.blinkist.com/en/nc/categories/productivity-and-time-management-en/books",
		"https://www.blinkist.com/en/nc/categories/mindfulness-and-happiness-en/books",
		"https://www.blinkist.com/en/nc/categories/religion-and-spirituality-en/books",
		"https://www.blinkist.com/en/nc/categories/biography-and-memoir-en/books",
		"https://www.blinkist.com/en/nc/categories/creativity-en/books",
	}
	// create a new collector
	c := colly.NewCollector(
		colly.MaxDepth(1),
		colly.AllowedDomains("www.blinkist.com"),
	)

	// authenticate
	err := c.Post("https://www.blinkist.com/en/nc/login/", map[string]string{"username": username, "password": password})
	if err != nil {
		log.Fatal(err)
	}

	// attach callbacks after login
	c.OnResponse(func(r *colly.Response) {
		log.Print("response received HTTP", r.StatusCode)
	})

	// on every a element which has href attribute call callback
	os.Create("temp")
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// if link starts with /en/books return from callback
		if !strings.HasPrefix(link, baseURL+"/en/books") {
			return
		}
		// print links to file
		f, err := os.OpenFile("temp", os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			panic(err)
		}

		defer f.Close()

		if _, err = f.WriteString(strings.Replace(link, "/en/books", "/en/nc/reader", -1) + "\n"); err != nil {
			panic(err)
		}
	})

	// start scraping
	for _, url := range categories {
		c.Visit(url)
	}

	// start removing duplicate entries
	os.Create(filename)
	line, _ := ioutil.ReadFile("temp")
	strLine := string(line)
	lines := strings.Split(strLine, "\n")
	resultSlice := removeDuplicates(lines)
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		log.Println(err)
	}
	for i := range resultSlice {
		f.Write([]byte(resultSlice[i] + "\n"))
	}
	f.Close()
	os.Remove("temp")
}

func removeDuplicates(s []string) []string {
	counter := 0
	m := make(map[string]bool)
	for _, item := range s {
		if _, ok := m[item]; ok {
			// duplicate item
			counter++
			fmt.Println(item, "is a duplicate", "#", counter)
		} else {
			m[item] = true
		}
	}

	var result []string
	for item := range m {
		result = append(result, item)
	}
	return result
}
