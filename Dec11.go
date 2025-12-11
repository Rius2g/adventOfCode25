package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

var branchInt int = 0

func main() {
	f, err := os.ReadFile("dec11Input.txt")
	if err != nil {
		log.Fatal(err)
	}

	//	dec11Part1(f)
	dec11Part2(f)

}

// this finds ALL branches, we need to eliminate those that DONT lead to "out"
func branchOut(continuePoint string, lines []string) {
	for i := 0; i < len(lines)-1; i++ {
		stri := strings.Split(lines[i], ":")
		if strings.Compare(stri[0], continuePoint) == 0 {
			branches := strings.Split(stri[1], " ")
			//could also trim remove space here, but who tf cares
			if len(branches) == 2 && branches[1] == "out" {
				branchInt++ //only count resolved branches that lead to out
				break
			}
			for j := 1; j < len(branches); j++ {
				branchOut(branches[j], lines)
			}
			break
		}
	}
}

func dec11Part1(f []byte) {
	lines := strings.Split(strings.TrimSpace(string(f)), "\n")

	entry := "you"

	for i := 0; i < len(lines)-1; i++ {
		stri := strings.Split(lines[i], ":")
		if strings.Compare(stri[0], entry) == 0 {
			branches := strings.Split(stri[1], " ")
			for j := 1; j < len(branches); j++ {
				branchOut(branches[j], lines)
			}
			break
		}
	}
	fmt.Println(branchInt)

}

// Adjacency map and memo for part 2
var adj map[string][]string
var memo map[string]int

func dec11Part2(f []byte) {
	lines := strings.Split(strings.TrimSpace(string(f)), "\n")

	// Build adjacency map once - O(1) lookups instead of O(n)
	adj = make(map[string][]string)
	for _, line := range lines {
		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			continue
		}
		node := parts[0]
		neighbors := strings.Fields(parts[1]) // splits on whitespace, handles leading space
		adj[node] = neighbors
	}

	// Memo cache: key = "node|visitedBitmask" -> count of valid paths
	memo = make(map[string]int)

	// Start from svr with visited=0 (neither dac nor fft seen yet)
	result := countPaths("svr", 0)
	fmt.Println(result)
}

// countPaths returns the number of paths from 'node' to "out" that visit both dac and fft
// visited is a bitmask: bit 0 = dac visited, bit 1 = fft visited
func countPaths(node string, visited int) int {
	// Update visited bitmask when we hit required nodes
	if node == "dac" {
		visited |= 1
	}
	if node == "fft" {
		visited |= 2
	}

	// Base case: reached "out"
	if node == "out" {
		if visited == 3 { // both dac (1) and fft (2) visited
			return 1
		}
		return 0
	}

	// Check memo
	key := fmt.Sprintf("%s|%d", node, visited)
	if v, ok := memo[key]; ok {
		return v
	}

	// Recurse through all neighbors
	total := 0
	for _, next := range adj[node] {
		total += countPaths(next, visited)
	}

	memo[key] = total
	return total
}
