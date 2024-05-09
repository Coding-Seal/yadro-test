package queue

import (
	"testing"
)

var testSample = []int{1, 2, 3, 4, 5, 6, 7, 8, 9}

func TestFixedQueue_Seq(t *testing.T) {
	q := NewFixedQueue[int](len(testSample))
	for _, i := range testSample {
		q.PushBack(i)
	}
	for _, i := range testSample {
		if q.Front() != i {
			t.Errorf("expect %d got %d", i, q.Front())
		}
		q.PopFront()
	}
}
func TestFixedQueue_Part(t *testing.T) {
	q := NewFixedQueue[int](len(testSample))

	for i := 0; i < len(testSample)/2; i++ {
		q.PushBack(testSample[i])
	}
	for i := 0; i < len(testSample)/2; i++ {
		if q.Front() != testSample[i] {
			t.Errorf("expect %d got %d", testSample[i], q.Front())
		}
		q.PopFront()
	}
	for i := len(testSample); i < len(testSample); i++ {
		q.PushBack(testSample[i])
	}
	for i := len(testSample); i < len(testSample); i++ {
		if q.Front() != testSample[i] {
			t.Errorf("expect %d got %d", testSample[i], q.Front())
		}
		q.PopFront()
	}
}
func TestFixedQueue_Remove_Middle(t *testing.T) {
	q := NewFixedQueue[int](len(testSample))
	for _, i := range testSample {
		q.PushBack(i)
	}
	q.Remove(testSample[len(testSample)/2])
	for _, i := range testSample {
		if i == testSample[len(testSample)/2] {
			continue
		}
		if q.Front() != i {
			t.Errorf("expect %d got %d", i, q.Front())
		}
		q.PopFront()
	}
}
func TestFixedQueue_Remove_Rear(t *testing.T) {
	q := NewFixedQueue[int](len(testSample))
	for _, i := range testSample {
		q.PushBack(i)
	}
	q.Remove(testSample[len(testSample)-1])
	for _, i := range testSample {
		if i == testSample[len(testSample)-1] {
			continue
		}
		if q.Front() != i {
			t.Errorf("expect %d got %d", i, q.Front())
		}
		q.PopFront()
	}
}
