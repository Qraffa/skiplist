package main

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

const (
	maxLevel  = 16
	skipListP = 1
)

type SkipList struct {
	header, tail *node
	length       int
	topLevel     int
}

type node struct {
	key     int
	value   interface{}
	forward []*node
}

func init() {
	rand.Seed(time.Now().Unix())
}

func NewSkipList() *SkipList {
	return &SkipList{
		header:   makeNode(maxLevel, 0, nil),
		tail:     nil,
		length:   0,
		topLevel: 0,
	}
}

func makeNode(level int, key int, val interface{}) *node {
	n := &node{}
	n.key = key
	n.value = val
	n.forward = make([]*node, level)
	return n
}

func (s *SkipList) Search(key int) (interface{}, error) {
	// start := s.header
	// var level = maxLevel
	// for ; level >= 0 && s.header.forward[level] != nil; level-- {
	// }
	// TODO comp
	x := s.header
	for i := s.topLevel - 1; i >= 0; i-- {
		for x.forward[i] != nil && x.forward[i].key < key {
			x = x.forward[i]
		}
	}
	x = x.forward[0]
	if x != nil && x.key == key {
		return x.value, nil
	}
	return nil, errors.New("not found")
}

func (s *SkipList) Insert(key int, val interface{}) {
	// update saves prev node
	update := make([]*node, maxLevel)
	x := s.header
	for i := s.topLevel - 1; i >= 0; i-- {
		for x.forward[i] != nil && x.forward[i].key < key {
			x = x.forward[i]
		}
		update[i] = x
	}
	x = x.forward[0]
	// update
	if x != nil && x.key == key {
		x.value = val
		return
	}
	level := randomLevel()
	// greater than top level, update the top level
	if level > s.topLevel {
		// add prev
		for i := s.topLevel; i < level; i++ {
			update[i] = s.header
		}
		s.topLevel = level
	}
	x = makeNode(level, key, val)
	// update pointer
	for i := 0; i < level; i++ {
		x.forward[i] = update[i].forward[i]
		update[i].forward[i] = x
	}
	s.length++
}

func (s *SkipList) Delete(key int) {
	// update saves prev node
	update := make([]*node, maxLevel)
	x := s.header
	for i := s.topLevel - 1; i >= 0; i-- {
		for x.forward[i] != nil && x.forward[i].key < key {
			x = x.forward[i]
		}
		update[i] = x
	}
	x = x.forward[0]
	// not exists
	if x == nil || x.key != key {
		return
	}
	// update pointer
	for i := 0; i < s.topLevel; i++ {
		// greater than x's level
		if update[i].forward[i] != x {
			break
		}
		update[i].forward[i] = x.forward[i]
	}
	s.length--
	// update top level
	for s.topLevel > 1 && s.header.forward[s.topLevel-1] == nil {
		s.topLevel--
	}
}

func (s *SkipList) Size() int {
	return s.length
}

func (s *SkipList) print() {
	img := make([][]int, s.topLevel)
	for i := 0; i < len(img); i++ {
		img[i] = make([]int, s.length+2)
	}
	km := make(map[int]int)
	x := s.header
	for i := s.topLevel - 1; i >= 0; i-- {
		var col int = 1
		for x.forward[i] != nil {
			img[s.topLevel-1-i][col] = x.forward[i].key
			km[x.forward[i].key] = len(x.forward[i].forward)
			x = x.forward[i]
			col++
		}
		x = s.header
	}
	fmt.Println(km)
	for i := 0; i < len(img)-1; i++ {
		for j := len(img[i]) - 2; j >= 0; j-- {
			v := img[i][j]
			if v != 0 && img[i][v] == 0 {
				img[i][j], img[i][v] = img[i][v], img[i][j]
			}
		}
	}
	for i := 0; i < len(img); i++ {
		for j := 0; j < len(img[0]); j++ {
			if j > 0 && j < len(img[0])-1 && img[i][j] == 0 {
				if img[i][j+1] == 0 {
					fmt.Printf("---\t")
				} else {
					fmt.Printf("-->\t")
				}
			} else {
				fmt.Printf("%3d\t", img[i][j])
			}
		}
		fmt.Println()
	}
}

func randomLevel() int {
	level := 1
	for rand.Int31()&0xFFFF < 0xFFFF>>skipListP && level < maxLevel {
		level++
	}
	return level
}
