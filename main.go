package main

import (
	"context"
	"image/color"
	"math/rand/v2"
	"time"

	llq "github.com/emirpasic/gods/queues/linkedlistqueue"
	pq "github.com/emirpasic/gods/queues/priorityqueue"
	"github.com/emirpasic/gods/utils"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Direction byte

const (
	Left = iota
	Right
	Up
	Down
)

type Cell struct {
	x      uint32
	y      uint32
	isWall bool
}

func NewCell(x, y uint32, isWall bool) *Cell {
	c := &Cell{
		x, y,
		isWall,
	}
	return c
}

var maze [][]*Cell
var visited [][]bool

const n = 20
const m = 20
const scale = 30
const framerate = 20
const frametime = time.Second / framerate

type Edge struct {
	a      *Cell
	b      *Cell
	weight uint32
}

func NewEdge(a, b *Cell, weight uint32) Edge {
	return Edge{a, b, weight}
}

func EdgeLess(a, b interface{}) int {
	ac := a.(Edge)
	bc := b.(Edge)
	return utils.UInt32Comparator(ac.weight, bc.weight)
}

func genMaze(ctx context.Context) error {

	current := maze[0][0]
	queue := pq.NewWith(EdgeLess)

	enc_adj := func(cell *Cell) {
		if !cell.isWall {
			return
		}
		cell.isWall = false
		if cell.x > 0 {
			queue.Enqueue(NewEdge(cell, maze[cell.y][cell.x-1], rand.Uint32()))
		}
		if cell.y > 0 {
			queue.Enqueue(NewEdge(cell, maze[cell.y-1][cell.x], rand.Uint32()))
		}
		if cell.x < m-1 {
			queue.Enqueue(NewEdge(cell, maze[cell.y][cell.x+1], rand.Uint32()))
		}
		if cell.y < n-1 {
			queue.Enqueue(NewEdge(cell, maze[cell.y+1][cell.x], rand.Uint32()))
		}
	}

	count_neighbors := func(cell *Cell) (count uint8) {
		if cell.x > 0 {
			if !maze[cell.y][cell.x-1].isWall {
				count++
			}
		}
		if cell.y > 0 {
			if !maze[cell.y-1][cell.x].isWall {
				count++
			}
		}
		if cell.x < m-1 {
			if !maze[cell.y][cell.x+1].isWall {
				count++
			}
		}
		if cell.y < n-1 {
			if !maze[cell.y+1][cell.x].isWall {
				count++
			}
		}
		/*
			if cell.x > 0 && cell.y > 0 {
				if !maze[cell.y-1][cell.x-1].isWall {
					count++
				}
			}
			if cell.y > 0 && cell.x < m-1 {
				if !maze[cell.y-1][cell.x+1].isWall {
					count++
				}
			}
			if cell.y < n-1 && cell.x > 0 {
				if !maze[cell.y+1][cell.x-1].isWall {
					count++
				}
			}
			if cell.y < n-1 && cell.x < m-1 {
				if !maze[cell.y+1][cell.x+1].isWall {
					count++
				}
			}*/
		return count
	}

	enc_adj(current)

	var (
		val any
		v   Edge
		ok  bool
	)

	for !queue.Empty() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if val, ok = queue.Dequeue(); !ok {
			continue
		}
		if v, ok = val.(Edge); !ok {
			continue
		}

		if (v.a.isWall != v.b.isWall) && (count_neighbors(v.b) < 2) {
			enc_adj(v.b)
			time.Sleep(frametime)
		}
	}
	return nil
}

var cb *Cell

