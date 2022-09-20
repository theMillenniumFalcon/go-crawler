package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type SEOData struct {
	URL             string
	Title           string
	H1              string
	MetaDescription string
	StatusCode      int
}

type DefaultParser struct {
}

type Parser interface {
	getSEOData(resp *http.Response) (SEOData, error)
}

var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Safari/604.1.38",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:56.0) Gecko/20100101 Firefox/56.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Safari/604.1.38",
}

func randomUserAgent() string {
	rand.Seed(time.Now().Unix())
	randNum := rand.Int() % len(userAgents)
	return userAgents[randNum]
}

func isSitemap(urls []string) ([]string, []string) {
	siteMapFiles := []string{}
	pages := []string{}

	for _, page := range urls {
		foundSitemap := strings.Contains(page, "xml")
		if foundSitemap == true {
			fmt.Println("Found sitemap", page)
			siteMapFiles = append(siteMapFiles, page)
		} else {
			pages = append(pages, page)
		}
	}

	return siteMapFiles, pages
}

func extractSiteMapURLsFunc(startURL string) []string {
	workList := make(chan []string)
	toCrawl := []string{}
	var n int
	n++

	go func() { workList <- []string{startURL} }()

	for ; n > 0; n-- {
		list := <-workList

		for _, link := range list {
			n++
			go func(link string) {
				response, err := makeRequest(link)
				if err != nil {
					log.Printf("Error retrieving URL: %s", link)
				}
				urls, _ := extractURLs(response)
				if err != nil {
					log.Printf("Error extracting document from response, URL: %s", link)
				}
				sitemapFiles, pages := isSitemap(urls)
				if sitemapFiles != nil {
					workList <- sitemapFiles
				}
				for _, page := range pages {
					toCrawl = append(toCrawl, page)
				}
			}(link)
		}
	}

	return toCrawl
}

func makeRequest(url string) (*http.Response, error) {
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", randomUserAgent())
	if err != nil {
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func scrapeURLs(urls []string, parser Parser, concurrency int) []SEOData {

}

func extractURLs(response *http.Response) ([]string, error) {
	doc, err := goquery.NewDocumentFromResponse(response)
	if err != nil {
		return nil, err
	}
	results := []string{}
	sel := doc.Find("loc")

	for i := range sel.Nodes {
		loc := sel.Eq(i)
		result := loc.Text()
		results = append(results, result)
	}

	return results, nil
}

func scrapePage(url string, parser Parser) (SEOData, error) {
	res, err := crawlPage(url)
	if err != nil {
		return SEOData{}, err
	}
	data, err := parser.getSEOData(res)
	if err != nil {
		return SEOData{}, err
	}

	return data, nil
}

func crawlPage(url string) {

}

func (d DefaultParser) getSEOData(resp *http.Response) (SEOData, error) {
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return SEOData{}, err
	}
	result := SEOData{}
	result.URL = resp.Request.URL.String()
	result.StatusCode = resp.StatusCode
	result.Title = doc.Find("title").First().Text()
	result.H1 = doc.Find("h1").First().Text()
	result.MetaDescription, _ = doc.Find("meta[name^=description]").Attr("content")

	return result, nil
}

func ScrapeSiteMapfunc(url string, parser Parser, concurrency int) []SEOData {
	results := extractSiteMapURLsFunc(url)
	res := scrapeURLs(results, parser, concurrency)
	return res
}

func main() {
	p := DefaultParser{}
	results := ScrapeSiteMapfunc("https://www.quicksprout.com/sitemap.xml", p, 10)
	for _, res := range results {
		fmt.Println(res)
	}
}
