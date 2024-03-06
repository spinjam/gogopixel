package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"path/filepath"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/solarlune/ldtkgo"
)

type TilesetLoader interface {
	LoadTileset(string) *ebiten.Image
}

// ----------------------------------------------------------------------------- DiskLoader struct
type DiskLoader struct {
	BasePath string
	Filter   ebiten.Filter
}

func NewDiskLoader(basePath string) *DiskLoader {
	return &DiskLoader{
		BasePath: basePath,
		Filter:   ebiten.FilterNearest,
	}
}

func (d *DiskLoader) LoadTileset(tilesetPath string) *ebiten.Image {
	fmt.Println("Loading Tileset ", tilesetPath)
	if img, _, err := ebitenutil.NewImageFromFile(filepath.Join(d.BasePath, tilesetPath)); err == nil {
		return img
	}
	return nil
}

// ----------------------------------------------------------------------------- RenderedLayer struct
type RenderedLayer struct {
	Image *ebiten.Image
	Layer *ldtkgo.Layer
}

// ---------------------------------------------------------------------------- Renderer struct
type Renderer struct {
	Tilesets       map[string]*ebiten.Image
	CurrentTileset string
	RenderedLayers []*RenderedLayer
	Offscreen      *ebiten.Image
	Loader         TilesetLoader
	camera         *Camera
	Buffer         *ebiten.Image
}

func NewRenderer(loader TilesetLoader, cam *Camera) *Renderer {

	//img, _, err := image.Decode(bytes.NewReader(images.Tiles_png))
	/*
		img, _, err := ebitenutil.NewImageFromFile("assets/tilesets/SunnyLand_by_Ansimuz-extended.png")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("tiles_png w=%d h=%d\n", img.Bounds().Dx(), img.Bounds().Dy())
	*/

	return &Renderer{
		Tilesets:       map[string]*ebiten.Image{},
		RenderedLayers: []*RenderedLayer{},
		Loader:         loader,
		camera:         cam,
	}
}

/*
func (er *EbitenRenderer) beginLayer(layer *ldtkgo.Layer, w, h int) {

	_, exists := er.Tilesets[layer.Tileset.Path]
	if !exists {
		fmt.Printf("Begin Layer > tileset path = %s\n", layer.Tileset.Path)
		er.Tilesets[layer.Tileset.Path] = er.Loader.LoadTileset(layer.Tileset.Path)
	}

	er.CurrentTileset = layer.Tileset.Path
	renderedImage := ebiten.NewImage(w, h)
	er.RenderedLayers = append(er.RenderedLayers, &RenderedLayer{Image: renderedImage, Layer: layer})
}
*/

func (er *Renderer) Load(level *ldtkgo.Level) {

	fmt.Println("---------------- Loading -------------------")
	fmt.Printf("LEVEL \tWidth=%d - Height=%d\n", level.Width, level.Height)
	for i := len(level.Layers) - 1; i >= 0; i-- {
		layer := level.Layers[i]
		fmt.Printf("LAYER [%s] >>> id=%s\n", layer.Type, layer.Identifier)
		fmt.Printf("\tgrid size=%d\n", layer.GridSize)
		fmt.Printf("\tcell wxh=%dx%d\n", layer.CellWidth, layer.CellHeight)
		fmt.Printf("\toffset x=%d y=%d\n", layer.OffsetX, layer.OffsetY)

		//er.beginLayer(layer, level.Width, level.Height)
		_, exists := er.Tilesets[layer.Tileset.Path]
		if !exists {
			// LDTK store external assets filenames with relative syntax
			newPath := "assets/" + strings.Replace(layer.Tileset.Path, "../", "", -1)
			fmt.Printf("Loading Tileset (path = %s)\n", newPath)
			tileimg, _, err := ebitenutil.NewImageFromFile(newPath)
			if err != nil {
				log.Fatal(err)
			}
			er.Tilesets[layer.Tileset.Path] = tileimg
		}
	}

	er.Offscreen = ebiten.NewImage(level.Width, level.Height)
	// only to test image buffer size
	er.Buffer = ebiten.NewImage(4096*2, 4096*2)
	fmt.Printf("buffer w=%d h=%d\n", er.Buffer.Bounds().Dx(), er.Buffer.Bounds().Dy())

	er.RenderOffscreen(level)
}

