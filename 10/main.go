package main

import (
	"fmt"
	"io"
	"os"
	"slices"
	"strconv"
	"strings"
)

const (
	blocked = -1
	start   = 0
	finish  = 9
)

type position struct {
	lat  int
	long int
}

func main() {
	m := new(topographicMap)
	if err := m.read(os.Stdin); err != nil {
		panic(err)
	}
	fmt.Printf("peek\n%+v\n", m.peek())
	fmt.Printf("%d", m.sumScore())
}

type topographicMap struct {
	data       [][]*node
	trailHeads []*node
	trailTails []*node
}

func (t *topographicMap) read(reader io.Reader) error {
	data, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	strdata := strings.Split(strings.TrimSpace(string(data)), "\n")
	t.trailHeads = []*node{}
	t.trailTails = []*node{}
	t.data = make([][]*node, len(strdata), len(strdata))
	for i, long := range strdata {
		longitude := make([]*node, len(long), len(long))
		for j, lat := range long {
			height, err := strconv.Atoi(string(lat))
			if err != nil {
				height = blocked
			}
			node := newNode(height)
			if i != 0 {
				node.attach(north, t.data[i-1][j])
			}
			if j != 0 {
				node.attach(west, longitude[j-1])
			}
			if height == start {
				t.trailHeads = append(t.trailHeads, node)
			}
			if height == finish {
				node.score = 1
				t.trailTails = append(t.trailTails, node)
			}
			longitude[j] = node
		}
		t.data[i] = longitude
	}
	t.backtrack()
	return nil
}

func rateDistinct(head *node) func(n *node) {
	return func(n *node) {
		if n.height == start && !slices.Contains(
			n.destinations,
			head,
		) {
			n.destinations = append(n.destinations, head)
			n.score += 1
		}
	}
}

func ratePath(n *node) {
	n.score += 1
}

func (t *topographicMap) backtrack() {
	for _, head := range t.trailTails {
		head.backtrack(ratePath)
	}
}

func (t topographicMap) sumScore() int {
	sum := 0
	for _, n := range t.trailHeads {
		sum += n.score
	}
	return sum
}

func (t *topographicMap) peek() string {
	var sb strings.Builder
	for i, lat := range t.data {
		repeats := len(lat)*2 + 1
		if i == 0 {
			sb.WriteString(fmt.Sprintf("\n%s\n", strings.Repeat("-", repeats)))
		}
		for j, long := range lat {
			if j == 0 {
				sb.WriteString("|")
			}
			if long.height == blocked {
				sb.WriteString(" |")
			} else {
				sb.WriteString(fmt.Sprintf("%d|", long.score))
			}
		}
		sb.WriteString(fmt.Sprintf("\n%s\n", strings.Repeat("-", repeats)))
	}
	return sb.String()
}

/*

[Task]

[Part 1]
What is the sum of the scores of all trailheads on your topographic map?

[Part 2]
rating is the number of distinct hiking trails

=======================================================================

[Definitions]

topographic map - task input, surrounding area

trailhead       - any position that starts one or more hiking trails

hiking trail    - any path that starts at height 0, ends at height 9
it always increases height by exactly 1 at each step
never includes diagonal steps, only NWSE

score           - trailhead ranking based on number of 9 height positions
reachable from trailhead via hiking trail

=======================================================================

[Examples]

[Part 1]

89010123
78121874
87430965
96549874
45678903
32019012
01329801
10456732

> 9 trailheads, score reading in order: 5,6,5,3,1,3,5,3,5
> total summed score 36

------------------------------------------------------------------------

...0...
...1...
...2...
6543456
7.....7
8.....8
9.....9

> score: 2

------------------------------------------------------------------------

..90..9
...1.98
...2..7
6543456
765.987
876....
987....

> score: 4

> two paths are leadign to same height at bottom left

------------------------------------------------------------------------

10..9..
2...8..
3...7..
4567654
...8..3
...9..2
.....01

> top trailhead has score of 1, bottom trailhead - 2

------------------------------------------------------------------------

[Part 2]

.....0.
..4321.
..5..2.
..6543.
..7..4.
..8765.
..9....

> 3

------------------------------------------------------------------------

..90..9
...1.98
...2..7
6543456
765.987
876....
987....

> 13

------------------------------------------------------------------------

89010123
78121874
87430965
96549874
45678903
32019012
01329801
10456732

> 81

*/
