package main

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
)

func main() {
	data, err := read(os.Stdin)
	if err != nil {
		panic(err)
	}
	mulDataList, err := parseData(cleanup(data))
	if err != nil {
		panic(err)
	}
	combinedMul := 0
	for _, mulData := range mulDataList {
		combinedMul += mul(mulData.a, mulData.b)
	}
	fmt.Print(combinedMul)
}

type mulDataContainer struct {
	a int
	b int
}

func cleanup(content []byte) []byte {
	pattern := regexp.MustCompile(`don't\(\).*?(\n|\z)?.*?do\(\)`)
	mid := pattern.ReplaceAllString(string(content), "")
	return []byte(regexp.MustCompile(`don't\(\).*`).ReplaceAllString(mid, ""))
}

func parseData(content []byte) ([]mulDataContainer, error) {
	pattern := regexp.MustCompile(`mul\((?P<left>\d+),\s?(?P<right>\d+)\)`)
	left := pattern.SubexpIndex("left")
	right := pattern.SubexpIndex("right")
	d := []mulDataContainer{}
	var e error
	for _, submatches := range pattern.FindAllSubmatch(content, -1) {
		a, err := strconv.Atoi(string(submatches[left]))
		if err != nil {
			e = err
			break
		}
		b, err := strconv.Atoi(string(submatches[right]))
		if err != nil {
			e = err
			break
		}
		d = append(d, mulDataContainer{a: a, b: b})
	}
	return d, e
}

func mul(a int, b int) int {
	return a * b
}

func read(r io.Reader) ([]byte, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return b, nil
}
