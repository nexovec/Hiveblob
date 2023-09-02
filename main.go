package main

import (
	"fmt"
	"log"
	"os"

	pixelui "github.com/dusk125/pixelui"
	"github.com/inkyblackness/imgui-go"
	"golang.org/x/image/colornames"

	// beep "github.com/faiface/beep"
	box2d "github.com/E4/box2d"
	pixelutils "github.com/dusk125/pixelutils"
	pixel "github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	pixelgl "github.com/faiface/pixel/pixelgl"

	cobra "github.com/spf13/cobra"
)

var VERSION = "v0.0.0-dev"
var WINDOW_TITLE = "Hiveblob"
var WINDOW_HEIGHT float64 = 1080
var WINDOW_ASPECT_RATIO float64 = 16.0 / 9.0

var VELOCITY_ITERATIONS int = 6
var POSITION_ITERATIONS int = 2

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
	ui.AddTTFFont("resources/fonts/bitstream_vera_mono/VeraMono.ttf", 48)
	// font := imgui.CurrentIO().Fonts().AddFontFromFileTTF("resources/fonts/bitstream_vera_mono/VeraMono.ttf", 12)
	// fmt.Printf("imgui.CurrentIO().Fonts(): %v\n", imgui.CurrentIO().Fonts())
	// ui.AddTTFFont()
	// imgui.PushFont(ui.fonts.)

	ticker := pixelutils.NewTicker(60)

	// set-up physics
	gravityVector := box2d.B2Vec2{
		X: 0.0,
		Y: 0.0}
	world := box2d.MakeB2World(gravityVector)

	groundBodyDef := box2d.MakeB2BodyDef()
	groundBodyDef.Position.Set(20.0, 10.0)

	groundShape := box2d.MakeB2PolygonShape()
	groundShape.SetAsBox(15.0, 3.0)

	groundFixture := box2d.MakeB2Fixture()
	groundFixture.SetDensity(0.0)
	groundFixture.M_shape = &groundShape

	groundBody := world.CreateBody(&groundBodyDef)
	groundBody.CreateFixture(&groundShape, 0.0)

	playerBodyDef := box2d.MakeB2BodyDef()
	// playerBodyDef.Type = box2d.B2BodyType.B2_dynamicBody
	playerBodyDef.Position.Set(23.0, 4.0)

	playerShape := box2d.MakeB2PolygonShape()
	playerShape.SetAsBox(1.0, 2.0)

	playerFixtureDef := box2d.MakeB2FixtureDef()
	playerFixtureDef.Density = 1.0
	playerFixtureDef.Friction = 0.3
	playerFixtureDef.Shape = &playerShape

	playerBody := world.CreateBody(&playerBodyDef)
	playerBody.CreateFixtureFromDef(&playerFixtureDef)

	// gameObjects := struct
	for !window.Closed() {

		// game update
		dt := ticker.Deltat()

		world.Step(dt, VELOCITY_ITERATIONS, POSITION_ITERATIONS)
		// game rendering
		window.Clear(colornames.Beige)
		// DEBUG: render a circle out of on the screen.
		// {
		// 	ctx := imdraw.New(nil)
		// 	ctx.Color = colornames.Black
		// 	ctx.Push(pixel.V(100.0, 0.0))
		// 	ctx.Circle(100.0, 100.0)
		// 	ctx.Draw(window)
		// }
		DRAW_PHYSICS_OBJ_OUTLINES := true
		if DRAW_PHYSICS_OBJ_OUTLINES {
			DEBUG_fixture := groundBody.GetFixtureList()
			if DEBUG_fixture.GetType() == box2d.B2Shape_Type.E_polygon {
				shape := DEBUG_fixture.GetShape()
				shape.Clone()
				polygon := shape.(*box2d.B2PolygonShape)
				vertices := polygon.M_vertices
				for i := 0; i < polygon.M_count; i++ {
					vertex := vertices[i]
					nextVertex := vertices[(i+1)%polygon.M_count]
					position := groundBody.GetWorldPoint(vertex)
					nextPosition := groundBody.GetWorldPoint(nextVertex)
					BOX2D_RENDERING_SCALE := 10.0
					p1 := pixel.V(position.X, position.Y).Scaled(BOX2D_RENDERING_SCALE)
					p2 := pixel.V(nextPosition.X, nextPosition.Y).Scaled(BOX2D_RENDERING_SCALE)
					// line := pixel.Line{
					// 	p1,
					// 	p2,
					// }
					ctx := imdraw.New(nil)
					ctx.Color = colornames.Black
					fmt.Println(p1, p2)
					ctx.Push(p1, p2)
					ctx.Line(5.0)
					ctx.Draw(window)
				}
			} else {
				panic("Not implemented!")
			}
		}

		// UI goes here
		ui.NewFrame()
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
