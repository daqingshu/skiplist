package skiplist

type Ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64 |
		~string
}

type Node[K Ordered, V any] struct {
	Key   K
	Value V
	Next  []*Node[K, V]
	Prev  *Node[K, V]
}

func newNode[K Ordered, V any](key K, value V, level uint32) *Node[K, V] {
	nd := &Node[K, V]{
		Key:   key,
		Value: value,
		Next:  make([]*Node[K, V], level),
		Prev:  nil,
	}
	for i := 0; i < int(level); i++ {
		nd.Next[i] = nil
	}
	return nd
}
