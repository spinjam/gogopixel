package main

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/math/f64"
)

type Camera struct {
	X, Y, Rot, Scale float64
	Width, Height    int
	Surface          *ebiten.Image
	// -------
	Viewport f64.Vec2
	Position f64.Vec2
}

func NewCamera(width, height int, x, y, rotation, zoom float64) *Camera {
	return &Camera{
		X:        x,
		Y:        y,
		Width:    width,
		Height:   height,
		Rot:      rotation,
		Scale:    zoom,
		Surface:  ebiten.NewImage(width, height),
		Viewport: f64.Vec2{0., 0.},
	}
}

func (c *Camera) SetPosition(x, y float64) *Camera {
	c.X = x
	c.Y = y
	return c
}

func (c *Camera) MovePosition(dx, dy float64) *Camera {
	c.X += dx
	c.Y += dy
	return c
}

func (c *Camera) Blit(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	w, h := c.Surface.Bounds().Dx(), c.Surface.Bounds().Dy()
	cx := float64(w) / 2.0
	cy := float64(h) / 2.0

	op.GeoM.Translate(-cx, -cy)
	op.GeoM.Scale(c.Scale, c.Scale)
	op.GeoM.Rotate(c.Rot)
	op.GeoM.Translate(cx*c.Scale, cy*c.Scale)

	screen.DrawImage(c.Surface, op)
}

func (c *Camera) Info() {
	w, h := c.Surface.Bounds().Dx(), c.Surface.Bounds().Dy()
	fmt.Println("--------------- Camera Info ----------------")
	fmt.Printf("w=%d - h=%d\n", w, h)
	//fmt.Println("--------------------------------------------")
}

// get screen coords from world coords
func (c *Camera) WorldToScreenCoords(x, y float64) (float64, float64) {
	// Extracts the width and height of the camera's viewport.
	w, h := c.Width, c.Height
	// Calculates the cosine and sine of the camera's rotation angle Rot. This will be used to rotate the coordinates.
	co := math.Cos(c.Rot)
	si := math.Sin(c.Rot)

	// Translates the given world coordinates by subtracting the camera's position X and Y.
	// This makes the camera the origin in the world space.
	x, y = x-c.X, y-c.Y
	// Rotates the translated coordinates using rotation matrices. This applies the camera's rotation to the coordinates.
	x, y = co*x-si*y, si*x+co*y

	// Scales the rotated coordinates by the camera's scale Scale and
	// then translates them to the center of the screen space
	// (adding half of the viewport width and height).
	// This effectively maps the rotated, scaled, and translated world coordinates to the screen coordinates.
	return x*c.Scale + float64(w)/2, y*c.Scale + float64(h)/2
}

func (c *Camera) ScreenToWorldCoords(x, y float64) (float64, float64) {
	// Extracts the width and height of the camera's viewport.
	w, h := c.Width, c.Height
	// Calculate the cosine and sine of the negative of the camera's rotation angle.
	// Negative rotation is used to reverse the effect of rotation on the coordinates.
	co := math.Cos(-c.Rot)
	si := math.Sin(-c.Rot)

	// Translate the coordinates to the center of the viewport and scale them
	// according to the camera's scale.
	x, y = (x-float64(w)/2)/c.Scale, (y-float64(h)/2)/c.Scale
	x, y = co*x-si*y, si*x+co*y

	return x + c.X, y + c.Y
}

/*
func (c *Camera) Reset() {
	c.Position[0] = 0.
	c.Position[1] = 0.
}

func (c *Camera) viewportCenter() f64.Vec2 {
	return f64.Vec2{
		c.Viewport[0] * 0.5,
		c.Viewport[1] * 0.5,
	}
}
*/
