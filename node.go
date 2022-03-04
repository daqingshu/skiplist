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
	Level []*Level[K, V]
	Prev  *Node[K, V]
}

type Level[K Ordered, V any] struct {
	Next *Node[K, V]
}

func newNode[K Ordered, V any](level uint32, key K, value V) *Node[K, V] {
	nd := &Node[K, V]{
		Key:   key,
		Value: value,
		Level: make([]*Level[K, V], level),
		Prev:  nil,
	}
	for i := 0; i < int(level); i++ {
		nd.Level[i] = &Level[K, V]{Next: nil}
	}
	return nd
}
