package client

import (
	"time"
	"yadro-test/queue"
)

const (
	arrived = 1
	waiting = 2
	atTable = 3
)

type status struct {
	Status          int
	TableNumber     int // 0 if client
	OccupiedTableAt time.Time
}

type Store struct {
	clients map[string]status
	queue   *queue.FixedQueue[string]
}

func NewClientStore(queueLength int) *Store {
	return &Store{
		clients: make(map[string]status),
		queue:   queue.NewFixedQueue[string](queueLength),
	}
}
func (s *Store) Come(clientName string) error {
	if _, ok := s.clients[clientName]; ok {
		return ErrAlreadyExists
	}
	s.clients[clientName] = status{Status: arrived}
	return nil
}
func (s *Store) Wait(clientName string) error {
	st, ok := s.clients[clientName]
	if !ok {
		return ErrNotFound
	}
	if st.Status != arrived {
		return ErrShouldNotWait
	}
	if s.queue.Full() {
		return ErrQueueIsFull
	}
	st.Status = waiting
	s.clients[clientName] = st
	s.queue.PushBack(clientName)
	return nil
}
func (s *Store) TakeTable(clientName string, tableNum int, t time.Time) error {
	st, ok := s.clients[clientName]
	if !ok {
		return ErrNotFound
	}
	if st.Status == waiting {
		s.queue.Remove(clientName)
	}
	st.Status = atTable
	st.TableNumber = tableNum
	st.OccupiedTableAt = t
	s.clients[clientName] = st
	return nil
}
func (s *Store) Leave(clientName string, t time.Time) (tableNum int, dur time.Duration, err error) {
	st, ok := s.clients[clientName]
	if !ok {
		return 0, 0, ErrNotFound
	}
	switch st.Status {
	case arrived:
	case waiting:
		s.queue.Remove(clientName)
	case atTable:
		dur = t.Sub(st.OccupiedTableAt)
		tableNum = st.TableNumber
	}
	delete(s.clients, clientName)
	return tableNum, dur, nil
}
func (s *Store) FirstInLine() (string, error) {
	if s.queue.Empty() {
		return "", ErrQueueIsEmpty
	}
	clientName := s.queue.Front()
	s.queue.PopFront()
	return clientName, nil
}
func (s *Store) GetNames() []string {
	res := make([]string, 0, len(s.clients))
	for name := range s.clients {
		res = append(res, name)
	}
	return res
}
