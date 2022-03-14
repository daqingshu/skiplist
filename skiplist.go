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
		head:   newNode(k, v, maxLevel),
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

	x = s.head
	for i := int(s.level) - 1; i >= 0; i-- {
		for x.Next[i] != nil &&
			x.Next[i].Key < key {
			x = x.Next[i]
		}
		update[i] = x
	}
	x = update[0]
	if x.Next[0] != nil {
		if x.Next[0].Key == key {
			return x.Next[0]
		}
	}
	level := randomLevel()
	if level > s.level {
		for i := int(s.level); i < int(level); i++ {
			update[i] = s.head
		}
		s.level = level
	}
	newNode := newNode(key, value, level)
	for i := 0; i < int(level); i++ {
		newNode.Next[i] = update[i].Next[i]
		update[i].Next[i] = newNode
	}
	if update[0] == s.head {
		newNode.Prev = nil
	} else {
		newNode.Prev = update[0]
	}
	if newNode.Next[0] != nil {
		newNode.Next[0].Prev = newNode
	} else {
		s.tail = newNode
	}
	s.length++
	return newNode
}

func (s *Skiplist[K, V]) Search(key K) *V {
	if s == nil {
		return nil
	}
	x := s.head
	for i := int(s.level) - 1; i >= 0; i-- {
		for x.Next[i] != nil &&
			x.Next[i].Key < key {
			x = x.Next[i]
		}
	}
	x = x.Next[0]
	if x == nil {
		return nil
	}
	if x.Key == key {
		return &x.Value
	}
	return nil
}

func (s *Skiplist[K, V]) deleteNode(x *Node[K, V], update *[maxLevel]*Node[K, V]) {
	for i := 0; i < int(s.level); i++ {
		if update[i].Next[i] == nil {
			continue
		}
		if update[i].Next[i] != x {
			break
		}
		update[i].Next[i] = x.Next[i]
	}
	if x.Next[0] != nil {
		x.Next[0].Prev = x.Prev
	} else {
		s.tail = x.Prev
	}
	for {
		if s.level > 1 && s.head.Next[s.level-1] == nil {
			s.level--
		} else {
			break
		}
	}
	s.length--
}

func (s *Skiplist[K, V]) Delete(key K) *V {
	if s == nil {
		return nil
	}
	var update [maxLevel]*Node[K, V]
	x := s.head
	for i := int(s.level) - 1; i >= 0; i-- {
		for x.Next[i] != nil &&
			x.Next[i].Key < key {
			x = x.Next[i]
		}
		update[i] = x
	}
	x = x.Next[0]
	if x != nil && x.Key == key {
		s.deleteNode(x, &update)
		return &x.Value
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
