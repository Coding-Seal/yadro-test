package client

import "errors"

var (
	ErrAlreadyExists = errors.New("client already exists")
	ErrNotFound      = errors.New("client not found")
	ErrShouldNotWait = errors.New("client should not wait")
	ErrQueueIsFull   = errors.New("client queue is full")
	ErrQueueIsEmpty  = errors.New("client queue is empty")
)
