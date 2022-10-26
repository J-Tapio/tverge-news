package main

import (
	"context"
	"log"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// Article struct
type TvergeArticle struct {
	Category     string `json:"category"`
	CategoryLink string `json:"categoryLink"`
	Title        string `json:"title"`
	ArticleDate  string `json:"date"`
	Author       string `json:"author"`
	URL          string `json:"URL"`
	Img          string `json:"img"`
	ImgSrcSet    string `json:"imgSrcSet"`
}

var sourceURL = "https://www.theverge.com"

var categoryXmlURL = map[string]string{
	"Tech":          "https://www.theverge.com/rss/tech/index.xml",
	"Science":       "https://www.theverge.com/rss/science/index.xml",
	"Reviews":       "https://www.theverge.com/rss/reviews/index.xml",
	"Entertainment": "https://www.theverge.com/rss/entertainment/index.xml",
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) domready(ctx context.Context) {
	go startScraper(a)
}

func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	dialog, err := runtime.MessageDialog(ctx, runtime.MessageDialogOptions{
		Type:    runtime.QuestionDialog,
		Title:   "Quit?",
		Message: "Are you sure you want to quit?",
	})

	if err != nil {
		return false
	}
	return dialog != "Yes"
}

// Binding to frontend
func (a *App) Latest() []TvergeArticle {
	return tvergeArticles
}

func outputToTvergeArticles(c <-chan TvergeArticle) {
	for {
		article := <-c
		tvergeArticles = append(tvergeArticles, article)
	}
}

func startScraper(app *App) {
	for {
		// Remove old articles
		if len(tvergeArticles) > 0 {
			tvergeArticles = []TvergeArticle{}
		}

		log.Println("Fetchin latest news from The Verge")
		// Channel per site request / goroutine
		fromTech := make(chan TvergeArticle, 10)
		fromScience := make(chan TvergeArticle, 10)
		fromReviews := make(chan TvergeArticle, 10)
		fromEnt := make(chan TvergeArticle, 10)
		toVergeArticles := make(chan TvergeArticle, 10)

		fromTechOpen := true
		fromScienceOpen := true
		fromReviewsOpen := true
		fromEntOpen := true

		go scrapeTheVergeXML(fromTech, categoryXmlURL["Tech"], "Tech")
		go scrapeTheVergeXML(fromScience, categoryXmlURL["Science"], "Science")
		go scrapeTheVergeXML(fromReviews, categoryXmlURL["Reviews"], "Reviews")
		go scrapeTheVergeXML(fromEnt, categoryXmlURL["Entertainment"], "Entertainment")
		go outputToTvergeArticles(toVergeArticles)

		for fromTechOpen || fromScienceOpen || fromReviewsOpen || fromEntOpen {
			select {
			case techArticle, open := <-fromTech:
				{
					if open {
						toVergeArticles <- techArticle
					} else {
						fromTechOpen = false
					}
				}
			case scienceArticle, open := <-fromScience:
				{
					if open {
						toVergeArticles <- scienceArticle
					} else {
						fromScienceOpen = false
					}
				}
			case reviewArticle, open := <-fromReviews:
				{
					if open {
						toVergeArticles <- reviewArticle
					} else {
						fromReviewsOpen = false
					}
				}
			case entArticle, open := <-fromEnt:
				{
					if open {
						toVergeArticles <- entArticle
					} else {
						fromEntOpen = false
					}
				}
			}
		}

		runtime.EventsEmit(app.ctx, "news")
		time.Sleep(1 * time.Hour)
	}
}