/*
	 ------------------------------------------------------------------------------
		Clear rendered layers (unused)
*/
func (er *Renderer) Clear() {
	for _, layer := range er.RenderedLayers {
		layer.Image.Dispose()
	}
	er.RenderedLayers = []*RenderedLayer{}
}

/*
	 ------------------------------------------------------------------------------
		Render the level (all layers flatten) offscreen
*/
func (er *Renderer) RenderOffscreen(level *ldtkgo.Level) {

	// disegno i layer in ordine inverso
	for i := len(level.Layers) - 1; i >= 0; i-- {

		layer := level.Layers[i]
		// er.Tilesets[layer.Tileset.Path] = tileimg
		tileimg := er.Tilesets[layer.Tileset.Path]
		// fmt.Printf("layer = %s\n", layer.Identifier)

		switch layer.Type {
		case ldtkgo.LayerTypeAutoTile:
			fallthrough
		case ldtkgo.LayerTypeIntGrid:
			fallthrough
		case ldtkgo.LayerTypeTile:
			//opacity := 1.0
			if tiles := layer.AllTiles(); len(tiles) > 0 {
				for _, tileData := range tiles {
					//fmt.Printf("%d-", tileData.ID)
					tilex := tileData.Position[0]
					tiley := tileData.Position[1]
					//fmt.Printf("tile x=%d y=%d\n", tilex, tiley)
					//rect := image.Rect(tilex, tiley, tilex+16, tiley+16)
					rect := image.Rect(tileData.Src[0], tileData.Src[1], tileData.Src[0]+layer.GridSize, tileData.Src[1]+layer.GridSize)
					tileimg := tileimg.SubImage(rect).(*ebiten.Image)

					opt := &ebiten.DrawImageOptions{}
					opt.GeoM.Translate(float64(-layer.GridSize/2), float64(-layer.GridSize/2))
					if tileData.FlipX() {
						opt.GeoM.Scale(-1, 1)
					}
					if tileData.FlipY() {
						opt.GeoM.Scale(1, -1)
					}
					opt.GeoM.Translate(float64(layer.GridSize/2), float64(layer.GridSize/2))

					opt.GeoM.Translate(float64(tilex), float64(tiley))
					if i == 1 {
						opt.ColorScale.ScaleAlpha(0.4)
					}
					er.Offscreen.DrawImage(tileimg, opt)

					//fmt.Printf("%v+", rect)
					//tile := er.Tilesets[er.CurrentTileset].SubImage(rect).(*ebiten.Image)
					//opt := &ebiten.DrawImageOptions{}
					//er.RenderedLayers[len(er.RenderedLayers)-1].Image.DrawImage(tile, opt)
				}
			}
			//fmt.Println(".")
		}
	}
}

/*
	 ------------------------------------------------------------------------------
		Render the level (all layers flatten) previously rendered offscreen
*/
func (er *Renderer) Render(screen *ebiten.Image, level *ldtkgo.Level) {

	er.camera.Surface.Clear()
	er.camera.Surface.Fill(color.RGBA{255, 128, 128, 255})

	area := image.Rect(
		int(er.camera.X), int(er.camera.Y), level.Width, level.Height)
	areaimg := er.Offscreen.SubImage(area).(*ebiten.Image)
	opt := &ebiten.DrawImageOptions{}
	screen.DrawImage(areaimg, opt)
}

/* ------------------------------------------------------------------------------
 */
