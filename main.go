package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
)

//go:embed frontend/dist
var assets embed.FS

var tvergeArticles []TvergeArticle

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:            "tverge-news",
		Width:            1024,
		MinWidth:         768,
		MaxWidth:         1280,
		Height:           768,
		Assets:           assets,
		BackgroundColour: &options.RGBA{R: 0, G: 0, B: 0, A: 1},
		OnStartup:        app.startup,
		OnDomReady:       app.domready,
		OnBeforeClose:    app.beforeClose,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
