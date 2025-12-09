package main

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
)

func Dec7() {
	Dec7Part1()
	Dec7Part2()
}

var count int = 0

func Dec7Part1() {

	//easy solution? we have an array of column indexes were the beams pass through, each time we reach splitter we modify and continue

	//read the file
	f, err := os.ReadFile("dec7Input.txt")
	if err != nil {
		log.Fatal(err)
	}
	//parse the file contents

	lines := strings.Split(string(f), "\n")
	//lets start by finding S in the line and setting that
	beamIdxs := make([]int, 0)
	for startIdx, col := range lines[0] {
		if col == 'S' {
			beamIdxs = append(beamIdxs, startIdx) //we append the start one
		}
	}
	for i := 1; i <= len(lines)-1; i++ {
		//nice here we have number of lines
		//now we simply start here, we see if a "^" is directly below any indexes in beamIdx, and we modify the array
		checkLine(lines[i], &beamIdxs)

	}

	fmt.Println(count)
}

func checkLine(col string, beamIdxs *[]int) {
	for i, entry := range col {
		if entry != '^' {
			continue
		}

		// only split if i is already a beam index
		if !slices.Contains(*beamIdxs, i) {
			continue
		}

		left := i - 1
		right := i + 1
		count++

		*beamIdxs = removeValue(*beamIdxs, i) //remove the original value, for the 2 edge case we will add it back in following lines

		if !slices.Contains(*beamIdxs, left) {
			*beamIdxs = append(*beamIdxs, left)
		}
		if !slices.Contains(*beamIdxs, right) {
			*beamIdxs = append(*beamIdxs, right)
		}
	}

}

func removeValue(xs []int, value int) []int {
	for i, v := range xs {
		if v == value {
			return append(xs[:i], xs[i+1:]...)
		}
	}
	return xs // value not found
}

func Dec7Part2() {
	//read the file
	f, err := os.ReadFile("dec7Input.txt")
	if err != nil {
		log.Fatal(err)
	}
	//parse the file contents

	lines := strings.Split(string(f), "\n")
	//lets start by finding S in the line and setting that
	startIdx := 0
	for idx, col := range lines[0] {
		if col == 'S' {
			startIdx = idx
			break
		}
	}
	timelineCounts := map[int]int{}
	timelineCounts[startIdx] = 1
	for i := 1; i <= len(lines)-1; i++ {
		//nice here we have number of lines
		//now we simply start here, we see if a "^" is directly below any indexes in beamIdx, and we modify the array
		timelineCounts = propagateRow(lines[i], timelineCounts)
	}

	total := 0
	for _, k := range timelineCounts {
		total += k
	}
	fmt.Println(total)

}

func propagateRow(colLine string, counts map[int]int) map[int]int {
	width := len(colLine)
	if width == 0 {
		return counts
	}

	next := make(map[int]int)

	for col, k := range counts {
		if k == 0 { //check bounds for last line
			continue
		}
		if col < 0 || col >= width {
			continue
		}

		ch := colLine[col]

		if ch == '^' { //timelineSplit
			next[col-1] += k
			next[col+1] += k
		} else { //no timeline split (only 1 timeline for this index)
			next[col] += k
		}
	}

	return next
}