func bfs(ctx context.Context, start *Cell) error {
	queue := llq.New()
	enc_adj := func(cell *Cell) {
		if cell.x > 0 {
			if !maze[cell.y][cell.x-1].isWall && !visited[cell.y][cell.x-1] {
				queue.Enqueue(maze[cell.y][cell.x-1])
			}
		}
		if cell.y > 0 {
			if !maze[cell.y-1][cell.x].isWall && !visited[cell.y-1][cell.x] {
				queue.Enqueue(maze[cell.y-1][cell.x])
			}
		}
		if cell.x < m-1 {
			if !maze[cell.y][cell.x+1].isWall && !visited[cell.y][cell.x+1] {
				queue.Enqueue(maze[cell.y][cell.x+1])
			}
		}
		if cell.y < n-1 {
			if !maze[cell.y+1][cell.x].isWall && !visited[cell.y+1][cell.x] {
				queue.Enqueue(maze[cell.y+1][cell.x])
			}
		}
	}

	var (
		current *Cell
		v       interface{}
		ok      bool
	)
	queue.Enqueue(start)

	for !queue.Empty() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		if v, ok = queue.Dequeue(); !ok {
			break
		}
		current = v.(*Cell)
		cb = v.(*Cell)
		visited[current.y][current.x] = true
		enc_adj(current)
		time.Sleep(frametime)
	}
	return nil
}

type Palette struct {
	black    color.RGBA
	darkBlue color.RGBA
	blue     color.RGBA
	white    color.RGBA
}

func NewPalette() Palette {
	return Palette{
		black:    rl.GetColor(0x000000ff),
		darkBlue: rl.GetColor(0x013d58ff),
		blue:     rl.GetColor(0x53dfffff),
		white:    rl.GetColor(0xdaf6ffff),
	}
}

var palette Palette

func Draw() {
	rl.BeginDrawing()
	var c color.RGBA
	for y := uint32(0); y < n; y++ {
		for x := uint32(0); x < m; x++ {
			if maze[y][x].isWall {
				c = palette.black
			} else {
				c = palette.darkBlue
			}
			rl.DrawRectangle(
				int32(x*scale),
				int32(y*scale),
				scale,
				scale,
				c,
			)
		}
	}
	for y := uint32(0); y < n; y++ {
		for x := uint32(0); x < m; x++ {
			if visited[y][x] {
				rl.DrawRectangle(
					int32(x*scale),
					int32(y*scale),
					scale,
					scale,
					palette.blue,
				)

			}
		}
	}
	if cb != nil {
		rl.DrawRectangle(
			int32(cb.x*scale),
			int32(cb.y*scale),
			scale,
			scale,
			palette.white,
		)
	}
	rl.EndDrawing()
}

func initMaze() {
	maze = make([][]*Cell, n)
	for y := uint32(0); y < n; y++ {
		maze[y] = make([]*Cell, m)
		for x := uint32(0); x < m; x++ {
			maze[y][x] = NewCell(x, y, true)
		}
	}
}
func initVisited() {
	cb = nil
	visited = make([][]bool, n)
	for y := uint32(0); y < n; y++ {
		visited[y] = make([]bool, m)
		for x := uint32(0); x < m; x++ {
			visited[y][x] = false
		}
	}
}

func main() {
	var (
		ctx    context.Context
		cancel context.CancelFunc
	)
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	palette = NewPalette()
	rl.InitWindow(m*scale, n*scale, "main")
	defer rl.CloseWindow()

	rl.SetTargetFPS(framerate)

	startRun := func() {
		if cancel != nil {
			cancel()
		}

		initMaze()
		initVisited()

		ctx, cancel = context.WithCancel(context.Background())

		startCell := maze[0][0]

		go func(c context.Context, start *Cell) {
			if err := genMaze(c); err != nil {
				println("genMaze finished or cancelled:", err.Error())
				return
			}

			select {
			case <-c.Done():
				println("cancelled before bfs")
				return
			default:
			}

			if err := bfs(c, start); err != nil {
				println("bfs finished or cancelled:", err.Error())
				return
			}
		}(ctx, startCell)
	}

	startRun()

	for !rl.WindowShouldClose() {
		if rl.IsKeyPressed(rl.KeyR) {
			startRun()
		}
		Draw()
	}
}
