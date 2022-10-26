package main

import (
	"log"
	"regexp"
	"strings"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
)

func parseXmlImgSrc(input string) (imgSrc string) {
	srcRe, err := regexp.Compile(`src="(.*)"`)
	if err != nil {
		log.Println("Error with regex")
		return imgSrc
	}

	return srcRe.FindStringSubmatch(input)[1]
}

func scrapeTheVergeXML(c chan<- TvergeArticle, URL, path string) {
	scraper := colly.NewCollector()
	extensions.RandomUserAgent(scraper)

	scraper.OnRequest(func(r *colly.Request) {
		log.Println("Visiting: ", r.URL.String())
		log.Println("User-Agent: ", r.Headers.Values("User-Agent"))
	})

	scraper.OnError(func(r *colly.Response, err error) {
		log.Println("Error while scraping: ", err)
		close(c)
	})

	scraper.OnXML("//entry", func(x *colly.XMLElement) {
		tvergeArticle := TvergeArticle{}
		//Category
		tvergeArticle.Category = path
		tvergeArticle.CategoryLink = sourceURL + "/" + strings.ToLower(path)
		//Title
		tvergeArticle.Title = x.ChildText("title")
		//Date
		tvergeArticle.ArticleDate = x.ChildText("published")
		//Author
		tvergeArticle.Author = x.ChildText("author/name")
		//Article URL
		tvergeArticle.URL = x.ChildAttr("link", "href")
		//Image
		tvergeArticle.Img = parseXmlImgSrc(x.ChildText("content"))
		tvergeArticle.ImgSrcSet = ""

		c <- tvergeArticle
	})

	scraper.OnScraped(func(r *colly.Response) {
		log.Println("Finished with scraping XML: ", r.Request.URL)
		close(c)
	})

	scraper.Visit(URL)
}
