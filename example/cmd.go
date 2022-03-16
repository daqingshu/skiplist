package main

import (
	"fmt"
	"time"

	"github.com/daqingshu/skiplist"
)

func main() {
	begin := time.Now().UnixNano()
	l := skiplist.NewSkiplist[int, uint32]()
	for i := 0; i < 10; i++ {
		n := l.Search(i)
		if n != nil {
			fmt.Println("get", i)
		} else {
			fmt.Println("can not find ", i)
		}
	}
	for i := 0; i < 1000000; i++ {
		l.Insert(i, uint32(i))
	}

	l.Delete(55)

	for i := 59; i >= 50; i-- {
		n := l.Search(i)
		if n != nil {
			fmt.Println("get", i)
		} else {
			fmt.Println("can not find ", i)
		}
	}
	end := time.Now().UnixNano()
	elspsedTime := (end - begin) / 1000000
	fmt.Printf("total used %v ms\n", elspsedTime)
	// fmt.Println("Press the Enter Key to exit")
	// fmt.Scanln()
}
