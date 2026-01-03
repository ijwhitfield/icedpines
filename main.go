package main

import (
	"fmt"
	"image/color"
	"math"
	"os"
	"path/filepath"
	"slices"
	"unsafe"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Scores struct {
	fastestTime float64
	lowest      float32
	fewestHits  int16
	wins        int16
}

var scores = Scores{}

func loadScores() {
	scoreFilename := resources.dir + "scores"
	file, err := os.Open(scoreFilename)
	if err != nil {
		scores = Scores{lowest: startingHeight, fastestTime: -1, fewestHits: -1, wins: 0}
		return
	}
	countNeeded := unsafe.Sizeof(scores)
	buffer := make([]uint8, countNeeded)
	count, err := file.Read(buffer)
	if uintptr(count) < countNeeded {
		scores = Scores{lowest: startingHeight, fastestTime: -1, fewestHits: -1, wins: 0}
		return
	}
	ptr := (*byte)(unsafe.Pointer(&scores))
	slc := unsafe.Slice(ptr, countNeeded)
	copy(slc, buffer)
}

func saveScores() {
	scoreFilename := resources.dir + "scores"
	file, err := os.Create(scoreFilename)
	if err != nil {
		rl.TraceLog(rl.LogError, "Failed to save scores to file")
		return
	}
	ptr := (*byte)(unsafe.Pointer(&scores))
	countNeeded := unsafe.Sizeof(scores)
	var data []byte = unsafe.Slice(ptr, countNeeded)
	count, err := file.Write(data)
	if uintptr(count) < countNeeded || err != nil {
		rl.TraceLog(rl.LogError, "Failed to save scores to file")
	}
}

type Resources struct {
	dir            string
	font           rl.Font
	fontBig        rl.Font
	background     [1]AnimSource
	bear           [10]AnimSource
	crap           [6]AnimSource
	rock           [3]AnimSource
	health         [3]AnimSource
	icons          [3]AnimSource
	menu           [2]AnimSource
	penguin        [4]AnimSource
	penguinIce     [1]AnimSource
	pole           [1]AnimSource
	snowball       [2]AnimSource
	trap           [2]AnimSource
	trees          [2]AnimSource
	treeIce        [1]AnimSource
	music          rl.Music
	musicMenu      rl.Music
	slideCenter    rl.Music
	slideSide      rl.Music
	boost          rl.Sound
	click          rl.Sound
	iceBreak       rl.Sound
	iced           rl.Sound
	impact         rl.Sound
	item           rl.Sound
	meatBreak      rl.Sound
	meatDead       rl.Sound
	penguinSquawk  rl.Sound
	rockBreak      rl.Sound
	scoop          rl.Sound
	snowballImpact rl.Sound
	snowballReady  rl.Sound
	snowballThrow  rl.Sound
	trapClosing    rl.Sound
	treeBreak      rl.Sound
	win            rl.Sound
}

var resources = Resources{}

const leftAnimIndex int32 = 0
const trapOpenAnimIndex int32 = 0
const centerAnimIndex int32 = 1
const trapClosedAnimIndex int32 = 1
const rightAnimIndex int32 = 2
const shockedAnimIndex int32 = 3
const leftGrabAnimIndex int32 = 3
const leftThrowAnimIndex int32 = 4
const centerGrabAnimIndex int32 = 5
const centerThrowAnimIndex int32 = 6
const rightGrabAnimIndex int32 = 7
const rightThrowAnimIndex int32 = 8
const hurtAnimIndex int32 = 9

func loadResources() {
	resources.dir, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	resources.dir += "/resources/"

	/* FONTS */
	runes := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!.?-:")
	resources.font = rl.LoadFontEx(resources.dir+"Steak Melt.otf", 48, runes, int32(len(runes)))
	resources.fontBig = rl.LoadFontEx(resources.dir+"Steak Melt.otf", 144, runes, int32(len(runes)))

	/* ANIMATIONS */
	resources.background = [1]AnimSource(makeAnimSources([]string{"background.png"}))
	resources.bear = [10]AnimSource(makeAnimSources([]string{
		"bearLeft.png",
		"bearCenter.png",
		"bearRight.png",
		"bearLeftGrab.png",
		"bearLeftThrow.png",
		"bearCenterGrab.png",
		"bearCenterThrow.png",
		"bearRightGrab.png",
		"bearRightThrow.png",
		"bearHurt.png",
	}))
	resources.crap = [6]AnimSource(makeAnimSources([]string{"crap1.png", "crap2.png", "crap3.png", "crap4.png", "crap5.png", "crap6.png"}))
	resources.rock = [3]AnimSource(makeAnimSources([]string{"boulder1.png", "boulder2.png", "boulder3.png"}))
	resources.health = [3]AnimSource(makeAnimSources([]string{"healthDead.png", "healthAlive.png", "healthFrame.png"}))
	resources.icons = [3]AnimSource(makeAnimSources([]string{"iconAltitude.png", "iconSpeed.png", "iconClock.png"}))
	resources.menu = [2]AnimSource(makeAnimSources([]string{"menu1.png", "menu2.png"}))
	resources.penguin = [4]AnimSource(makeAnimSources([]string{"penguinLeft.png", "penguinCenter.png", "penguinRight.png", "penguinShocked.png"}))
	resources.penguinIce = [1]AnimSource(makeAnimSources([]string{"penguinIce.png"}))
	resources.pole = [1]AnimSource(makeAnimSources([]string{"pole.png"}))
	resources.snowball = [2]AnimSource(makeAnimSources([]string{"snowball1.png", "snowball2.png"}))
	resources.trap = [2]AnimSource(makeAnimSources([]string{"trapOpen.png", "trapClosed.png"}))
	resources.trees = [2]AnimSource(makeAnimSources([]string{"tree1.png", "tree2.png"}))
	resources.treeIce = [1]AnimSource(makeAnimSources([]string{"treeIce.png"}))

	/* SOUNDS */
	resources.music = rl.LoadMusicStream(resources.dir + "audio/music.ogg")
	resources.musicMenu = rl.LoadMusicStream(resources.dir + "audio/musicMenu.ogg")
	resources.boost = loadSound("boost.ogg")
	resources.click = loadSound("click.ogg")
	resources.iceBreak = loadSound("iceBreak.ogg")
	resources.iced = loadSound("iced.ogg")
	resources.impact = loadSound("impact.ogg")
	resources.item = loadSound("item.ogg")
	resources.meatBreak = loadSound("meatBreak.ogg")
	resources.meatDead = loadSound("meatDead.ogg")
	resources.penguinSquawk = loadSound("penguinSquawk.ogg")
	resources.rockBreak = loadSound("rockBreak.ogg")
	resources.scoop = loadSound("scoop.ogg")
	resources.slideCenter = rl.LoadMusicStream(resources.dir + "audio/slideCenter.ogg")
	resources.slideSide = rl.LoadMusicStream(resources.dir + "audio/slideSide.ogg")
	resources.snowballImpact = loadSound("snowballImpact.ogg")
	resources.snowballReady = loadSound("snowballReady.ogg")
	resources.snowballThrow = loadSound("snowballThrow.ogg")
	resources.trapClosing = loadSound("trapClosing.ogg")
	resources.treeBreak = loadSound("treeBreak.ogg")
	resources.win = loadSound("win.ogg")

}

func loadSound(filename string) rl.Sound {
	return rl.LoadSound(fmt.Sprint(resources.dir, "audio/", filename))
}

func pauseSounds() {
	rl.PauseSound(resources.boost)
	// HACK never pause click
	// rl.PauseSound(resources.click)
	rl.PauseSound(resources.iceBreak)
	rl.PauseSound(resources.iced)
	rl.PauseSound(resources.impact)
	rl.PauseSound(resources.item)
	rl.PauseSound(resources.meatBreak)
	rl.PauseSound(resources.meatDead)
	rl.PauseSound(resources.penguinSquawk)
	rl.PauseSound(resources.rockBreak)
	rl.PauseSound(resources.scoop)
	rl.PauseMusicStream(resources.slideCenter)
	rl.PauseMusicStream(resources.slideSide)
	rl.PauseSound(resources.snowballImpact)
	rl.PauseSound(resources.snowballReady)
	rl.PauseSound(resources.snowballThrow)
	rl.PauseSound(resources.trapClosing)
	rl.PauseSound(resources.treeBreak)
	rl.PauseSound(resources.win)
}

func resumeSounds() {
	rl.ResumeSound(resources.boost)
	rl.ResumeSound(resources.click)
	rl.ResumeSound(resources.iceBreak)
	rl.ResumeSound(resources.iced)
	rl.ResumeSound(resources.impact)
	rl.ResumeSound(resources.item)
	rl.ResumeSound(resources.meatBreak)
	rl.ResumeSound(resources.meatDead)
	rl.ResumeSound(resources.penguinSquawk)
	rl.ResumeSound(resources.rockBreak)
	rl.ResumeSound(resources.scoop)
	rl.ResumeMusicStream(resources.slideCenter)
	rl.ResumeMusicStream(resources.slideSide)
	rl.ResumeSound(resources.snowballImpact)
	rl.ResumeSound(resources.snowballReady)
	rl.ResumeSound(resources.snowballThrow)
	rl.ResumeSound(resources.trapClosing)
	rl.ResumeSound(resources.treeBreak)
	rl.ResumeSound(resources.win)
}

const windowWidth int32 = 600
const windowHeight int32 = 800

const playerWidth int32 = 50
const playerAcceleration float32 = 200

const entitysMaxCount int32 = 512
const entitysPlayerIndex int32 = 1

/* BEHAVIORS */
const bExists uint64 = 1 << 0
const bIced uint64 = 1 << 2
const bSkier uint64 = 1 << 3
const bCanBeIced uint64 = 1 << 4
const bEarnsPoints uint64 = 1 << 5
const bDynamic uint64 = 1 << 6
const bSolid uint64 = 1 << 7
const bCausesIce uint64 = 1 << 8
const bDropsItem uint64 = 1 << 9
const bExplodesOnDeath uint64 = 1 << 10
const bExplosion uint64 = 1 << 11 // once all dots have expired, destroy this entity
const bLow uint64 = 1 << 12
const bHigh uint64 = 1 << 13
const bSmashEverything uint64 = 1 << 14
const bInvincible uint64 = 1 << 15

type Timer struct {
	time float32
	max  float32
}

func (timer *Timer) reset() {
	timer.time = timer.max
}

type Entity struct {
	x             float32 // side to side position on the mountain, center is zero
	y             float32 // altitude from bottom of mountain, so this goes down with time
	vx, vy        float32 // velocity
	width, height float32 // this is purely visually for now, hitbox defined lower
	rotationSpeed float32
	behavior      uint64
	hp, hpMax     int32
	damage        int32
	wishSpeed     float32
	centerX       float32 // center of where a skier wants to be
	hitbox        rl.Rectangle
	deathSound    rl.Sound
	explosionKind uint32
	attackTimer   Timer
	invulnTimer   Timer
	boostTimer    Timer
	smashTimer    Timer
	shockedTimer  Timer
	snowTimer     Timer
	iceTexture    rl.Texture2D
	anim          AnimState
	dots          []Dot
	flipped       bool
}

type Camera struct {
	x, y           float32 // same coordinate system as entities
	vx, vy         float32
	shakeMagnitude float32
	shakeX, shakeY float32
}

type Input struct {
	move   rl.Vector2
	pause  bool
	action bool
	mute   bool
}

const dotNothing uint32 = 0
const dotSnow uint32 = 1
const dotIce uint32 = 2
const dotBlood uint32 = 3
const dotTree uint32 = 4
const dotRock uint32 = 5
const dotBoost uint32 = 6

type Dot struct {
	x, y, z    float32 // z in this case means up and down. 0 is ground
	vx, vy, vz float32
	expiry     float64
	kind       uint32
}

func createEmptyDot() Dot {
	return Dot{}
}

func (entity *Entity) addTrail(x, y, z, vx, vy, vz float32, now, expiry float64, count uint32, dotKind uint32) {
	for count > 0 {
		for i := range len(entity.dots) {
			if count <= 0 {
				return
			}
			dot := &entity.dots[i]

			if dot.kind == dotNothing || dot.expiry < now {
				// vx := float32(rl.GetRandomValue(-100, 100))
				// vy := float32(rl.GetRandomValue(-100, 100))
				newSnow := Dot{
					x:      x,
					y:      y,
					z:      z,
					vx:     vx,
					vy:     vy,
					vz:     vz,
					expiry: expiry,
					kind:   dotKind,
				}
				*dot = newSnow
				count -= 1
			}
		}
		old := entity.dots
		newLen := len(entity.dots) * 2
		if newLen == 0 {
			newLen = 1
		}
		entity.dots = make([]Dot, newLen)
		copy(entity.dots, old)
	}

}

const itemHealth int32 = 0
const itemBoost int32 = 1

type HealthBar struct {
	fullness       float32
	shakeMagnitude float32
	shakeX, shakeY float32
}

type Game struct {
	playTime          float64
	skierPoints       float32
	treePoints        float32
	rockPoints        float32
	trapPoints        float32
	outerTreePoints   float32
	crapPoints        float32
	furthestY         float32
	lastBarrierY      float32
	deathTimer        Timer
	skierTimer        Timer
	hpShakeTimer      Timer
	boostTimer        Timer
	notificationTimer Timer
	notificationText  string
	hits              int16
	menuSelection     int32
	menuOpen          bool
	quit              bool
	finished          bool
	muted             bool
	musicVolume       float32
	musicMenuVolume   float32
	camera            Camera
	healthBar         HealthBar
	input             Input
	entitys           [entitysMaxCount]Entity
}

func main() {
	game := Game{}
	initGame(&game)
	for !rl.WindowShouldClose() && !game.quit {
		updateDraw(&game)
	}
}

const clippingPlane float32 = 10
const viewDistance float32 = 2500
const hillWidth float32 = 2000
const barrierDistance float32 = 100
const cameraFollowDistance float32 = 150

// const startingHeight float32 = 30000

const startingHeight float32 = 300000

func cameraProjectRectangle(camera Camera, input rl.Rectangle) (rl.Rectangle, bool) {
	output := rl.Rectangle{}
	// determine y distance relative to the camera
	yDiff := (input.Y + input.Height) - camera.y
	yDiff *= -1 // since we're looking downwards, flip the y value

	// if y distance out of range, return zero sized rectangle and a false for visible
	if yDiff > viewDistance || yDiff < clippingPlane {
		return rl.Rectangle{}, false
	}
	// determine x distance relative to the camera
	xDiff := (input.X + input.Width/2) - camera.x
	// determine scale value based on y distance
	scale := cameraFollowDistance / yDiff
	// calculate output
	output.Width = input.Width * scale
	output.Height = input.Height * scale
	output.X = (xDiff*scale - output.Width/2) + float32(windowWidth)/2
	//https://www.desmos.com/calculator/lutldqk9dn
	output.Y = float32(windowHeight) - (500 - (500*105)/(yDiff)) - output.Height

	output.X -= camera.shakeX * scale * camera.shakeMagnitude
	output.Y -= camera.shakeY * scale * camera.shakeMagnitude

	return output, true
}

func cameraProjectDot(camera Camera, y float32, dot Dot) (pos rl.Vector2, size float32) {
	// determine y distance relative to the camera
	yDiff := (y + 10) - camera.y
	yDiff *= -1 // since we're looking downwards, flip the y value

	// determine x distance relative to the camera
	xDiff := (dot.x + 10/2) - camera.x
	// determine scale value based on y distance
	scale := cameraFollowDistance / yDiff
	// calculate output
	size = 10 * scale
	pos.X = (xDiff*scale - size/2) + float32(windowWidth)/2
	//https://www.desmos.com/calculator/lutldqk9dn
	pos.Y = float32(windowHeight) - (500 - (500*105)/(yDiff)) - size - (dot.z * scale)

	pos.X -= camera.shakeX * scale * camera.shakeMagnitude
	pos.Y -= camera.shakeY * scale * camera.shakeMagnitude

	return
}

func aabbCollisionCheck(r1 rl.Rectangle, r2 rl.Rectangle) bool {
	return r1.X+r1.Width > r2.X && r1.X < r2.X+r2.Width && r1.Y+r1.Height > r2.Y && r1.Y < r2.Y+r2.Height
}

func aabbCollision(r1 rl.Rectangle, r2 rl.Rectangle) rl.Rectangle {
	if aabbCollisionCheck(r1, r2) {
		x := max(r1.X-r2.X, r2.X-r1.X)
		y := max(r1.Y-r2.Y, r2.Y-r1.Y)
		return rl.Rectangle{
			X:      x,
			Y:      y,
			Width:  min(r1.X+r1.Width-r2.X, r2.X+r2.Width-r1.X),
			Height: min(r1.Y+r1.Height-r2.Y, r2.Y+r2.Height-r1.Y),
		}
	}
	return rl.Rectangle{}
}

func YSort(a indexYPair, b indexYPair) int {
	if a.y < b.y {
		return -1
	}
	if a.y > b.y {
		return 1
	}
	return 0
}

type indexYPair struct {
	index int32
	y     float32
}

func drawText(str string, x float32, y float32) {
	rl.DrawTextEx(resources.font, str, rl.Vector2{X: x, Y: y}, 24, 2, color.RGBA{0, 0, 0, 255})
}

func drawTextBig(str string, x float32, y float32) {
	rl.DrawTextEx(resources.fontBig, str, rl.Vector2{X: x, Y: y}, 64, 2, color.RGBA{0, 0, 0, 255})
}

func measureText(str string) float32 {
	return rl.MeasureTextEx(resources.font, str, 24, 2).X
}

func measureTextBig(str string) float32 {
	return rl.MeasureTextEx(resources.font, str, 64, 2).X
}

func drawTexture(texture rl.Texture2D, dst rl.Rectangle) {
	rl.DrawTexturePro(texture, rl.Rectangle{X: 0, Y: 0, Width: float32(texture.Width), Height: float32(texture.Height)}, dst, rl.Vector2{X: 0, Y: 0}, 0, rl.White)
}
func drawTextureRotating(texture rl.Texture2D, dst rl.Rectangle, rotation float32) {
	dst.X += dst.Width / 2
	dst.Y += dst.Height / 2
	rl.DrawTexturePro(texture, rl.Rectangle{X: 0, Y: 0, Width: float32(texture.Width), Height: float32(texture.Height)}, dst, rl.Vector2{X: float32(dst.Width / 2), Y: float32(dst.Height / 2)}, rotation, rl.White)
}

func drawTextureFlipped(texture rl.Texture2D, dst rl.Rectangle) {
	rl.DrawTexturePro(texture, rl.Rectangle{X: 0, Y: 0, Width: -float32(texture.Width), Height: float32(texture.Height)}, dst, rl.Vector2{X: 0, Y: 0}, 0, rl.White)
}

type AnimSource struct {
	texture      rl.Texture2D
	width        float32
	height       float32
	timePerFrame float32
	frameCount   int32
}

func makeAnimSources(filenames []string) []AnimSource {
	sources := make([]AnimSource, len(filenames))
	for i, filename := range filenames {
		filename = fmt.Sprint(resources.dir, "sprites/", filename)
		texture := rl.LoadTexture(filename)
		source := AnimSource{
			texture:      texture,
			width:        float32(texture.Width),
			height:       float32(texture.Height),
			timePerFrame: 0,
			frameCount:   1,
		}
		sources[i] = source
	}
	return sources
}

type AnimState struct {
	sources     []AnimSource
	timeStarted float64
	activeIndex int32
}

var colorBlack = color.RGBA{0, 0, 0, 255}
var colorWhite = color.RGBA{255, 255, 255, 255}
var colorLightGrey = color.RGBA{175, 191, 210, 255}
var colorBrown3 = color.RGBA{116, 63, 57, 255}
var colorDarkGrey = color.RGBA{79, 103, 129, 255}
var colorDarkRed = color.RGBA{158, 40, 54, 255}
var colorLightRed = color.RGBA{229, 59, 68, 255}
var colorLightBlue = color.RGBA{44, 231, 244, 255}
var colorYellow = color.RGBA{255, 231, 98, 255}

func draw(game Game) {
	rl.BeginDrawing()

	/* BACKGROUND */
	rl.ClearBackground(colorWhite)
	drawTexture(resources.background[0].texture, rl.Rectangle{0, 0, float32(windowWidth), 400})
	/* ENTITIES */
	indices := [entitysMaxCount]indexYPair{}
	for i := range entitysMaxCount {
		indices[i] = indexYPair{i, game.entitys[i].y}

	}
	slices.SortFunc(indices[:], YSort)
	fogLayer := 0
	fogLayerCount := 20
	fogLayerDepth := float32(50)
	for i := range entitysMaxCount {
		entity := &game.entitys[indices[i].index]
		if fogLayer < fogLayerCount && game.camera.y-entity.y < viewDistance-fogLayerDepth*float32(fogLayer+1) {
			// rl.DrawRectangle(0, 250, windowWidth, windowHeight, color.RGBA{255, 0, 0, 50})
			rl.DrawRectangle(0, 250, windowWidth, windowHeight, color.RGBA{255, 255, 255, 50})
			fogLayer += 1
		}
		if entity.anim.sources != nil {
			preProjection := rl.Rectangle{
				X:      entity.x - entity.width/2,
				Y:      entity.y - entity.height,
				Width:  entity.width,
				Height: entity.height,
			}
			postProjection, visible := cameraProjectRectangle(game.camera, preProjection)
			if visible {
				if entity.invulnTimer.time > 0 && int32(entity.invulnTimer.time*5)%2 == 1 {
					continue
				}
				if entity.anim.activeIndex >= 0 && entity.anim.activeIndex < int32(len(entity.anim.sources)) {
					anim := entity.anim.sources[entity.anim.activeIndex]
					if entity.flipped {
						drawTextureFlipped(anim.texture, postProjection)
					} else if entity.rotationSpeed > 0 {
						drawTextureRotating(anim.texture, postProjection, float32(game.playTime)*entity.rotationSpeed)
					} else {
						drawTexture(anim.texture, postProjection)
					}
				} else {
					rl.DrawRectangleRec(postProjection, color.RGBA{255, 0, 255, 255})
				}
				if entity.hasBehavior(bIced) {
					drawTexture(entity.iceTexture, postProjection)
				}
			}
		}
		/* DOTS */
		if entity.dots != nil {
			for _, dot := range entity.dots {
				color := color.RGBA{}
				switch dot.kind {
				case dotSnow:
					color = colorLightGrey
				case dotBlood:
					color = colorLightRed
				case dotIce:
					color = colorLightBlue
				case dotRock:
					color = colorDarkGrey
				case dotTree:
					color = colorBrown3
				case dotBoost:
					color = colorYellow
				}
				if color.A != 0 {
					pos, size := cameraProjectDot(game.camera, dot.y, dot)
					scale := float32(dot.expiry-game.playTime) * 2
					scale = min(1, scale)
					size *= scale
					rl.DrawCircleV(pos, size/2, color)
				}
			}
		}
	}
	/* UI */
	if game.menuOpen {
		frameDuration := float32(2.0 / 3.0)
		rl.GetMusicTimePlayed(resources.music)
		// rl.DrawRectangleRec(rl.Rectangle{X: 0, Y: float32(windowHeight/2 - 4), Width: 224, Height: 182}, rl.White)
		title := "iced birds"
		measureTextBig(title)
		drawTextBig(title, float32(windowWidth)/2-measureTextBig(title)/2, 120)

		menuY := float32(470)
		snowballWidth := 24

		strNew := "new run"
		strQuit := "quit game"
		strLow := fmt.Sprintf("Lowest: %d m", int32(scores.lowest/100))
		strQuick := fmt.Sprintf("Quickest: %d s", int32(scores.fastestTime))
		strWin := fmt.Sprintf("Wins: %d", scores.wins)
		width := float32(237.5 - 16)
		width = max(width, measureText(strNew)+float32(snowballWidth+8))
		width = max(width, measureText(strQuit)+float32(snowballWidth+8))
		width = max(width, measureText(strLow))
		width = max(width, measureText(strQuick))
		width = max(width, measureText(strWin))
		width += 16
		var height float32
		if scores.wins > 0 {
			height = 150
		} else {
			height = 90
		}
		rl.DrawRectangleRec(rl.Rectangle{X: 12, Y: menuY - 8, Width: width + 8, Height: height + 8}, rl.Black)
		rl.DrawRectangleRec(rl.Rectangle{X: 16, Y: menuY - 4, Width: width, Height: height}, rl.White)

		rl.DrawRectangleRec(rl.Rectangle{X: 20, Y: menuY, Width: width - (float32(snowballWidth) + 8) - 4, Height: 24}, colorLightGrey)
		drawText(strNew, 20+(width-(float32(snowballWidth)+8)-4)/2-measureText(strNew)/2, menuY)
		rl.DrawRectangleRec(rl.Rectangle{X: 20, Y: menuY + 30, Width: width - (float32(snowballWidth) + 8) - 4, Height: 24}, colorLightGrey)
		drawText(strQuit, 20+(width-(float32(snowballWidth)+8)-4)/2-measureText(strQuit)/2, menuY+30)
		/* CURSOR */
		drawTextureRotating(resources.snowball[0].texture, rl.Rectangle{X: 20 + width - (float32(snowballWidth) + 4) - 4, Y: menuY + float32(game.menuSelection)*30, Width: float32(snowballWidth), Height: float32(snowballWidth)}, snowballRotationSpeed*float32(game.playTime))

		/* SCORES */
		// if scores.wins > 0 {
		// 	rl.DrawRectangleRec(rl.Rectangle{X: float32(windowWidth - 400), Y: float32(windowHeight / 2), Width: 400, Height: 200}, rl.White)
		// 	str := fmt.Sprintf("Lowest altitude: %d", int32(scores.lowest/100))
		// 	width := measureText(str)
		// 	drawText(str, float32(windowWidth)-width-24, float32(windowHeight/2))
		// 	str = fmt.Sprintf("Fastest time: %d", int32(scores.fastestTime))
		// 	width = measureText(str)
		// 	drawText(str, float32(windowWidth)-width-24, float32(windowHeight/2+30))
		// 	str = fmt.Sprintf("Fewest hits: %d", scores.fewestHits)
		// 	width = measureText(str)
		// 	drawText(str, float32(windowWidth)-width-24, float32(windowHeight/2+60))
		// 	str = fmt.Sprintf("Total wins: %d", scores.wins)
		// 	width = measureText(str)
		// 	drawText(str, float32(windowWidth)-width-24, float32(windowHeight/2+90))
		// }
		if scores.wins > 0 {
			drawText(strLow, 24, float32(menuY+60))
			drawText(strQuick, 24, float32(menuY+90))
			drawText(strWin, 24, float32(menuY+120))
		} else {
			str := fmt.Sprintf("Lowest: %d m", int32(scores.lowest/100))
			drawText(str, 24, float32(menuY+60))
		}
		frame := int32(rl.GetMusicTimePlayed(resources.music)/frameDuration) % 2
		drawTexture(resources.menu[frame].texture, rl.Rectangle{0, 0, float32(windowWidth), float32(windowHeight)})

	} else {
		player := game.entitys[entitysPlayerIndex]
		rl.DrawRectangleRec(rl.Rectangle{24, 20, 125, 75}, colorBlack)
		rl.DrawRectangleRec(rl.Rectangle{25, 21, 123, 73}, colorWhite)
		rl.DrawRectangleRec(rl.Rectangle{26, 22, 121, 71}, colorBlack)
		rl.DrawRectangleRec(rl.Rectangle{27, 23, 119, 69}, colorWhite)

		drawTexture(resources.icons[0].texture, rl.Rectangle{32, 30, 24, 24})
		drawText(fmt.Sprintf("%d", int32(game.furthestY/100)), 60, 30)
		drawTexture(resources.icons[1].texture, rl.Rectangle{32, 60, 24, 24})
		drawText(fmt.Sprintf("%d", int32(-player.vy)/10), 60, 60)
		timeText := fmt.Sprintf("%d", int32(game.playTime))
		timeX := float32(windowWidth)/2 - measureText(timeText)/2 + 15
		drawTexture(resources.icons[2].texture, rl.Rectangle{timeX - 30, 20, 24, 24})
		drawText(timeText, timeX, 20)

		health := game.healthBar.fullness
		texture := resources.health[0].texture
		x := game.healthBar.shakeX * game.healthBar.shakeMagnitude
		y := game.healthBar.shakeY * game.healthBar.shakeMagnitude
		glowSize := (1-health)*8.0 - 2
		rl.DrawRectangleRec(rl.Rectangle{X: x + float32(windowWidth-125-24) - glowSize, Y: y + 20 - glowSize, Width: 125 + glowSize*2, Height: 75 + glowSize*2}, colorLightRed)
		dst := rl.Rectangle{X: x + float32(windowWidth-24) - 125, Y: y + 20, Width: 125 * (1 - health), Height: 75}
		rl.DrawTexturePro(texture, rl.Rectangle{X: 0, Y: 0, Width: float32(texture.Width) * (1 - health), Height: float32(texture.Height)}, dst, rl.Vector2{X: 0, Y: 0}, 0, rl.White)
		texture = resources.health[1].texture
		dst = rl.Rectangle{X: x + float32(windowWidth-24) - 125*health, Y: y + 20, Width: 125 * health, Height: 75}
		rl.DrawTexturePro(texture, rl.Rectangle{X: float32(texture.Width) * (1 - health), Y: 0, Width: float32(texture.Width) * health, Height: float32(texture.Height)}, dst, rl.Vector2{X: 0, Y: 0}, 0, rl.White)
		drawTexture(resources.health[2].texture, rl.Rectangle{X: x + float32(windowWidth-125-24), Y: y + 20, Width: 125, Height: 75})

		// for i := range player.hp {
		// 	drawTexture(resources.heart[0].texture, rl.Rectangle{X: float32(windowWidth - (i+1)*30), Y: 20, Width: 25, Height: 25})
		// }

		if game.notificationTimer.time > 0 {
			textWidth := rl.MeasureTextEx(resources.font, game.notificationText, 24, 2).X
			drawText(game.notificationText, float32(windowWidth/2)-textWidth/2, float32(windowHeight/2-24))
		}

		/* DEBUG OVERLAY */
		if showOverlay {
			for i := range entitysMaxCount {
				entity := game.entitys[i]
				hitbox := entity.getHitbox()
				hitbox.X += float32(windowWidth)/2 - game.camera.x
				hitbox.Y += float32(windowHeight)/2 - game.camera.y
				rl.DrawRectangleRec(hitbox, color.RGBA{0, 0, 255, 255})
			}

		}
	}

	rl.EndDrawing()
}

const walking bool = false
const showOverlay bool = false

func updateInput(input *Input) {
	input.move = rl.Vector2{}
	if rl.IsKeyDown(rl.KeyUp) || rl.IsKeyDown(rl.KeyW) {
		input.move.Y += -1.0
	}
	if rl.IsKeyDown(rl.KeyDown) || rl.IsKeyDown(rl.KeyS) {
		input.move.Y += 1.0
	}
	if rl.IsKeyDown(rl.KeyLeft) || rl.IsKeyDown(rl.KeyA) {
		input.move.X += -1.0
	}
	if rl.IsKeyDown(rl.KeyRight) || rl.IsKeyDown(rl.KeyD) {
		input.move.X += 1.0
	}
	input.move = rl.Vector2Normalize(input.move)
	input.pause = rl.IsKeyPressed(rl.KeyEscape)
	input.action = rl.IsKeyPressed(rl.KeySpace) || rl.IsKeyPressed(rl.KeyEnter) || rl.IsKeyPressed(rl.KeyZ) || rl.IsKeyPressed(rl.KeyX)
	input.mute = rl.IsKeyPressed(rl.KeyM)
}

func getFirstEmptyEntity(entitys []Entity) *Entity {
	for i := range entitys {
		entity := &entitys[i]
		if entity.behavior == 0 {
			return entity
		}
	}
	return nil
}

func addTree(y float32, entitys []Entity) bool {
	slot := getFirstEmptyEntity(entitys)
	if slot == nil {
		return false
	}
	var newObstacle Entity
	treeIndex := rl.GetRandomValue(0, int32(len(resources.trees)-1))
	flipped := rl.GetRandomValue(0, 1) == 0
	for {
		x := float32(rl.GetRandomValue(-int32(hillWidth)/2, int32(hillWidth)/2))
		yRand := float32(rl.GetRandomValue(-300, 0))
		newObstacle = Entity{
			x:             x,
			y:             y + yRand,
			width:         200,
			height:        400,
			hitbox:        rl.Rectangle{X: -50, Y: -25 / 2, Width: 100, Height: 25},
			behavior:      bExists | bCanBeIced | bSolid | bExplodesOnDeath,
			hp:            1,
			damage:        1,
			explosionKind: dotTree,
			deathSound:    resources.treeBreak,
			iceTexture:    resources.treeIce[0].texture,
			flipped:       flipped,
			anim:          AnimState{sources: resources.trees[:], activeIndex: treeIndex},
		}
		collided := false
		for i := range entitysMaxCount {
			entity := entitys[i]
			box1 := rl.Rectangle{X: entity.x - 10, Y: entity.y - 25, Width: 20, Height: 50}
			box2 := rl.Rectangle{X: newObstacle.x - 10, Y: newObstacle.y - 25, Width: 20, Height: 50}

			if aabbCollisionCheck(box1, box2) {
				collided = true
				break
			}
		}
		if !collided {
			break
		}
	}
	*slot = newObstacle
	return true
}

func addRock(y float32, entitys []Entity) bool {
	slot := getFirstEmptyEntity(entitys)
	if slot == nil {
		return false
	}
	var newObstacle Entity
	rockIndex := rl.GetRandomValue(0, int32(len(resources.rock)-1))
	flipped := rl.GetRandomValue(0, 1) == 0
	for {
		x := float32(rl.GetRandomValue(-int32(hillWidth)/2, int32(hillWidth)/2))
		yRand := float32(rl.GetRandomValue(-300, 0))
		newObstacle = Entity{
			x:             x,
			y:             y + yRand,
			width:         200,
			height:        200,
			hitbox:        rl.Rectangle{X: -80, Y: -25 / 2, Width: 160, Height: 25},
			behavior:      bExists | bSolid | bExplodesOnDeath,
			hp:            1,
			damage:        1,
			explosionKind: dotRock,
			deathSound:    resources.rockBreak,
			flipped:       flipped,
			anim:          AnimState{sources: resources.rock[:], activeIndex: rockIndex},
		}
		collided := false
		for i := range entitysMaxCount {
			entity := entitys[i]
			box1 := rl.Rectangle{X: entity.x - 10, Y: entity.y - 25, Width: 20, Height: 50}
			box2 := rl.Rectangle{X: newObstacle.x - 10, Y: newObstacle.y - 25, Width: 20, Height: 50}

			if aabbCollisionCheck(box1, box2) {
				collided = true
				break
			}
		}
		if !collided {
			break
		}
	}
	*slot = newObstacle
	return true
}

func addTrap(y float32, entitys []Entity) bool {
	slot := getFirstEmptyEntity(entitys)
	if slot == nil {
		return false
	}
	var newObstacle Entity
	flipped := rl.GetRandomValue(0, 1) == 0
	for {
		x := float32(rl.GetRandomValue(-int32(hillWidth)/2, int32(hillWidth)/2))
		yRand := float32(rl.GetRandomValue(-300, 0))
		newObstacle = Entity{
			x:        x,
			y:        y + yRand,
			width:    200,
			height:   150,
			hitbox:   rl.Rectangle{X: -80, Y: -25 / 2, Width: 160, Height: 25},
			behavior: bExists | bLow | bInvincible,
			hp:       100,
			damage:   100,
			flipped:  flipped,
			anim:     AnimState{sources: resources.trap[:], activeIndex: trapOpenAnimIndex},
		}
		collided := false
		for i := range entitysMaxCount {
			entity := entitys[i]
			box1 := rl.Rectangle{X: entity.x - 10, Y: entity.y - 25, Width: 20, Height: 50}
			box2 := rl.Rectangle{X: newObstacle.x - 10, Y: newObstacle.y - 25, Width: 20, Height: 50}

			if aabbCollisionCheck(box1, box2) {
				collided = true
				break
			}
		}
		if !collided {
			break
		}
	}
	*slot = newObstacle
	return true
}

func addOuterTree(y float32, entitys []Entity) bool {
	var x float32
	for {
		x = float32(rl.GetRandomValue(-int32(hillWidth), int32(hillWidth)))
		if !(x >= -float32(hillWidth)/2-50 && x <= float32(hillWidth)/2+50) {
			break
		}
	}
	y += float32(rl.GetRandomValue(-300, 0))
	treeIndex := rl.GetRandomValue(0, int32(len(resources.trees)-1))
	flipped := rl.GetRandomValue(0, 1) == 0
	slot := getFirstEmptyEntity(entitys)
	if slot != nil {
		*slot = Entity{
			x:        x,
			y:        y,
			width:    400,
			height:   800,
			behavior: bExists | bInvincible,
			hp:       100,
			damage:   0,
			flipped:  flipped,
			anim:     AnimState{sources: resources.trees[:], activeIndex: treeIndex},
		}
		return true
	}
	return false
}

func addCrap(y float32, entitys []Entity) bool {
	slot := getFirstEmptyEntity(entitys)
	if slot == nil {
		return false
	}
	var newObstacle Entity
	crapIndex := rl.GetRandomValue(0, int32(len(resources.crap)-1))
	flipped := rl.GetRandomValue(0, 1) == 0
	for {
		x := float32(rl.GetRandomValue(-int32(hillWidth), int32(hillWidth)))
		yRand := float32(rl.GetRandomValue(-300, 0))
		newObstacle = Entity{
			x: x,
			y: y + yRand,
			// width:    150,
			// height:   20,
			width:    100,
			height:   50,
			behavior: bExists | bLow | bInvincible,
			hp:       100,
			damage:   0,
			flipped:  flipped,
			anim:     AnimState{sources: resources.crap[:], activeIndex: crapIndex},
		}
		collided := false
		for i := range entitysMaxCount {
			entity := entitys[i]
			box1 := rl.Rectangle{X: entity.x - 10, Y: entity.y - 25, Width: 20, Height: 50}
			box2 := rl.Rectangle{X: newObstacle.x - 10, Y: newObstacle.y - 25, Width: 20, Height: 50}

			if aabbCollisionCheck(box1, box2) {
				collided = true
				break
			}
		}
		if !collided {
			break
		}
	}
	*slot = newObstacle
	return true
}
func addSkier(y float32, entitys []Entity) bool {
	x := float32(rl.GetRandomValue(-int32(hillWidth)/2+300, int32(hillWidth)/2-300))
	vx := float32(800)
	goLeft := rl.GetRandomValue(0, 1) == 1
	if goLeft {
		vx *= -1
	}
	slot := getFirstEmptyEntity(entitys)
	if slot != nil {
		*slot = Entity{
			x:             x,
			y:             y,
			centerX:       x,
			width:         100,
			height:        80,
			wishSpeed:     500,
			vy:            -500,
			vx:            vx,
			hitbox:        rl.Rectangle{X: -50, Y: -25, Width: 100, Height: 50},
			hp:            1,
			damage:        1,
			deathSound:    resources.meatBreak,
			explosionKind: dotBlood,
			shockedTimer:  Timer{0, 0.25},
			snowTimer:     Timer{0, 0.05},
			iceTexture:    resources.penguinIce[0].texture,
			anim:          AnimState{sources: resources.penguin[:]},
			behavior:      bExists | bSkier | bCanBeIced | bDropsItem | bExplodesOnDeath,
		}
		return true
	}
	return false
}

const snowballSpeed float32 = 1000
const snowballRotationSpeed float32 = 600

func addSnowball(x float32, y float32, vy float32, entitys []Entity) bool {
	slot := getFirstEmptyEntity(entitys)
	if slot != nil {
		ballIndex := rl.GetRandomValue(0, int32(len(resources.snowball)-1))
		*slot = Entity{
			x:             x,
			y:             y,
			vy:            vy,
			width:         50,
			height:        50,
			hitbox:        rl.Rectangle{X: -25, Y: -25, Width: 50, Height: 50},
			hp:            1,
			rotationSpeed: snowballRotationSpeed,
			explosionKind: dotSnow,
			deathSound:    resources.snowballImpact,
			anim:          AnimState{sources: resources.snowball[:], activeIndex: ballIndex},
			behavior:      bExists | bDynamic | bSolid | bCausesIce | bExplodesOnDeath | bHigh,
		}
		return true
	}
	return false
}

func addBarriers(y float32, entitys []Entity) {
	leftDone := false
	for i := range entitys {
		entity := &entitys[i]
		if entity.behavior == 0 {
			if !leftDone {
				*entity = Entity{
					x:        -hillWidth / 2,
					y:        y,
					hp:       100, // TODO get rid of this once i fix the transparency debug thing
					width:    20,
					height:   100,
					behavior: bExists | bInvincible,
					anim:     AnimState{sources: resources.pole[:]},
				}
				leftDone = true
			} else {
				*entity = Entity{
					x:        hillWidth / 2,
					y:        y,
					hp:       100, // TODO get rid of this once i fix the transparency debug thing
					width:    20,
					height:   100,
					behavior: bExists,
					anim:     AnimState{sources: resources.pole[:]},
				}
				break
			}
		}
	}
}

const boostTime float32 = 5
const boostSpeed float32 = 3000

func addPlayer(entitys []Entity) {
	entitys[entitysPlayerIndex] = Entity{
		vy:            -300,
		y:             startingHeight,
		width:         100,
		height:        150,
		hitbox:        rl.Rectangle{X: -float32(playerWidth) / 2, Y: -25 / 2, Width: float32(playerWidth), Height: 25},
		hp:            3,
		hpMax:         3,
		damage:        3,
		deathSound:    resources.meatDead,
		invulnTimer:   Timer{0, 3},
		wishSpeed:     700,
		attackTimer:   Timer{0, 1},
		boostTimer:    Timer{0, boostTime},
		snowTimer:     Timer{0, 0.01},
		smashTimer:    Timer{0, boostTime + 1},
		anim:          AnimState{sources: resources.bear[:]},
		explosionKind: dotBlood,
		behavior:      bExists | bEarnsPoints | bDynamic | bSolid | bExplodesOnDeath,
	}
}

func createEmpty() Entity {
	entity := Entity{}
	return entity
}

func (entity *Entity) addDamage(damage int32) {
	entity.hp -= damage
	if entity.hp > 0 {
		entity.invulnTimer.reset()
	}
	// if entity.hp <= 0 {
	// 	*entity = createEmpty()
	// }
}

func (entity Entity) getHitbox() rl.Rectangle {
	return rl.Rectangle{
		X:      entity.x + entity.hitbox.X,
		Y:      entity.y + entity.hitbox.Y,
		Width:  entity.hitbox.Width,
		Height: entity.hitbox.Height,
	}
}

func abs(a float32) float32 {
	if a < 0 {
		return -a
	}
	return a
}

func (entity Entity) hasBehavior(flags uint64) bool {
	return (^entity.behavior & flags) == 0
}

func tryIce(e1 *Entity, e2 *Entity) {
	if e1.hasBehavior(bCausesIce) && e2.hasBehavior(bCanBeIced) {
		e2.behavior |= bIced
		e1.addDamage(e2.damage)
		rl.PlaySound(resources.iced)
		//*e1 = createEmpty()

	}

}

func tryDamage(e1 *Entity, e2 *Entity) bool {
	miss := (e1.hasBehavior(bHigh) && e2.hasBehavior(bLow)) || (e1.hasBehavior(bLow) && e2.hasBehavior(bHigh))
	if e1.invulnTimer.time <= 0 && e2.damage > 0 && !e2.hasBehavior(bIced) && !miss {
		e1.addDamage(e2.damage)
		e1.wishSpeed *= 0.75
		if &e2.anim.sources[0] == &resources.trap[0] {
			rl.PlaySound(resources.trapClosing)
			e2.anim.activeIndex = trapClosedAnimIndex
		}
		return true
	}
	return false
}

/* returns which item was given */
func (entity *Entity) giveRandomItem() int32 {
	var item int32
	if entity.hp < entity.hpMax {
		item = itemHealth
	} else {
		item = itemBoost
	}
	switch item {
	case itemHealth:
		entity.hp += 1
	case itemBoost:
		entity.vy = -boostSpeed
		entity.vx = 0
		entity.behavior |= bInvincible
		entity.boostTimer.reset()
		entity.behavior |= bSmashEverything
		entity.smashTimer.reset()
		rl.StopSound(resources.scoop)
		entity.attackTimer.time = 0
		entity.anim.activeIndex = centerAnimIndex
	}
	rl.PlaySound(resources.item)
	return item
}

func tryDeath(entity *Entity, vy float32, now float64) {
	if entity.hp <= 0 {
		if entity.hasBehavior(bIced) {
			rl.PlaySound(resources.iceBreak)
		}
		rl.PlaySound(entity.deathSound)
		if entity.hasBehavior(bExplodesOnDeath) {
			entity.vx = 0
			entity.vy = 0
			entity.explode(vy, now)
			entity.anim.sources = nil
			entity.behavior = bExplosion
		}
	}
}

func (entity *Entity) explode(vy float32, now float64) {
	entity.dots = make([]Dot, 30)
	for i := range entity.dots {
		dot := &entity.dots[i]
		kind := entity.explosionKind
		if entity.hasBehavior(bIced) && i%2 == 0 {
			kind = dotIce
		}
		*dot = Dot{
			x:      entity.x,
			y:      entity.y,
			z:      entity.height / 2,
			vx:     float32(rl.GetRandomValue(-800, 800)),
			vy:     vy,
			vz:     float32(rl.GetRandomValue(-800, 800)),
			expiry: now + 1,
			kind:   kind,
		}
	}
}

const skierAcceleration float32 = 1000
const cameraScrollSpeed float32 = 300
const trackYPerfectly bool = true

func (camera *Camera) track(x, y float32, speed float32, frameTime float32) {
	diffx := x - camera.x
	vx := diffx * speed
	if abs(vx*frameTime) > abs(diffx) {
		camera.x = x
	} else {
		camera.x += vx * frameTime
	}
	if trackYPerfectly {
		camera.y = y
	} else {
		diffy := y - camera.y
		vy := diffy * speed
		if abs(vy*frameTime) > abs(diffy) {
			camera.y = y
		} else {
			camera.y += vy * frameTime
		}
	}
}

func update(game *Game) {
	// we always hit 60fps, actually getting frame time only causes crazy stuff to happen on stalls
	// frameTime := rl.GetFrameTime()
	frameTime := float32(1.0 / 60.0)
	player := &game.entitys[entitysPlayerIndex]
	playerMomentum := player.vy

	rl.UpdateMusicStream(resources.music)
	rl.UpdateMusicStream(resources.musicMenu)
	rl.UpdateMusicStream(resources.slideCenter)
	rl.UpdateMusicStream(resources.slideSide)

	updateInput(&game.input)

	if game.input.pause && (player.hp > 0 || game.deathTimer.time > 0) {
		game.menuOpen = !game.menuOpen
		if game.menuOpen {
			pauseSounds()
		} else {
			resumeSounds()
		}
	}
	if game.deathTimer.time <= 0 && player.hp <= 0 {
		game.menuOpen = true
		pauseSounds()
	}
	if game.input.mute {
		if game.muted {
			game.muted = false
		} else {
			game.muted = true
		}
	}

	if !(game.menuOpen && player.hp > 0) {
		game.playTime += float64(frameTime)
		// game.musicVolume = min(1, game.musicVolume+frameTime)
		// game.musicMenuVolume = max(0, game.musicMenuVolume-frameTime)
		game.musicVolume = min(1, 1-(1-game.musicVolume)*0.9)
		game.musicMenuVolume = max(0, game.musicMenuVolume*0.9)
		wasScooping := player.attackTimer.time > 0
		player.attackTimer.time -= frameTime
		if wasScooping && player.attackTimer.time <= 0 {
			rl.PlaySound(resources.snowballReady)
		}
		wasBoosting := player.boostTimer.time > 0
		player.boostTimer.time -= frameTime
		if wasBoosting && player.boostTimer.time <= 0 {
			player.invulnTimer.reset()
			game.skierTimer.reset()
			player.behavior &^= bInvincible
			rl.StopSound(resources.boost)
		}
		wasSmashing := player.smashTimer.time > 0
		player.smashTimer.time -= frameTime
		if wasSmashing && player.smashTimer.time <= 0 {
			player.behavior &^= bSmashEverything
		}
		game.deathTimer.time -= frameTime
		game.skierTimer.time -= frameTime
		game.camera.shakeMagnitude -= frameTime * 200
		game.camera.shakeMagnitude = max(0, game.camera.shakeMagnitude)
		game.healthBar.shakeMagnitude -= frameTime * 200
		game.healthBar.shakeMagnitude = max(0, game.healthBar.shakeMagnitude)
		game.notificationTimer.time -= frameTime
		/* PLAYER */
		if player.hp > 0 {
			if player.boostTimer.time <= 0 {
				if game.input.action && player.attackTimer.time <= 0 {
					if addSnowball(player.x, player.y-50, player.vy-snowballSpeed, game.entitys[:]) {
						player.attackTimer.reset()
						rl.PlaySound(resources.snowballThrow)
					}
				}
				if game.input.move.X > 0 {
					if player.attackTimer.time > player.attackTimer.max*0.8 {
						player.anim.activeIndex = rightThrowAnimIndex
					} else if player.attackTimer.time > 0 {
						player.anim.activeIndex = rightGrabAnimIndex
					} else {
						player.anim.activeIndex = rightAnimIndex
					}
					player.vx = player.wishSpeed * 2
				} else if game.input.move.X < 0 {
					if player.attackTimer.time > player.attackTimer.max*0.8 {
						player.anim.activeIndex = leftThrowAnimIndex
					} else if player.attackTimer.time > 0 {
						player.anim.activeIndex = leftGrabAnimIndex
					} else {
						player.anim.activeIndex = leftAnimIndex
					}
					player.vx = -player.wishSpeed * 2
				} else {
					if player.attackTimer.time > player.attackTimer.max*0.8 {
						player.anim.activeIndex = centerThrowAnimIndex
					} else if player.attackTimer.time > 0 {
						player.anim.activeIndex = centerGrabAnimIndex
					} else {
						player.anim.activeIndex = centerAnimIndex
					}
					player.vx = 0
				}
				if player.vx == 0 {
					if !rl.IsMusicStreamPlaying(resources.slideCenter) {

						rl.PlayMusicStream(resources.slideCenter)
					}
					rl.StopMusicStream(resources.slideSide)
				} else {
					if !rl.IsMusicStreamPlaying(resources.slideSide) {
						rl.PlayMusicStream(resources.slideSide)
					}
					rl.StopMusicStream(resources.slideCenter)
				}
				if player.vy > 0 {
					player.anim.activeIndex = hurtAnimIndex
				}
				if rl.GetRandomValue(0, 700) < int32(abs(player.vy)) {
					if player.snowTimer.time <= 0 {
						x := float32(rl.GetRandomValue(-10, 10))
						y := float32(rl.GetRandomValue(0, -30))
						z := float32(rl.GetRandomValue(0, 0))
						vx := float32(rl.GetRandomValue(-100, 100))
						vz := float32(rl.GetRandomValue(-100, 100))

						player.addTrail(player.x+x, player.y+y, z, vx, 0, 300+vz, game.playTime, game.playTime+.5, 1, dotSnow)
						player.snowTimer.reset()
					}
				}
				if player.boostTimer.time > 0 {
					player.vy = -boostSpeed
					game.camera.shakeMagnitude = 100
				}
				/* SOUND */
				if player.attackTimer.time > player.attackTimer.max*0.8 {
				} else if player.attackTimer.time > 0 {
					if !rl.IsSoundPlaying(resources.scoop) {
						rl.PlaySound(resources.scoop)
					}
				} else {
					if rl.IsSoundPlaying(resources.scoop) {
						rl.StopSound(resources.scoop)
					}
				}

				if walking {
					player.vy = game.input.move.Y * 100
				} else {
					if player.vy > -player.wishSpeed {
						player.vy -= playerAcceleration * frameTime
						player.vy = max(player.vy, -player.wishSpeed)
					} else if player.vy < -player.wishSpeed {
						player.vy += 5 * playerAcceleration * frameTime
						player.vy = min(player.vy, -player.wishSpeed)
					}
					// if player.vy < -player.wishSpeed {
					// 	diff := -player.wishSpeed - player.vy
					// 	if diff < 200*frameTime {
					// 		player.vy = -player.wishSpeed
					// 	} else {
					// 		player.vy += 10 * frameTime
					// 	}
					// }
					// player.vy = max(player.vy, -player.wishSpeed)
					player.x = min(player.x, hillWidth/2-float32(playerWidth))
					player.x = max(player.x, -hillWidth/2+float32(playerWidth))
					player.wishSpeed += max(0, -player.vy*frameTime) / 100
				}
			} else {
				if player.snowTimer.time <= 0 {
					for range 10 {
						x := float32(rl.GetRandomValue(-int32(player.width/2), int32(player.width/2)))
						y := float32(rl.GetRandomValue(0, -30))
						z := float32(rl.GetRandomValue(int32(player.height*0.2), int32(player.height*0.8)))
						vx := float32(rl.GetRandomValue(-100, 100))
						// vy := float32(rl.GetRandomValue(-100, 100))

						player.addTrail(player.x+x, player.y+y, z, vx, -2000, 0, game.playTime, game.playTime+.5, 1, dotBoost)
					}
					player.snowTimer.reset()
				}
			}
		}

		/* BARRIERS */
		if game.camera.y-viewDistance <= game.lastBarrierY-barrierDistance {
			addBarriers(game.lastBarrierY-barrierDistance, game.entitys[:])
			game.lastBarrierY -= barrierDistance
		}

		/* SPAWNING */
		treeDifficulty := float32(0)
		rockDifficulty := float32(0)
		trapDifficulty := float32(0)
		levelDistance := startingHeight / 3
		if player.y > levelDistance*2 {
			treeDifficulty = (levelDistance*3 - player.y) / levelDistance
			rockDifficulty = 0
			trapDifficulty = 0
		} else if player.y > levelDistance {
			rockDifficulty = (levelDistance*2 - player.y) / (levelDistance * 2)
			treeDifficulty = 1 - rockDifficulty
			trapDifficulty = 0
		} else if player.y > 0 {
			rockDifficulty = (levelDistance - player.y) / (levelDistance * 3)
			trapDifficulty = (levelDistance - player.y) / (levelDistance * 3)
			treeDifficulty = 1 - rockDifficulty - trapDifficulty
		} else {
			treeDifficulty = 0
			rockDifficulty = 0
			trapDifficulty = 1
		}
		if treeDifficulty > 0 {
			treeCost := 50 / treeDifficulty
			for game.treePoints > treeCost {
				if addTree(game.camera.y-viewDistance, game.entitys[:]) {
					game.treePoints -= treeCost
				} else {
					break
				}
			}
		}
		if rockDifficulty > 0 {
			rockCost := 50 / rockDifficulty
			for game.rockPoints > rockCost {
				if addRock(game.camera.y-viewDistance, game.entitys[:]) {
					game.rockPoints -= rockCost
				} else {
					break
				}
			}
		}
		if trapDifficulty > 0 {
			trapCost := 100 / trapDifficulty
			for game.trapPoints > trapCost {
				if addTrap(game.camera.y-viewDistance, game.entitys[:]) {
					game.trapPoints -= trapCost
				} else {
					break
				}
			}
		}
		for game.crapPoints > 50 {
			if addCrap(game.camera.y-viewDistance, game.entitys[:]) {
				game.crapPoints -= 50
			} else {
				break
			}
		}
		for game.outerTreePoints > 25 {
			if addOuterTree(game.camera.y-viewDistance, game.entitys[:]) {
				game.outerTreePoints -= 25
			} else {
				break
			}
		}

		skierCount := 0
		for i := range entitysMaxCount {
			entity := &game.entitys[i]
			if entity.hasBehavior(bSkier) {
				skierCount += 1
			}
		}
		if game.skierTimer.time <= 0 && skierCount < 2 && player.boostTimer.time <= 0 {
			addSkier(game.camera.y-viewDistance, game.entitys[:])
			game.skierTimer.reset()
		}

		/* BASIC LOOP */
		for i := range entitysMaxCount {
			entity := &game.entitys[i]

			/* MOVE */
			if !entity.hasBehavior(bIced) {

				if entity.hp <= 0 {
					entity.vy *= 1 - 0.5*frameTime
					entity.vx *= 1 - 0.5*frameTime
				} else if entity.hasBehavior(bSkier) {
					/* X */
					if entity.x < entity.centerX {
						entity.vx += skierAcceleration * frameTime
					} else {
						entity.vx -= skierAcceleration * frameTime
					}
					if entity.vx < -400 {
						entity.anim.activeIndex = leftAnimIndex
					} else if entity.vx < 400 {
						entity.anim.activeIndex = centerAnimIndex
					} else {
						entity.anim.activeIndex = rightAnimIndex
					}
					/* Y */
					if player.boostTimer.time <= 0 {
						entity.y = min(player.y-50, entity.y)
					}
					if entity.y == player.y-50 && abs(player.x-entity.x) < 200 {
						entity.wishSpeed = player.wishSpeed + 400
						rl.PlaySound(resources.penguinSquawk)
						entity.shockedTimer.reset()
					} else {
						entity.wishSpeed -= 25 * frameTime
						entity.wishSpeed = max(500, entity.wishSpeed)
					}
					entity.vy = -entity.wishSpeed
					if entity.shockedTimer.time > 0 {
						entity.anim.activeIndex = shockedAnimIndex
					}
					if entity.snowTimer.time <= 0 {
						x := float32(rl.GetRandomValue(-10, 10))
						y := float32(rl.GetRandomValue(0, -30))
						z := float32(rl.GetRandomValue(0, 0))
						vx := float32(rl.GetRandomValue(-100, 100))
						vz := float32(rl.GetRandomValue(-100, 100))

						entity.addTrail(entity.x+x, entity.y+y, z, vx, 0, 300+vz, game.playTime, game.playTime+2, 1, dotSnow)
						entity.snowTimer.reset()
					}
				}
				entity.y += entity.vy * frameTime
				entity.x += entity.vx * frameTime
			}
			/* DOTS */
			dotsLiving := false
			for i := range len(entity.dots) {
				dot := &entity.dots[i]
				if dot.expiry < game.playTime {
					*dot = createEmptyDot()
				} else {
					dotsLiving = true
					dot.vz -= 1000 * frameTime // gravity
					dot.x += dot.vx * frameTime
					dot.y += dot.vy * frameTime
					dot.z += dot.vz * frameTime
					if dot.z <= 0 {
						dot.z = 0
						dot.vz *= -1
					}
				}
			}

			/* DESPAWN */
			if entity != player {
				if entity.hasBehavior(bExplosion) {
					if !dotsLiving {
						*entity = createEmpty()
					}
				} else {
					if entity.y > game.camera.y+50 || entity.y < game.camera.y-3000 {
						*entity = createEmpty()
					}
				}
			}

			/* TIMERS */
			entity.invulnTimer.time -= frameTime
			entity.shockedTimer.time -= frameTime
			entity.snowTimer.time -= frameTime

		}

		/* COLLISIONS */
		for i1 := range entitysMaxCount {
			e1 := &game.entitys[i1]
			if !e1.hasBehavior(bDynamic|bSolid) || e1.hp <= 0 {
				continue
			}
			for i2 := range entitysMaxCount {
				if i1 == i2 {
					continue
				}
				e2 := &game.entitys[i2]
				if e2.hp <= 0 {
					continue
				}
				collision := aabbCollision(e1.getHitbox(), e2.getHitbox())
				if collision.Width != 0 && collision.Height != 0 {
					if e1 == player && e2.hasBehavior(bSolid) {
						game.camera.shakeMagnitude += 30
					}
					tryIce(e1, e2)
					tryIce(e2, e1)

					if !e1.hasBehavior(bSmashEverything) && e2.hasBehavior(bSolid) && !e2.hasBehavior(bIced) {
						displacement1 := -e1.vy / abs(e1.vy-e2.vy) * collision.Height
						e1.y += displacement1
						displacement2 := -e2.vy / abs(e2.vy-e1.vy) * collision.Height
						e2.y += displacement2
						e1.vy = rl.Clamp(-e1.vy, -100, 100)
						e2.vy = rl.Clamp(-e2.vy, -100, 100)
						if math.IsNaN(float64(e1.y)) || math.IsNaN(float64(e1.x)) {
							panic("NaN position")
						}

					}

					if !e1.hasBehavior(bInvincible) {
						damaged := tryDamage(e1, e2)
						if e1 == player && damaged {
							game.hits += 1
							game.healthBar.shakeMagnitude += 50
							rl.PlaySound(resources.impact)
						}
					}
					if !e2.hasBehavior(bInvincible) {
						damaged := tryDamage(e2, e1)
						if e2 == player && damaged {
							game.hits += 1
							game.healthBar.shakeMagnitude += 50
							rl.PlaySound(resources.impact)
						}
					}
					if (e1 == player || e2 == player) && player.hp <= 0 {
						scores.lowest = min(scores.lowest, player.y)
						rl.StopSound(resources.scoop)
						rl.StopMusicStream(resources.slideCenter)
						rl.StopMusicStream(resources.slideSide)
						saveScores()
						game.deathTimer.reset()
					}
					if (e1.hasBehavior(bDropsItem) && e1.hp <= 0) || (e2.hasBehavior(bDropsItem) && e2.hp <= 0) {
						item := player.giveRandomItem()
						switch item {
						case itemHealth:
							game.healthBar.shakeMagnitude += 50
							which := rl.GetRandomValue(0, 5)
							var txt string
							switch which {
							case 0:
								txt = "DELICIOUS!"
							case 1:
								txt = "DELECTABLE!"
							case 2:
								txt = "SCRUMPTIOUS!"
							case 3:
								txt = "YUMMY!"
							case 4:
								txt = "MMMM!"
							case 5:
								txt = "TASTY!"
							}
							game.notificationText = txt
						case itemBoost:
							rl.PlaySound(resources.boost)
							game.notificationText = "BOOST!"
						}
						game.notificationTimer.reset()
						game.skierTimer.reset()
					}
					var vy float32
					if e2 == player {
						vy = playerMomentum
					} else {
						vy = 0
					}
					tryDeath(e1, vy, game.playTime)
					if e1 == player {
						vy = playerMomentum

					} else {
						vy = 0
					}
					tryDeath(e2, vy, game.playTime)
					if !e1.hasBehavior(bDynamic|bSolid) || e1.hp <= 0 {
						break
					}
				}
			}
		}

	} else {
		game.musicMenuVolume = min(1, 1-(1-game.musicMenuVolume)*0.9)
		game.musicVolume = max(0, game.musicVolume*0.9)
	}
	if game.muted {
		rl.SetMusicVolume(resources.music, 0)
		rl.SetMusicVolume(resources.musicMenu, 0)
	} else {
		rl.SetMusicVolume(resources.musicMenu, game.musicMenuVolume)
		rl.SetMusicVolume(resources.music, game.musicVolume)
	}
	{
		health := float32(max(0, player.hp)) / float32(player.hpMax)
		speed := float32(3)
		diff := health - game.healthBar.fullness
		hpChange := diff * speed
		if abs(hpChange*frameTime) > abs(diff) {
			game.healthBar.fullness = health
		} else {
			game.healthBar.fullness += hpChange * frameTime
		}
	}
	/* CAMERA */
	if player.hp > 0 || (player.hp <= 0 && game.deathTimer.time > 0) {
		game.camera.track(player.x, player.y+cameraFollowDistance, 40, frameTime)
	} else if player.hp <= 0 && game.deathTimer.time <= 0 {
		game.camera.track(0, game.camera.y-cameraScrollSpeed*frameTime, 5, frameTime)
	}
	if int32(game.playTime*200)%2 == 0 {
		game.camera.shakeX = float32(rl.GetRandomValue(-int32(50), int32(50))) / 100
		game.camera.shakeY = float32(rl.GetRandomValue(-int32(50), int32(50))) / 100
		game.healthBar.shakeX = float32(rl.GetRandomValue(-int32(50), int32(50))) / 100
		game.healthBar.shakeY = float32(rl.GetRandomValue(-int32(50), int32(50))) / 100
	}
	/* SPAWNING POINTS */
	if game.camera.y < game.furthestY {
		pointsAdded := game.furthestY - game.camera.y
		game.furthestY = game.camera.y
		game.skierPoints += pointsAdded
		game.outerTreePoints += pointsAdded
		game.crapPoints += pointsAdded
		game.treePoints += pointsAdded
		game.trapPoints += pointsAdded
		game.rockPoints += pointsAdded
	}
	/* WIN IF WINNING */
	if player.y <= 0 && !game.finished {
		game.finished = true
		rl.PlaySound(resources.win)
		game.notificationText = "FINISHED!\nNow playing endless mode..."
		game.notificationTimer.reset()
		scores.wins += 1
		if scores.fastestTime <= -1 {
			scores.fastestTime = game.playTime
		} else {
			scores.fastestTime = min(scores.fastestTime, game.playTime)
		}
		if scores.fewestHits <= -1 {
			scores.fewestHits = game.hits
		} else {
			scores.fewestHits = min(scores.fewestHits, game.hits)
		}
		saveScores()
	}
	/* MENU INPUT */
	if game.menuOpen {
		if game.input.move.Y > 0 {
			prev := game.menuSelection
			game.menuSelection = min(1, game.menuSelection+1)
			if game.menuSelection != prev {
				rl.PlaySound(resources.click)
			}
		} else if game.input.move.Y < 0 {
			prev := game.menuSelection
			game.menuSelection = max(0, game.menuSelection-1)
			if game.menuSelection != prev {
				rl.PlaySound(resources.click)
			}
		}
		if game.input.action {
			switch game.menuSelection {
			case 0:
				reset(game)
			case 1:
				game.quit = true
			}
			rl.PlaySound(resources.click)
		}
	}
}

func updateDraw(game *Game) {
	update(game)
	draw(*game)
}

func reset(game *Game) {
	muted := game.muted
	*game = Game{}
	// game.playTime = rl.GetTime()
	addPlayer(game.entitys[:])
	player := &game.entitys[entitysPlayerIndex]
	game.camera.x = player.x
	game.camera.y = player.y + cameraFollowDistance
	game.deathTimer.max = 3
	game.skierTimer.max = 5
	game.hpShakeTimer.max = 2
	game.boostTimer.max = boostTime
	game.notificationTimer.max = 2
	game.healthBar.fullness = 1
	game.furthestY = startingHeight
	game.muted = muted
	if !game.muted {
		game.musicVolume = 1
	}

	y := player.y
	for ; y > player.y-viewDistance; y -= barrierDistance {
		addBarriers(y, game.entitys[:])
	}
	game.lastBarrierY = y + barrierDistance
}

func initGame(game *Game) {
	reset(game)
	player := &game.entitys[entitysPlayerIndex]
	player.hp = 0
	game.camera.y -= cameraFollowDistance
	game.menuOpen = true

	rl.InitWindow(windowWidth, windowHeight, "iced birds")
	rl.InitAudioDevice()
	rl.SetExitKey(0)
	rl.SetTargetFPS(60)
	loadResources()
	loadScores()

	rl.PlayMusicStream(resources.music)
	rl.SetMusicVolume(resources.musicMenu, 0)
	rl.PlayMusicStream(resources.musicMenu)
}
