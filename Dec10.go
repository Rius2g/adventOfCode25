package main

import (
	"fmt"
	"log"
	"math/bits"
	"os"
	"strconv"
	"strings"
)

func main() {
	dec10Part1()
	dec10part2()
}

func dec10Part1() {
	f, err := os.ReadFile("dec10Input.txt")
	if err != nil {
		log.Fatal(err)
	}

	var count int = 0

	lines := strings.Split(strings.TrimSpace(string(f)), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		rawWanted := parts[0]                        // e.g. "[###.]"
		wantedInner := strings.Trim(rawWanted, "[]") // e.g. "###."

		if len(wantedInner) > 63 {
			// We only support up to 63 lamps with a single uint64 (bit 0..62).
			continue
		}

		// collect only button tokens "(...)" and ignore the trailing "{...}"
		var buttons []string
		for _, p := range parts[1:] {
			if strings.HasPrefix(p, "{") {
				break // stop when we reach the stats block
			}
			buttons = append(buttons, p) // "(0,1,2)" etc
		}

		targetMask := parseTargetToMask(wantedInner)
		buttonMasks := make([]uint64, len(buttons))
		for i, b := range buttons {
			buttonMasks[i] = parseButtonMaskToUint(b)
		}

		minPresses, _ := findMinPressesBitwise(targetMask, buttonMasks)

		if minPresses == -1 {
			fmt.Printf("no solution with given buttons. Line: %q\n",
				line)
		} else {
			count += minPresses
		}
	}

	fmt.Println(count)
}

// parseTargetToMask converts a pattern like "###." into a uint64 bitmask.
// Bit i is 1 if position i is '#', 0 if '.'
func parseTargetToMask(s string) uint64 {
	var mask uint64
	for i, ch := range s {
		if ch == '#' {
			mask |= 1 << uint(i)
		}
	}
	return mask
}

// parseButtonMaskToUint parses "(0,1,3)" into a uint64 mask with bits 0,1,3 set.
func parseButtonMaskToUint(button string) uint64 {
	var mask uint64
	inner := strings.Trim(button, "()")
	if inner == "" {
		return mask
	}
	parts := strings.Split(inner, ",")
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		idx, err := strconv.Atoi(p)
		if err != nil {
			continue
		}
		if idx < 0 || idx >= 63 {
			// We decided to support [0,62] only
			continue
		}
		mask ^= 1 << uint(idx) // ^= in case same index appears twice in one button
	}
	return mask
}

// findMinPressesBitwise finds minimal presses with bitwise masks.
// target: desired lamp configuration as uint64
// buttons: each button's toggle mask as uint64
// returns minimal number of presses and which button indices were used.
func findMinPressesBitwise(target uint64, buttons []uint64) (int, []int) {
	n := len(buttons)
	if n == 0 {
		if target == 0 {
			return 0, []int{}
		}
		return -1, nil
	}

	totalCombos := 1 << n

	best := -1
	var bestCombo []int

	for subset := 0; subset < totalCombos; subset++ {
		var state uint64 = 0 // start from all dots

		// apply each button whose bit is 1 in subset
		for b := 0; b < n; b++ {
			if subset&(1<<b) != 0 {
				state ^= buttons[b] // toggle
			}
		}

		if state == target {
			presses := bits.OnesCount(uint(subset))
			if best == -1 || presses < best {
				best = presses
				// record indices of pressed buttons
				bestCombo = bestCombo[:0]
				for b := 0; b < n; b++ {
					if subset&(1<<b) != 0 {
						bestCombo = append(bestCombo, b)
					}
				}
			}
		}
	}

	return best, bestCombo
}

