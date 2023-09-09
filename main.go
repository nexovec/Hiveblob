package main

import (
	"fmt"
	"log"
	"os"

	pixelui "github.com/dusk125/pixelui"
	imgui "github.com/inkyblackness/imgui-go"
	colornames "golang.org/x/image/colornames"

	// beep "github.com/faiface/beep"
	box2d "github.com/E4/box2d"
	pixelutils "github.com/dusk125/pixelutils"

	// cb "github.com/emirpasic/gods/queues/circularbuffer"

	pixel "github.com/faiface/pixel"
	imdraw "github.com/faiface/pixel/imdraw"
	pixelgl "github.com/faiface/pixel/pixelgl"
	cobra "github.com/spf13/cobra"
)

var VERSION = "v0.0.0-dev"
var WINDOW_TITLE = "Hiveblob"
var WINDOW_HEIGHT float64 = 1080
var WINDOW_ASPECT_RATIO float64 = 16.0 / 9.0

var DRAW_PHYSICS_OBJ_OUTLINES = true
var VELOCITY_ITERATIONS int = 6
var POSITION_ITERATIONS int = 2
var SPRINT_FORCE float64 = 1400.0
var BOX2D_RENDERING_SCALE float64 = 10.0

type LaunchOptions struct {
	verbose        bool
	displayVersion bool
}

