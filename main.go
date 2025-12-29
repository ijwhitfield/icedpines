package main

import (
	"fmt"
	"image/color"
	"math"
	"slices"

	rl "github.com/gen2brain/raylib-go/raylib"
)

/*

	only one scene
	there's a menu layer that's either open or closed
	the title is displayed above
	there's two, new run, quit game

	if menu is open and player is alive, pause time and don't simulate
	if menu is open and player is dead, slowly roll camera forward
	menu opens 3 seconds after player is dead
	while menu is open, UI is not shown

*/

type Resources struct {
	font     rl.Font
	penguin  [3]AnimSource
	bear     [3]AnimSource
	pole     [1]AnimSource
	trees    [2]AnimSource
	snowball [1]AnimSource
	heart    [1]AnimSource
}

var resources = Resources{}

const leftAnimIndex int32 = 0
const centerAnimIndex int32 = 1
const rightAnimIndex int32 = 2
const deadAnimIndex int32 = 3

func loadResources() {
	runes := []rune("-abcdefghijklmnopqrstuvwxyz:0123456789")
	resources.font = rl.LoadFontEx("resources/Snowstorm Black.otf", 48, runes, int32(len(runes)))
	resources.penguin = [3]AnimSource(makeAnimSources([]string{"resources/penguinLeft.png", "resources/penguinCenter.png", "resources/penguinRight.png"}))
	resources.bear = [3]AnimSource(makeAnimSources([]string{"resources/bearLeft.png", "resources/bearCenter.png", "resources/bearRight.png"}))
	resources.pole = [1]AnimSource(makeAnimSources([]string{"resources/pole.png"}))
	resources.trees = [2]AnimSource(makeAnimSources([]string{"resources/tree1.png", "resources/tree2.png"}))
	resources.snowball = [1]AnimSource(makeAnimSources([]string{"resources/snowball.png"}))
	resources.heart = [1]AnimSource(makeAnimSources([]string{"resources/heart.png"}))
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

type Timer struct {
	time float32
	max  float32
}

func (timer *Timer) reset() {
	timer.time = timer.max
}

type Entity struct {
	/* PHYSICS */
	x             float32 // side to side position on the mountain, center is zero
	y             float32 // altitude from bottom of mountain, so this goes down with time
	vx, vy        float32 // velocity
	width, height float32 // this is purely visually for now, hitbox defined lower

	/* GAMEPLAY */
	behavior    uint64
	hp          int32
	damage      int32
	wishSpeed   float32
	centerX     float32 // center of where a skier wants to be
	hitbox      rl.Rectangle
	attackTimer Timer
	invulnTimer Timer

	/* ANIMATIONS */
	anim AnimState
}

type Camera struct {
	x, y float32 // same coordinate system as entities
}

type Input struct {
	move     rl.Vector2
	pause    bool
	snowball bool
	item     bool
}

type Game struct {
	playTime             float64
	skierPoints          float32
	obstaclePoints       float32
	outerObstaclesPoints float32
	furthestY            float32
	lastBarrierY         float32
	deathTimer           Timer
	skierTimer           Timer
	menuSelection        int32
	menuOpen             bool
	quit                 bool
	camera               Camera
	input                Input
	entitys              [entitysMaxCount]Entity
}

func main() {
	game := Game{}
	initGame(&game)
	for !rl.WindowShouldClose() && !game.quit {
		updateDraw(&game)
	}
}

const clippingPlane float32 = 10
const viewDistance float32 = 2000
const hillWidth float32 = 2000
const barrierDistance float32 = 100
const cameraFollowDistance float32 = 150
const startingHeight float32 = 100000

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

	return output, true
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

func drawTexture(texture rl.Texture2D, dst rl.Rectangle) {
	rl.DrawTexturePro(texture, rl.Rectangle{X: 0, Y: 0, Width: float32(texture.Width), Height: float32(texture.Height)}, dst, rl.Vector2{X: 0, Y: 0}, 0, rl.White)
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

func draw(game Game) {
	rl.BeginDrawing()

	/* BACKGROUND */
	rl.ClearBackground(rl.RayWhite)
	/* ENTITIES */
	indices := [entitysMaxCount]indexYPair{}
	for i := range entitysMaxCount {
		indices[i] = indexYPair{i, game.entitys[i].y}

	}
	slices.SortFunc(indices[:], YSort)
	for i := range entitysMaxCount {
		entity := &game.entitys[indices[i].index]
		if entity.anim.sources != nil {
			preProjection := rl.Rectangle{
				X:      entity.x - entity.width/2,
				Y:      entity.y - entity.height,
				Width:  entity.width,
				Height: entity.height,
			}
			postProjection, visible := cameraProjectRectangle(game.camera, preProjection)
			if visible {
				// texture := animGetTexture(entity.animIndex, entity.animStartTime)
				// rl.DrawTexturePro(texture, src, postProjection, origin, 0, color.White.RGBA())
				if entity.invulnTimer.time > 0 && int32(entity.invulnTimer.time*5)%2 == 0 {
					continue
				}
				if entity.anim.activeIndex >= 0 && entity.anim.activeIndex < int32(len(entity.anim.sources)) {
					anim := entity.anim.sources[entity.anim.activeIndex]
					drawTexture(anim.texture, postProjection)
				} else {
					rl.DrawRectangleRec(postProjection, color.RGBA{255, 0, 255, 255})
				}
				if entity.hasBehavior(bIced) {
					rl.DrawRectangleRec(postProjection, color.RGBA{0, 255, 255, 100})
				}
				// rl.DrawRectangleRec(postProjection, color)
			}
		}
	}
	/* UI */
	if game.menuOpen {
		drawText("iced birds", 24, float32(windowHeight/2))

		rl.DrawRectangleRec(rl.Rectangle{X: 20, Y: float32(windowHeight/2 + 30), Width: 200, Height: 24}, color.RGBA{175, 175, 175, 255})
		drawText("new run", 24, float32(windowHeight/2+30))
		rl.DrawRectangleRec(rl.Rectangle{X: 20, Y: float32(windowHeight/2 + 60), Width: 200, Height: 24}, color.RGBA{175, 175, 175, 255})
		drawText("quit game", 24, float32(windowHeight/2+60))

		drawTexture(resources.snowball[0].texture, rl.Rectangle{X: 240, Y: float32(windowHeight/2 + (game.menuSelection+1)*30), Width: 24, Height: 24})
	} else {
		player := game.entitys[entitysPlayerIndex]
		drawText(fmt.Sprintf("altitude: %d", int32(game.furthestY/100)), 24, 20)
		drawText(fmt.Sprintf("speed: %d", int32(-player.vy)), 24, 50)

		if player.y <= 0 {
			drawText("you win", 24, float32(windowHeight/2))
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
const sidewaysCollision bool = false

func updateInput(input *Input) {
	input.move = rl.Vector2{}
	if rl.IsKeyDown(rl.KeyUp) {
		input.move.Y += -1.0
	}
	if rl.IsKeyDown(rl.KeyDown) {
		input.move.Y += 1.0
	}
	if rl.IsKeyDown(rl.KeyLeft) {
		input.move.X += -1.0
	}
	if rl.IsKeyDown(rl.KeyRight) {
		input.move.X += 1.0
	}
	input.move = rl.Vector2Normalize(input.move)
	input.pause = rl.IsKeyPressed(rl.KeyEscape)
	input.snowball = rl.IsKeyPressed(rl.KeyX)
	input.item = rl.IsKeyPressed(rl.KeyZ)
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

func addObstacle(y float32, entitys []Entity) bool {
	x := float32(rl.GetRandomValue(-int32(hillWidth)/2, int32(hillWidth)/2))
	y += float32(rl.GetRandomValue(-300, 0))
	treeIndex := rl.GetRandomValue(0, int32(len(resources.trees)-1))
	slot := getFirstEmptyEntity(entitys)
	if slot != nil {
		*slot = Entity{
			x:        x,
			y:        y,
			width:    200,
			height:   400,
			hitbox:   rl.Rectangle{X: -50, Y: -25 / 2, Width: 100, Height: 25},
			behavior: bExists | bCanBeIced | bSolid,
			hp:       1,
			damage:   1,
			anim:     AnimState{sources: resources.trees[:], activeIndex: treeIndex},
		}
		return true
	}
	return false
}

func addOuterObstacle(y float32, entitys []Entity) bool {
	x := float32(rl.GetRandomValue(-int32(hillWidth), int32(hillWidth)))
	for x >= -float32(hillWidth)/2 && x <= float32(hillWidth)/2 {
		x = float32(rl.GetRandomValue(-int32(hillWidth), int32(hillWidth)))
	}
	y += float32(rl.GetRandomValue(-300, 0))
	treeIndex := rl.GetRandomValue(0, int32(len(resources.trees)-1))
	slot := getFirstEmptyEntity(entitys)
	if slot != nil {
		*slot = Entity{
			x:        x,
			y:        y,
			width:    400,
			height:   800,
			behavior: bExists,
			hp:       1,
			damage:   1,
			anim:     AnimState{sources: resources.trees[:], activeIndex: treeIndex},
		}
		return true
	}
	return false
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
			x:        x,
			y:        y,
			centerX:  x,
			width:    80,
			height:   60,
			vy:       -500,
			vx:       vx,
			hitbox:   rl.Rectangle{X: -25, Y: -25, Width: 50, Height: 50},
			hp:       1,
			anim:     AnimState{sources: resources.penguin[:]},
			behavior: bExists | bSkier | bCanBeIced,
		}
		return true
	}
	return false
}

const snowballSpeed float32 = 1000

func addSnowball(x float32, y float32, vy float32, entitys []Entity) bool {
	slot := getFirstEmptyEntity(entitys)
	if slot != nil {
		*slot = Entity{
			x:        x,
			y:        y,
			vy:       vy,
			width:    50,
			height:   50,
			hitbox:   rl.Rectangle{X: -25, Y: -25, Width: 50, Height: 50},
			hp:       1,
			anim:     AnimState{sources: resources.snowball[:]},
			behavior: bExists | bDynamic | bSolid | bCausesIce,
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
					hp:       1, // TODO get rid of this once i fix the transparency debug thing
					width:    20,
					height:   100,
					behavior: bExists,
					anim:     AnimState{sources: resources.pole[:]},
				}
				leftDone = true
			} else {
				*entity = Entity{
					x:        hillWidth / 2,
					y:        y,
					hp:       1, // TODO get rid of this once i fix the transparency debug thing
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

func addPlayer(entitys []Entity) {
	entitys[entitysPlayerIndex] = Entity{
		y:           startingHeight,
		width:       80,
		height:      150,
		hitbox:      rl.Rectangle{X: -float32(playerWidth) / 2, Y: -25 / 2, Width: float32(playerWidth), Height: 25},
		hp:          3,
		damage:      3,
		invulnTimer: Timer{0, 3},
		wishSpeed:   500,
		attackTimer: Timer{0, 1},
		anim:        AnimState{sources: resources.bear[:]},
		behavior:    bExists | bEarnsPoints | bDynamic | bSolid,
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
	} else {
		entity.anim.activeIndex = deadAnimIndex
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

const skierAcceleration float32 = 1000
const cameraScrollSpeed float32 = 300

func update(game *Game) {
	frameTime := rl.GetFrameTime()
	player := &game.entitys[entitysPlayerIndex]
	updateInput(&game.input)
	player.attackTimer.time -= frameTime
	game.deathTimer.time -= frameTime
	game.skierTimer.time -= frameTime

	if game.input.pause && (player.hp > 0 || game.deathTimer.time > 0) {
		game.menuOpen = !game.menuOpen
	}
	if game.deathTimer.time <= 0 && player.hp <= 0 {
		game.menuOpen = true
	}

	if !(game.menuOpen && player.hp > 0) {
		game.playTime += float64(frameTime)

		/* PLAYER */
		if player.hp > 0 {
			if game.input.snowball && player.attackTimer.time <= 0 {
				if addSnowball(player.x, player.y-50, player.vy-snowballSpeed, game.entitys[:]) {
					player.attackTimer.reset()
				}
			}
			if game.input.move.X > 0 {
				player.anim.activeIndex = rightAnimIndex

				player.vx = player.wishSpeed * 2
			} else if game.input.move.X < 0 {
				player.anim.activeIndex = leftAnimIndex
				player.vx = -player.wishSpeed * 2
			} else {
				player.anim.activeIndex = centerAnimIndex
				player.vx = 0
			}
			if walking {
				player.vy = game.input.move.Y * 100
			} else {
				player.vy -= playerAcceleration * frameTime
				player.vy = max(player.vy, -player.wishSpeed)
			}
		}

		/* BARRIERS */
		if game.camera.y-viewDistance <= game.lastBarrierY-barrierDistance {
			addBarriers(game.lastBarrierY-barrierDistance, game.entitys[:])
			game.lastBarrierY -= barrierDistance
		}

		/* SPAWNING */
		for game.obstaclePoints > 50 {
			if addObstacle(game.camera.y-viewDistance, game.entitys[:]) {
				game.obstaclePoints -= 50
			}
		}
		for game.outerObstaclesPoints > 25 {
			if addOuterObstacle(game.camera.y-viewDistance, game.entitys[:]) {
				game.outerObstaclesPoints -= 25
			}
		}

		skierCount := 0
		for i := range entitysMaxCount {
			entity := &game.entitys[i]
			if entity.hasBehavior(bSkier) {
				skierCount += 1
			}
		}
		if game.skierTimer.time <= 0 && skierCount == 0 {
			if addSkier(game.camera.y-viewDistance, game.entitys[:]) {
				game.skierTimer.reset()
			}
		}

		/* BASIC LOOP */
		for i := range entitysMaxCount {
			entity := &game.entitys[i]

			/* MOVE */
			if !entity.hasBehavior(bIced) {
				entity.y += entity.vy * frameTime
				entity.x += entity.vx * frameTime
				if entity.behavior&bEarnsPoints != 0 {
					entity.x = min(entity.x, hillWidth/2-float32(playerWidth))
					entity.x = max(entity.x, -hillWidth/2+float32(playerWidth))
					entity.wishSpeed += max(0, -entity.vy*frameTime) / 100
				}
				if entity.hp <= 0 {
					entity.vy *= 1 - 0.5*frameTime
					entity.vx *= 1 - 0.5*frameTime
				} else if entity.hasBehavior(bSkier) {
					if entity.x < entity.centerX {
						entity.vx += skierAcceleration * frameTime
					} else {
						entity.vx -= skierAcceleration * frameTime
					}
					if entity.vx < -100 {
						entity.anim.activeIndex = leftAnimIndex
					} else if entity.vx < 100 {
						entity.anim.activeIndex = centerAnimIndex
					} else {
						entity.anim.activeIndex = rightAnimIndex
					}
					entity.y = min(player.y-150, entity.y)
				}
			}

			/* DESPAWN */
			if (entity.y > game.camera.y+50 || entity.y < game.camera.y-3000) && i != entitysPlayerIndex {
				*entity = createEmpty()
			}

			entity.invulnTimer.time -= frameTime

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
					if e1.hasBehavior(bCausesIce) && e2.hasBehavior(bCanBeIced) {
						e2.behavior |= bIced
						*e1 = createEmpty()
					}
					if e2.hasBehavior(bCausesIce) && e1.hasBehavior(bCanBeIced) {
						e1.behavior |= bIced
						*e2 = createEmpty()
					}
					if e2.hasBehavior(bSolid) && !e2.hasBehavior(bIced) {
						if sidewaysCollision {
							timeX := collision.Width / abs(e1.vx-e2.vx)
							timeY := collision.Height / abs(e1.vy-e2.vy)
							if timeX < timeY {
								displacement1 := -e1.vx / abs(e1.vx-e2.vx) * collision.Width
								e1.x += displacement1
								displacement2 := -e2.vx / abs(e2.vx-e1.vx) * collision.Width
								e2.x += displacement2
								e1.vx = rl.Clamp(-e1.vx, -100, 100)
								e2.vx = rl.Clamp(-e2.vx, -100, 100)
							} else {
								displacement1 := -e1.vy / abs(e1.vy-e2.vy) * collision.Height
								e1.y += displacement1
								displacement2 := -e2.vy / abs(e2.vy-e1.vy) * collision.Height
								e2.y += displacement2
								e1.vy = rl.Clamp(-e1.vy, -100, 100)
								e2.vy = rl.Clamp(-e2.vy, -100, 100)
							}
						}
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

					if e1.invulnTimer.time <= 0 && e2.damage > 0 && !e2.hasBehavior(bIced) {
						e1.addDamage(e2.damage)
						if e1 == player {
							e1.wishSpeed *= 0.75
							if player.hp <= 0 {
								game.deathTimer.reset()
							}
						}
					}
					if e2.invulnTimer.time <= 0 && e1.damage > 0 && !e1.hasBehavior(bIced) {
						e2.addDamage(e1.damage)
						if e2 == player {
							e2.wishSpeed *= 0.75
							if player.hp <= 0 {
								game.deathTimer.reset()
							}
						}
					}

				}

			}
		}
	}
	if player.hp > 0 || (player.hp <= 0 && game.deathTimer.time > 0) {
		game.camera.x = game.entitys[entitysPlayerIndex].x
		game.camera.y = game.entitys[entitysPlayerIndex].y + cameraFollowDistance
	} else if player.hp <= 0 && game.deathTimer.time <= 0 {
		game.camera.x = 0
		game.camera.y += -cameraScrollSpeed * frameTime
	}
	if game.camera.y < game.furthestY {
		pointsAdded := game.furthestY - game.camera.y
		game.furthestY = game.camera.y
		game.skierPoints += pointsAdded
		game.outerObstaclesPoints += pointsAdded
		game.obstaclePoints += pointsAdded

	}
	if game.menuOpen {
		if game.input.move.Y > 0 {
			game.menuSelection = min(1, game.menuSelection+1)
		} else if game.input.move.Y < 0 {
			game.menuSelection = max(0, game.menuSelection-1)
		}
		if game.input.snowball {
			switch game.menuSelection {
			case 0:
				reset(game)
			case 1:
				game.quit = true
			}
		}
	}
}

func updateDraw(game *Game) {
	update(game)
	draw(*game)
}

func reset(game *Game) {
	*game = Game{}
	// game.playTime = rl.GetTime()
	addPlayer(game.entitys[:])
	player := &game.entitys[entitysPlayerIndex]
	game.camera.x = player.x
	game.camera.y = player.y + cameraFollowDistance
	game.deathTimer.max = 3
	game.skierTimer.max = 2
	game.furthestY = startingHeight

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
	rl.SetExitKey(0)
	rl.SetTargetFPS(60)
	loadResources()
}
