package main

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

type Point2D struct {
	X, Y int
}

type GridCoord struct {
	X, Y int
}

func main() {
	//Dec9Part1()
	Dec9Part2()
}

func Dec9Part1() {
	//read the file
	f := parseFile("dec9Input.txt")

	//parse the file contents

	largest := 0
	for i := 0; i < len(f)-1; i++ {
		row1 := f[i]
		for j := i + 1; j < len(f); j++ {
			row2 := f[j]
			dx := absInt(row1.X - row2.X)
			dy := absInt(row1.Y - row2.Y)
			area := (dx + 1) * (dy + 1)
			largest = max(largest, area)

		}
	}

	fmt.Println(largest)
	//iterate over double loop, abs value after - then *

}

func absInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func Dec9Part2() {
	points := parseFile("dec9Input.txt")
	if len(points) < 2 {
		fmt.Println(0)
		return
	}

	// Coordinate Compression and Boundary Marking
	allowed, xMap, yMap, gridPointsCompressed := buildAllowedGridCompressed(points)

	H_prime := len(allowed)
	W_prime := len(allowed[0])

	// Extract sorted coordinate lists
	sortedX := make([]int, W_prime)
	for x, i := range xMap {
		sortedX[i] = x
	}
	sortedY := make([]int, H_prime)
	for y, i := range yMap {
		sortedY[i] = y
	}

	pArea := make([][]int, H_prime+1)
	for y := range pArea {
		pArea[y] = make([]int, W_prime+1)
	}

	for cy := range H_prime {
		dY := int(1)
		if cy < H_prime-1 {
			dY = sortedY[cy+1] - sortedY[cy]
		}

		rowSum := int(0)
		for cx := range W_prime {
			// Physical width (X-coordinate difference)
			dX := int(1)
			if cx < W_prime-1 {
				dX = sortedX[cx+1] - sortedX[cx]
			}

			cellArea := int(0)
			if allowed[cy][cx] {
				cellArea = dX * dY
			}

			rowSum += cellArea
			pArea[cy+1][cx+1] = rowSum + pArea[cy][cx+1]
		}
	}

	rectPhysicalArea := func(cx1, cy1, cx2, cy2 int) int {
		return pArea[cy2+1][cx2+1] -
			pArea[cy1][cx2+1] -
			pArea[cy2+1][cx1] +
			pArea[cy1][cx1]
	}

	largestArea := int(0)
	n := len(points)

	for i := range n {
		for j := i + 1; j < n; j++ {
			p1_orig := points[i]
			p2_orig := points[j]
			p1_comp := gridPointsCompressed[i]
			p2_comp := gridPointsCompressed[j]

			cx1, cx2 := p1_comp.X, p2_comp.X
			cy1, cy2 := p1_comp.Y, p2_comp.Y
			if cx1 > cx2 {
				cx1, cx2 = cx2, cx1
			}
			if cy1 > cy2 {
				cy1, cy2 = cy2, cy1
			}

			x1_orig, x2_orig := p1_orig.X, p2_orig.X
			y1_orig, y2_orig := p1_orig.Y, p2_orig.Y
			if x1_orig > x2_orig {
				x1_orig, x2_orig = x2_orig, x1_orig
			}
			if y1_orig > y2_orig {
				y1_orig, y2_orig = y2_orig, y1_orig
			}

			currentArea := (x2_orig - x1_orig + 1) * (y2_orig - y1_orig + 1)

			if currentArea <= largestArea {
				continue
			}

			// match the total physical area of the rectangle?
			allowedArea := rectPhysicalArea(cx1, cy1, cx2, cy2)

			if allowedArea == currentArea {
				largestArea = currentArea
			}
		}
	}

	fmt.Println(largestArea)
}

func collectCompressedCoords(points []Point2D) ([]int, []int) {
	uniqueX := make(map[int]struct{})
	uniqueY := make(map[int]struct{})

	for _, p := range points {
		uniqueX[p.X] = struct{}{}
		uniqueY[p.Y] = struct{}{}
		// Include coordinates around the point to represent the "gaps" for flood-fill
		uniqueX[p.X-1] = struct{}{}
		uniqueX[p.X+1] = struct{}{}
		uniqueY[p.Y-1] = struct{}{}
		uniqueY[p.Y+1] = struct{}{}
	}

	var sortedX []int
	for x := range uniqueX {
		sortedX = append(sortedX, x)
	}
	slices.Sort(sortedX)

	var sortedY []int
	for y := range uniqueY {
		sortedY = append(sortedY, y)

	}

	slices.Sort(sortedY)

	return sortedX, sortedY
}

