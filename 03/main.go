package main

import (
	"fmt"
	"io"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	data, err := read(os.Stdin)
	if err != nil {
		panic(err)
	}
	mulDataList, err := parseData([]byte(purgeBetween(string(data), "don't()", "do()")))
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

func purgeBetween(content string, start string, end string) string {
	removed := content
	p1 := -1
	p2 := -1
	for i := range removed {
		if p1 == -1 {
			for j, jChar := range start {
				next := i + j
				if next >= len(removed) || string((removed[next])) != string(jChar) {
					p1 = -1
					break
				}
				if p1 != -1 {
					continue
				}
				p1 = i
			}
			continue
		}

		if p1 == -1 || p1+len(start)-1 >= i {
			continue
		}

		if p2 == -1 {
			for j, jChar := range end {
				next := i + j
				if next >= len(removed) || string(removed[next]) != string(jChar) {
					p2 = -1
					break
				}
				p2 = next
			}
			continue
		}

		if p2 == -1 {
			continue
		}

		removed = replaceBetween(removed, p1, p2, "#")
		p1 = -1
		p2 = -1
	}
	if p1 != -1 {
		removed = replaceBetween(removed, p1, len(removed)-1, "#")
	}
	return removed
}

func replaceBetween(original string, p1 int, p2 int, replacement string) string {
	start := original[0:p1]
	end := original[p2:]
	return start + strings.Repeat(replacement, int(math.Abs(float64(p2-p1)))) + end
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
