package main

import (
	"image/color"
	"slices"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const windowWidth int32 = 600
const windowHeight int32 = 800

const playerWidth int32 = 50
const playerSpeed float32 = 300

const entitysMaxCount int32 = 512
const entitysPlayerIndex int32 = 1

/* BEHAVIORS */
const bExists uint64 = 1 << 0
const bIced uint64 = 1 << 2
const bBooster uint64 = 1 << 3
const bCanBeIced uint64 = 1 << 4

type Entity struct {
	/* PHYSICS */
	x             float32 // side to side position on the mountain, center is zero
	y             float32 // altitude from bottom of mountain, so this goes down with time
	vx, vy        float32 // velocity
	width, height float32 // this is purely visually for now, hitbox defined lower

	/* GAMEPLAY */
	behavior     uint64
	hp           int32
	damage       int32
	hitboxWidth  float32
	hitboxHeight float32

	/* ANIMATIONS */
	color color.RGBA
	// these are unused until we have art in the game
	animIndex     uint32
	animStartTime float64
}

type Camera struct {
	x, y float32 // same coordinate system as entities
}

type Input struct {
	move  rl.Vector2
	reset bool
}

type Game struct {
	playTime       float64
	playerPoints   float32
	obstaclePoints float32
	lastBarrierY   float32
	camera         Camera
	input          Input
	entitys        [entitysMaxCount]Entity
}

func main() {
	game := Game{}
	initGame(&game)
	for !rl.WindowShouldClose() {
		updateDraw(&game)
	}
}

const clippingPlane float32 = 10
const viewDistance float32 = 2000
const hillWidth float32 = 2000
const barrierDistance float32 = 100
const cameraFollowDistance float32 = 150

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
	return r1.X > r2.X && r1.X < r2.X+r2.Width && r1.Y > r2.Y && r1.Y < r2.Y+r2.Height
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
		if entity.animIndex != 0 {
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
				rl.DrawRectangleRec(postProjection, entity.color)
			}
		}
	}
	/* UI */

	rl.EndDrawing()
}