func buildAllowedGridCompressed(points []Point2D) ([][]bool, map[int]int, map[int]int, []GridCoord) {
	sortedX, sortedY := collectCompressedCoords(points)

	xMap := make(map[int]int)
	for i, x := range sortedX {
		xMap[x] = i
	}
	yMap := make(map[int]int)
	for i, y := range sortedY {
		yMap[y] = i
	}

	W_prime := len(sortedX)
	H_prime := len(sortedY)

	boundary := make([][]bool, H_prime)
	for y := range boundary {
		boundary[y] = make([]bool, W_prime)
	}

	n := len(points)
	for i := range n {
		p1 := points[i]
		p2 := points[(i+1)%n]

		cx1, cx2 := xMap[p1.X], xMap[p2.X]
		cy1, cy2 := yMap[p1.Y], yMap[p2.Y]

		if cx1 > cx2 {
			cx1, cx2 = cx2, cx1
		}
		if cy1 > cy2 {
			cy1, cy2 = cy2, cy1
		}

		if p1.X == p2.X {
			cx := xMap[p1.X]
			for cy := cy1; cy <= cy2; cy++ {
				boundary[cy][cx] = true
			}
		} else if p1.Y == p2.Y {
			cy := yMap[p1.Y]
			for cx := cx1; cx <= cx2; cx++ {
				boundary[cy][cx] = true
			}
		} else {
			log.Fatalf("Adjacent points not axis-aligned: %v %v", p1, p2)
		}
	}

	type coord struct{ y, x int }

	// Grid is H_prime x W_prime, extended grid for flood fill is (H_prime+2) x (W_prime+2)
	visited := make([][]bool, H_prime+2)
	for y := range visited {
		visited[y] = make([]bool, W_prime+2)
	}

	q := []coord{{0, 0}} // Start outside, top-left corner
	visited[0][0] = true

	dirs := [...]coord{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}

	for head := 0; head < len(q); head++ {
		c := q[head]
		for _, d := range dirs {
			ny, nx := c.y+d.y, c.x+d.x

			// Check bounds of extended grid
			if ny < 0 || ny > H_prime+1 || nx < 0 || nx > W_prime+1 {
				continue
			}
			if visited[ny][nx] {
				continue
			}

			// Check if we hit the boundary (fence) in the main compressed grid area (index 1 to H_prime/W_prime)
			if ny >= 1 && ny <= H_prime && nx >= 1 && nx <= W_prime {
				if boundary[ny-1][nx-1] {
					continue // Cannot pass through a boundary tile
				}
			}

			visited[ny][nx] = true
			q = append(q, coord{ny, nx})
		}
	}

	allowed := make([][]bool, H_prime)
	for y := range allowed {
		allowed[y] = make([]bool, W_prime)
	}

	for cy := range H_prime {
		for cx := range W_prime {
			if boundary[cy][cx] {
				allowed[cy][cx] = true
				continue
			}
			// If not boundary and not visited from the outside flood-fill, it is interior (allowed)
			if !visited[cy+1][cx+1] { // Use +1 offset for the extended visited grid
				allowed[cy][cx] = true
			}
		}
	}

	gridPointsCompressed := make([]GridCoord, len(points))
	for i, p := range points {
		gridPointsCompressed[i] = GridCoord{
			X: xMap[p.X],
			Y: yMap[p.Y],
		}
	}

	return allowed, xMap, yMap, gridPointsCompressed
}

func parseFile(path string) []Point2D {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")

	points := make([]Point2D, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Split(line, ",")
		if len(parts) != 2 {
			log.Fatalf("invalid line: %q", line)
		}

		// Use ParseInt for int
		x, err1 := strconv.ParseInt(strings.TrimSpace(parts[0]), 10, 64)
		y, err2 := strconv.ParseInt(strings.TrimSpace(parts[1]), 10, 64)
		if err1 != nil || err2 != nil {
			log.Fatalf("error parsing coords in %q", line)
		}

		points = append(points, Point2D{X: int(x), Y: int(y)})
	}

	if len(points) == 0 {
		log.Fatal("no points found")
	}
	return points
}
