package main

import (
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"slices"
	"strings"
)

type coordinates struct {
	longitude int
	latitude  int
}

type antinode struct {
	frequency string
	coordinates
}

type _state struct {
	antennas map[string][]coordinates
	city     []string
}

func main() {
	input, err := readInput()
	if err != nil {
		panic(err)
	}
	state := getInitialState(input)
	antinodes := filterUniqueCoordinatesInBounds(
		findAntinodes(state.antennas),
		len(state.city[0])-1,
	)
	visual := flag.Bool("v", false, "")
	flag.Parse()
	if *visual {
		printAll(state, antinodes)
	}
	fmt.Printf(
		"%d\n",
		len(antinodes),
	)
}

func printAll(state _state, antinodes []antinode) {
	out := slices.Clone(state.city)

	anti := map[string]bool{}
	for _, val := range antinodes {
		fmt.Printf("%+v\n", val)
		anti[fmt.Sprintf("%d:%d", val.longitude, val.latitude)] = true
	}

	for longitude, row := range out {
		runes := []rune(row)
		for latitude := range runes {
			if anti[fmt.Sprintf("%d:%d", longitude, latitude)] {
				runes[latitude] = rune('#')
			}
		}
		fmt.Printf("%s\n", string(runes))
	}

}

func filterUniqueCoordinatesInBounds(
	data []antinode,
	maxCoords int,
) []antinode {
	seen := map[string]bool{}
	return slices.DeleteFunc(
		slices.Clone(data),
		func(e antinode) bool {
			coordstring := fmt.Sprintf(
				"%d:%d",
				e.longitude,
				e.latitude,
			)
			if 0 > compare(e.longitude, e.latitude, math.Min) ||
				maxCoords < compare(e.longitude, e.latitude, math.Max) ||
				seen[coordstring] {
				return true
			}
			seen[coordstring] = true
			return false
		},
	)
}

func findAntinodes(antennas map[string][]coordinates) []antinode {
	antinodes := []antinode{}
	for frequency, data := range antennas {
		antinodes = append(
			antinodes,
			calculateCoordinates(data[0], data[1:], frequency)...,
		)
	}
	return antinodes
}

func calculateCoordinates(
	target coordinates,
	data []coordinates,
	frequency string,
) []antinode {
	if len(data) == 1 {
		return calculateAntinodes(target, data[0], frequency)
	}
	res := []antinode{}
	for _, value := range data {
		res = append(res, calculateAntinodes(target, value, frequency)...)
	}
	res = append(res, calculateCoordinates(data[0], data[1:], frequency)...)
	return res
}

func calculateAntinodes(a coordinates, b coordinates, frequency string) []antinode {
	diffLongitude := diff(a.longitude, b.longitude)
	diffLatitude := diff(a.latitude, b.latitude)
	res := []antinode{}
	for i := -500; i <= 500; i++ { // cursed sizing, but whatever
		if i == 0 {
			continue
		}
		bottom := antinode{
			frequency:   frequency,
			coordinates: coordinates{},
		}
		top := antinode{
			frequency:   frequency,
			coordinates: coordinates{},
		}
		if a.longitude > b.longitude {
			bottom.coordinates.longitude = a.longitude + diffLongitude*i
			top.coordinates.longitude = b.longitude - diffLongitude*i
		} else {
			bottom.coordinates.longitude = a.longitude - diffLongitude*i
			top.coordinates.longitude = b.longitude + diffLongitude*i
		}
		if a.latitude > b.latitude {
			bottom.coordinates.latitude = a.latitude + diffLatitude*i
			top.coordinates.latitude = b.latitude - diffLatitude*i
		} else {
			bottom.coordinates.latitude = a.latitude - diffLatitude*i
			top.coordinates.latitude = b.latitude + diffLatitude*i
		}
		res = append(res, bottom, top)
	}
	return res
}

func diff(a int, b int) int {
	return int(math.Abs(float64(a - b)))
}

func compare(a int, b int, comparator func(float64, float64) float64) int {
	return int(comparator(float64(a), float64(b)))
}

func getInitialState(input string) _state {
	antennas := map[string][]coordinates{}
	city := strings.Split(input, "\n")
	for longitude, row := range city {
		for latitude, rune := range row {
			char := string(rune)
			if char == "." {
				continue
			}
			antennas[char] = append(
				antennas[char],
				coordinates{longitude: longitude, latitude: latitude},
			)
		}
	}
	return _state{antennas: antennas, city: city}
}

func readInput() (string, error) {
	data, err := io.ReadAll(os.Stdin)
	return string(data), err
}
