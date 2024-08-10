package main

import (
	"fmt"

	ray "github.com/gen2brain/raylib-go/raylib"
)

type Point struct {
	x int32
	y int32
}

type Cell struct {
	alive     bool
	neighbors int32
	nextState bool
}

const cellSize = 24

var (
	direction     []Point
	grid          map[Point]Cell
	lastPos       Point
	gametick      int32   = 20
	zoomHoldSpeed float64 = 30
	tickHoldSpeed float64 = 60
	windowWidth   int32   = 1280
	windowHeight  int32   = 720
	middleVec     ray.Vector2
)

func init() {
	grid = make(map[Point]Cell)
	direction = []Point{
		{x: 1, y: 1},
		{x: 1, y: 0},
		{x: 1, y: -1},

		{x: 0, y: 1},
		{x: 0, y: -1},

		{x: -1, y: 1},
		{x: -1, y: 0},
		{x: -1, y: -1},
	}
	middleVec = ray.Vector2{X: float32(windowWidth / 2), Y: float32(windowHeight / 2)}
}

func inputCells(camera ray.Camera2D) {
	if ray.IsMouseButtonDown(ray.MouseButtonLeft) {
		mousePos := ray.GetScreenToWorld2D(ray.GetMousePosition(), camera)
		pos := Point{
			x: int32(mousePos.X) / cellSize,
			y: int32(mousePos.Y) / cellSize,
		}

		if ray.IsMouseButtonPressed(ray.MouseButtonLeft) || (pos != lastPos) {
			lastPos = pos

			if cell, ok := grid[pos]; ok {
				cell.alive = true
				grid[pos] = cell
				checkNeighbor(pos, cell.alive)
			} else {
				cell = Cell{true, 0, false}
				grid[pos] = cell
				checkNeighbor(pos, cell.alive)
			}
		}

	}

	if ray.IsMouseButtonDown(ray.MouseRightButton) {
		mousePos := ray.GetScreenToWorld2D(ray.GetMousePosition(), camera)
		pos := Point{
			x: int32(mousePos.X) / cellSize,
			y: int32(mousePos.Y) / cellSize,
		}

		if cell, ok := grid[pos]; ok && cell.alive {
			cell.alive = false
			grid[pos] = cell

			checkNeighbor(pos, false)
			for point, cell := range grid {
				if !cell.alive && cell.neighbors == 0 {
					delete(grid, point)
				}
			}
		}
	}
}

func drawBoard() {
	for pos, cell := range grid {
		color := ray.DarkGreen
		if !cell.alive {
			color = ray.LightGray
		}
		ray.DrawRectangle(pos.x*cellSize, pos.y*cellSize, cellSize-1, cellSize-1, color)
	}
}

func checkNeighbor(point Point, state bool) {
	change := int32(1)
	if !state {
		change = -1
	}

	var pos Point
	for _, dir := range direction {
		pos = point
		pos.x += dir.x
		pos.y += dir.y

		if cell, ok := grid[pos]; ok {
			cell.neighbors += change
			grid[pos] = cell
		} else {
			cell = Cell{false, 1, false}
			grid[pos] = cell
		}
	}
}

func tickStep() {
	// calculate next state
	for point, cell := range grid {
		if cell.alive {
			cell.nextState = cell.neighbors >= 2 && cell.neighbors <= 3
			grid[point] = cell
		} else if !cell.alive {
			cell.nextState = cell.neighbors == 3
			grid[point] = cell
		}
	}

	// change all the nextState
	for point, cell := range grid {
		if cell.nextState != cell.alive {
			cell.alive = cell.nextState
			grid[point] = cell
			checkNeighbor(point, cell.alive)
		}
	}

	// delete the died one
	for point, cell := range grid {
		if !cell.alive && cell.neighbors == 0 {
			delete(grid, point)
		}
	}
}