var launchOptions = struct {
	verbose     bool
	showVersion bool
	showFPS     bool
	fullscreen  bool
}{}

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

	transform := pixel.IM.Scaled(pixel.ZV, BOX2D_RENDERING_SCALE)

	// set-up physics
	gravityVector := box2d.B2Vec2{
		X: 0.0,
		Y: -100.0}
	world := box2d.MakeB2World(gravityVector)
	lookforwardWorld := box2d.MakeB2World(gravityVector)
	bodies := []*box2d.B2Body{}
	var (
		playerBody     *box2d.B2Body
		projectileBody *box2d.B2Body
	)
	{
		groundBodyDef := box2d.MakeB2BodyDef()
		groundBodyDef.Position.Set(40.0, 10.0)

		groundShape := box2d.MakeB2PolygonShape()
		groundShape.SetAsBox(35.0, 3.0)

		groundFixtureDef := box2d.MakeB2FixtureDef()
		groundFixtureDef.Density = 0.0
		groundFixtureDef.Friction = 1.0
		groundFixtureDef.Shape = &groundShape

		groundBody := world.CreateBody(&groundBodyDef)
		groundBody.CreateFixtureFromDef(&groundFixtureDef)
		bodies = append(bodies, groundBody)

		playerBodyDef := box2d.MakeB2BodyDef()
		playerBodyDef.Type = box2d.B2BodyType.B2_dynamicBody
		playerBodyDef.Position.Set(23.0, 23.0)

		playerShape := box2d.MakeB2PolygonShape()
		playerShape.SetAsBox(1.0, 2.0)

		playerFixtureDef := box2d.MakeB2FixtureDef()
		playerFixtureDef.Density = 1.0
		playerFixtureDef.Friction = 0.8
		playerFixtureDef.Shape = &playerShape

		playerB2Body := world.CreateBody(&playerBodyDef)
		playerB2Body.CreateFixtureFromDef(&playerFixtureDef)
		playerB2Body.SetFixedRotation(true)
		bodies = append(bodies, playerB2Body)
		playerBody = playerB2Body

		projectileBodyDef := box2d.MakeB2BodyDef()
		projectileBodyDef.Type = box2d.B2BodyType.B2_dynamicBody
		projectileBodyDef.Position.Set(23.0, 25.0)

		projectileShape := box2d.MakeB2PolygonShape()
		projectileShape.SetAsBox(1.0, 1.0)

		projectileBody = lookforwardWorld.CreateBody(&projectileBodyDef)
		// projectileBody := world.CreateBody(&projectileBodyDef)
		// shape, ok := projectileShape.(box2d.B2ShapeInterface)
		// if ok != nil {
		// 	panic("oof")
		// }
		projectileBody.CreateFixture(&projectileShape, 1.0)
		bodies = append(bodies, projectileBody)
	}

	BUFFERED_PROJECTILE_POSITIONS := 10
	PROJECTILE_SIMULATION_SKIP_IN_TICKS := 6
	impulse := box2d.MakeB2Vec2(300.0, 255.0)
	projectileBody.ApplyLinearImpulseToCenter(impulse, true)
	projectilePositionsBuffer := []box2d.B2Vec2{}
	for i := 0; i < BUFFERED_PROJECTILE_POSITIONS; i++ {
		timestep := ticker.TargetFrametime().Seconds() * float64(PROJECTILE_SIMULATION_SKIP_IN_TICKS)
		lookforwardWorld.Step(timestep, 6, 2)
		projectilePositionsBuffer = append(projectilePositionsBuffer, projectileBody.GetPosition())

		// lookforwardWorld.Step(0.1, 6, 2)
	}

	// simulate projectile path
	// staticBodies := []*box2d.B2Body{}
	// simulatedStaticBodies := []*box2d.B2Body{}
	// BUFFERED_PROJECTILE_POSITIONS := 50
	// projectileBodyTemporalBuffer := cb.New(BUFFERED_PROJECTILE_POSITIONS)
	// {
	// 	// // filter out static box2d bodies
	// 	// for _, body := range bodies {
	// 	// 	// fmt.Println(body.GetType() == box2d.B2BodyType.B2_staticBody)
	// 	// 	if body.GetType() == box2d.B2BodyType.B2_staticBody {
	// 	// 		staticBodies = append(staticBodies, body)
	// 	// 	}
	// 	// }
	// 	// log.Println("Static body count: ", len(staticBodies))
	// 	futureSimWorld := box2d.MakeB2World(gravityVector)

	// 	// populate static bodies into a separate box2d world
	// 	// for _, body := range staticBodies {
	// 	// 	bodyDef := box2d.MakeB2BodyDef()
	// 	// 	bodyDef.Type = body.GetType()
	// 	// 	bodyDef.Position = body.GetPosition()

	// 	// 	body := futureSimWorld.CreateBody(&bodyDef)
	// 	// 	for fixture := body.GetFixtureList(); fixture != nil; fixture.GetNext() {
	// 	// 		// TODO: test
	// 	// 		body.CreateFixture(fixture.GetShape(), fixture.GetDensity())
	// 	// 	}
	// 	// 	simulatedStaticBodies = append(simulatedStaticBodies, body)
	// 	// }

	// 	// projectileShape := box2d.MakeB2CircleShape()
	// 	// projectileShape.M_radius = 0.3
	// 	// projectileShape := box2d.MakeB2PolygonShape()
	// 	// projectileShape.SetAsBox(1.0, 1.0)

	// 	// projectileBody := futureSimWorld.CreateBody(&projectileBodyDef)
	// 	// projectileBody.CreateFixture(&projectileShape, 1.0)
	// 	// projectileBody.GetFixtureList().SetFriction(0.1) // FIXME: only does one fixture
	// 	// projectileBody.M_xf.Set(box2d.B2Vec2{25.0, 25.0}, 0.0)
	// 	// bodies = append(bodies, projectileBody)
	// 	// v := box2d.MakeB2Vec2(0.0, 255.0)
	// 	// projectileBody.SetLinearVelocity(v)
	// 	// projectileBody.ApplyLinearImpulseToCenter(v, true)

	// 	// 	// simulate a projectile
	// 	// 	// @nocheckin fix projectile rendering first
	// 	// 	for i := 0; i < BUFFERED_PROJECTILE_POSITIONS; i++ {
	// 	// 		projectileBodyDef.Position = projectileBody.GetPosition()
	// 	// 		b := futureSimWorld.CreateBody(&projectileBodyDef)
	// 	// 		b.CreateFixture(&projectileShape, 1.0)
	// 	// 		projectileBodyTemporalBuffer.Enqueue(projectileBody)

	// 	// 		tgtFrametime := ticker.TargetFrametime()
	// 	// 		futureSimWorld.Step(float64(tgtFrametime), 6, 2)
	// 	// 	}
	// }

	for !window.Closed() {

		// game update
		ticker.Tick()
		dt := ticker.Deltat()
		// log.Printf("Tick delta time: %5f", dt)

		// handle controls
		if window.Pressed(pixelgl.KeyD) {
			playerBody.ApplyForce(box2d.MakeB2Vec2(SPRINT_FORCE, 0.0), playerBody.GetWorldCenter(), true)
		}
		if window.Pressed(pixelgl.KeyA) {
			playerBody.ApplyForce(box2d.MakeB2Vec2(-SPRINT_FORCE, 0.0), playerBody.GetWorldCenter(), true)
		}
		if window.JustPressed(pixelgl.KeyW) {
			playerBody.ApplyLinearImpulseToCenter(box2d.MakeB2Vec2(0.0, 255.0), true)
		}

		world.Step(dt, VELOCITY_ITERATIONS, POSITION_ITERATIONS)

		// game rendering
		window.Clear(colornames.Beige)
		if DRAW_PHYSICS_OBJ_OUTLINES {
			ctx := imdraw.New(nil)
			ctx.SetMatrix(transform)
			// render box2d bounding boxes
			ctx.Color = colornames.Indianred
			for _, body := range bodies {
				aabb := body.GetFixtureList().GetAABB(0)
				startX, startY := aabb.LowerBound.X, aabb.LowerBound.Y
				endX, endY := aabb.UpperBound.X, aabb.UpperBound.Y
				ctx.Push(pixel.V(startX, startY), pixel.V(endX, endY))
				ctx.Rectangle(0.0)
			}
			// render box2d body outlines
			ctx.Color = colornames.Black
			for _, body := range bodies {
				// for _, body := range simulatedStaticBodies { // FIXME: crashes, why
				drawBox2dBodyOutline(body, ctx)
			}
			for _, position := range projectilePositionsBuffer {
				pos := pixel.Vec{X: position.X, Y: position.Y}
				ctx.Push(pos)
				ctx.Circle(0.2, 0.0)
			}
			// for it := projectileBodyTemporalBuffer.Iterator(); it.Next(); {
			// 	body := it.Value()
			// 	drawBox2dBodyOutline(body.(*box2d.B2Body), ctx)
			// }
			ctx.Draw(window)

		}

		// UI goes here
		ui.NewFrame()
		imgui.Begin("Base")
		// FIXME: https://github.com/dusk125/pixelui/issues/15
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

// NOTE: don't forget to call ctx.draw(window) to blit to your window
func drawBox2dBodyOutline(body *box2d.B2Body, ctx *imdraw.IMDraw) {
	DEBUG_fixture := body.GetFixtureList()
	if DEBUG_fixture.GetType() == box2d.B2Shape_Type.E_polygon {
		shape := DEBUG_fixture.GetShape()
		polygon := shape.(*box2d.B2PolygonShape)
		vertices := polygon.M_vertices
		// ctx.Reset()
		for i := 0; i < polygon.M_count; i++ {
			vertex := vertices[i]
			nextVertex := vertices[(i+1)%polygon.M_count]
			position := body.GetWorldPoint(vertex)
			nextPosition := body.GetWorldPoint(nextVertex)
			p1 := pixel.V(position.X, position.Y)
			p2 := pixel.V(nextPosition.X, nextPosition.Y)
			ctx.Push(p1, p2)
			ctx.Line(5.0 / BOX2D_RENDERING_SCALE)
		}
	} else if DEBUG_fixture.GetType() == box2d.B2Shape_Type.E_circle {
		ctx.Circle(DEBUG_fixture.GetShape().GetRadius(), 5.0)
	} else {
		panic("Not implemented!")
	}
}

func main() {
	fmt.Println("LAUNCHING THE GAME.")
	log.SetOutput(os.Stdout)
	// parse executable arguments
	rootCmd := cobra.Command{
		Use:     "hiveblob [--version] [--verbose|--v] [--fullscreen|--f] [--fps]",
		Short:   "Hiveblob the game development build executable",
		Example: "hiveblob --version",
		Run: func(cmd *cobra.Command, args []string) {
			if launchOptions.showVersion {
				log.Default().Println("version:", VERSION)
			}
			pixelgl.Run(run)
		},
	}
	{
		rootCmd.Flags().BoolVarP(&launchOptions.verbose, "verbose", "v", false, "print more logs")
		rootCmd.Flags().BoolVar(&launchOptions.showVersion, "version", false, "show version number")
		rootCmd.Flags().BoolVar(&launchOptions.showFPS, "fps", false, "show frames per second in the corner")
	}
	if err := rootCmd.Execute(); err != nil {
		log.Default().Fatalln(err)
	}
}
