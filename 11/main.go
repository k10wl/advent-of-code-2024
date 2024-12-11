package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

var (
	visual = false
	blinks = 0
)

const invalid = -1

func main() {
	pebbles, err := readPebbles(os.Stdin)
	if err != nil {
		panic(err)
	}
	observer := new(person)
	observer.observe(pebbles)
	fmt.Printf("%+v", observer.blink(blinks))
}

func init() {
	v := flag.Bool("visual", false, "")
	b := flag.Int("blink", 0, "")
	flag.Parse()
	visual = *v
	blinks = *b
}

var cache = make(map[cached]int)

type cached struct {
	value int
	n     int
}

func act(value int, n int) int {
	if n == 0 {
		return 1
	}
	c := cached{value: value, n: n}
	if val, ok := cache[c]; ok {
		return val
	}
	if value == 0 {
		res := act(1, n-1)
		cache[c] = res
		return res
	}
	if asString := fmt.Sprintf("%d", value); len(asString)%2 == 0 {
		split := len(asString) / 2
		fmt.Println(asString)
		left, err := strconv.Atoi(asString[split:])
		if err != nil {
			panic(err)
		}
		right, err := strconv.Atoi(asString[:split])
		if err != nil {
			panic(err)
		}
		res := act(left, n-1) + act(right, n-1)
		cache[c] = res
		return res
	}
	res := act(value*2024, n-1)
	cache[c] = res
	return res
}

type person struct {
	observing []int
}

func (p *person) observe(items []int) {
	p.observing = items
}

func (p person) peek() string {
	var sb strings.Builder
	for _, pebble := range p.observing {
		sb.WriteString(fmt.Sprintf("%d ", pebble))
	}
	return strings.TrimSpace(sb.String())
}

func (p person) blink(n int) int {
	res := 0
	for _, entity := range p.observing {
		res += act(entity, n)
	}
	return res
}

func readPebbles(r io.Reader) ([]int, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	out := []int{}
	for _, str := range strings.Split(string(data), " ") {
		asInt, err := strconv.Atoi(str)
		if err != nil {
			return nil, err
		}
		out = append(out, asInt)
	}
	return out, nil
}

/*
--- [Day 11 : Plutonian Pebbles](https://adventofcode.com/2024/day/11) ---

Part 1: Count how many stones will you have after blinking 25 times?

stones are arranged in a perfectly straight line, each stone has a **number**
engraved on it. Sometimes the number changes. Other times the stone might
**split in two**, causing other stone to shift over a bit to make room in line

Every time you blink, the stones each **simultaneously** change according to
rules:
- if stone is 0, it is replaced by 1.
- if number has even number of digits, it is replaced by two **stones**.
  the left half of the digits are engraved on the new left stone, and the
  right half on the new right stone. New numbers don't keep extra zeros
  1000 would become 10 | 0.
- if none of the other rules apply the stune is replaced by a new stone;
  the old stone's number multiplied by 2024 is engraved on the new stone.

Order is preserved after changes, they stay on their perfectly straight line


--------------------------------[examples]-------------------------------------

0 1 10 99 999

blink once - 1 2024 1 0 9 9 2021976
7 stones

-------------------------------------------------------------------------------

125 17

6 blinks - 2097446912 14168 4048 2 0 2 4 40 48 2024 40 48 80 96 2 8 6 7 6 0 3 2
22 stones

92 0 286041 8034 34394 795 8 2051489


*/
