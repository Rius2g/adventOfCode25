package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

func Dec5() {
	////	Dec5Part1()
	Dec5part2()
}

func Dec5Part1() {

	//what is the plan, split the file into two (on the empty line)
	//store the ranges as min-max, check num x if >= [i].min AND x <= [i].max ==> fresh

	var count int = 0

	//read the file
	f, err := os.ReadFile("dec5Input.txt")
	if err != nil {
		log.Fatal(err)
	}

	// Normalize and split on the empty line
	content := strings.ReplaceAll(string(f), "\r\n", "\n") // handle Windows endings
	parts := strings.SplitN(strings.TrimSpace(content), "\n\n", 2)
	if len(parts) != 2 {
		panic("expected exactly one empty line separator")
	}

	// [0] = ranges before empty line
	freshRanges := strings.Split(strings.TrimSpace(parts[0]), "\n")

	// [1] = lines after empty line
	questioned := strings.Split(strings.TrimSpace(parts[1]), "\n")

	//parse the file contents

	//split file into FRESH part [0] AND non-fresh part[1]

	//loop over questioned part
	//check against fresh part
	for i := 0; i <= len(questioned)-1; i++ {
		//how to do this shit? IF 0 we read only + 1, if else we read +1 and -1 lines as well
		//then we loop through and check the grid around
		//reset each loop
		qNum, err := strconv.Atoi(strings.TrimSpace(questioned[i]))
		if err != nil {
			log.Fatal("no valid qnum")
		}
		for j := 0; j < len(freshRanges)-1; j++ {
			rang := strings.Split(freshRanges[j], "-")
			if len(rang) != 2 {
				log.Fatal("invalid range")
			}

			minNum, err1 := strconv.Atoi(strings.TrimSpace(rang[0]))
			maxNum, err2 := strconv.Atoi(strings.TrimSpace(rang[1]))

			if err1 != nil || err2 != nil {
				log.Fatal("Invalid range")
			}
			if checkNumInRange(qNum, minNum, maxNum) == 1 {
				count++
				break

			}
		}

	}

	fmt.Println(count)

}

type Range struct {
	minimum int
	maximum int
}

func Dec5part2() {
	var count int = 0

	//read the file
	f, err := os.ReadFile("dec5Input.txt")
	if err != nil {
		log.Fatal(err)
	}

	// Normalize and split on the empty line
	content := strings.ReplaceAll(string(f), "\r\n", "\n") // handle Windows endings
	parts := strings.SplitN(strings.TrimSpace(content), "\n\n", 2)
	if len(parts) != 2 {
		panic("expected exactly one empty line separator")
	}

	// [0] = ranges before empty line
	freshRanges := strings.Split(strings.TrimSpace(parts[0]), "\n")

	ranges := make([]Range, 0, len(freshRanges))
	for _, r := range freshRanges {
		p := strings.Split(strings.TrimSpace(r), "-")
		if len(p) != 2 {
			log.Fatal("bad range")
		}
		mi, _ := strconv.Atoi(p[0])
		mx, _ := strconv.Atoi(p[1])
		ranges = append(ranges, Range{mi, mx})
	}

	sort.Slice(ranges, func(i, j int) bool {
		return ranges[i].minimum < ranges[j].minimum
	}) //sorted, now we check for overlap and MERGE

	mergedRanges := mergeRanges(ranges)

	//now just count?

	for i := range mergedRanges {
		count += (mergedRanges[i].maximum - mergedRanges[i].minimum) + 1
	}

	fmt.Println(count)

}

func mergeRanges(ranges []Range) []Range {
	if len(ranges) == 0 {
		return ranges
	}
	merged := []Range{ranges[0]}

	for i := 1; i < len(ranges); i++ {
		last := &merged[len(merged)-1]
		cur := ranges[i]
		if cur.minimum <= last.maximum {
			if cur.maximum > last.maximum {
				last.maximum = cur.maximum
			}
		} else {
			merged = append(merged, cur)
		}

	}
	return merged
}

func checkNumInRange(num, minNum, maxNum int) int {

	if num >= minNum && num <= maxNum {
		return 1
	}
	return 0
}
