package main

import (
	"log"
	"regexp"
	"strings"
	
	"github.com/gocolly/colly"
)

func category(c string, path string) string {
	// If category not available somehow - category set from url path
	if c == "" {
		return path
	}
	// For some reason scraping returns eg. 'WebWeb'
	// Return only one occurence of category word
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
		return img, srcSet
	}
	imgSrcSetMatch := re.FindStringSubmatch(source)
	srcSet = imgSrcSetMatch[1]
	// Clean 'amp;' from the strings
	if len(imgSrcSetMatch) > 0 {
		srcSet = strings.Replace(srcSet, "amp;", "", -1)
	}

	// imgSrc
	re3, err := regexp.Compile(`src="(.*?)"\sdecoding`)
	if err != nil {
		log.Println("Error with regex: ", err)
		return img, srcSet
	}

	imgSrcMatch := re3.FindStringSubmatch(source)
	img = imgSrcMatch[1]
	if len(imgSrcMatch) > 0 {
		imgUncleaned := "https://www.theverge.com" + imgSrcMatch[1]
		img = strings.Replace(imgUncleaned, "amp;", "", -1)
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

		if tvergeArticle != (TvergeArticle{}) && tvergeArticle.URL != "https://www.theverge.com" && tvergeArticle.Img != "https://www.theverge.com" {
			c <- tvergeArticle
		}
	})

	scraper.OnScraped(func(r *colly.Response) {
		log.Println("Finished with scraping: ", r.Request.URL)
		close(c)
	})

	scraper.Visit(URL)
}
