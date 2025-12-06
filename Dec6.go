package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

func main() {
	Dec6Part1()
	Dec6Part2()
}

func Dec6Part1() {
	var count int = 0

	//split into columns not rows?

	lines, err := readLinesAsFields("dec6Input.txt")
	if err != nil {
		fmt.Print(err)
		return
	}
	if len(lines) < 2 {
		fmt.Println("Need at least one data row and one operator row")
		return
	}
	opRow := lines[len(lines)-1]
	numCols := len(opRow)
	for col := 0; col < numCols; col++ {
		res, err := computeColumn("dec6Input.txt", col)
		if err != nil {
			fmt.Println(err)
			continue
		}

		count += res
	}
	fmt.Println(count)

}

func readLinesAsFields(fileName string) ([][]string, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var lines [][]string

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		lines = append(lines, fields)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func computeColumn(fileName string, col int) (int, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	var lines [][]string
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		lines = append(lines, fields)
	}
	if err := scanner.Err(); err != nil {
		return 0, err
	}

	if len(lines) < 2 {
		return 0, fmt.Errorf("need at least one data row and one operator row")
	}

	opRow := lines[len(lines)-1]

	// include col and length in error
	if col < 0 || col >= len(opRow) {
		return 0, fmt.Errorf("column %d out of range for operator row (len=%d)", col, len(opRow))
	}
	op := opRow[col]

	var values []int
	for i := 0; i < len(lines)-1; i++ {
		row := lines[i]
		if col >= len(row) {
			return 0, fmt.Errorf("column %d out of range for data row %d (len=%d)", col, i, len(row))
		}
		v, err := strconv.Atoi(row[col])
		if err != nil {
			return 0, fmt.Errorf("failed to parse int in row %d, col %d (%q): %w", i, col, row[col], err)
		}
		values = append(values, v)
	}

	if len(values) == 0 {
		return 0, fmt.Errorf("no values found in column %d", col)
	}

	//reconstruct the values in the proper format here //we start from the back

	result := values[0]
	for i := 1; i < len(values); i++ {
		switch op {
		case "*":
			result *= values[i]
		case "+":
			result += values[i]
		default:
			return 0, fmt.Errorf("unsupported operator %q in column %d", op, col)
		}
	}

	return result, nil
}

func Dec6Part2() {
	grid, err := readGrid2("dec6Input.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	if len(grid) == 0 {
		fmt.Println("empty grid")
		return
	}

	rows := len(grid)
	cols := len(grid[0])

	total := 0

	var curNums []int
	var curOp rune

	flushProblem := func() { //adding function
		if len(curNums) == 0 {
			return
		}

		acc := curNums[0]
		for i := 1; i < len(curNums); i++ {
			switch curOp {
			case '+':
				acc += curNums[i]
			case '*':
				acc *= curNums[i]
			default:
				panic(fmt.Sprintf("unknown op"))
			}
		}
		total += acc
		curNums = nil
		curOp = 0
	}

	//parsing part

	isSpaceCol := func(c int) bool {
		for r := 0; r < rows; r++ {
			if grid[r][c] != ' ' {
				return false
			}
		}
		return true
	}

	for c := cols - 1; c >= 0; c-- {
		if isSpaceCol(c) {
			flushProblem()
			continue
		}

		var sb strings.Builder
		for r := 0; r < rows-1; r++ {
			ch := grid[r][c]
			if unicode.IsDigit(ch) {
				sb.WriteRune(ch)
			}
		}
		if sb.Len() == 0 {
			continue
		}
		val, err := strconv.Atoi(sb.String())
		if err != nil {
			fmt.Println(err)
			return
		}
		curNums = append(curNums, val)

		bottom := grid[rows-1][c]
		if bottom == '+' || bottom == '*' {
			curOp = bottom
		}
	}
	flushProblem()

	fmt.Println(total)

}

func readGrid2(filename string) ([][]rune, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)

	maxLen := 0 //need to write the remaining empty strings

	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
		if len(line) > maxLen {
			maxLen = len(line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	grid := make([][]rune, len(lines))
	for i, line := range lines {
		if len(line) < maxLen {
			line = line + strings.Repeat(" ", maxLen-len(line)) //make all the lines equally long by adding whitespaces
		}
		grid[i] = []rune(line)
	}
	return grid, nil

}
