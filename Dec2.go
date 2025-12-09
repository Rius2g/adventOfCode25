package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func Dec2() {

	f, err := os.ReadFile("dec2Input.txt")
	if err != nil {
		log.Fatal(err)
	}

	// remove trailing newline or spaces
	content := strings.TrimSpace(string(f))
	total := 0

	ranges := strings.Split(content, ",")

	for _, r := range ranges {
		r = strings.TrimSpace(r)
		parts := strings.SplitN(r, "-", 2)
		if len(parts) != 2 {
			log.Fatalf("invalid range: %q", r)
		}

		startStr := strings.TrimSpace(parts[0])
		endStr := strings.TrimSpace(parts[1])

		start, err := strconv.Atoi(startStr)
		if err != nil {
			log.Fatalf("could not convert start %q: %v", startStr, err)
		}

		end, err := strconv.Atoi(endStr)
		if err != nil {
			log.Fatalf("could not convert end %q: %v", endStr, err)
		}

		total += sumInvalidInRange(start, end)
	}
	fmt.Println(total)

}

func sumInvalidInRange(start, stop int) int {
	sum := 0
	for n := start; n <= stop; n++ {
		if checkNum(n) {
			sum += n
		}
	}
	return sum
}

func checkNum(num int) bool {
	s := strconv.Itoa(num)
	n := len(s)

	// Try every possible pattern length
	for L := 1; L <= n/2; L++ {
		// String length must be a multiple of L
		if n%L != 0 {
			continue
		}

		reps := n / L
		if reps < 2 {
			continue
		}

		pattern := s[0:L]
		ok := true

		for i := 1; i < reps; i++ {
			if s[i*L:(i+1)*L] != pattern {
				ok = false
				break
			}
		}

		if ok {
			return true
		}
	}

	return false
}
