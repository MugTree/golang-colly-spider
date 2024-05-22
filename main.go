package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

func main() {

	domainsFlag := flag.String("domains", "", "Comma seperated list of domains to crawl")
	keywordsFlag := flag.String("keywords", "", "Comma seperated list of keywords to check")

	flag.Parse()

	domains := strings.Split(*domainsFlag, ",")
	keywords := strings.Split(*keywordsFlag, ",")

	if len(domains) == 0 || len(keywords) == 0 {
		log.Fatal("You must provide at least one domain and one keyword")
	}

	for _, startURL := range domains {

		foundPages := []string{}

		startURLParsed, err := url.Parse(startURL)

		if err != nil {
			log.Fatalf("Failed to parse start URL: %v", err)
		}

		startHostname := startURLParsed.Hostname()

		c := colly.NewCollector()

		c.OnHTML("a[href]", func(e *colly.HTMLElement) {

			link := e.Attr("href")
			absLink := e.Request.AbsoluteURL(link)
			linkParsed, err := url.Parse(absLink)

			if err != nil {
				log.Println("Failed to parse link:", absLink)
				return
			}

			linkHostname := linkParsed.Hostname()
			if linkHostname != startHostname {
				return
			}

			e.Request.Visit(absLink)

		})

		c.OnHTML("html", func(e *colly.HTMLElement) {

			pageURL := e.Request.URL.String()
			pageText := e.Text

			for _, keyword := range keywords {
				if strings.Contains(strings.ToLower(pageText), strings.ToLower(keyword)) {
					foundPages = append(foundPages, pageURL)
					break
				}
			}
		})

		c.OnRequest(func(r *colly.Request) {
			fmt.Println("Visiting: ", r.URL)
		})

		c.OnScraped(func(r *colly.Response) {
			fmt.Println("Finished scraping", startHostname)
		})

		c.Visit(startURL)

		c.Wait()

		dirName := "public"

		newpath := filepath.Join(".", dirName)
		err = os.MkdirAll(newpath, os.ModePerm)

		if err != nil {
			log.Fatal(err)
		}

		file, err := os.Create("./" + dirName + "/" + startHostname + "_" + fmt.Sprint(time.Now().Unix()) + ".txt")
		if err != nil {
			log.Fatal("Could not create file:", err)
		}
		defer file.Close()

		for _, page := range foundPages {
			fmt.Println(page)
			file.WriteString(page + "\n")
		}
	}

}
