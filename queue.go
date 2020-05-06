package queue

import (
	"sync"
)

//TODO:
// * Check signal
// * Check indexing

//---------------------------------------------------------------------------------------------------
// Queue Interface
//---------------------------------------------------------------------------------------------------

//Queue provides a finite fifo queue interface that provides a signal for event driven code
type Queue interface {
	//Close is provides cleanup
	Close()
	//GetSignal returns a signal channel that can be monitored if something is enqueue (event driven)
	GetSignal() (signal <-chan struct{})
	//Dequeue dequeues a single (last) element
	Dequeue() (element interface{}, underflow bool)
	//DequeueMultiple dequeues a multiple (last) element
	DequeueMultiple(count int) (elements []interface{}, underflow bool)
	//Enqueue enqueues a single element
	Enqueue(element interface{}) (overflow bool)
	//EnqueueMultiple enqueues multiple element
	EnqueueMultiple(elements []interface{}) (overflow bool)
	//Flush will flush the queue and return all elements that were in it
	Flush() (elements []interface{})
	//Peek allows for peeking at all elements in queue
	Peek() (elements []interface{}, empty bool)
	//PeekHead allows for peek at last element
	PeekHead() (element interface{}, empty bool)
	//PeekTail allows for peek at first element
	PeekTail() (element interface{}, empty bool)
	//GetSize will return the size (max size of the queue)
	GetSize() (size int)
	//GetLength will return the current length of the queue
	GetLength() (len int)
}

//---------------------------------------------------------------------------------------------------
// queue (pointer implmentation of the Queue Interface)
//---------------------------------------------------------------------------------------------------

//queue provides a pointer implementation of Queue
type queue struct {
	sync.Mutex
	data   []interface{} //the elements inside the queue
	size   int           //the max size of the queue
	index  int           //current index of the queue (this is also avoids len())
	signal chan struct{} //signal to notify that element has been enqueued
}

//NewQueue returns a new queue
func NewQueue(size int) interface {
	Queue
} {
	//Check if size is valid
	if size <= 0 {
		size = DefaultSize
	}
	data := make([]interface{}, size)
	signal := make(chan struct{})
	return &queue{
		data:   data,
		size:   size,
		index:  -1,
		signal: signal,
	}
}

//ensure the implementation
var _ Queue = &queue{}

//Close is provides cleanup
func (q *queue) Close() {
	q.Lock()
	defer q.Unlock()
	//cleanup
	q.signal, q.data = nil, nil
	q.size, q.index = 0, 0
	return
}

//GetSignal returns a signal channel that can be monitored if something is enqueue (event driven)
func (q *queue) GetSignal() (signal <-chan struct{}) {
	q.Lock()
	defer q.Unlock()
	//signal
	signal = q.signal
	return
}

//Dequeue dequeues a single (last) element
func (q *queue) Dequeue() (element interface{}, underflow bool) {
	q.Lock()
	defer q.Unlock()
	//dequeue
	underflow, element = q.dequeue()

	return
}

//DequeueMultiple dequeues a multiple (last) element
func (q *queue) DequeueMultiple(count int) (elements []interface{}, underflow bool) {
	q.Lock()
	defer q.Unlock()
	//dequeue based on count
	for i := 0; i < count; i++ {
		var element interface{}
		//dequeue
		underflow, element = q.dequeue()
		// check for underflow
		if underflow {
			return
		}
		elements = append(elements, element)
	}
	return
}

//Enqueue enqueues a single element
func (q *queue) Enqueue(element interface{}) (overflow bool) {
	q.Lock()
	defer q.Unlock()
	//enqueue
	overflow = q.enqueue(element)
	//trigger signal
	q.triggerSignal()
	return
}

//EnqueueMultiple enqueues multiple element
func (q *queue) EnqueueMultiple(elements []interface{}) (overflow bool) {
	q.Lock()
	defer q.Unlock()
	//enqueue each element
	for _, element := range elements {
		if overflow = q.enqueue(element); overflow {
			//overflow
			return
		}
		q.triggerSignal()
	}
	return
}

//Flush will flush the queue and return all elements that were in it
func (q *queue) Flush() (elements []interface{}) {
	q.Lock()
	defer q.Unlock()
	//flush
	elements = q.flush()
	return
}

//TODO: Check the empty case
func (q *queue) Peek() (elements []interface{}, empty bool) {
	q.Lock()
	defer q.Unlock()
	//Check if queue is empty
	if q.checkifEmpty() {
		empty = true
		return
	}
	//Peak by reversing data
	for i := q.index; i > -1; i-- {
		elements = append(elements, q.data[i])
	}
	return
}

//PeekHead allows for peek at last element
func (q *queue) PeekHead() (element interface{}, empty bool) {
	q.Lock()
	defer q.Unlock()
	//Check if queue is empty
	if q.checkifEmpty() {
		empty = true
		return
	}
	//Peak head
	element = q.data[q.index]
	return
}

//PeekTail allows for peek at first element
func (q *queue) PeekTail() (element interface{}, empty bool) {
	q.Lock()
	defer q.Unlock()
	//Check if queue is empty
	if q.checkifEmpty() {
		empty = true
		return
	}
	//Peak tail
	element = q.data[0]
	return
}

//GetSize will return the size (max size of the queue)
func (q *queue) GetSize() (size int) {
	q.Lock()
	defer q.Unlock()
	size = q.size
	return
}

//GetLength will return the current length of the queue
func (q *queue) GetLength() (len int) {
	q.Lock()
	defer q.Unlock()
	len = q.index + 1
	return
}

//---------------------------------------------------------------------------------------------------
// private queue methods
//---------------------------------------------------------------------------------------------------

//triggerSignal will send the signal that elemnent(s) have been enqueued
func (q *queue) triggerSignal() {
	select {
	case q.signal <- struct{}{}:
	default:
		//WHOAOAOAOAOA
	}
}

//shift shifts the elements up by one
func (q *queue) shift() {
	l := len(q.data)
	q.data = append(q.data[l-1:], q.data[:l-1]...)
	return
}

//checkifEmpty will check if the queue is empty
func (q *queue) checkifEmpty() (empty bool) {
	empty = q.index < 0
	return
}

//checkifFull will check if the queue is full
func (q *queue) checkifFull() (full bool) {
	full = q.index+1 >= q.size
	return
}

//enqueue performs the enqueue logic
func (q *queue) enqueue(element interface{}) (overflow bool) {
	//Check if queue is full (overflow)
	if q.checkifFull() {
		overflow = true
		return
	}
	//shift data
	q.shift()
	//insert element (use zero)
	q.data[0] = element
	//increment index
	q.index = q.index + 1
	return
}

//dequeue performs the dequeue logic
func (q *queue) dequeue() (underflow bool, element interface{}) {
	//Check if queue is empty (underflow)
	if q.checkifEmpty() {
		underflow = true
		return
	}
	//Get the last value in
	element = q.data[q.index-1]
	//clear element from data
	q.data[q.index-1] = nil
	//decrement index
	q.index = q.index - 1
	return
}

//flush will flush the queue and return all elements
func (q *queue) flush() (elements []interface{}) {
	elements = q.data
	//reset index
	q.index = -1
	//clear data
	q.data = make([]interface{}, q.size)
	return
}
