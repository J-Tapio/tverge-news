package main

import (
	"context"
	"fmt"
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

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	go startScraper()
}

func (a *App) domready(ctx context.Context) {}

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

// For some reason The Verge mixes up news within different sections, eg.
// Entertainment news article might be found within Science section
// Hence the need for duplicate removal.
func removeDuplicates(articles []TvergeArticle) []TvergeArticle {
	var uniqueArticles []TvergeArticle
	// Make map of [string]bool - in this case the url should be unique
	keys := make(map[string]bool)
	// Loop articles
	for _, article := range articles {
		if keys[article.URL] {
			continue
		} else {
			keys[article.URL] = true
			uniqueArticles = append(uniqueArticles, article)
		}
	}
	return uniqueArticles
}

func outputToTvergeArticles(c <-chan TvergeArticle) {
	for {
		article := <-c
		tvergeArticles = append(tvergeArticles, article)
	}
}

func startScraper() {
	for {
		// If articles - remove old
		if len(tvergeArticles) > 0 {
			tvergeArticles = []TvergeArticle{}
		}
		fmt.Println("Fetchin latest news from The Verge")
		// Channel per site request
		fromTech := make(chan TvergeArticle, 10)
		fromScience := make(chan TvergeArticle, 10)
		fromReviews := make(chan TvergeArticle, 10)
		fromEnt := make(chan TvergeArticle, 10)
		toVergeArticles := make(chan TvergeArticle, 10)

		fromTechOpen := true
		fromScienceOpen := true
		fromReviewsOpen := true
		fromEntOpen := true

		go scrapeTheVerge(fromTech, "https://www.theverge.com/tech", "Tech")
		go scrapeTheVerge(fromScience, "https://www.theverge.com/science", "Science")
		go scrapeTheVerge(fromReviews, "https://www.theverge.com/reviews", "Reviews")
		go scrapeTheVerge(fromEnt, "https://www.theverge.com/entertainment", "Entertainment")
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
		tvergeArticles = removeDuplicates(tvergeArticles)
		time.Sleep(1 * time.Hour)
	}
}

