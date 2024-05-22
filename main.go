package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
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

	customHeaders := map[string]string{
		"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		"Accept-Language": "en-US,en;q=0.9",
		"Referer":         "http://google.com",
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
			for key, value := range customHeaders {
				r.Headers.Set(key, value)
			}

			delay := time.Duration(rand.Intn(3000-1000)+1000) * time.Millisecond
			time.Sleep(delay)
			fmt.Println("Added a small " + fmt.Sprint(delay) + "ms delay")
			fmt.Println("Visiting: ", r.URL)
		})

		c.OnResponse(func(r *colly.Response) {
			fmt.Println(r.Request.URL)
			fmt.Println(r.StatusCode)
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
