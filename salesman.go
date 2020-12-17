package main

import (
	"fmt"
	"time"
)

var graph = [][]int{
	{0, 72, 47, 61, 82, 40, 71, 88, 34, 51, 20, 21, 4},
	{72, 0, 10, 28, 16, 10, 41, 10, 60, 47, 94, 86, 10},
	{47, 10, 0, 23, 16, 99, 80, 77, 26, 5, 18, 51, 87},
	{61, 28, 23, 0, 83, 98, 26, 29, 0, 98, 5, 83, 80},
	{82, 16, 16, 83, 0, 15, 32, 18, 68, 15, 52, 66, 61},
	{40, 10, 99, 98, 15, 0, 11, 66, 60, 8, 25, 80, 14},
	{71, 41, 80, 26, 32, 11, 0, 15, 74, 79, 26, 88, 77},
	{88, 10, 77, 29, 18, 66, 15, 0, 94, 88, 6, 51, 73},
	{34, 60, 26, 0, 68, 60, 74, 94, 0, 94, 81, 24, 41},
	{51, 47, 5, 98, 15, 8, 79, 88, 94, 0, 8, 84, 41},
	{20, 94, 18, 5, 52, 25, 26, 6, 81, 8, 0, 85, 65},
	{21, 86, 51, 83, 66, 80, 88, 51, 24, 84, 85, 0, 86},
	{4, 10, 87, 80, 61, 14, 77, 73, 41, 41, 65, 86, 0},
}

const startingCity = 0

func testThat(currentCity int, left []int, lastToVisit int, path []int, verbose int) (int, []int) {
	if verbose > -1 {
		fmt.Println("Testing", currentCity)
	}
	var cost int
	var bestWay []int
	if len(left) == 1 {
		//fmt.Println("Returning this values:", graph[currentCity][lastToVisit], append(path, lastToVisit))
		return graph[currentCity][lastToVisit], append(path, lastToVisit)
	}
	for currIndex, nextCity := range left {
		if nextCity != currentCity && nextCity != lastToVisit {
			value := graph[currentCity][nextCity]
			//fmt.Println("Testing branch:", currentCity, left, value, nextCity, currIndex, cost, bestWay, path)
			newLeft := make([]int, len(left))
			if len(left) > 2 {
				copy(newLeft, left)
				lastIndex := len(newLeft) - 1
				newLeft[currIndex] = newLeft[lastIndex]
				newLeft[lastIndex] = 0
				newLeft = newLeft[:lastIndex]
			} else if len(left) > 1 {
				newLeft = left[currIndex-1 : currIndex]
			}
			//fmt.Println("New left:", newLeft)
			possibleCost, possibleWay := testThat(nextCity, newLeft, lastToVisit, append(path, currentCity), verbose-1)
			if cost == 0 || (possibleCost+value) < cost {
				cost = value + possibleCost
				bestWay = possibleWay
			}

		}
	}
	return cost, bestWay
}

func main() {
	before := time.Now().UnixNano()
	left := make([]int, len(graph))
	for i := range left {
		left[i] = i
	}
	cost, bestWay := testThat(startingCity, left, startingCity, []int{}, 1)
	fmt.Println("Best:")
	fmt.Println(cost)
	fmt.Println(bestWay)
	after := time.Now().UnixNano()
	diff := after - before
	fmt.Println("Took:", diff/1000, "ms")
}
