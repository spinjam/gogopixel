package main

import (
	"fmt"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/yohamta/ganim8/v2"
)

type Dir int

const (
	Dir_Right Dir = iota
	Dir_Left
)

// ------------------------------------------------
type PlayerState int

const (
	Player_Idle PlayerState = iota
	Player_Run
	Player_Jump
	Player_Climb
)

// String - Creating common behavior - give the type a String function
func (d PlayerState) String() string {
	return [...]string{"Idle", "Run", "Jump", "Climb"}[d]
}

// EnumIndex - Creating common behavior - give the type a EnumIndex functio
func (d PlayerState) EnumIndex() int {
	return int(d)
}

// ----------------------------------------------
const (
	MoveDx = 1.5
	FrameW = 32
	FrameH = 32
)

type Player struct {
	x, y      float64
	velocity  Vec2D[float64]
	dir       Dir
	state     PlayerState
	grounded  bool
	images    map[PlayerState]*ebiten.Image
	anims     map[PlayerState]*ganim8.Animation
	curr_anim *ganim8.Animation
}

func NewPlayer() *Player {
	p := &Player{}
	// init "map"s for images and animations
	p.images = map[PlayerState]*ebiten.Image{}
	p.anims = map[PlayerState]*ganim8.Animation{}

	var err error
	p.images[Player_Idle], _, err = ebitenutil.NewImageFromFile("assets/hero/Pink Man/Idle (32x32).png")
	if err != nil {
		log.Fatal(err)
	}
	p.images[Player_Run], _, err = ebitenutil.NewImageFromFile("assets/hero/Pink Man/Run (32x32).png")
	if err != nil {
		log.Fatal(err)
	}
	p.images[Player_Jump], _, err = ebitenutil.NewImageFromFile("assets/hero/Pink Man/Jump (32x32).png")
	if err != nil {
		log.Fatal(err)
	}
	p.images[Player_Climb], _, err = ebitenutil.NewImageFromFile("assets/hero/Pink Man/Wall Jump (32x32).png")
	if err != nil {
		log.Fatal(err)
	}

	idleGrid := ganim8.NewGrid(FrameW, FrameH, 640, 192, 0, 0, 0)
	p.anims[Player_Idle] = ganim8.New(p.images[Player_Idle], idleGrid.Frames("1-10", "1-2", "1-4", 3), time.Millisecond*60)

	runGrid := ganim8.NewGrid(FrameW, FrameH, 640, 64, 0, 0, 0)
	p.anims[Player_Run] = ganim8.New(p.images[Player_Run], runGrid.Frames("1-10", 1), time.Millisecond*60)

	jumpGrid := ganim8.NewGrid(FrameW, FrameH, 832, 64, 0, 0, 0)
	p.anims[Player_Jump] = ganim8.New(p.images[Player_Jump], jumpGrid.Frames("1-13", 1), time.Millisecond*60)

	climbGrid := ganim8.NewGrid(FrameW, FrameH, 640, 128, 0, 0, 0)
	p.anims[Player_Climb] = ganim8.New(p.images[Player_Climb], climbGrid.Frames("1-10", 1, "1-2", 2), time.Millisecond*60)

	p.state = Player_Run
	p.grounded = true
	p.curr_anim = p.anims[p.state]

	p.x = 50.0
	p.y = 150.0
	p.velocity = Vec2D[float64]{0., 0.}
	p.dir = Dir_Right

	return p
}

func (p *Player) Move(dx, dy float64) {
	// Get the delta time
	//dt := 1 / ebiten.ActualTPS()
	if dx > 0 {
		p.dir = Dir_Right
	} else if dx < 0 {
		p.dir = Dir_Left
	}
	p.velocity.X = (dx * MoveDx)
	p.velocity.Y = (dy * MoveDx)
}

func (p *Player) CycleAnim() {
	if p.state < Player_Climb {
		p.state += 1
		fmt.Println("cycle anim")
	} else {
		p.state = Player_Idle
	}
}

func (p *Player) Update() error {

	// speed := rotationPerSecond / float64(ebiten.TPS())
	p.x += p.velocity.X
	p.y += p.velocity.Y

	if p.velocity.X > 0 {
		p.velocity.X -= MoveDx
	} else if p.velocity.X < 0 {
		p.velocity.X += MoveDx
	}

	if p.velocity.Y > 0 {
		p.velocity.Y -= MoveDx
	}

	p.curr_anim = p.anims[p.state]
	p.curr_anim.Update()

	return nil
}

func (p *Player) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(p.x, p.y)
	// // The paramters are x, y, rotate (in radian), scaleX, scaleY
	// originX, originY.
	sx := 1.0
	px := p.x
	if p.dir == Dir_Left {
		sx = -1.
		px += 64 // (64 / 2)
	} else if p.dir == Dir_Right {
		sx = 1.
	}
	//opt.GeoM.Translate(float64(-layer.GridSize/2), float64(-layer.GridSize/2))
	p.curr_anim.Draw(screen, ganim8.DrawOpts(px, p.y, 0, sx, 1.0))
}
