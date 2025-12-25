package main

import (
	"fmt"
	"image/color"
	"slices"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const windowWidth int32 = 600
const windowHeight int32 = 800

const playerWidth int32 = 50
const playerAcceleration float32 = 200

const entitysMaxCount int32 = 512
const entitysPlayerIndex int32 = 1

/* BEHAVIORS */
const bExists uint64 = 1 << 0
const bIced uint64 = 1 << 2
const bBooster uint64 = 1 << 3
const bCanBeIced uint64 = 1 << 4
const bEarnsPoints uint64 = 1 << 5

type Entity struct {
	/* PHYSICS */
	x             float32 // side to side position on the mountain, center is zero
	y             float32 // altitude from bottom of mountain, so this goes down with time
	vx, vy        float32 // velocity
	width, height float32 // this is purely visually for now, hitbox defined lower

	/* GAMEPLAY */
	behavior      uint64
	hp            int32
	damage        int32
	wishSpeed     float32
	hitbox        rl.Rectangle
	invulnTime    float32
	invulnTimeMax float32

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
	furthestY      float32
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
				if entity.invulnTime > 0 && int32(entity.invulnTime*5)%2 == 0 {
					continue
				}
				color := entity.color
				if entity.hp <= 0 {
					color.A /= 2
				}
				rl.DrawRectangleRec(postProjection, color)
			}
		}
	}

	/* UI */
	player := game.entitys[entitysPlayerIndex]
	rl.DrawText(fmt.Sprintf("Altitude: %.0f", game.furthestY/100), 20, 20, 24, color.RGBA{0, 0, 0, 255})
	rl.DrawText(fmt.Sprintf("Speed: %.0f", -player.vy), 20, 50, 24, color.RGBA{0, 0, 0, 255})
	if player.y <= 0 {
		rl.DrawText("YOU WIN", 20, windowHeight/2, 24, color.RGBA{0, 0, 0, 255})
	}
	// rl.DrawText("Points: 0", 20, 20, 24, color.RGBA{0, 0, 0, 255})

	/* DEBUG OVERLAY */
	if showOverlay {
		for i := range entitysMaxCount {
			entity := game.entitys[i]
			hitbox := getHitbox(entity)
			hitbox.X += float32(windowWidth)/2 - game.camera.x
			hitbox.Y += float32(windowHeight)/2 - game.camera.y
			rl.DrawRectangleRec(hitbox, color.RGBA{0, 0, 255, 255})
		}
	}

	rl.EndDrawing()
}

const walking bool = false
const showOverlay bool = false

func updateInput(input *Input) {
	input.move = rl.Vector2{}
	if walking {
		if rl.IsKeyDown(rl.KeyUp) {
			input.move.Y += -1.0
		}
		if rl.IsKeyDown(rl.KeyDown) {
			input.move.Y += 1.0
		}
	}
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
		x:         x,
		y:         y,
		width:     100,
		height:    300,
		hitbox:    rl.Rectangle{-50, -25 / 2, 100, 25},
		behavior:  bExists | bCanBeIced,
		hp:        1,
		damage:    1,
		animIndex: 2,
		color:     color.RGBA{0, 255, 0, 255},
	}
	return entity
}

func createBarrier(x float32, y float32) Entity {
	entity := Entity{
		x:         x,
		y:         y,
		hp:        1, // TODO get rid of this once i fix the transparency debug thing
		width:     50,
		height:    50,
		behavior:  bExists,
		animIndex: 3,
		color:     color.RGBA{0, 0, 0, 255},
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
		X:      entity.x + entity.hitbox.X,
		Y:      entity.y + entity.hitbox.Y,
		Width:  entity.hitbox.Width,
		Height: entity.hitbox.Height,
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
		player.vx = game.input.move.X * player.wishSpeed * 2
		if walking {
			player.vy = game.input.move.Y * 100
		} else {
			player.vy -= playerAcceleration * frameTime
			player.vy = max(player.vy, -player.wishSpeed)
		}
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
	// TODO this only spawns one thing per frame, even if multiple are possible
	if game.obstaclePoints > 50 {
		slot := getFirstEmptyEntity(game.entitys[:entitysMaxCount])
		if slot != nil {
			*slot = createObstacle(player.y - viewDistance)
			game.obstaclePoints -= 50
		}
	}

	/* BASIC LOOP */
	for i := range entitysMaxCount {
		entity := &game.entitys[i]

		/* MOVE */
		entity.y += entity.vy * frameTime
		entity.x += entity.vx * frameTime
		if entity.behavior&bEarnsPoints != 0 {

			entity.x = min(entity.x, hillWidth/2-float32(playerWidth))
			entity.x = max(entity.x, -hillWidth/2+float32(playerWidth))
			if entity.y < game.furthestY {
				pointsAdded := game.furthestY - entity.y
				game.furthestY = entity.y
				game.playerPoints += pointsAdded
				game.obstaclePoints += pointsAdded
				entity.wishSpeed += max(0, -entity.vy*frameTime) / 100
			}

		}

		/* DESPAWN */
		if entity.y > game.camera.y+50 {
			*entity = createEmpty()
		}

		entity.invulnTime -= frameTime

	}

	/* COLLISIONS */
	for i1 := range entitysMaxCount {
		e1 := &game.entitys[i1]
		if e1.hp <= 0 {
			continue
		}
		for i2 := i1 + 1; i2 < entitysMaxCount; i2++ {
			e2 := &game.entitys[i2]
			if e2.hp <= 0 {
				continue
			}
			// TODO resolve position as well, using the velocities to determine who moves and how far
			/*
				if e1 is dynamic, iterate over *all* other entities and resolve any collisions
			*/
			if aabbCollisionCheck(getHitbox(*e1), getHitbox(*e2)) {
				if e1.invulnTime <= 0 {
					doDamage(e1, e2.damage)
					e1.invulnTime = e1.invulnTimeMax
				}
				if e2.invulnTime <= 0 {
					doDamage(e2, e1.damage)
					e2.invulnTime = e2.invulnTimeMax
				}
				e1.vx = rl.Clamp(-e1.vx, -100, 100)
				e1.vy = rl.Clamp(-e1.vy, -100, 100)
				e2.vx = rl.Clamp(-e2.vx, -100, 100)
				e2.vy = rl.Clamp(-e2.vy, -100, 100)
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
		y:             startingHeight,
		width:         float32(playerWidth),
		height:        float32(playerWidth),
		hitbox:        rl.Rectangle{-float32(playerWidth) / 2, -25 / 2, float32(playerWidth), 25},
		hp:            3,
		damage:        3,
		invulnTimeMax: 3,
		wishSpeed:     500,
		color:         color.RGBA{255, 0, 0, 255},
		animIndex:     1,
		behavior:      bExists | bEarnsPoints,
	}
}

func reset(game *Game) {
	*game = Game{}
	// game.playTime = rl.GetTime()
	player := &game.entitys[entitysPlayerIndex]
	*player = createPlayer()
	game.camera.x = player.x
	game.camera.y = player.y + cameraFollowDistance

	game.furthestY = startingHeight

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