func updateInput(input *Input) {
	input.move = rl.Vector2{}
	// if rl.IsKeyDown(rl.KeyUp) {
	// 	input.move.Y += 1.0
	// }
	// if rl.IsKeyDown(rl.KeyDown) {
	// 	input.move.Y += -1.0
	// }
	if rl.IsKeyDown(rl.KeyLeft) {
		input.move.X += -1.0
	}
	if rl.IsKeyDown(rl.KeyRight) {
		input.move.X += 1.0
	}
	input.move = rl.Vector2Normalize(input.move)
	input.reset = rl.IsKeyPressed(rl.KeyR)
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

func createObstacle(y float32) Entity {
	x := float32(rl.GetRandomValue(-int32(hillWidth)/2, int32(hillWidth)/2))
	y += float32(rl.GetRandomValue(-300, 0))
	entity := Entity{
		x:            x,
		y:            y,
		width:        100,
		height:       300,
		hitboxWidth:  float32(100),
		hitboxHeight: float32(50),
		behavior:     bExists | bCanBeIced,
		hp:           1,
		damage:       1,
		animIndex:    2,
		color:        color.RGBA{0, 255, 0, 255},
	}
	return entity
}

func createBarrier(x float32, y float32) Entity {
	entity := Entity{
		x:            x,
		y:            y,
		width:        50,
		height:       50,
		hitboxWidth:  float32(50),
		hitboxHeight: float32(50),
		behavior:     bExists,
		animIndex:    3,
		color:        color.RGBA{0, 0, 0, 255},
	}
	return entity
}

func createEmpty() Entity {
	entity := Entity{}
	return entity
}

func doDamage(entity *Entity, damage int32) {
	entity.hp -= damage
	// if entity.hp <= 0 {
	// 	*entity = createEmpty()
	// }
}

func getHitbox(entity Entity) rl.Rectangle {
	return rl.Rectangle{
		X:      entity.x - entity.hitboxWidth/2,
		Y:      entity.y - entity.hitboxHeight,
		Width:  entity.hitboxWidth,
		Height: entity.hitboxHeight,
	}
}

func update(game *Game) {
	/* TIME */
	frameTime := rl.GetFrameTime()
	game.playTime += float64(rl.GetFrameTime())

	/* INPUTS */
	updateInput(&game.input)
	if game.input.reset {
		reset(game)
	}

	/* PLAYER MOVEMENT */
	player := &game.entitys[entitysPlayerIndex]
	if player.hp > 0 {
		player.vx = game.input.move.X * playerSpeed * 2
		player.x = min(player.x, hillWidth/2)
		player.x = max(player.x, -hillWidth/2)
		//player.vy = game.input.move.Y * playerSpeed
		player.vy = -playerSpeed
	} else {
		player.vy *= 1 - 0.5*frameTime
		player.vx *= 1 - 0.5*frameTime
	}

	/* BARRIERS */
	if player.y-viewDistance <= game.lastBarrierY-barrierDistance {
		slot := getFirstEmptyEntity(game.entitys[:])
		*slot = createBarrier(-hillWidth/2, game.lastBarrierY-barrierDistance)
		slot = getFirstEmptyEntity(game.entitys[:])
		*slot = createBarrier(hillWidth/2, game.lastBarrierY-barrierDistance)
		game.lastBarrierY -= barrierDistance
	}

	/* SPAWNING */
	game.obstaclePoints += frameTime * 500
	if game.obstaclePoints > 100 {
		slot := getFirstEmptyEntity(game.entitys[:entitysMaxCount])
		if slot != nil {
			*slot = createObstacle(player.y - viewDistance)
			game.obstaclePoints -= 100
		}
	}

	/* MOVE AND DESTROY OLD ENTITIES */
	for i := range entitysMaxCount {
		entity := &game.entitys[i]

		/* MOVE */
		entity.y += entity.vy * frameTime
		entity.x += entity.vx * frameTime

		/* DESPAWN */
		if entity.y > game.camera.y+50 {
			*entity = createEmpty()
		}
	}

	/* COLLISIONS */
	for i1 := range entitysMaxCount {
		e1 := &game.entitys[i1]
		for i2 := i1 + 1; i2 < entitysMaxCount; i2++ {
			e2 := &game.entitys[i2]
			if aabbCollisionCheck(getHitbox(*e1), getHitbox(*e2)) {
				doDamage(e1, e2.damage)
				doDamage(e2, e1.damage)
			}

		}
	}

	/* CAMERA */
	game.camera.x = game.entitys[entitysPlayerIndex].x
	game.camera.y = game.entitys[entitysPlayerIndex].y + cameraFollowDistance
}

func updateDraw(game *Game) {
	update(game)
	draw(*game)
}

func createPlayer() Entity {
	return Entity{
		width:        float32(playerWidth),
		height:       float32(playerWidth),
		hitboxWidth:  float32(playerWidth),
		hitboxHeight: float32(playerWidth),
		hp:           3,
		vy:           -100,
		color:        color.RGBA{255, 0, 0, 255},
		animIndex:    1,
		behavior:     bExists,
	}
}

func reset(game *Game) {
	*game = Game{}
	// game.playTime = rl.GetTime()
	player := &game.entitys[entitysPlayerIndex]
	*player = createPlayer()
	game.camera.x = player.x
	game.camera.y = player.y + cameraFollowDistance

	y := player.y
	for ; y > player.y-viewDistance; y -= barrierDistance {
		slot := getFirstEmptyEntity(game.entitys[:])
		*slot = createBarrier(-hillWidth/2, y)
		slot = getFirstEmptyEntity(game.entitys[:])
		*slot = createBarrier(hillWidth/2, y)
	}
	game.lastBarrierY = y + barrierDistance
}

func initGame(game *Game) {
	reset(game)

	rl.InitWindow(windowWidth, windowHeight, "iced pines")
	rl.SetTargetFPS(60)
}
