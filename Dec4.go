package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func Dec4() {
	dec4part1()
	dec4part2()
}

func dec4part1() {
	//sliding 3x3 window

	var count int = 0

	//read the file
	f, err := os.ReadFile("dec4Input.txt")
	if err != nil {
		log.Fatal(err)
	}
	//parse the file contents

	lines := strings.Split(string(f), "\n")
	for i := 0; i <= len(lines)-1; i++ {
		//how to do this shit? IF 0 we read only + 1, if else we read +1 and -1 lines as well
		//then we loop through and check the grid around
		//reset each loop
		var line1 string = ""
		var line2 string = ""
		var line3 string = ""
		line2 = lines[i]
		if i == len(lines)-1 {
			//only read backwards, not forwards
			line1 = lines[i-1] //not forward
		} else if i == 0 {
			line3 = lines[i+1]

		} else { //read both froward and backward
			line1 = lines[i-1] //backward
			line3 = lines[i+1] //forward

		}
		count += countAdjecent(line1, line2, line3)

	}

	fmt.Println(count)

}

//create an array of these

func dec4part2() {
	grid, err := readGrid("dec4Input.txt")
	if err != nil {
		log.Fatal(err)
	}

	total := 0
	removed := true

	for removed {
		removed = false
		for row := 0; row < len(grid); row++ {
			added := countAndMarkAdjacent(grid, row)
			if added > 0 {
				removed = true
				total += added
			}
		}
	}

	fmt.Println(total)
}

func readGrid(path string) ([][]rune, error) {
	f, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	rawLines := strings.Split(strings.TrimRight(string(f), "\n"), "\n")
	grid := make([][]rune, 0, len(rawLines))
	for _, line := range rawLines {
		grid = append(grid, []rune(line))
	}
	return grid, nil
}

func countAndMarkAdjacent(grid [][]rune, row int) int {
	rollsAccessible := 0
	for col, char := range grid[row] {
		if char == '@' {
			count := 0
			count += checkLineIndexesRunes(grid[row], []int{col - 1, col + 1})
			if row > 0 {
				count += checkLineIndexesRunes(grid[row-1], []int{col - 1, col, col + 1})
			}
			if row < len(grid)-1 {
				count += checkLineIndexesRunes(grid[row+1], []int{col - 1, col, col + 1})
			}

			if count < 4 {
				rollsAccessible++
				grid[row][col] = 'x' // THIS now sticks in memory
			}
		}
	}
	return rollsAccessible
}

func checkLineIndexesRunes(line []rune, index []int) int {
	if len(line) == 0 {
		return 0
	}
	count := 0
	for _, idx := range index {
		if idx >= 0 && idx < len(line) {
			if line[idx] == '@' {
				count++
			}
		}
	}
	return count
}

func countAdjecent(line1, line2, line3 string) int {
	//loop over line2 AKA the line we are on
	rollsAccessible := 0
	for i, char := range line2 {
		if char == '@' {
			count := 0

			count += checklineIndexs(line2, []int{i - 1, i + 1})
			count += checklineIndexs(line1, []int{i - 1, i + 1, i})
			count += checklineIndexs(line3, []int{i - 1, i + 1, i})
			//found paper, need top check IF accessible
			//we can always count forward and backwards
			//check other lines, first we chekck if the other lines extist

			if count < 4 {
				rollsAccessible++
			}
		}
	}
	return rollsAccessible
}

func checklineIndexs(line string, index []int) int {
	if line == "" {
		return 0
	}
	count := 0
	for _, idx := range index {
		if idx >= 0 && idx <= len(line)-1 {
			if line[idx] == '@' {
				count++
			}
		}

	}

	return count
}