func main() {
	ray.InitWindow(windowWidth, windowHeight, "Go Of Life")
	defer ray.CloseWindow()

	camera := ray.Camera2D{
		Offset:   ray.Vector2Zero(),
		Zoom:     1,
		Rotation: 0,
	}

	autoTickStep := false
	tick := false

	var fps int32 = 144
	ray.SetTargetFPS(fps)

	lastTime := ray.GetTime()
	var lastKey int
	var lastKeyTime float64
	var lastKeyTimeHold float64

	for !ray.WindowShouldClose() {
		inputCells(camera)

		gameTime := ray.GetTime()
		if ray.IsKeyPressed(ray.KeyS) {
			tick = true
		}
		if ray.IsKeyDown(ray.KeyEnter) {
			tick = true
		}

		if ray.IsKeyDown(ray.KeyR) {
			grid = make(map[Point]Cell)
		}

		// Game tick
		if ray.IsKeyPressed(ray.KeyLeftBracket) || (ray.IsKeyDown(ray.KeyLeftBracket) && lastKey != 0 && gameTime-lastKeyTime > 0.5) {
			if lastKey != ray.KeyLeftBracket {
				lastKey = ray.KeyLeftBracket
				lastKeyTime = gameTime
			}

			if gameTime-lastKeyTimeHold > 1/tickHoldSpeed {
				gametick -= 5
				if gametick < 5 {
					gametick = 5
				}
				lastKeyTimeHold = gameTime
			}
		}

		if ray.IsKeyPressed(ray.KeyRightBracket) || (ray.IsKeyDown(ray.KeyRightBracket) && lastKey != 0 && gameTime-lastKeyTime > 0.5) {
			if lastKey != ray.KeyRightBracket {
				lastKey = ray.KeyRightBracket
				lastKeyTime = gameTime
			}

			if gameTime-lastKeyTimeHold > 1/tickHoldSpeed {
				gametick += 5
				lastKeyTimeHold = gameTime
			}
		}

		// Zoom
		if ray.IsKeyPressed(ray.KeyMinus) || (ray.IsKeyDown(ray.KeyMinus) && lastKey != 0 && gameTime-lastKeyTime > 0.5) {
			if lastKey != ray.KeyMinus {
				lastKey = ray.KeyMinus
				lastKeyTime = gameTime
			}

			if gameTime-lastKeyTimeHold > 1/zoomHoldSpeed {
				camera.Target = ray.GetScreenToWorld2D(middleVec, camera)
				camera.Offset = middleVec
				camera.Zoom *= 0.75

				lastKeyTimeHold = gameTime
			}
		}

		if ray.IsKeyPressed(ray.KeyEqual) || (ray.IsKeyDown(ray.KeyEqual) && lastKey != 0 && gameTime-lastKeyTime > 0.5) {
			if lastKey != ray.KeyEqual {
				lastKey = ray.KeyEqual
				lastKeyTime = gameTime
			}

			if gameTime-lastKeyTimeHold > 1/zoomHoldSpeed {
				camera.Target = ray.GetScreenToWorld2D(middleVec, camera)
				camera.Offset = middleVec
				camera.Zoom /= 0.75

				lastKeyTimeHold = gameTime
			}

		}

		if ray.IsKeyReleased(ray.KeyRightBracket) || ray.IsKeyReleased(ray.KeyLeftBracket) || ray.IsKeyReleased(ray.KeyMinus) || ray.IsKeyReleased(ray.KeyEqual) {
			lastKey = 0
		}

		if ray.IsKeyPressed(ray.KeySpace) {
			autoTickStep = !autoTickStep
		}

		if wheelMove := ray.GetMouseWheelMoveV(); wheelMove != ray.Vector2Zero() {
			camera.Offset.X += wheelMove.X * cellSize
			camera.Offset.Y += wheelMove.Y * cellSize
		}

		ray.BeginDrawing()

		if gameTime-lastTime > 1/float64(gametick) {
			lastTime = gameTime
			if tick || autoTickStep {
				tickStep()
				tick = false
			}

		}

		ray.BeginMode2D(camera)
		ray.ClearBackground(ray.RayWhite)

		drawBoard()
		ray.EndMode2D()
		ray.DrawText(fmt.Sprint("Game Tick: ", gametick), 0, 0, 32, ray.Black)
		if !autoTickStep {
			ray.DrawText(fmt.Sprint("Paused"), 0, 33, 32, ray.Black)
		}

		ray.EndDrawing()

	}
}
