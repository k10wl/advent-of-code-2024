package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"slices"
	"strconv"
	"strings"
)

type memory struct{ id int }

var visual bool

func main() {
	v := flag.Bool("v", false, "")
	flag.Parse()
	visual = *v
	memReader := new(memoryReader)
	data, err := memReader.read(os.Stdin)
	if err != nil {
		panic(err)
	}

	worker := new(fileMover)
	organized := worker.organize(data)

	sum, err := checksum(organized)
	if err != nil {
		panic(err)
	}
	if visual {
		fmt.Printf("%+v\n", inline(data))
		fmt.Printf("%+v\n", inline(organized))
	}
	fmt.Printf("%d\n", sum)
}

func inline(data []memory) string {
	var sb strings.Builder
	for _, mem := range data {
		if mem.id == -1 {
			sb.WriteString(".|")
			continue
		}
		sb.WriteString(fmt.Sprintf("%d|", mem.id))
	}
	return sb.String()
}

type memoryReader struct{}

func (mem memoryReader) read(r io.Reader) ([]memory, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	m := []memory{}
	for i, b := range data {
		asInt, err := strconv.Atoi(string(b))
		if err != nil {
			return nil, err
		}
		if i%2 != 0 {
			for i := 0; i < asInt; i++ {
				m = append(m, memory{id: -1})
			}
			continue
		}
		for j := 0; j < asInt; j++ {
			m = append(m, memory{id: i / 2})
		}
	}

	return m, nil
}

type fragmentator struct {
	p1   int
	p2   int
	data []memory
}

func (f *fragmentator) organize(data []memory) []memory {
	defer f.reset()
	f.data = slices.Clone(data)
	f.p1 = 0
	f.p2 = len(data) - 1
	for {
		p1 := f.findP1()
		p2 := f.findP2()
		if p1 > p2 {
			break
		}
		f.p1 = p1
		f.p2 = p2
		f.swap()
	}
	return f.data
}

func (f *fragmentator) reset() {
	f.p1 = -1
	f.p2 = -1
}

func (f *fragmentator) findP1() int {
	for i := f.p1; i < len(f.data); i++ {
		if id := f.data[i].id; id == -1 {
			return i
		}
	}
	return -1
}

func (f *fragmentator) findP2() int {
	for i := f.p2; i > 0; i-- {
		if id := f.data[i].id; id != -1 {
			return i
		}
	}
	return -1
}

func (f *fragmentator) swap() {
	f.data[f.p1], f.data[f.p2] = f.data[f.p2], f.data[f.p1]
	f.p1 += 1
	f.p2 -= 1
}

type fileMover struct {
	data               []memory
	cachedMemoryBounds map[int][]*memoryBounds
	filePointer        int
	memoryPointer      int
	attempted          map[int]bool
}

type memoryBounds struct {
	left  int
	right int
}

func (f *fileMover) organize(data []memory) []memory {
	f.reset()
	defer f.reset()
	f.data = slices.Clone(data)
	for {
		fileBounds := f.findFile()
		if fileBounds == nil {
			if visual {
				fmt.Printf("file bounds are nil\n")
			}
			break
		}
		memoryBounds := f.findFree(fileBounds.right - fileBounds.left + 1)
		if visual {
			fmt.Printf("fileBounds: %v\n", fileBounds)
			fmt.Printf("memoryBounds: %v\n", memoryBounds)
			fmt.Printf("before: %v\n", inline(f.data))
		}
		if memoryBounds == nil || memoryBounds.right > fileBounds.left {
			continue
		}
		f.swap(fileBounds, memoryBounds)
		if visual {
			fmt.Printf("after:  %v\n", inline(f.data))
		}
		if f.filePointer == 0 {
			break
		}
	}
	return f.data
}

func (f *fileMover) retrieveCached(n int) *memoryBounds {
	bounds, ok := f.cachedMemoryBounds[n]
	if !ok {
		return nil
	}
	bound := bounds[0]
	f.cachedMemoryBounds[n] = bounds[1:]
	return bound
}

func (f *fileMover) cache(n int, bounds *memoryBounds) {
	f.cachedMemoryBounds[n] = append(f.cachedMemoryBounds[n], bounds)
}

func (f *fileMover) findFree(n int) *memoryBounds {
	// cached := f.retrieveCached(n)
	// if cached != nil {
	// 	return cached
	// }
	if visual {
		fmt.Printf("size lookup: %v\n", n)
	}
	for left, mem := range f.data {
		if mem.id != -1 {
			continue
		}
		size := 1
		for j := left + 1; j < len(f.data)-1; j++ {
			if f.data[j].id != -1 {
				break
			}
			size++
		}
		free := &memoryBounds{left: left, right: left + size - 1}
		if visual {
			fmt.Printf("free: %v | %d\n", free, size)
			fmt.Printf("found size size: %v\n", size)
		}
		if size >= n {
			return free
		}
		f.cache(size, free)
	}
	return nil
}

func (f *fileMover) findFile() *memoryBounds {
	if f.filePointer == -1 {
		f.filePointer = len(f.data) - 1
	}
	var mem memoryBounds
	if visual {
		fmt.Printf("all atempted: %+v\n", f.attempted)
	}
	for i := f.filePointer; i > 0; i-- {
		if f.data[i].id == -1 || f.attempted[f.data[i].id] {
			continue
		}
		j := i - 1
		for ; j > -1; j-- {
			if f.data[j].id != f.data[i].id {
				break
			}
		}
		f.attempted[f.data[i].id] = true
		mem.left = j + 1
		mem.right = i
		if visual {
			fmt.Printf("file: %v\n", mem)
		}
		f.filePointer = j
		return &mem
	}
	return nil
}

func (f *fileMover) reset() {
	f.cachedMemoryBounds = map[int][]*memoryBounds{}
	f.data = nil
	f.filePointer = -1
	f.attempted = map[int]bool{}
}

func (f *fileMover) swap(file *memoryBounds, free *memoryBounds) {
	for i := 0; i <= file.right-file.left; i++ {
		if visual {
			fmt.Printf("swapping: %v\n", i)
		}
		f.data[file.left+i], f.data[free.left+i] = f.data[free.left+i], f.data[file.left+i]
	}
}

func checksum(data []memory) (int, error) {
	sum := 0
	for i, mem := range data {
		if mem.id == -1 {
			continue
		}
		sum += mem.id * i
	}
	return sum, nil
}
