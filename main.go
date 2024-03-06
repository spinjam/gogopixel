package main

import (
	"fmt"
	_ "image/png"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/solarlune/ldtkgo"
)

const (
	ScreenW = 1280
	ScreenH = 720
)

// -------------------------------------------------------------
type Game struct {
	player         *Player
	LDTKProject    *ldtkgo.Project
	EbitenRenderer *Renderer
	CurrentLevel   int
	time           int64
	camera         *Camera
}

func NewGame() *Game {

	cam := NewCamera(ScreenW, ScreenH, 0, 0, 0, 1.0)
	cam.Info()

	g := &Game{
		player: NewPlayer(),
		camera: cam,
	}

	var err error
	g.LDTKProject, err = ldtkgo.Open("assets/map/map1.ldtk")
	if err != nil {
		panic(err)
	}
	fmt.Printf("LTDK JSON ver = %s\n", g.LDTKProject.JSONVersion)

	fmt.Println("--- Tilesets")
	for i, tileset := range g.LDTKProject.Tilesets {
		fmt.Printf("%d: %d - Tileset id = %s - path = %s\n", i, tileset.ID, tileset.Identifier, tileset.Path)
	}
	fmt.Println("--- Levels")
	for i, level := range g.LDTKProject.Levels {
		fmt.Printf("%d: level %s\n", i, level.Identifier)
	}
	g.CurrentLevel = 0

	g.EbitenRenderer = NewRenderer(NewDiskLoader(""), cam)
	g.EbitenRenderer.Load(g.LDTKProject.Levels[g.CurrentLevel])

	g.time = 0

	return g
}

/*
// repeatingKeyPressed return true when key is pressed considering the repeat state.
func repeatingKeyPressed(key ebiten.Key) bool {
	const (
		delay    = 30
		interval = 3
	)
	d := inpututil.KeyPressDuration(key)
	if d == 1 {
		return true
	}
	if d >= delay && (d-delay)%interval == 0 {
		return true
	}
	return false
}
*/

func (g *Game) Update() error {
	// --- move camera
	if ebiten.IsKeyPressed(ebiten.Key2) {
		//camX := g.camera.X
		g.camera.MovePosition(5.0, 0)
		//g.EbitenRenderer.MoveCamera(5, 0)
	}
	if ebiten.IsKeyPressed(ebiten.Key1) {
		g.camera.MovePosition(-5.0, 0)
		//g.EbitenRenderer.MoveCamera(-5, 0)
	}
	// --- move player
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.player.Move(1.0, 0.0)
		//g.EbitenRenderer.MoveCamera(1, 0)
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.player.Move(-1.0, 0.0)
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.player.Move(0., -1.)
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.player.Move(0., 1.)
	}
	/*
		if repeatingKeyPressed(ebiten.KeyRight) {
			g.player.Move(1.0, 0)
		}
		if repeatingKeyPressed(ebiten.KeyLeft) {
			g.player.Move(-1.0, 0)
		}
	*/

	if inpututil.IsKeyJustPressed(ebiten.KeyA) {
		g.player.CycleAnim()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		os.Exit(0)
	}
	g.player.Update()
	g.time += 1

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	//g.RenderLevel(screen)
	level := g.LDTKProject.Levels[g.CurrentLevel]
	g.EbitenRenderer.Render(screen, level)

	g.player.Draw(screen)

	//screen.Fill(color.RGBA{0x33, 0x33, 0x33, 0xff})
	if (g.time / 60) > 5.0 {
		ebitenutil.DebugPrint(screen, "Ebiten Engine (after 5 sec)")
	}
	/*
		for _, layer := range g.EbitenRenderer.RenderedLayers {
			fmt.Println("draw layer ", layer.Layer.Identifier)
			screen.DrawImage(layer.Image, &ebiten.DrawImageOptions{})
		}
	*/
}

func (g *Game) RenderLevel(screen *ebiten.Image) {

	level := g.LDTKProject.Levels[g.CurrentLevel]
	g.EbitenRenderer.Render(screen, level)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenW / 2, ScreenH / 2
}

/*
var assets embed.FS

func readImage(file string) image.Image {
	b, _ := assets.ReadFile(file)
	return bytes2Image(&b)
}

func bytes2Image(rawImage *[]byte) image.Image {
	img, format, error := image.Decode(bytes.NewReader(*rawImage))
	if error != nil {
		log.Fatal("Bytes2Image Failed: ", format, error)
	}
	return img
}
*/

func main() {
	ebiten.SetWindowSize(ScreenW, ScreenH)
	ebiten.SetWindowTitle("Goblit")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
