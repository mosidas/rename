package main

import (
	"embed"
	"os"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/out
var assets embed.FS

func main() {
	// Get command-line arguments (excluding program name)
	argsWithoutProg := os.Args[1:]

	// Create an instance of the app structure
	app := NewApp()

	// Set initial files if provided via command-line
	if len(argsWithoutProg) > 0 {
		app.SetInitialFiles(argsWithoutProg)
	}

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "file rename",
		Width:  1524,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
		SingleInstanceLock: &options.SingleInstanceLock{
			UniqueId: "com.wails.rename",
			OnSecondInstanceLaunch: func(data options.SecondInstanceData) {
				// When a second instance is launched, load files in the existing instance
				if len(data.Args) > 1 {
					files := data.Args[1:] // Skip program name
					app.LoadFilesFromSecondInstance(files)
				}
			},
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
