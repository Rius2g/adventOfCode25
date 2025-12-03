package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	part1()
	part2()

}

func part1() {
	var count int = 0

	//read the file
	f, err := os.ReadFile("dec3Input.txt")
	if err != nil {
		log.Fatal(err)
	}
	//parse the file contents

	lines := strings.Split(string(f), "\n")
	for i := range len(lines) - 1 {
		line := lines[i]

		intLine := digitsIntoTokens(line)
		//iterate over string and find largest and second largest int,
		//be careful about the order they are in
		var num1 int8
		var num2 int8

		var biggestSum int8 = convertTo2IntegerNum(num1, num2)
		for j := 0; j < len(intLine)-1; j++ { //first num does NOT go to the end
			//double loop this bitc-1h
			num1 = int8(intLine[j])
			for k := j + 1; k < len(intLine); k++ { //second num goes to the end
				num2 = int8(intLine[k])
				sumNum := convertTo2IntegerNum(num1, num2)
				if sumNum > biggestSum {
					biggestSum = sumNum
				}

			} //we loop second num first (moving that, becuase it migth free up idx)
		}
		//
		fmt.Println(biggestSum)
		count += int(biggestSum)
	}

	fmt.Println(count)

}

func part2() {
	var count int = 0

	//read the file
	f, err := os.ReadFile("dec3Input.txt")
	if err != nil {
		log.Fatal(err)
	}

	lines := strings.Split(string(f), "\n")
	for i := range len(lines) - 1 {
		line := lines[i]

		intLine := digitsIntoTokens(line)
		//each is 100 ints long, we want the 12 best

		//greedy window algo
		//we loop through 12 times, we find the LARGEST num we can for the given index AS LONG as there are enough nums after to fill!
		nums := make([]int, 12) //pre-alloc 12
		K := 12                 //amount of ints needed
		n := len(intLine)
		start := 0
		for pos := 0; pos < K; pos++ {
			//we loop throgh find big num as long as 11 - i spots left we can select it
			remaining := K - pos
			maxStart := n - remaining //cant start later than this as it would be not enough left
			bestIdx := start          //start on first available
			bestVal := intLine[start]
			for k := start + 1; k <= maxStart; k++ { //loop the inner window for the max allowed start
				if intLine[k] > bestVal {
					bestVal = intLine[k]
					bestIdx = k
				}
			}
			nums[pos] = bestVal
			start = bestIdx + 1
		}
		count += convertTo12DigistNum(nums)
	}

	fmt.Println(count)
}

func digitsIntoTokens(s string) []int {
	nums := make([]int, 0, len(s))
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			continue
		}
		nums = append(nums, int(s[i]-'0'))
	}
	return nums
}

func convertTo12DigistNum(nums []int) int {
	sum := 0
	for _, d := range nums {
		sum = sum*10 + d
	}
	return sum
}

func convertTo2IntegerNum(num1, num2 int8) int8 {

	return (num1*10 + num2)
}
