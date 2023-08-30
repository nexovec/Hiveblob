package main

import (
	"fmt"
	"log"
	"os"

	pixelui "github.com/dusk125/pixelui"
	// beep "github.com/faiface/beep"
	pixelutils "github.com/dusk125/pixelutils"
	pixel "github.com/faiface/pixel"
	pixelgl "github.com/faiface/pixel/pixelgl"
	colornames "golang.org/x/image/colornames"

	cobra "github.com/spf13/cobra"
)

var VERSION = "v0.0.0-dev"
var WINDOW_TITLE = "Hiveblob"

type LaunchOptions struct {
	verbose        bool
	displayVersion bool
}

var isVerbose bool
var displayVersion bool

func run() {
	conf := pixelgl.WindowConfig{
		Title:     WINDOW_TITLE,
		Bounds:    pixel.R(0, 0, 1280, 720),
		Resizable: false,
		VSync:     true,
	}
	window, err := pixelgl.NewWindow(conf)
	if err != nil {
		panic("Could not create a window")
	}

	ui := pixelui.NewUI(window, 0)
	defer ui.Destroy()

	ticker := pixelutils.NewTicker(60)
	for !window.Closed() {
		window.Clear(colornames.Beige)
		ui.NewFrame()

		// UI goes here

		ui.Draw(window)

		window.Update()
		ticker.Wait()
	}
}

func main() {
	fmt.Println("LAUNCHING THE GAME.")
	// runGame()
	// panic("")
	log.SetOutput(os.Stdout)
	// parse executable arguments
	rootCmd := cobra.Command{
		Use:     "hiveblob [--version|--v]",
		Short:   "Hiveblob the game development build executable",
		Example: "hiveblob --server",
		Run: func(cmd *cobra.Command, args []string) {
			if displayVersion {
				log.Default().Println("version:", VERSION)
			}
			pixelgl.Run(run)
		},
	}
	rootCmd.Flags().BoolVar(&isVerbose, "verbose", false, "hiveblob --verbose")
	rootCmd.Flags().BoolVar(&displayVersion, "version", false, "hiveblob --version")
	if err := rootCmd.Execute(); err != nil {
		log.Default().Fatalln(err)
	}
}
