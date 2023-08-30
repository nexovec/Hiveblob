package main

import (
	"fmt"
	pixelui "github.com/dusk125/pixelui"
	"github.com/inkyblackness/imgui-go"
	"log"
	"os"
	// beep "github.com/faiface/beep"
	pixelutils "github.com/dusk125/pixelutils"
	pixel "github.com/faiface/pixel"
	pixelgl "github.com/faiface/pixel/pixelgl"
	colornames "golang.org/x/image/colornames"

	cobra "github.com/spf13/cobra"
)

var VERSION = "v0.0.0-dev"
var WINDOW_TITLE = "Hiveblob"
var WINDOW_HEIGHT float64 = 1080
var WINDOW_ASPECT_RATIO float64 = 16.0 / 9.0

type LaunchOptions struct {
	verbose        bool
	displayVersion bool
}

var isVerbose bool
var displayVersion bool

func run() {
	windowX := float64(0)
	windowY := float64(0)
	conf := pixelgl.WindowConfig{
		Title:     WINDOW_TITLE,
		Bounds:    pixel.R(windowX, windowY, windowX+WINDOW_HEIGHT*WINDOW_ASPECT_RATIO, windowY+WINDOW_HEIGHT),
		Resizable: false,
		VSync:     true,
	}
	window, err := pixelgl.NewWindow(conf)
	if err != nil {
		panic("Could not create a window")
	}

	ui := pixelui.NewUI(window, 0)
	defer ui.Destroy()

	// font := imgui.CurrentIO().Fonts().AddFontDefault()
	// ui.AddTTFFont("resources/fonts/bitstream_vera_mono/VeraMono.ttf", 48)
	// font := imgui.CurrentIO().Fonts().AddFontFromFileTTF("resources/fonts/bitstream_vera_mono/VeraMono.ttf", 12)
	// ui.AddTTFFont()
	// imgui.PushFont(font)

	ticker := pixelutils.NewTicker(60)
	for !window.Closed() {
		window.Clear(colornames.Beige)
		ui.NewFrame()

		// UI goes here
		imgui.Begin("Base")
		// imgui.PushFont(font)
		imgui.Text("Placeholder text")
		// imgui.PopFont()
		imgui.End()

		imgui.ShowDemoWindow(nil)

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
