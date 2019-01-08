package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/gocolly/colly"
)

var baseURL = "https://www.blinkist.com"

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
	err := c.Post("https://www.blinkist.com/en/nc/login/", map[string]string{"username": "XXX", "password": "XXX"})
	if err != nil {
		log.Fatal(err)
	}

	// attach callbacks after login
	c.OnResponse(func(r *colly.Response) {
		log.Print("response received HTTP", r.StatusCode)
	})

	// on every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// if link starts with /en/books return from callback
		if !strings.HasPrefix(link, baseURL+"/en/books") {
			return
		}
		// print link
		fmt.Println(link)
		// visit link found on page
		e.Request.Visit(link)
	})

	// start scraping
	for _, url := range categories {
		c.Visit(url)
	}
}
