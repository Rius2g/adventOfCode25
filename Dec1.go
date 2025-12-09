package main

import (
	//"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func Dec1() {

	//need to read the prefix (L OR R), convert the rest of the string to int
	//need to calculate, we are going from 0-99, how many 0 does it land on?
	//bounds, take number +- new number

	var start int = 50
	var count int = 0

	//read the file
	f, err := os.ReadFile("dec1Input.txt")
	if err != nil {
		log.Fatal(err)
	}
	//parse the file contents

	lines := strings.Split(string(f), "\n")
	for i := range len(lines) - 1 {
		prefix := string(lines[i][0])
		num, err := strconv.Atoi(string(lines[i][1:]))
		if err != nil {
			log.Fatalf("could not fetch number")
		}

		var countIncrease int
		switch prefix {
		case "L":
			{
				start, countIncrease = getPositivt(start, num)
			}
		case "R":
			{
				start, countIncrease = decreaseSubHundred(start, num)
			}
		}
		count += countIncrease
	}
	log.Println(count)
}

func getPositivt(start int, movement int) (int, int) {
	newPos := start - movement
	for newPos < 0 {
		newPos += 100
	}
	countIncrease := 0
	if start == 0 {
		countIncrease = movement / 100
	} else if start <= movement {
		countIncrease = (movement-start)/100 + 1
	}
	return newPos, countIncrease
}

func decreaseSubHundred(start int, movement int) (int, int) {
	newPos := start + movement
	countIncrease := newPos / 100
	newPos = newPos % 100
	return newPos, countIncrease
}
