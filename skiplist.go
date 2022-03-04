package skiplist

import (
	"math"
	"sync/atomic"

	"github.com/valyala/fastrand"
)

const (
	maxLevel = 32
	pValue   = 1 / math.E
)

var (
	probabilities [maxLevel]uint32
)

func init() {
	p := float64(1.0)
	for i := 0; i < maxLevel; i++ {
		probabilities[i] = uint32(float64(math.MaxUint32) * p)
		p *= pValue
	}
}

type Skiplist[K Ordered, V any] struct {
	head   *Node[K, V]
	tail   *Node[K, V]
	length uint32
	level  uint32
}

func NewSkiplist[K Ordered, V any]() *Skiplist[K, V] {
	var k K
	var v V
	skl := &Skiplist[K, V]{
		head:   newNode(maxLevel, k, v),
		tail:   nil,
		length: 0,
		level:  1,
	}
	skl.head.Prev = nil

	return skl
}

func (s *Skiplist[K, V]) Insert(key K, value V) *Node[K, V] {
	if s == nil {
		return nil
	}
	var update [maxLevel]*Node[K, V]
	var x *Node[K, V]
	var i int
	x = s.head
	for i = int(s.level) - 1; i >= 0; i-- {
		for {
			if x.Level[i].Next != nil &&
				x.Level[i].Next.Key <= key {
				x = x.Level[i].Next
			} else {
				break
			}
		}
		update[i] = x
	}
	level := randomLevel()
	if level > s.level {
		for i = int(s.level); i < int(level); i++ {
			update[i] = s.head
		}
		s.level = level
	}
	x = newNode(level, key, value)
	for i = 0; i < int(level); i++ {
		x.Level[i].Next = update[i].Level[i].Next
		update[i].Level[i].Next = x
	}
	if update[0] == s.head {
		x.Prev = nil
	} else {
		x.Prev = update[0]
	}
	if x.Level[0].Next != nil {
		x.Level[0].Next.Prev = x
	} else {
		s.tail = x
	}
	s.length++
	return x
}

func (s *Skiplist[K, V]) Search(key K) *V {
	if s == nil {
		return nil
	}
	x := s.head
	var i int
	for i = int(s.level) - 1; i >= 0; i-- {
		for {
			if x.Level[i].Next != nil &&
				x.Level[i].Next.Key <= key {
				x = x.Level[i].Next
			} else {
				break
			}
		}
		if x.Key == key {
			return &x.Value
		}
	}
	return nil
}

func (s *Skiplist[K, V]) Level() uint32 {
	return atomic.LoadUint32(&s.level)
}

func randomLevel() uint32 {
	rnd := fastrand.Uint32()

	h := uint32(1)
	for h < maxLevel && rnd <= probabilities[h] {
		h++
	}

	return h
}
