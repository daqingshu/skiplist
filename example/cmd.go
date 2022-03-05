package main

import (
	"fmt"

	"gitee.com/daqingshu/skiplist"
)

func main() {
	l := skiplist.NewSkiplist[int, uint32]()
	for i := 0; i < 10; i++ {
		l.Insert(i, uint32(i))
	}
	var k = 5
	v := l.Search(k)
	fmt.Println(*v)
	v = l.Delete(k)
	fmt.Println(*v)

	s := l.Search(k)
	if s != nil {
		fmt.Println(*s)
	} else {
		fmt.Printf("%v is not in list\n", k)
	}
	fmt.Println("Press the Enter Key to exit")
	fmt.Scanln()
}
