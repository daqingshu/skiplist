package main

import (
	"fmt"

	"gitee.com/daqingshu/skiplist"
	"github.com/valyala/fastrand"
)

func main() {
	l := skiplist.NewSkiplist[int, uint32]()
	for i := 0; i < 100000; i++ {
		l.Insert(i, fastrand.Uint32())
	}

	v := l.Search(10000)
	fmt.Println(*v)
	fmt.Println("Press the Enter Key to exit")
	fmt.Scanln()
}
