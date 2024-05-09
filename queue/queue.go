package queue

type FixedQueue[K comparable] struct {
	items                       []K
	front, rear, capacity, size int
}

func NewFixedQueue[K comparable](capacity int) *FixedQueue[K] {
	return &FixedQueue[K]{
		items:    make([]K, capacity),
		capacity: capacity,
		front:    1,
		rear:     0,
		size:     0,
	}
}
func (q *FixedQueue[K]) Full() bool {
	return q.size == q.capacity
}
func (q *FixedQueue[K]) Empty() bool {
	return q.size == 0
}
func (q *FixedQueue[K]) Front() K {
	if q.Empty() {
		panic("queue underflow")
	}
	return q.items[q.front]
}
func (q *FixedQueue[K]) PushBack(k K) {
	if q.Full() {
		panic("queue overflow")
	}
	q.rear = (q.rear + 1) % q.capacity
	q.items[q.rear] = k
	q.size++
}
func (q *FixedQueue[K]) PopFront() {
	if q.Empty() {
		panic("queue underflow")
	}
	q.front = (q.front + 1) % q.capacity
	q.size--
}

func (q *FixedQueue[K]) Remove(k K) {
	for i := q.front; i != q.rear; i = (i + 1) % q.capacity {
		item := q.items[i]
		if item == k {
			q.shiftElementsAfter(i)
		}
	}
	if q.items[q.rear] == k {
		if q.rear-1 >= 0 {
			q.rear = q.rear - 1
		} else {
			q.rear = q.capacity - 1
		}
	}
}
func (q *FixedQueue[K]) shiftElementsAfter(i int) {
	j := (i + 1) % q.capacity
	for j != q.rear {
		q.items[i] = q.items[j]
		i = j
		j = (j + 1) % q.capacity
	}
	q.items[i] = q.items[j]
	q.rear = i

}
