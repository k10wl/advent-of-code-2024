package main

import (
	"fmt"
	"io"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
)

func main() {
	data := readToData(os.Stdin)
	safeCount := 0
	for _, v := range data {
		at := failAt(v)
		if at == -1 {
			safeCount++
			continue
		}
		for i := range v {
			retry := slices.Clone(v)
			at = failAt(slices.Delete(retry, i, i+1))
			if at == -1 {
				safeCount++
				break
			}
		}
	}
	fmt.Printf("%d", safeCount)
}

func failAt(data []int) int {
	at := -1
	progressionVector := 0
	for index, current := range data {
		if index == 0 {
			continue
		}
		previous := data[index-1]
		if math.Abs(float64(previous-current)) > 3 {
			at = index
			break
		}
		if index == 1 {
			progressionVector = getDirection(previous, current)
		} else if progressionVector != getDirection(previous, current) {
			at = index
			break
		}
	}
	return at
}

func getDirection(a int, b int) int {
	return int(math.Max(-1, math.Min(1, float64(b-a))))
}

func readToData(r io.Reader) [][]int {
	b, err := io.ReadAll(r)
	if err != nil {
		panic(err)
	}
	return readRows(getRows(string(b)))
}

func getRows(all string) []string {
	return strings.Split(strings.TrimSpace(string(all)), "\n")
}

func readRows(all []string) [][]int {
	data := make([][]int, len(all), len(all))
	for i, v := range all {
		data[i] = readData(v)
	}
	return data
}

func readData(row string) []int {
	str := strings.Split(row, " ")
	data := make([]int, len(str), len(str))
	for i, v := range str {
		intData, err := strconv.Atoi(v)
		if err != nil {
			panic(err)
		}
		data[i] = intData
	}
	return data
}