func dec10part2() {
	f, err := os.ReadFile("dec10Input.txt")
	if err != nil {
		log.Fatal(err)
	}

	lines := strings.Split(strings.TrimSpace(string(f)), "\n")
	totalPresses := 0

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		// Parse target { ... }
		var targetStr string
		var buttonStrs []string
		for _, p := range parts {
			if strings.HasPrefix(p, "{") {
				targetStr = p
			} else if strings.HasPrefix(p, "(") {
				buttonStrs = append(buttonStrs, p)
			}
		}

		if targetStr == "" {
			continue // skip incomplete lines
		}

		target := parseTarget(targetStr)
		buttonIndices := make([][]int, len(buttonStrs))
		for i, b := range buttonStrs {
			buttonIndices[i] = parseButtonIndices(b)
		}

		presses := solveJoltageILP(target, buttonIndices)

		if presses != -1 {
			totalPresses += presses
		} else {
			// To help debugging, we only print if a known "good" puzzle fails
			// fmt.Printf("No solution found for line: %s\n", line)
		}
	}

	fmt.Println("Total minimal button presses:", totalPresses)
}

// parseTarget parses "{3,5,4,7}" into []int{3,5,4,7}
func parseTarget(s string) []int {
	inner := strings.Trim(s, "{}")
	if inner == "" {
		return []int{}
	}
	parts := strings.Split(inner, ",")
	res := make([]int, len(parts))
	for i, p := range parts {
		p = strings.TrimSpace(p)
		v, err := strconv.Atoi(p)
		if err != nil {
			log.Fatalf("invalid target value %q in %q", p, s)
		}
		res[i] = v
	}
	return res
}

// parseButtonIndices parses "(1,3)" into []int{1,3}
func parseButtonIndices(s string) []int {
	inner := strings.Trim(s, "()")
	if inner == "" {
		return []int{}
	}
	parts := strings.Split(inner, ",")
	var res []int
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		v, err := strconv.Atoi(p)
		if err != nil {
			log.Fatalf("invalid button index %q in %q", p, s)
		}
		res = append(res, v)
	}
	return res
}

// solveJoltageILP uses optimized brute force with pruning
func solveJoltageILP(target []int, buttons [][]int) int {
	n := len(target)
	m := len(buttons)

	if m == 0 {
		for _, t := range target {
			if t != 0 {
				return -1
			}
		}
		return 0
	}

	// Build coefficient matrix
	coeff := make([][]int, n)
	for i := 0; i < n; i++ {
		coeff[i] = make([]int, m)
	}
	for j, btn := range buttons {
		for _, i := range btn {
			if i < n {
				coeff[i][j] = 1
			}
		}
	}

	minPresses := -1
	presses := make([]int, m)
	current := make([]int, n)

	var search func(btnIdx int, currentSum int)
	search = func(btnIdx int, currentSum int) {
		// Prune if current sum already >= best
		if minPresses != -1 && currentSum >= minPresses {
			return
		}

		// Check if any counter exceeds target (infeasible)
		for i := 0; i < n; i++ {
			if current[i] > target[i] {
				return
			}
		}

		if btnIdx == m {
			// Check if we've hit all targets
			for i := 0; i < n; i++ {
				if current[i] != target[i] {
					return
				}
			}
			if minPresses == -1 || currentSum < minPresses {
				minPresses = currentSum
			}
			return
		}

		// Calculate max useful presses for this button
		// Can't exceed the remaining deficit for any counter it affects
		maxPress := -1
		for i := 0; i < n; i++ {
			if coeff[i][btnIdx] == 1 {
				remaining := target[i] - current[i]
				if maxPress == -1 || remaining < maxPress {
					maxPress = remaining
				}
			}
		}
		if maxPress < 0 {
			maxPress = 0
		}

		// Try different press counts for this button
		for p := 0; p <= maxPress; p++ {
			presses[btnIdx] = p
			// Update current counters
			for i := 0; i < n; i++ {
				current[i] += coeff[i][btnIdx] * p
			}
			search(btnIdx+1, currentSum+p)
			// Restore current counters
			for i := 0; i < n; i++ {
				current[i] -= coeff[i][btnIdx] * p
			}
		}
	}

	search(0, 0)
	return minPresses
}
