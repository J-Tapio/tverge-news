package main

import (
	"github.com/gocolly/colly"
	"log"
	"regexp"
	"strings"
)

//TODO: Quite possibly just a flaw in my selector - category somehow comes as a duplicate, eg. "WebWeb"??

func category(c string, path string) string {
	if c == "" {
		return path
	}
	firstChar := rune(c[0])
	for i, character := range c {
		if firstChar == character {
			c = c[i:]
		}
	}
	return c
}

func categoryLink(link string, path string) string {
	if link == "" {
		return "https://www.theverge.com/" + strings.ToLower(path)
	}
	return "https://www.theverge.com" + link
}

// Regex and formatting for image src and srcSet values
func scrapeImageSrc(source string) (img, srcSet string) {
	
	// srcSet
	re, err := regexp.Compile(`srcSet="(.*?)"\ssrc`)
	if err != nil {
		log.Println("Error with regex: ", err)
		return "", ""
	}
	re2, err := regexp.Compile(`/_next`)
	if err != nil {
		log.Println("Error with regex: ", err)
		return "", ""
	}

	imgSrcSetMatch := re.FindStringSubmatch(source)
	// Replace srcSet strings with URL
	// Finally clean 'amp;' - unicode from the strings
	if len(imgSrcSetMatch) > 0 {
		srcSetUncleaned := re2.ReplaceAllString(imgSrcSetMatch[1], "https://www.theverge.com/_next")
		srcSet = strings.Replace(srcSetUncleaned, "amp;", "", -1)
	} else {
		srcSet = ""
	}

	// src
	re3, err := regexp.Compile(`src="(.*?)"\sdecoding`)
	if err != nil {
		log.Println("Error with regex: ", err)
		return "", ""
	}

	imgSrcMatch := re3.FindStringSubmatch(source)
	if len(imgSrcMatch) > 0 {
		imgUncleaned := "https://www.theverge.com" + imgSrcMatch[1]
		img = strings.Replace(imgUncleaned, "amp;", "", -1)
	} else {
		img = ""
	}

	return img, srcSet
}

func scrapeTheVerge(c chan<- TvergeArticle, URL string, path string) {
	scraper := colly.NewCollector()
	scraper.OnRequest(func(r *colly.Request) {
		log.Println("Visiting: ", r.URL.String())
	})

	scraper.OnError(func(r *colly.Response, err error) {
		log.Println("Error while scraping: ", err)
		close(c)
	})

	scraper.OnHTML(".duet--content-cards--content-card.group", func(h *colly.HTMLElement) {
		var tvergeArticle TvergeArticle

		// Category
		tvergeArticle.Category = category(h.ChildText("div:first-of-type>span>a"), path)
		tvergeArticle.CategoryLink = categoryLink(h.ChildAttr("div:first-of-type>span>a", "href"), path)
		// Title
		tvergeArticle.Title = h.ChildText("div:first-of-type div:first-of-type h2 a")
		// Date
		tvergeArticle.ArticleDate = h.ChildText("div:first-of-type div:first-of-type span span:last-child")
		// Author
		tvergeArticle.Author = h.ChildText("div:first-of-type div:first-of-type span> span:first-child>a")
		// Article URL
		tvergeArticle.URL = "https://www.theverge.com" + h.ChildAttr("div:first-of-type div:first-of-type h2 a", "href")
		// Image
		tvergeArticle.Img, tvergeArticle.ImgSrcSet = scrapeImageSrc(h.ChildText("div:last-of-type div:last-of-type .block a span img+noscript"))

		if tvergeArticle != (TvergeArticle{})  && tvergeArticle.URL != "https://www.theverge.com" && tvergeArticle.Img != "https://www.theverge.com" {
			c <- tvergeArticle
		}
	})

	scraper.OnScraped(func(r *colly.Response) {
		log.Println("Finished with scraping: ", r.Request.URL)
		close(c)
	})

	scraper.Visit(URL)
}
