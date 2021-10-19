package main

import (
	"fmt"
	"os"
	"runtime/pprof"
	"testing"
)

func TestRandomLevel(t *testing.T) {
	ls := make([]int, 32+1)
	for i := 0; i < 1000; i++ {
		l := randomLevel()
		ls[l]++
	}
	for i, v := range ls {
		fmt.Println(i, float64(v)/1000)
	}
}

func TestSkiplist(t *testing.T) {
	sl := NewSkipList()
	sl.Insert(1, 1)
	sl.Insert(2, 2)

	v, err := sl.Search(1)
	if err != nil {
		panic(err)
	}
	fmt.Println(v.(int))
	v, err = sl.Search(2)
	if err != nil {
		panic(err)
	}
	fmt.Println(v.(int))
}

func TestDel(t *testing.T) {
	sl := NewSkipList()
	for i := 1; i <= 100; i++ {
		sl.Insert(i, i)
	}
	fmt.Println(sl.Size())
	for i := 1; i <= 200; i++ {
		sl.Delete(i)
	}
	fmt.Println(sl.Size())
}

func TestImage(t *testing.T) {
	sl := NewSkipList()
	for i := 1; i <= 16; i++ {
		sl.Insert(i, i)
	}
	sl.print()
}

// go test -bench=BenchmarkInsert -test.run=BenchmarkInsert -benchmem
func BenchmarkInsert(b *testing.B) {
	f, _ := os.Create("cpu_profile")
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()
	b.ReportAllocs()
	sl := NewSkipList()
	for i := 0; i < b.N; i++ {
		sl.Insert(i, i)
	}
}

// go test -bench=BenchmarkDel -test.run=BenchmarkDel -benchmem -benchtime=1000000x
func BenchmarkDel(b *testing.B) {
	f, _ := os.Create("mem_profile")
	b.ReportAllocs()
	sl := NewSkipList()
	for i := 0; i < b.N; i++ {
		sl.Insert(i, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sl.Delete(i)
	}
	defer pprof.WriteHeapProfile(f)
}