func (er *Renderer) RenderLevel(screen *ebiten.Image, level *ldtkgo.Level) {

	er.Clear()

	//opt := &ebiten.DrawImageOptions{}
	opt2 := &ebiten.DrawImageOptions{}
	opt2.GeoM.Translate(100., 200.)
	//er.RenderedLayers[0].Image.DrawImage(screen, opt)

	//optArea := &ebiten.DrawImageOptions{}

	//screen.DrawImage(er.tilesImage, opt)
	// fmt.Printf("-----------\n")
	// disegno i layer in ordine inverso
	for i := len(level.Layers) - 1; i >= 0; i-- {

		layer := level.Layers[i]
		tileimg := er.Tilesets[layer.Tileset.Path]

		// fmt.Printf("layer = %s\n", layer.Identifier)
		switch layer.Type {
		case ldtkgo.LayerTypeAutoTile:
			fallthrough
		case ldtkgo.LayerTypeIntGrid:
			fallthrough
		case ldtkgo.LayerTypeTile:
			//opacity := 1.0
			if tiles := layer.AllTiles(); len(tiles) > 0 {
				for _, tileData := range tiles {
					//fmt.Printf("%d-", tileData.ID)
					tilex := tileData.Position[0]
					tiley := tileData.Position[1]
					//fmt.Printf("tile x=%d y=%d\n", tilex, tiley)
					//rect := image.Rect(tilex, tiley, tilex+16, tiley+16)
					rect := image.Rect(tileData.Src[0], tileData.Src[1], tileData.Src[0]+layer.GridSize, tileData.Src[1]+layer.GridSize)
					tileimg := tileimg.SubImage(rect).(*ebiten.Image)

					opt := &ebiten.DrawImageOptions{}
					opt.GeoM.Translate(float64(-layer.GridSize/2), float64(-layer.GridSize/2))
					if tileData.FlipX() {
						opt.GeoM.Scale(-1, 1)
					}
					if tileData.FlipY() {
						opt.GeoM.Scale(1, -1)
					}
					opt.GeoM.Translate(float64(layer.GridSize/2), float64(layer.GridSize/2))

					opt.GeoM.Translate(float64(tilex), float64(tiley))
					if i == 1 {
						opt.ColorScale.ScaleAlpha(0.4)
					}
					screen.DrawImage(tileimg, opt)
					//er.Offscreen.DrawImage(tileimg, opt)

					//fmt.Printf("%v+", rect)
					//tile := er.Tilesets[er.CurrentTileset].SubImage(rect).(*ebiten.Image)
					//opt := &ebiten.DrawImageOptions{}
					//er.RenderedLayers[len(er.RenderedLayers)-1].Image.DrawImage(tile, opt)
				}
			}
			//fmt.Println(".")
		}
	}

	/*
		for _, layer := range level.Layers {

			switch layer.Type {

			// IntGrids get autotiles automatically
			case ldtkgo.LayerTypeIntGrid:
				fallthrough
			case ldtkgo.LayerTypeAutoTile:
				fallthrough
			case ldtkgo.LayerTypeTile:

				if tiles := layer.AllTiles(); len(tiles) > 0 {

					er.beginLayer(layer, level.Width, level.Height)

					for _, tileData := range tiles {
						// er.renderTile(tile.Position[0]+layer.OffsetX, tile.Position[1]+layer.OffsetY, tile.Src[0], tile.Src[1], layer.GridSize, layer.GridSize, tile.Flip)

						// Subimage the Tile from the Tileset
						tile := er.Tilesets[er.CurrentTileset].SubImage(
							image.Rect(tileData.Src[0], tileData.Src[1], tileData.Src[0]+layer.GridSize, tileData.Src[1]+layer.GridSize)).(*ebiten.Image)

						opt := &ebiten.DrawImageOptions{}

						// We have to offset the tile to be centered before flipping
						opt.GeoM.Translate(float64(-layer.GridSize/2), float64(-layer.GridSize/2))

						// Handle flipping; first bit in byte is horizontal flipping, second is vertical flipping.

						if tileData.FlipX() {
							opt.GeoM.Scale(-1, 1)
						}
						if tileData.FlipY() {
							opt.GeoM.Scale(1, -1)
						}

						// Undo offsetting
						opt.GeoM.Translate(float64(layer.GridSize/2), float64(layer.GridSize/2))

						// Move tile to final position; note that slightly unlike LDtk, layer offsets in LDtk-Go are added directly into the final tiles' X and Y positions. This means that with this renderer,
						// if a layer's offset pushes tiles outside of the layer's render Result image, they will be cut off. On LDtk, the tiles are still rendered, of course.
						opt.GeoM.Translate(float64(tileData.Position[0]+layer.OffsetX), float64(tileData.Position[1]+layer.OffsetY))

						// Finally, draw the tile to the Result image.
						er.RenderedLayers[len(er.RenderedLayers)-1].Image.DrawImage(tile, opt)
					}
				}
			}
		}

		// Reverse sort the layers when drawing because in LDtk, the numbering order is from top-to-bottom, but the drawing order is from bottom-to-top.
		sort.Slice(er.RenderedLayers, func(i, j int) bool {
			return i > j
		})
	*/

}

func (er *Renderer) MoveCamera(dx int, dy int) {
	/*
		if er.camera.X >= 0 {
			er.camera.X += dx
		}
		if er.camera.Y >= 0 {
			er.camera.Y += dy
		}
	*/
}
