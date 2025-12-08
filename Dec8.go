package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Point3D struct {
	X, Y, Z int
}

type Edge struct {
	I, J int
	Dist float64
}

// Disjoint Set Union (Union-Find) to track circuits.
type DSU struct {
	parent []int
	size   []int
}

func NewDSU(n int) *DSU {
	parent := make([]int, n)
	size := make([]int, n)
	for i := range n {
		parent[i] = i
		size[i] = 1
	}
	return &DSU{parent: parent, size: size}
}

func (d *DSU) Find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.Find(d.parent[x])
	}
	return d.parent[x]
}

func (d *DSU) Union(a, b int) bool {
	ra := d.Find(a) //find parent for the circuit they are in
	rb := d.Find(b)
	if ra == rb { //if same parent, same circuit
		return false // already same circuit
	}
	if d.size[ra] < d.size[rb] {
		ra, rb = rb, ra
	}
	d.parent[rb] = ra
	d.size[ra] += d.size[rb]
	return true
}

func Distance3D(p1, p2 Point3D) float64 {
	dx := float64(p2.X - p1.X)
	dy := float64(p2.Y - p1.Y)
	dz := float64(p2.Z - p1.Z)
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

func main() {
	points := parseInput("dec8Input.txt")
	edges := buildEdges(points)

	// Sort edges by distance, then by indices for stability
	sort.Slice(edges, func(a, b int) bool {
		if edges[a].Dist == edges[b].Dist {
			if edges[a].I == edges[b].I {
				return edges[a].J < edges[b].J
			}
			return edges[a].I < edges[b].I
		}
		return edges[a].Dist < edges[b].Dist
	})

	// Part 1 (what you already did): 1000 connections, product of 3 largest circuits
	// You can comment this out if you’re done with part 1.
	dec8part1(points, edges)

	// Part 2: keep connecting until everything is one circuit;
	// report the pair that finished the job.
	dec8part2(points, edges)
}

func parseInput(path string) []Point3D {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")

	points := make([]Point3D, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Split(line, ",")
		if len(parts) != 3 {
			log.Fatalf("invalid line: %q", line)
		}

		x, err1 := strconv.Atoi(strings.TrimSpace(parts[0]))
		y, err2 := strconv.Atoi(strings.TrimSpace(parts[1]))
		z, err3 := strconv.Atoi(strings.TrimSpace(parts[2]))
		if err1 != nil || err2 != nil || err3 != nil {
			log.Fatalf("error parsing coords in %q", line)
		}

		points = append(points, Point3D{X: x, Y: y, Z: z})
	}

	if len(points) == 0 {
		log.Fatal("no points found")
	}
	return points
}

func buildEdges(points []Point3D) []Edge {
	n := len(points)
	var edges []Edge
	for i := 0; i < n-1; i++ {
		for j := i + 1; j < n; j++ {
			d := Distance3D(points[i], points[j])
			edges = append(edges, Edge{I: i, J: j, Dist: d})
		}
	}
	return edges
}

// Part 1, for reference: 1000 shortest connections overall.
func dec8part1(points []Point3D, edges []Edge) {
	n := len(points)
	dsu := NewDSU(n)

	connectionsToTake := min(1000, len(edges))

	for k := range connectionsToTake {
		e := edges[k]
		_ = dsu.Union(e.I, e.J) // even if already same circuit, still counts as a connection
	}

	// collect circuit sizes
	componentSeen := make(map[int]bool)
	var circuitSizes []int
	for i := range n {
		root := dsu.Find(i)
		if !componentSeen[root] {
			componentSeen[root] = true
			circuitSizes = append(circuitSizes, dsu.size[root])
		}
	}

	sort.Slice(circuitSizes, func(i, j int) bool {
		return circuitSizes[i] > circuitSizes[j]
	})

	if len(circuitSizes) < 3 {
		log.Printf("Part 1: need at least 3 circuits, got %d; sizes=%v", len(circuitSizes), circuitSizes)
		return
	}

	result := circuitSizes[0] * circuitSizes[1] * circuitSizes[2]
	fmt.Println("Part 1 – product of 3 largest circuits:", result)
}

// Part 2: keep connecting until we get down to a single circuit.
// The last merging edge is the one we care about.
func dec8part2(points []Point3D, edges []Edge) {
	n := len(points)
	dsu := NewDSU(n)
	components := n

	var finalEdge *Edge

	for _, e := range edges {
		if dsu.Union(e.I, e.J) {
			components--
			if components == 1 {
				// This is the first time everything is connected.
				finalEdge = &e
				break
			}
		}
	}

	if finalEdge == nil {
		log.Fatal("Part 2: never reached a single circuit – something is off")
	}

	p1 := points[finalEdge.I]
	p2 := points[finalEdge.J]
	product := p1.X * p2.X

	fmt.Printf("Part 2 – final connecting pair indices: %d and %d\n", finalEdge.I, finalEdge.J)
	fmt.Printf("Part 2 – final connecting pair coords: (%d,%d,%d) and (%d,%d,%d)\n",
		p1.X, p1.Y, p1.Z, p2.X, p2.Y, p2.Z,
	)
	fmt.Printf("Part 2 – product of X coords: %d\n", product)
}
